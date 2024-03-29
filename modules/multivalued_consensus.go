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
	// BCAnswer - Channel to receive the answer from BC
	BCAnswer = make(map[int]chan uint)
)

// MultiValuedConsensus - The method that is called to initiate the MVC module
func MultiValuedConsensus(mvcid int, v []byte) {
	// START Variables initialization
	init := make(map[int][]byte)
	vect := make(map[int][]byte)
	initMutex := sync.RWMutex{}
	vectMutex := sync.RWMutex{}
	BCAnswer[mvcid] = make(chan uint, 1)

	if _, in := messenger.MvcChannel[mvcid]; !in {
		messenger.MvcChannel[mvcid] = make(chan struct {
			MvcMessage types.MvcMessage
			From       int
		})
	}
	// END Variables initialization

	/* ----------------------------------- Task 1 ----------------------------------- */
	go func() {
		init[variables.ID] = v
		rbMVC(ComputeUniqueIdentifier(mvcid, 1), types.NewMvcMessage(mvcid, "INIT", v, nil))

		for { // Wait until at least (n-f) INIT messages
			initMutex.Lock()
			if len(init) >= (variables.N - variables.F) {
				initMutex.Unlock()
				break
			}
			initMutex.Unlock()
		}

		// Fill vector with values in init else DEFAULT and calculate w value
		initMutex.Lock()
		vector := fillVector(init)
		initMutex.Unlock()
		w := calculateW(vector)
		logger.OutLogger.Print(mvcid, ".MVC: vector-", vector, " --> ", w, "\n")

		vect[variables.ID] = w
		rbMVC(ComputeUniqueIdentifier(mvcid, 2), types.NewMvcMessage(mvcid, "VECT", w, vector))

		for { // Wait until at least (n-f) valid VECT messages
			vectMutex.Lock()
			if len(vect) >= (variables.N - variables.F) {
				vectMutex.Unlock()
				break
			}
			vectMutex.Unlock()
		}

		// Fill vectorW with values in vect else DEFAULT and calculate BC input value
		vectMutex.Lock()
		vectorW := fillVector(vect)
		vectMutex.Unlock()
		bVal := calculateBinaryValue(vectorW)
		logger.OutLogger.Print(mvcid, ".MVC: vectorW-", vectorW, " --> ", bVal, "\n")

		go BinaryConsensus(mvcid, bVal)
		c := <-BCAnswer[mvcid]

		if c == 0 {
			logger.OutLogger.Print(mvcid, ".MVC: decide-", variables.DEFAULT, "\n")
			MVCAnswer[mvcid] <- variables.DEFAULT
			return
		}

		for {
			counter, dict := findOccurrences(fillVector(vect))
			for k, v := range counter {
				if v >= (variables.N - (2 * variables.F)) {
					logger.OutLogger.Print(mvcid, ".MVC: decide-", dict[k], "\n")
					MVCAnswer[mvcid] <- dict[k]
					return
				}
			}
		}
	}()

	/* ----------------------------------- Task 2 ----------------------------------- */
	go func() {
		for message := range messenger.MvcChannel[mvcid] {
			if message.MvcMessage.Type == "INIT" {
				if _, in := init[message.From]; in {
					continue // Only one value can be received from each process
				}
				initMutex.Lock()
				init[message.From] = message.MvcMessage.Value
				initMutex.Unlock()

			} else if message.MvcMessage.Type == "VECT" {
				if _, in := vect[message.From]; in {
					continue // Only one value can be received from each process
				}

				if checkVectValidity(message.MvcMessage, init) { // Accept only valid VECT msgs
					vectMutex.Lock()
					vect[message.From] = message.MvcMessage.Value
					vectMutex.Unlock()
				}
			}
		}
	}()
}

func rbMVC(id int, mvcMessage types.MvcMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(mvcMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	msg := types.NewMessage(w.Bytes(), "MVC")
	w = new(bytes.Buffer)
	encoder = gob.NewEncoder(w)
	err = encoder.Encode(msg)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	go ReliableBroadcast(id, "MVC", w.Bytes())
}

/* -------------------------------- Helper Functions -------------------------------- */

func checkVectValidity(message types.MvcMessage, init map[int][]byte) bool {
	for key, val := range message.Vector {
		if bytes.Equal(val, variables.DEFAULT) {
			continue
		}
		if !bytes.Equal(init[key], val) {
			return false
		}
	}

	val := calculateW(message.Vector)

	return bytes.Equal(val, message.Value)
}

func fillVector(array map[int][]byte) map[int][]byte {
	vector := make(map[int][]byte, variables.N)
	for i := 0; i < variables.N; i++ {
		if _, in := array[i]; in {
			vector[i] = array[i]
		} else {
			vector[i] = variables.DEFAULT
		}
	}
	return vector
}

func calculateW(vector map[int][]byte) []byte {
	counter, dict := findOccurrences(vector)

	w := variables.DEFAULT
	count := 0
	for k, v := range counter {
		if v >= (variables.N-(2*variables.F)) && count == 0 {
			w = dict[k]
			count = v
		} else if v >= (variables.N-(2*variables.F)) && v > count {
			w = dict[k]
			count = v
		} else if v >= (variables.N-(2*variables.F)) && v == count &&
			bytes.Compare(w, dict[k]) == -1 {
			w = dict[k]
		}
	}
	return w
}

func calculateBinaryValue(vector map[int][]byte) uint {
	counter, _ := findOccurrences(vector)

	if len(counter) > 1 {
		return 0
	}

	if counter[0] >= (variables.N - (2 * variables.F)) {
		return 1
	}
	return 0
}

func findOccurrences(vector map[int][]byte) (map[int]int, map[int][]byte) {
	counter := make(map[int]int)
	dict := make(map[int][]byte)
	for _, val := range vector {
		if bytes.Equal(val, variables.DEFAULT) {
			continue
		}
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
