package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Message struct {
	ID      int64
	Payload string
}

func main() {
	var netBuffer bytes.Buffer

	enc := gob.NewEncoder(&netBuffer)

	if errEncode := enc.Encode(Message{
		ID:      1,
		Payload: "xxx",
	}); errEncode != nil {
		log.Fatal("error:", errEncode)
	}
}
