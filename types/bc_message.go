package types

import (
	"BFTWithoutSignatures/logger"
	"bytes"
	"encoding/gob"
)

// BcMessage - Binary Consensus message struct
type BcMessage struct {
	Tag   int
	Value uint
}

// NewBcMessage - Creates a new Bc message
func NewBcMessage(tag int, value uint) BcMessage {
	return BcMessage{Tag: tag, Value: value}
}

// GobEncode - Binary Consensus message encoder
func (bcm BcMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(bcm.Tag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(bcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Binary Consensus message decoder
func (bcm *BcMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&bcm.Tag)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&bcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
