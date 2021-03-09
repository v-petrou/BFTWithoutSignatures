package types

import (
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/variables"
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

// Equals - Checks if client messages are equal
func (cm *ClientMessage) Equals(cmsg *ClientMessage) bool {
	return (cm.Client == cmsg.Client) && (cm.Value == cmsg.Value) &&
		(cm.TimeStamp.Equal(cmsg.TimeStamp)) && (cm.Ack == cmsg.Ack)
}

// ------------------------------------------------------------------------------------ //

// Reply struct
type Reply struct {
	TimeStamp time.Time
	Client    int
	ID        int
	Result    string
}

// NewReplyMessage - Creates a new Client message
func NewReplyMessage(client int, value rune) Reply {
	return Reply{TimeStamp: time.Now(), Client: client, ID: variables.ID, Result: string(value)}
}

// Equals - Checks if replies are equal
func (rep *Reply) Equals(reply *Reply) bool {
	return (rep.Client == reply.Client) && (rep.ID == reply.ID) && (rep.Result == reply.Result)
}
