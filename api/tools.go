package api

import (
	"encoding/json"
	"net/http"
)

func marshalResponse(data any) ([]byte, error) {
	return json.Marshal(data)
}

func writeResponseString(rw http.ResponseWriter, statusCode int, payload string) {
	writeResponseBytes(rw, statusCode, []byte(payload))
}

func writeResponseBytes(rw http.ResponseWriter, statusCode int, payload []byte) {
	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(payload)
}
