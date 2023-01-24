package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func ParseBody(r *http.Request, x any) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func GenerateReferenceId() string {
	return uuid.NewV4().String()
}
