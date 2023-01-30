package gateways

import (
	"bytes"
	"encoding/json"
	"github.com/montaser55/two-factor-authentication-service/pkg/models/requests"
	"io"
	"log"
	"net/http"
)

var wmsBaseUrl = "http://3.36.21.109:10270/wallet-management-service-1.0"

func VerifyPin(request requests.CredentialRequest) {
	jsonValue, _ := json.Marshal(request)
	resp, err := http.Post(wmsBaseUrl+"/api/user-pin/verify-pin", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Panic("Could not Verify Pin")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		m := map[string]string{}
		err := json.Unmarshal(body, &m)
		if err != nil {
			log.Panic("Could not Verify Pin")
			return
		}
		log.Println(m)
	}

}
