package handlers

import "encoding/json"

type response struct {
	Message string `json:"message"`
}

func (r response) toByteArray() []byte {
	b, _ := json.Marshal(r)
	return b
}

func newResponse(msg string) []byte {
	return response{Message: msg}.toByteArray()
}
