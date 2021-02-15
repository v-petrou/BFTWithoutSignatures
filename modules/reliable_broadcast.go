package modules

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/types"
	"bytes"
	"encoding/gob"
)

// ReliableBroadcast -
func ReliableBroadcast(rbid int, mType string, value []byte) {
	sendToAll(types.NewRbMessage(rbid, "init", mType, value))

	// START Variables initialization
	// initial := make(map[int][]byte, variables.N-1)
	// echo := make(map[int][]byte, variables.N-1)
	// ready := make(map[int][]byte, variables.N-1)

	// for message := range messenger.RbChannel[rbid] {
	// }
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
