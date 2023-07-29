package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

// GobEncode ..
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
