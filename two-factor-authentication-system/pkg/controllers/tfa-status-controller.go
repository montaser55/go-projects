package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
	userId, _ := strconv.ParseInt(mux.Vars(r)["userId"], 10, 64)
	userTfaInfo := models.GetUserTfaInfoByUserId(userId)
	tfaStatusResponse := &responses.TfaStatusResponse{}
	if userTfaInfo != nil {
		if userTfaInfo.Sms {
			tfaStatusResponse.IsEnabled = true
			tfaStatusResponse.TfaChannelType = enums.SMS
		} else if userTfaInfo.App {
			tfaStatusResponse.IsEnabled = true
			tfaStatusResponse.TfaChannelType = enums.APP
		}
	}

	res, _ := json.Marshal(tfaStatusResponse)
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
		Period:      uint(utils.GetExpiryTimeInSeconds(request.TfaChannelType)),
	})
	otp, _ := totp.GenerateCodeCustom(key.Secret(), time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(utils.GetExpiryTimeInSeconds(request.TfaChannelType)),
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	//config.RedisConnect()
	log.Printf(otp)
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
		validateOtpExpiration(twoFactorAuthenticationContext.OtpGenerationTime, utils.GetExpiryTimeInSeconds(enums.SMS))
	}
	validate, _ := totp.ValidateCustom(request.Otp, twoFactorAuthenticationContext.SecretKey, time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(utils.GetExpiryTimeInSeconds(request.TfaChannelType)),
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

func validateTfaDisableInitRequest(request requests.TfaDisableInitRequest, userTfaInfo *models.UserTfaInfo) {
	if userTfaInfo == nil {
		log.Panic("TFA is already disabled")
	}
	if (userTfaInfo.Sms && request.TfaChannelType != enums.SMS) || (userTfaInfo.App && request.TfaChannelType != enums.APP) {
		log.Panic("Invalid request")
	}
}

func validateTfaDisableRequest(request requests.TfaDisableRequest, userTfaInfo *models.UserTfaInfo) {
	tfaChannelType := request.TfaChannelType

	if tfaChannelType == enums.SMS && request.ReferenceId == "" {
		log.Panic("Reference ID not provided")
	}

	if userTfaInfo == nil || (!userTfaInfo.Sms && !userTfaInfo.App) {
		log.Panic("TFA is already disabled")
	}

	switch tfaChannelType {
	case enums.APP:
		if !userTfaInfo.App {
			log.Panic("TFA is already disabled")
		}
	case enums.SMS:
		if !userTfaInfo.Sms {
			log.Panic("TFA is already disabled")
		}
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
	res.ExpiryTimeInSeconds = utils.GetExpiryTimeInSeconds(request.TfaChannelType)
	return res
}

func buildTfaDisableInitResponse(referenceId string, expiryTimeInSeconds int) responses.TfaDisableInitResponse {
	res := responses.TfaDisableInitResponse{}
	res.ReferenceId = referenceId
	res.ExpiryTimeInSeconds = expiryTimeInSeconds
	return res
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

//func PostTfa(w http.ResponseWriter, r *http.Request) {
//	userTfaInfo := &models.UserTfaInfo{}
//	utils.ParseBody(r, userTfaInfo)
//	models.CreateUserTfaInfo(userTfaInfo)
//	res, _ := json.Marshal(userTfaInfo)
//	w.Header().Set(headerName, headerValue)
//	w.WriteHeader(http.StatusOK)
//	w.Write(res)
//}

func InitDisableTfa(w http.ResponseWriter, r *http.Request) {
	request := &requests.TfaDisableInitRequest{}
	utils.ParseBody(r, request)
	userTfaInfo := models.GetUserTfaInfoByUserId(request.UserId)
	validateTfaDisableInitRequest(*request, userTfaInfo)

	tfaChannelType := request.TfaChannelType

	key, _ := totp.Generate(totp.GenerateOpts{
		Period: uint(utils.GetExpiryTimeInSeconds(request.TfaChannelType)),
	})

	secretKey := ""
	interval := utils.GetExpiryTimeInSeconds(tfaChannelType)
	if userTfaInfo.Sms {
		secretKey := key.Secret()
		generatedOtp, _ := totp.GenerateCodeCustom(secretKey, time.Now().UTC(), totp.ValidateOpts{
			Period:    uint(interval),
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		log.Println("Generated OTP:", generatedOtp)
		// todo: send SMS
	} else if userTfaInfo.App {
		secretKey = userTfaInfo.SecretKey
	}

	referenceId := utils.GenerateReferenceId()

	redisClient := config.GetRedisClient()
	twoFactorAuthenticationContext := buildTwoFactorAuthenticationContext(request.UserId, secretKey, request.TfaChannelType)
	twoFactorAuthenticationContextJson, _ := json.Marshal(&twoFactorAuthenticationContext)
	redisClient.Set(config.GetContext(), referenceId, twoFactorAuthenticationContextJson, 0)

	tfaDisableInitResponse := buildTfaDisableInitResponse(referenceId, interval)
	res, _ := json.Marshal(&tfaDisableInitResponse)
	w.Header().Set(headerName, headerValue)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DisableTfa(w http.ResponseWriter, r *http.Request) {
	request := &requests.TfaDisableRequest{}
	utils.ParseBody(r, request)
	userTfaInfo := models.GetUserTfaInfoByUserId(request.UserId)
	validateTfaDisableRequest(*request, userTfaInfo)

	tfaChannelType := request.TfaChannelType

	redisClient := config.GetRedisClient()
	twoFactorAuthenticationContext := &contexts.TwoFactorAuthenticationContext{}
	str, err := redisClient.Get(config.GetContext(), request.ReferenceId).Result()
	if err != nil {
		log.Panic("ReferenceId not found in redis")
	}
	json.Unmarshal([]byte(str), twoFactorAuthenticationContext)
	if tfaChannelType == enums.SMS {
		validateOtpExpiration(twoFactorAuthenticationContext.OtpGenerationTime, utils.GetExpiryTimeInSeconds(enums.SMS))
	}
	validate, _ := totp.ValidateCustom(request.Otp, twoFactorAuthenticationContext.SecretKey, time.Now().UTC(), totp.ValidateOpts{
		Period:    uint(utils.GetExpiryTimeInSeconds(tfaChannelType)),
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	log.Printf("%v", validate)
	if !validate {
		log.Panic("Invalid Otp")
	}

	switch tfaChannelType {
	case enums.APP:
		userTfaInfo.App = false
	case enums.SMS:
		userTfaInfo.Sms = false
	default:
		log.Panic("Invalid OTP channel")
	}
	models.UpdateUserTfaInfo(userTfaInfo)
}
