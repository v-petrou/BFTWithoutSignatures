package types

import (
	"BFTWithoutSignatures/logger"
	"bytes"
	"encoding/gob"
)

// RbMessage - Reliable Broadcast message struct
type RbMessage struct {
	Rbid    int
	Tag     string // (init, echo, ready)
	Type    string // What is the type of the value (MVC, VC, ABC)
	Process int
	Value   []byte
}

// NewRbMessage - Creates a new Rb message
func NewRbMessage(rbid int, tag string, t string, process int, value []byte) RbMessage {
	return RbMessage{Rbid: rbid, Tag: tag, Type: t, Process: process, Value: value}
}

// GobEncode - Reliable Broadcast message encoder
func (rbm RbMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(rbm.Rbid)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rbm.Tag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rbm.Type)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rbm.Process)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(rbm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Reliable Broadcast message decoder
func (rbm *RbMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&rbm.Rbid)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rbm.Tag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rbm.Type)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rbm.Process)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&rbm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
