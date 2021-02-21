package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
	"BFTWithoutSignatures/variables"
	"bytes"
	"encoding/gob"
	"sync"
)

var (
	binValues = make(map[int][]uint)
	mutex     = sync.RWMutex{}
)

// BinaryConsensus - The method that is called to initiate the BC module
func BinaryConsensus(bcid int, initVal uint) {
	est := initVal
	for round := 1; ; round++ {
		id := ComputeUniqueIdentifier(bcid, round)
		logger.OutLogger.Print(id, ".BC: bcid-", bcid, "  round-", round, "\n")

		// BV_broadcast of the est value of the round
		go BvBroadcast(id, est)

		// Wait until not empty binValues
		for {
			mutex.Lock()
			if len(binValues[id]) != 0 {
				mutex.Unlock()
				break
			}
			mutex.Unlock()
		}

		broadcast("AUX", types.NewBcMessage(id, binValues[id][0]))

		// START Variables initialization
		values := make([]uint, 0)
		rec := make(map[int]uint)
		rec[variables.ID] = binValues[id][0]

		count := make(map[uint]int, 2)
		count[0], count[1] = 0, 0
		count[rec[variables.ID]]++

		if _, in := messenger.BcChannel[id]; !in {
			messenger.BcChannel[id] = make(chan struct {
				BcMessage types.BcMessage
				From      int
			})
		}
		// END Variables initialization

		for message := range messenger.BcChannel[id] {
			if _, in := rec[message.From]; in {
				continue // Only one value can be received from each process
			}

			rec[message.From] = message.BcMessage.Value
			count[message.BcMessage.Value]++

			// Wait until (n-t) AUX messages with the same value v
			if (count[0] >= (variables.N-variables.F) && inList(0, binValues[id])) &&
				(count[1] >= (variables.N-variables.F) && inList(1, binValues[id])) {
				values = append(values, 0)
				values = append(values, 1)
			} else if count[0] >= (variables.N-variables.F) && inList(0, binValues[id]) {
				values = append(values, 0)
			} else if count[1] >= (variables.N-variables.F) && inList(1, binValues[id]) {
				values = append(values, 1)
			}

			if len(values) != 0 {
				coin := random(id)
				logger.OutLogger.Print(id, ".BC:  vals-", values, "  coin-", coin, "\n")

				if len(values) == 2 {
					est = coin
				} else if len(values) == 1 && values[0] == coin {
					logger.OutLogger.Print(id, ".BC:  decide-", values[0], "\n")
					decide(bcid, values[0])
					return
				} else if len(values) == 1 && values[0] != coin {
					est = values[0]
				}

				break
			}
		}
	}

}

// BvBroadcast - Implements the BV_broadcast functionality
func BvBroadcast(identifier int, initVal uint) {
	// START variables initialization
	broadcasted := make(map[uint]bool, 2)
	broadcasted[0], broadcasted[1] = false, false

	received := make(map[int]int, (variables.N - 1))
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}
		received[i] = 0
	}

	counter := make(map[uint]int, 2)
	counter[0], counter[1] = 0, 0
	counter[initVal]++

	mutex.Lock()
	binValues[identifier] = make([]uint, 0, 2)
	mutex.Unlock()
	// END variables initialization

	// Broadcast initial value
	broadcast("EST", types.NewBcMessage(identifier, initVal))
	broadcasted[initVal] = true

	if _, in := messenger.BvbChannel[identifier]; !in {
		messenger.BvbChannel[identifier] = make(chan struct {
			BcMessage types.BcMessage
			From      int
		})
	}

	for message := range messenger.BvbChannel[identifier] {
		tag := message.BcMessage.Tag
		val := message.BcMessage.Value
		if received[message.From] < 2 { // Max 2 msgs can be accepted from other servers
			received[message.From]++
			counter[val]++
		}

		if counter[val] >= (variables.F+1) && !broadcasted[val] {
			broadcast("EST", types.NewBcMessage(tag, val))
			broadcasted[val] = true
		}

		if counter[val] >= ((2*variables.F)+1) && !inList(val, binValues[tag]) {
			mutex.Lock()
			binValues[tag] = append(binValues[tag], val)
			mutex.Unlock()
		}

		logger.OutLogger.Print(tag, ".BC:  bin_values-", binValues[tag], "\n")
	}

}

func broadcast(tag string, bcMessage types.BcMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(bcMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	var message types.Message
	if tag == "EST" {
		message = types.NewMessage(w.Bytes(), "BVB")
	} else if tag == "AUX" {
		message = types.NewMessage(w.Bytes(), "BC")
	} else {
		logger.ErrLogger.Fatal("Wrong message type in Binary Consensus Broadcast")
	}
	messenger.Broadcast(message)
}

// TODO: implement a more Byzantine Tolerant Common-Coin algorithm
func random(id int) uint {
	return uint(id % 2)
}

func decide(id int, value uint) {
	BCAnswer[id] <- value
}

/* -------------------------------- Helper Functions -------------------------------- */

// Checks if element a exists in list
func inList(a uint, list []uint) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ComputeUniqueIdentifier - Creates a unique num from (bcid,round) pair (Cantor's pairing func)
func ComputeUniqueIdentifier(a int, b int) int {
	res := (a * a) + (3 * a) + (2 * a * b) + b + (b * b)
	res = res / 2
	return res
}
