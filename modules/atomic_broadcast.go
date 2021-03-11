package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
	"BFTWithoutSignatures/variables"
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"log"
	"sort"
	"sync"
)

var (
	aid        int
	num        int
	received   map[int]map[int][]byte
	rDelivered [][]byte
	delMutex   = sync.RWMutex{}

	// VCAnswer - Channel to receive the answer from VC
	VCAnswer = make(map[int]chan map[int][]byte)
)

// InitiateAtomicBroadcast - The method that is called to initiate the ABC module
func InitiateAtomicBroadcast() {
	aid = 0
	num = 0
	received = make(map[int]map[int][]byte, (variables.N - 1))
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}
		received[i] = make(map[int][]byte)
	}
	rDelivered = make([][]byte, 0)

	go ReliableBroadcastAbc()

	go abcTask1()
	go abcTask2()
}

// AtomicBroadcast - The method that is called to broadcast a new ABC value
func AtomicBroadcast(m []byte) {
	rbABC(num, types.NewAbcMessage(num, m))

	delMutex.Lock()
	rDelivered = append(rDelivered, m)
	delMutex.Unlock()

	num++
}

func abcTask1() {
	for {
		for { // Wait until not empty R_delivered
			delMutex.Lock()
			if len(rDelivered) != 0 {
				delMutex.Unlock()
				break
			}
			delMutex.Unlock()
		}

		// Build the vector with the hashes of messages in R_delivered
		delMutex.Lock()
		h := hashMessages(rDelivered)
		delMutex.Unlock()

		VCAnswer[aid] = make(chan map[int][]byte, 1)
		w := new(bytes.Buffer)
		err := gob.NewEncoder(w).Encode(h)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Print(aid, ".ABC hash-", h, " -> VC\n")

		// Call VC and retrieve the answer
		go VectorConsensus(aid, w.Bytes())
		vc := <-VCAnswer[aid]
		x := make(map[int][][]byte)

		for k, v := range vc {
			if len(v) == 0 {
				continue
			}
			var temp [][]byte
			r := bytes.NewBuffer(v)
			err = gob.NewDecoder(r).Decode(&temp)
			if err != nil {
				logger.ErrLogger.Fatal(err)
			}
			x[k] = temp
		}

		aDelivered := make([][]byte, 0)

		// Wait until messages with hash in at least f+1 cells in X are in R_delivered
		count, dict := countHashes(x)
		for k, v := range count {
			if v >= (variables.F + 1) {
				for {
					val, in := checkIfDelivered(rDelivered, dict[k])
					if in {
						aDelivered = append(aDelivered, val)
						break
					}
				}
			}
		}

		// Sort messages in aDeliver and then deliver them
		// TODO: probably put them in a channel
		sort.Slice(aDelivered, func(i, j int) bool {
			return bytes.Compare(aDelivered[i], aDelivered[j]) < 0
		})

		logger.OutLogger.Print(aid, ".ABC aDelivered-", aDelivered, "\n")
		log.Print(variables.ID, " | ", aid, ".ABC aDelivered-", aDelivered, "\n")

		// Remove from R_delivered the values that have been already delivered
		for _, b := range aDelivered {
			delMutex.Lock()
			for i, v := range rDelivered {
				if bytes.Compare(b, v) == 0 {
					rDelivered = append(rDelivered[:i], rDelivered[i+1:]...)
					break
				}
			}
			delMutex.Unlock()
		}

		logger.OutLogger.Print(aid, ".ABC len-", len(rDelivered), " -> aid++\n")
		aid++
	}
}

func abcTask2() {
	for message := range messenger.AbcChannel {
		if _, in := received[message.From][message.AbcMessage.Num]; in {
			continue // Only one value can be received from each process
		}
		received[message.From][message.AbcMessage.Num] = message.AbcMessage.Value

		delMutex.Lock()
		rDelivered = append(rDelivered, message.AbcMessage.Value)
		delMutex.Unlock()
	}
}

func rbABC(id int, abcMessage types.AbcMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(abcMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	msg := types.NewMessage(w.Bytes(), "ABC")
	w = new(bytes.Buffer)
	encoder = gob.NewEncoder(w)
	err = encoder.Encode(msg)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	SendRBInit(id, w.Bytes())
}

func hashMessages(rDelivered [][]byte) [][]byte {
	h := make([][]byte, 0)
	for _, v := range rDelivered {
		hasher := sha512.New()
		hasher.Write(v)
		h = append(h, hasher.Sum(nil))
	}
	return h
}

func countHashes(vector map[int][][]byte) (map[int]int, map[int][]byte) {
	counter := make(map[int]int)
	dict := make(map[int][]byte)

	for _, val := range vector {
		for _, x := range val {
			key := len(dict)
			for k, v := range dict {
				if bytes.Compare(v, x) == 0 {
					key = k
					break
				}
			}
			dict[key] = x
			counter[key] = counter[key] + 1
		}
	}

	return counter, dict
}

func checkIfDelivered(rDelivered [][]byte, val []byte) ([]byte, bool) {
	for _, v := range rDelivered {
		hasher := sha512.New()
		hasher.Write(v)
		if bytes.Compare(hasher.Sum(nil), val) == 0 {
			return v, true
		}
	}
	return nil, false
}
