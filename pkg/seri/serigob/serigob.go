package serigob

import (
	"bytes"
	"encoding/gob"
	"main/pkg/seri"
)

type MGob struct{}

func (g MGob) Encode(m seri.Message) ([]byte, error) {
	var result bytes.Buffer

	err := gob.NewEncoder(&result).Encode(m)

	return result.Bytes(), err
}

func (g MGob) Decode(payload []byte) (seri.Message, error) {
	var result seri.Message

	err := gob.NewDecoder(bytes.NewReader(payload)).Decode(&result)

	return result, err
}
