package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
	"BFTWithoutSignatures/variables"
	"bytes"
	"encoding/gob"
	"log"
)

var (
	// DEFAULT - The default value
	DEFAULT = []byte("")

	// BCAnswer - Channel to receive the answer from BC
	BCAnswer = make(map[int]chan uint)
)

// MultiValuedConsensus - The method that is called to initiate the MVC module
func MultiValuedConsensus(cid int, v []byte) {
	// Initialize variables
	init := make(map[int]types.MvcMessage)
	vect := make(map[int]types.MvcMessage)

	BCAnswer[cid] = make(chan uint, 1)

	/* ----------------------------------- Task 1 ----------------------------------- */
	go func() {
		init[variables.ID] = types.NewMvcMessage(cid, "INIT", v, nil)
		broadcastMVC(init[variables.ID]) // Replace with RB if needed

		// Wait until at least (n-f) INIT messages
		for {
			if len(init) >= (variables.N - variables.F) {
				break
			}
		}

		// Fill vector with values received in init else DEFAULT
		vector := fillVector(init)

		w := calculateW(vector)
		logger.OutLogger.Print(cid, ".MVC:\n\tinit-", init, "\n\tvector-", vector, " --> ", w, "\n")

		vect[variables.ID] = types.NewMvcMessage(cid, "VECT", w, vector)
		broadcastMVC(vect[variables.ID])

		// Wait until at least (n-f) valid VECT messages
		for {
			if len(vect) >= (variables.N - variables.F) {
				break
			}
		}

		// Fill vectorW with values received in vect else DEFAULT
		vectorW := fillVector(vect)

		bVal := calculateBinaryValue(vectorW)
		logger.OutLogger.Print(cid, ".MVC:\n\tvect-", vect, "\n\tvector-", vectorW, " --> ", bVal, "\n")

		go BinaryConsensus(cid, bVal)
		c := <-BCAnswer[cid]

		if c == 0 {
			// Probably i will put the result in a chan or something like this to pass it to VC
			logger.OutLogger.Print(cid, ".MVC  decide-", DEFAULT, "\n")
			log.Println(variables.ID, "|", cid, ".MVC  decide-", DEFAULT)
			return
		}

		for {
			counter, dict := findOccurrences(fillVector(vect))
			for k, v := range counter {
				if v >= (variables.N - (2 * variables.F)) {
					// Probably i will put the result in a chan or something like this to pass it to VC
					logger.OutLogger.Print(cid, ".MVC  decide-", dict[k], "\n")
					log.Println(variables.ID, "|", cid, ".MVC  decide-", dict[k])
					return
				}
			}
		}
	}()

	/* ----------------------------------- Task 2 ----------------------------------- */
	go func() {
		if _, in := messenger.MvcChannel[cid]; !in {
			messenger.MvcChannel[cid] = make(chan struct {
				MvcMessage types.MvcMessage
				From       int
			})
		}

		for message := range messenger.MvcChannel[cid] {
			if message.MvcMessage.Type == "INIT" {
				// If already received from that process continue
				if _, in := init[message.From]; in {
					continue
				}
				init[message.From] = message.MvcMessage

			} else if message.MvcMessage.Type == "VECT" {
				// If already received from that process continue
				if _, in := vect[message.From]; in {
					continue
				}
				// If a valid VECT message add it in the map
				if checkVectValidity(message.MvcMessage, init) {
					vect[message.From] = message.MvcMessage
				}
			}
		}
	}()
}

func broadcastMVC(mvcMessage types.MvcMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(mvcMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	messenger.Broadcast(types.NewMessage(w.Bytes(), "MVC"))
}

/* -------------------------------- Helper Functions -------------------------------- */

func checkVectValidity(message types.MvcMessage, init map[int]types.MvcMessage) bool {
	for key, val := range message.Vector {
		if bytes.Compare(val, DEFAULT) == 0 {
			continue
		}
		if bytes.Compare(init[key].Value, val) != 0 {
			return false
		}
	}

	val := calculateW(message.Vector)
	if bytes.Compare(val, message.Value) != 0 {
		return false
	}
	return true
}

func fillVector(array map[int]types.MvcMessage) map[int][]byte {
	vector := make(map[int][]byte, variables.N)
	for i := 0; i < variables.N; i++ {
		if _, in := array[i]; in {
			vector[i] = array[i].Value
		} else {
			vector[i] = DEFAULT
		}
	}
	return vector
}

func calculateW(vector map[int][]byte) []byte {
	counter, dict := findOccurrences(vector)

	w := DEFAULT
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
		if bytes.Compare(val, DEFAULT) == 0 {
			continue
		}
		key := len(dict)
		for k, v := range dict {
			if bytes.Compare(v, val) == 0 {
				key = k
				break
			}
		}
		dict[key] = val
		counter[key] = counter[key] + 1
	}

	return counter, dict
}
