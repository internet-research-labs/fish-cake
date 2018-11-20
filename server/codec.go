package server

import (
	"encoding/json"
)

type WireMessage struct {
	Type string      `json:"type"`
	Blob interface{} `json:"blob"`
}

func EncodeWireMessage(t string, decoded interface{}) []byte {
	wrapped := WireMessage{
		Type: t,
		Blob: decoded,
	}

	encoded, err := json.Marshal(wrapped)

	if err != nil {
		return encoded
	}
	return encoded
}
