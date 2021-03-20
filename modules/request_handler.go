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
	// Delivered - Channel to receive delivered messages from ABC
	Delivered = make(chan struct {
		Id    int
		Value [][]byte
	})

	// Array - The array that has to be in consensus
	Array = make([]rune, 0)
)

// RequestHandler - The module that handles requests received from clients and replies to them
func RequestHandler() {
	// Accepts the requests from the clients and calls ABC
	go func() {
		for message := range messenger.RequestChannel {
			AtomicBroadcast(message)
		}
	}()

	// Gets the delivered result from ABC, appends it in the Array and replies to the client
	go func() {
		for message := range Delivered {
			for _, v := range message.Value {
				var m types.ClientMessage
				buffer := bytes.NewBuffer(v)
				decoder := gob.NewDecoder(buffer)
				err := decoder.Decode(&m)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				Array = append(Array, m.Value)

				go messenger.ReplyClient(types.NewReplyMessage(m.Num), m.Cid)
			}

			logger.OutLogger.Print(message.Id, ".REQH: array-", Array, "\n")
			log.Print(variables.ID, " | ", message.Id, ".REQH: array-", Array, "\n")
		}
	}()
}

// start := time.Now()
// log.Println(variables.ID, "|", "time-", time.Since(start))
