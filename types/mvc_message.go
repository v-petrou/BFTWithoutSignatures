package types

import (
	"BFTWithoutSignatures/logger"
	"bytes"
	"encoding/gob"
)

// MvcMessage - Multi-valued consensus message struct
type MvcMessage struct {
	Cid    int
	Type   string
	Value  []byte
	Vector map[int][]byte
}

// NewMvcMessage - Creates a new Mvc message
func NewMvcMessage(cid int, t string, value []byte, vector map[int][]byte) MvcMessage {
	return MvcMessage{Cid: cid, Type: t, Value: value, Vector: vector}
}

// GobEncode - Multi-valued consensus message encoder
func (mvcm MvcMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(mvcm.Cid)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(mvcm.Type)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(mvcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(mvcm.Vector)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Multi-valued consensus message decoder
func (mvcm *MvcMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&mvcm.Cid)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&mvcm.Type)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&mvcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&mvcm.Vector)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
