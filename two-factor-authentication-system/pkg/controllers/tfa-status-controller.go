package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/montaser55/two-factor-authentication-service/pkg/models"
	"github.com/montaser55/two-factor-authentication-service/pkg/models/requests"
	"github.com/montaser55/two-factor-authentication-service/pkg/utils"
)

const headerName = "Content-Type"
const headerValue = "pkglication/json"

func CheckTfa(w http.ResponseWriter, r *http.Request) {
	userTfaInfos := models.GetAllUserTfaInfos()
	res, _ := json.Marshal(userTfaInfos)
	w.Header().Set(headerName, headerValue)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GenerateSecret(w http.ResponseWriter, r *http.Request) {
	request := &requests.GenerateSecretRequest{}
	utils.ParseBody(r, request)
	log.Printf("%v", request)
	validateGenerateSecretRequest(*request)

}

func validateGenerateSecretRequest(request requests.GenerateSecretRequest) {
	log.Printf("%v", request.UserId)
}

func PostTfa(w http.ResponseWriter, r *http.Request) {
	UserTfaInfo := &models.UserTfaInfo{}
	utils.ParseBody(r, UserTfaInfo)
	newUserTfaInfo := UserTfaInfo.CreateUserTfaInfo()
	res, _ := json.Marshal(newUserTfaInfo)
	w.Header().Set(headerName, headerValue)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
