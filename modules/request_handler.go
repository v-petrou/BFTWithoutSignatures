package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/variables"
	"log"
	"time"
)

var (
	// Delivered - Channel to receive delivered messages from ABC
	Delivered = make(chan [][]byte)
)

// RequestHandler - The module that handles requests received from clients and replies to them
func RequestHandler() {
	start := time.Now()

	go func() {
		for message := range messenger.RequestChannel {
			messenger.ReplyClient([]byte("ACK"), message.Client)
			//AtomicBroadcast([]byte(string(message.Value)))

			log.Println(variables.ID, "|", "time-", time.Since(start))
		}
	}()

	go func() {
		for message := range Delivered {
			for _, v := range message {
				logger.OutLogger.Println(v)
			}
		}
	}()
}
