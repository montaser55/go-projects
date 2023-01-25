package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/montaser55/two-factor-authentication-service/pkg/config"
	"github.com/montaser55/two-factor-authentication-service/pkg/models"
	"github.com/montaser55/two-factor-authentication-service/pkg/models/contexts"
	"github.com/montaser55/two-factor-authentication-service/pkg/models/requests"
	"github.com/montaser55/two-factor-authentication-service/pkg/models/responses"
	"github.com/montaser55/two-factor-authentication-service/pkg/utils"
	"github.com/montaser55/two-factor-authentication-service/pkg/utils/enums"
	"github.com/montaser55/two-factor-authentication-service/pkg/utils/gateways"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
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
	referenceId := utils.GenerateReferenceId()
	email := gateways.GetUserEmailInfo(request.UserId)
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "cryptrade",
		AccountName: email,
		Period:      uint(getExpiryTimeInSeconds(request.TfaChannelType)),
	})
	//config.RedisConnect()
	log.Printf("%v", config.GetRedisClient().Ping(config.GetContext()))
	redisClient := config.GetRedisClient()
	twoFactorAuthenticationContext := buildTwoFactorAuthenticationContext(request.UserId, key.Secret(), request.TfaChannelType)
	twoFactorAuthenticationContextJson, _ := json.Marshal(&twoFactorAuthenticationContext)
	redisClient.Set(config.GetContext(), referenceId, twoFactorAuthenticationContextJson, 0)

	generateSecretResponse := buildGenerateSecretResponse(*key, referenceId, *request)
	res, _ := json.Marshal(&generateSecretResponse)
	w.Header().Set(headerName, headerValue)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func EnableTfa(w http.ResponseWriter, r *http.Request) {
	request := &requests.TfaEnableRequest{}
	utils.ParseBody(r, request)
	validateTfaEnableRequest(*request)
	redisClient := config.GetRedisClient()
	twoFactorAuthenticationContext := &contexts.TwoFactorAuthenticationContext{}
	str, err := redisClient.Get(config.GetContext(), request.ReferenceId).Result()
	if err != nil {
		log.Panic("ReferenceId not found in redis")
	}
	json.Unmarshal([]byte(str), twoFactorAuthenticationContext)
	if request.TfaChannelType == enums.SMS {
		validateOtpExpiration(twoFactorAuthenticationContext.OtpGenerationTime, getExpiryTimeInSeconds(enums.SMS))
	}
	validate, _ := totp.ValidateCustom(request.Otp, twoFactorAuthenticationContext.SecretKey, time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(getExpiryTimeInSeconds(request.TfaChannelType)),
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	log.Printf("%v", validate)
	if !validate {
		log.Panic("Invalid Otp")
	}

	userTfaInfo := models.GetUserTfaInfoByUserId(request.UserId)
	if userTfaInfo == nil {
		userTfaInfoNew := buildUserTfaInfo(*twoFactorAuthenticationContext)
		models.CreateUserTfaInfo(userTfaInfoNew)
	} else {
		userTfaInfo.SecretKey = twoFactorAuthenticationContext.SecretKey
		userTfaInfo.App = twoFactorAuthenticationContext.TfaChannelType == enums.APP
		userTfaInfo.Sms = twoFactorAuthenticationContext.TfaChannelType == enums.SMS
		models.UpdateUserTfaInfo(userTfaInfo)
	}
}

func buildUserTfaInfo(twoFactorAuthenticationContext contexts.TwoFactorAuthenticationContext) *models.UserTfaInfo {
	userTfaInfo := &models.UserTfaInfo{}
	userTfaInfo.UserId = twoFactorAuthenticationContext.UserId
	userTfaInfo.App = twoFactorAuthenticationContext.TfaChannelType == enums.APP
	userTfaInfo.Sms = twoFactorAuthenticationContext.TfaChannelType == enums.SMS
	userTfaInfo.TryCounter = 0
	userTfaInfo.SecretKey = twoFactorAuthenticationContext.SecretKey
	return userTfaInfo
}

func validateOtpExpiration(otpGenerationTime int64, otpExpiryInSeconds int) {
	timeDifferenceInSeconds := int(time.Now().UnixMilli()-otpGenerationTime) / 1000
	if timeDifferenceInSeconds > otpExpiryInSeconds {
		log.Panic("Otp Expired")
	}
}

func validateTfaEnableRequest(request requests.TfaEnableRequest) {
	if request.UserId == 0 {
		log.Panic("UserId is not Provided")
	}
	if err := request.TfaChannelType.IsValid(); err != nil {
		log.Panic("Invalid Tfa Channel")
	}
	if request.ReferenceId == "" {
		log.Panic("ReferenceId not Provided")
	}
	if request.Otp == "" || len(request.Otp) < 6 {
		log.Panic("Otp Invalid")
	}
}

func buildTwoFactorAuthenticationContext(userId int64, secret string, tfaChannelType enums.TfaChannelType) contexts.TwoFactorAuthenticationContext {
	twoFactorAuthenticationContext := contexts.TwoFactorAuthenticationContext{}
	twoFactorAuthenticationContext.UserId = userId
	twoFactorAuthenticationContext.SecretKey = secret
	twoFactorAuthenticationContext.TfaChannelType = tfaChannelType
	twoFactorAuthenticationContext.OtpGenerationTime = time.Now().UnixMilli()
	return twoFactorAuthenticationContext
}

func buildGenerateSecretResponse(key otp.Key, referenceId string, request requests.GenerateSecretRequest) responses.GenerateSecretResponse {
	res := responses.GenerateSecretResponse{}
	res.PlainSecret = key.Secret()
	res.QrCodeMessage = key.String()
	res.ReferenceId = referenceId
	res.TfaChannelType = request.TfaChannelType
	res.ExpiryTimeInSeconds = getExpiryTimeInSeconds(request.TfaChannelType)
	return res
}

func getExpiryTimeInSeconds(tfaChannelType enums.TfaChannelType) int {
	if tfaChannelType == enums.SMS {
		return 180
	} else if tfaChannelType == enums.APP {
		return 30
	}
	return 0
}

func validateGenerateSecretRequest(request requests.GenerateSecretRequest) {
	if request.UserId == 0 {
		log.Panic("UserId is not Provided")
	}
	if err := request.TfaChannelType.IsValid(); err != nil {
		log.Panic("Invalid Tfa Channel")
	}
	userTfaInfo := models.GetUserTfaInfoByUserId(request.UserId)
	if userTfaInfo != nil && (userTfaInfo.Sms || userTfaInfo.App) {
		log.Panic("Tfa Already Activated")
	}
}

func PostTfa(w http.ResponseWriter, r *http.Request) {
	userTfaInfo := &models.UserTfaInfo{}
	utils.ParseBody(r, userTfaInfo)
	models.CreateUserTfaInfo(userTfaInfo)
	res, _ := json.Marshal(userTfaInfo)
	w.Header().Set(headerName, headerValue)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
