package types

import (
	"BFTWithoutSignatures/logger"
	"bytes"
	"encoding/gob"
)

// VcMessage - Vector consensus message struct
type VcMessage struct {
	Vcid  int
	Value []byte
}

// NewVcMessage - Creates a new VC message
func NewVcMessage(id int, value []byte) VcMessage {
	return VcMessage{Vcid: id, Value: value}
}

// GobEncode - Vector consensus message encoder
func (vcm VcMessage) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(vcm.Vcid)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = encoder.Encode(vcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return w.Bytes(), nil
}

// GobDecode - Vector consensus message decoder
func (vcm *VcMessage) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&vcm.Vcid)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	err = decoder.Decode(&vcm.Value)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return nil
}
