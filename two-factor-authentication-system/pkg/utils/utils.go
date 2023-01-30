package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func ParseBody(r *http.Request, x any) {
	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, x); err != nil {
			return
		}
	}
}

func ParseBodyReusable(r *http.Request, x any) {
	if body, err := io.ReadAll(r.Body); err == nil {
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if err := json.Unmarshal(body, x); err != nil {
			return
		}
	}
}

func GenerateReferenceId() string {
	return uuid.NewV4().String()
}
