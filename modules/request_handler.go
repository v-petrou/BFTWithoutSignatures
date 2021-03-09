package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
)

var (
	// Delivered - Channel to receive delivered messages from ABC
	Delivered = make(chan [][]byte)
)

// RequestHandler - The module that handles requests received from clients and replies to them
func RequestHandler() {
	go func() {
		for message := range messenger.RequestChannel {
			messenger.ReplyClient([]byte("ACK"), message.Client)
			//AtomicBroadcast([]byte(string(message.Value)))
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
