package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
	"BFTWithoutSignatures/variables"
	"bytes"
	"encoding/gob"
	"encoding/json"
)

var (
	// MVCAnswer - Channel to receive the answer from MVC
	MVCAnswer = make(map[int]chan []byte)
)

// VectorConsensus - The method that is called to initiate the VC module
func VectorConsensus(vcid int, initVal []byte) {
	// START Variables initialization
	received := make(map[int][]byte)
	received[variables.ID] = initVal

	if _, in := messenger.VcChannel[vcid]; !in {
		messenger.VcChannel[vcid] = make(chan struct {
			VcMessage types.VcMessage
			From      int
		})
	}
	// END Variables initialization

	// Reliable Broadcast the given value
	rbVC(vcid, types.NewVcMessage(vcid, initVal))

	for round := 0; ; round++ {
		for message := range messenger.VcChannel[vcid] {
			if _, in := received[message.From]; in {
				continue // Only one value can be received from each process
			}
			received[message.From] = message.VcMessage.Value

			// Wait until at least ((n-f)+r) INIT messages
			if len(received) == ((variables.N - variables.F) + round) {
				break
			}
		}

		// Built the vector with the values received
		vector := make(map[int][]byte, variables.N)
		for i := 0; i < variables.N; i++ {
			if _, in := received[i]; in {
				vector[i] = received[i]
			} else {
				vector[i] = variables.DEFAULT
			}
		}

		// Compute the MVC identifier and convert vector to bytes
		mvcid := ComputeUniqueIdentifier(vcid, round)
		MVCAnswer[mvcid] = make(chan []byte, 1)
		w, err := json.Marshal(vector)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		logger.OutLogger.Print(vcid, ".VC: len-", len(received), " vector-", vector, " --> MVC\n")

		go MultiValuedConsensus(mvcid, w)
		v := <-MVCAnswer[mvcid]

		// If MVC answer != DEFAULT, then decide this value, else go to next the round
		if !bytes.Equal(v, variables.DEFAULT) {
			var vect map[int][]byte
			err = json.Unmarshal(v, &vect)
			if err != nil {
				logger.ErrLogger.Fatal(err)
			}

			logger.OutLogger.Print(vcid, ".VC: decide-", vect, "\n")
			VCAnswer[vcid] <- vect
			return
		}
	}
}

func rbVC(id int, vcMessage types.VcMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(vcMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	msg := types.NewMessage(w.Bytes(), "VC")
	w = new(bytes.Buffer)
	encoder = gob.NewEncoder(w)
	err = encoder.Encode(msg)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	go ReliableBroadcast(id, "VC", w.Bytes())
}
