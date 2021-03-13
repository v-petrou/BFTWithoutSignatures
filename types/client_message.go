package types

import (
	"BFTWithoutSignatures/logger"
	"bytes"
	"encoding/gob"
	"time"
)

// ClientMessage - Client message struct
type ClientMessage struct {
	Client    int
	TimeStamp time.Time
	Value     rune
	Ack       bool
}

// NewClientMessage - Creates a new Client message
func NewClientMessage(client int, value rune) ClientMessage {
	return ClientMessage{Client: client, TimeStamp: time.Now(), Value: value, Ack: false}
}

// GobEncode - Client message encoder
func (cm ClientMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(cm.Client)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(cm.TimeStamp)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(cm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(cm.Ack)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Client message decoder
func (cm *ClientMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&cm.Client)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&cm.TimeStamp)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&cm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&cm.Ack)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
