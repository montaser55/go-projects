package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/montaser55/go-orm/pkg/utils"
	"github.com/montaser55/two-factor-authentication-service/pkg/models"
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

func PostTfa(w http.ResponseWriter, r *http.Request) {
	UserTfaInfo := &models.UserTfaInfo{}
	utils.ParseBody(r, UserTfaInfo)
	newUserTfaInfo := UserTfaInfo.CreateUserTfaInfo()
	res, _ := json.Marshal(newUserTfaInfo)
	w.Header().Set(headerName, headerValue)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
