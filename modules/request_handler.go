package modules

import (
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
)

var (
	// Delivered - Channel to receive delivered messages from ABC
	Delivered = make(chan types.Reply)
)

// RequestHandler - The module that handles requests received from clients and replies to them
func RequestHandler() {
	go func() {
		for message := range messenger.RequestChannel {
			AtomicBroadcast([]byte(string(message.Value)))
		}
	}()

	go func() {
		for message := range Delivered {
			messenger.BroadcastClients(message)
		}
	}()
}

// start := time.Now()
// log.Println(variables.ID, "|", "time-", time.Since(start))
