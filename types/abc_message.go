package types

import (
	"BFTWithoutSignatures/logger"
	"bytes"
	"encoding/gob"
)

// AbcMessage - Atomic Broadcast message struct
type AbcMessage struct {
	Num   int
	Value []byte
}

// NewAbcMessage - Creates a new ABC message
func NewAbcMessage(num int, value []byte) AbcMessage {
	return AbcMessage{Num: num, Value: value}
}

// GobEncode - Atomic Broadcast message encoder
func (abcm AbcMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(abcm.Num)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(abcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Atomic Broadcast message decoder
func (abcm *AbcMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&abcm.Num)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&abcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
