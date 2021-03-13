package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
	"BFTWithoutSignatures/variables"
	"bytes"
	"encoding/gob"
)

// ReliableBroadcast - The method that is called to initiate the RB module
func ReliableBroadcast(rbid int, mType string, initVal []byte) {
	// START Variables initialization
	initial := make(map[int][]byte, variables.N)
	echo := make(map[int]map[int][]byte, variables.N)
	ready := make(map[int]map[int][]byte, variables.N)
	sentEcho := make(map[int]bool, variables.N)
	sentReady := make(map[int]bool, variables.N)
	accepted := make(map[int]bool, variables.N)
	for i := 0; i < variables.N; i++ {
		echo[i] = make(map[int][]byte, variables.N)
		ready[i] = make(map[int][]byte, variables.N)
		sentEcho[i] = false
		sentReady[i] = false
		accepted[i] = false
	}

	if _, in := messenger.RbChannel[mType][rbid]; !in {
		messenger.RbChannel[mType][rbid] = make(chan struct {
			RbMessage types.RbMessage
			From      int
		})
	}
	// END Variables initialization

	// Step 0
	sendToAll(types.NewRbMessage(rbid, "INIT", mType, variables.ID, initVal))
	sendToAll(types.NewRbMessage(rbid, "ECHO", mType, variables.ID, initVal))

	initial[variables.ID] = initVal
	echo[variables.ID][variables.ID] = initVal
	sentEcho[variables.ID] = true
	ready[variables.ID][variables.ID] = initVal
	accepted[variables.ID] = true

	logger.OutLogger.Print(rbid, ".RB-", mType, ": INIT ", variables.ID, "\n")
	logger.OutLogger.Print(rbid, ".RB-", mType, ": INIT->ECHO ", variables.ID, "\n")

	for message := range messenger.RbChannel[mType][rbid] {
		tag := message.RbMessage.Tag
		instance := message.RbMessage.Process
		if tag == "INIT" {
			if _, in := initial[message.From]; message.From != instance || in {
				continue // Only one value can be received from each process
			}
			initial[instance] = message.RbMessage.Value
			sendToAll(types.NewRbMessage(rbid, "ECHO", mType, instance, initial[instance]))

			echo[instance][variables.ID] = initial[instance]
			sentEcho[instance] = true
			logger.OutLogger.Print(rbid, ".RB-", mType, ": INIT->ECHO ", instance, "\n")

		} else if tag == "ECHO" {
			if _, in := echo[instance][message.From]; in {
				continue // Only one value can be received from each process
			}
			echo[instance][message.From] = message.RbMessage.Value

			counter, dict := CountMessages(echo[instance])
			for k, v := range counter {
				if v >= ((variables.N+variables.F)/2) && !sentEcho[instance] { // Step 1
					sendToAll(types.NewRbMessage(rbid, "ECHO", mType, instance, dict[k]))

					echo[instance][variables.ID] = dict[k]
					sentEcho[instance] = true
					logger.OutLogger.Print(rbid, ".RB-", mType, ": ECHO->ECHO ", instance, "\n")

				} else if v >= ((variables.N+variables.F)/2) && !sentReady[instance] { // Step 2
					sendToAll(types.NewRbMessage(rbid, "READY", mType, instance, dict[k]))

					ready[instance][variables.ID] = dict[k]
					sentReady[instance] = true
					logger.OutLogger.Print(rbid, ".RB-", mType, ": ECHO->READY ", instance, "\n")
				}
			}

		} else if tag == "READY" {
			if _, in := ready[instance][message.From]; in {
				continue // Only one value can be received from each process
			}
			ready[instance][message.From] = message.RbMessage.Value

			counter, dict := CountMessages(ready[instance])
			for k, v := range counter {
				if v >= ((2*variables.F)+1) && !accepted[instance] { // Step 3 - Accept v
					go messenger.HandleMessage(dict[k])
					accepted[instance] = true
					logger.OutLogger.Print(rbid, ".RB-", mType, ": accept-", instance, "\n")

				} else if v >= (variables.F+1) && !sentEcho[instance] { // Step 1
					sendToAll(types.NewRbMessage(rbid, "ECHO", mType, instance, dict[k]))

					echo[instance][variables.ID] = dict[k]
					sentEcho[instance] = true
					logger.OutLogger.Print(rbid, ".RB-", mType, ": READY->ECHO ", instance, "\n")

				} else if v >= (variables.F+1) && !sentReady[instance] { // Step 2
					sendToAll(types.NewRbMessage(rbid, "READY", mType, instance, dict[k]))

					ready[instance][variables.ID] = dict[k]
					sentReady[instance] = true
					logger.OutLogger.Print(rbid, ".RB-", mType, ": READY->READY ", instance, "\n")
				}
			}
		}
	}
}

func sendToAll(rbMessage types.RbMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(rbMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	message := types.NewMessage(w.Bytes(), "RB")
	messenger.Broadcast(message)
}

// CountMessages - Counts the messages received from RB
func CountMessages(vector map[int][]byte) (map[int]int, map[int][]byte) {
	counter := make(map[int]int)
	dict := make(map[int][]byte)
	for _, val := range vector {
		key := len(dict)
		for k, v := range dict {
			if bytes.Equal(v, val) {
				key = k
				break
			}
		}
		dict[key] = val
		counter[key] = counter[key] + 1
	}

	return counter, dict
}
