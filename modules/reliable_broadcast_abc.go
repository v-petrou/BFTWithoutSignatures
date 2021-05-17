package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
	"BFTWithoutSignatures/variables"
	"bytes"
	"encoding/gob"
)

var (
	initial   = make(map[int]map[int][]byte, variables.N)         // instance, num
	echo      = make(map[int]map[int]map[int][]byte, variables.N) // instance, num, from
	ready     = make(map[int]map[int]map[int][]byte, variables.N)
	sentEcho  = make(map[int]map[int]bool, variables.N)
	sentReady = make(map[int]map[int]bool, variables.N)
	accepted  = make(map[int]map[int]bool, variables.N)
)

// SendRBInit - Sends the INIT message
func SendRBInit(num int, initVal []byte) {
	broadcastAll(types.NewRbMessage(num, "INIT", "ABC", variables.ID, initVal))
	broadcastAll(types.NewRbMessage(num, "ECHO", "ABC", variables.ID, initVal))

	initial[variables.ID][num] = initVal

	echo[variables.ID][num] = make(map[int][]byte)
	echo[variables.ID][num][variables.ID] = initVal
	sentEcho[variables.ID][num] = true

	ready[variables.ID][num] = make(map[int][]byte)
	ready[variables.ID][num][variables.ID] = initVal
	accepted[variables.ID][num] = true

	logger.OutLogger.Print(num, ".RB-ABC: INIT ", variables.ID, "\n")
	logger.OutLogger.Print(num, ".RB-ABC: INIT->ECHO ", variables.ID, "\n")
}

// ReliableBroadcastAbc - The method that is called to initiate the RB module for ABC
func ReliableBroadcastAbc() {
	for i := 0; i < variables.N; i++ {
		initial[i] = make(map[int][]byte)
		echo[i] = make(map[int]map[int][]byte)
		ready[i] = make(map[int]map[int][]byte)
		sentEcho[i] = make(map[int]bool)
		sentReady[i] = make(map[int]bool)
		accepted[i] = make(map[int]bool)
	}

	for message := range messenger.RbAbcChannel {
		tag := message.RbMessage.Tag
		instance := message.RbMessage.Process
		num := message.RbMessage.Rbid
		if tag == "INIT" {
			if _, in := initial[instance][num]; message.From != instance || in {
				continue // Only one value can be received from each process
			}
			if echo[instance][num] == nil {
				echo[instance][num] = make(map[int][]byte)
			}

			initial[instance][num] = message.RbMessage.Value
			broadcastAll(types.NewRbMessage(num, "ECHO", "ABC", instance, initial[instance][num]))

			echo[instance][num][variables.ID] = initial[instance][num]
			sentEcho[instance][num] = true
			logger.OutLogger.Print(num, ".RB-ABC: INIT->ECHO ", instance, "\n")

		} else if tag == "ECHO" {
			if _, in := echo[instance][num][message.From]; in {
				continue // Only one value can be received from each process
			}
			if echo[instance][num] == nil {
				echo[instance][num] = make(map[int][]byte)
			}
			if ready[instance][num] == nil {
				ready[instance][num] = make(map[int][]byte)
			}

			echo[instance][num][message.From] = message.RbMessage.Value

			counter, dict := CountMessages(echo[instance][num])
			for k, v := range counter {
				if v >= ((variables.N+variables.F)/2) && !sentEcho[instance][num] { // Step 1
					broadcastAll(types.NewRbMessage(num, "ECHO", "ABC", instance, dict[k]))

					echo[instance][num][variables.ID] = dict[k]
					sentEcho[instance][num] = true
					logger.OutLogger.Print(num, ".RB-ABC: ECHO->ECHO ", instance, "\n")

				} else if v >= ((variables.N+variables.F)/2) && !sentReady[instance][num] { // Step 2
					broadcastAll(types.NewRbMessage(num, "READY", "ABC", instance, dict[k]))

					ready[instance][num][variables.ID] = dict[k]
					sentReady[instance][num] = true
					logger.OutLogger.Print(num, ".RB-ABC: ECHO->READY ", instance, "\n")
				}
			}

		} else if tag == "READY" {
			if _, in := ready[instance][num][message.From]; in {
				continue // Only one value can be received from each process
			}
			if echo[instance][num] == nil {
				echo[instance][num] = make(map[int][]byte)
			}
			if ready[instance][num] == nil {
				ready[instance][num] = make(map[int][]byte)
			}

			ready[instance][num][message.From] = message.RbMessage.Value

			counter, dict := CountMessages(ready[instance][num])
			for k, v := range counter {
				if v >= ((2*variables.F)+1) && !accepted[instance][num] { // Step 3 - Accept v
					go messenger.HandleMessage(dict[k])
					accepted[instance][num] = true
					logger.OutLogger.Print(num, ".RB-ABC: accept-", instance, "\n")

				} else if v >= (variables.F+1) && !sentEcho[instance][num] { // Step 1
					broadcastAll(types.NewRbMessage(num, "ECHO", "ABC", instance, dict[k]))

					echo[instance][num][variables.ID] = dict[k]
					sentEcho[instance][num] = true
					logger.OutLogger.Print(num, ".RB-ABC: READY->ECHO ", instance, "\n")

				} else if v >= (variables.F+1) && !sentReady[instance][num] { // Step 2
					broadcastAll(types.NewRbMessage(num, "READY", "ABC", instance, dict[k]))

					ready[instance][num][variables.ID] = dict[k]
					sentReady[instance][num] = true
					logger.OutLogger.Print(num, ".RB-ABC: READY->READY ", instance, "\n")
				}
			}
		}
	}
}

func broadcastAll(rbMessage types.RbMessage) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(rbMessage)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	message := types.NewMessage(w.Bytes(), "RB_ABC")
	messenger.Broadcast(message)
}
