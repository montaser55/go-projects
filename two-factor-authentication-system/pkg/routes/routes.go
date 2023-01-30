package routes

import (
	"github.com/gorilla/mux"
	"github.com/montaser55/two-factor-authentication-service/pkg/controllers"
	"github.com/montaser55/two-factor-authentication-service/pkg/models"
	"github.com/montaser55/two-factor-authentication-service/pkg/models/requests"
	"github.com/montaser55/two-factor-authentication-service/pkg/utils"
	"github.com/montaser55/two-factor-authentication-service/pkg/utils/enums"
	"github.com/montaser55/two-factor-authentication-service/pkg/utils/gateways"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var RegisterRoutes = func(router *mux.Router) {
	contextRouter := router.PathPrefix("/two-factor-authentication-service").Subrouter()

	tfaStatusRouter := contextRouter.PathPrefix("/api/tfa-status").Subrouter()
	tfaStatusRouter.HandleFunc("/check/{userId}", controllers.CheckTfa).Methods("GET")
	tfaStatusRouter.HandleFunc("/generate-secret", controllers.GenerateSecret).Methods("POST")
	tfaStatusRouter.HandleFunc("/enable", controllers.EnableTfa).Methods("POST")
	tfaStatusRouter.HandleFunc("/init-disable", controllers.InitDisableTfa).Methods("POST")
	tfaStatusRouter.HandleFunc("/disable", controllers.DisableTfa).Methods("POST")

	authenticateRouter := contextRouter.PathPrefix("/api/authenticate").Subrouter()
	authenticateRouter.HandleFunc("/verify-otp", controllers.CheckTfa).Methods("POST")
	authenticateRouter.HandleFunc("/resend-otp", controllers.CheckTfa).Methods("POST")

	filterRouter := contextRouter.PathPrefix("").Subrouter()
	//filterRouter.HandleFunc("/api/execute-trade/user/generate-coin-release-tx-dispute", controllers.RouteApi).Methods("POST")
	//filterRouter.HandleFunc("/api/execute-trade/user/generate-coin-release-tx", controllers.RouteApi).Methods("POST")
	filterRouter.HandleFunc("/api/transaction/send-coins", controllers.CheckTfa).Methods("POST")
	//filterRouter.HandleFunc("/api/transaction/request-coin-redeem", controllers.RouteApi).Methods("POST")

	filterRouter.Use(credentialVerificationMiddleware, twoFactorAuthenticationMiddleware)

	router.Use(loggingMiddleware)

}

func credentialVerificationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentialRequest := &requests.CredentialRequest{}
		utils.ParseBodyReusable(r, credentialRequest)
		gateways.VerifyPin(*credentialRequest)
		next.ServeHTTP(w, r)
	})
}

func twoFactorAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentialRequest := &requests.CredentialRequest{}
		utils.ParseBodyReusable(r, credentialRequest)
		userTfaInfo := models.GetUserTfaInfoByUserId(credentialRequest.UserId)
		if userTfaInfo != nil && (userTfaInfo.Sms || userTfaInfo.App) {
			handleTwoFactorAuthentication(w, r, *userTfaInfo)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleTwoFactorAuthentication(w http.ResponseWriter, r *http.Request, userTfaInfo models.UserTfaInfo) {
	_ = deduceTransactionType(r.URL)

}

func deduceTransactionType(url *url.URL) enums.TransactionType {
	if strings.HasSuffix(url.String(), "/send-coins") {
		return enums.SEND_COIN
	} else if strings.HasSuffix(url.String(), "/request-coin-redeem") {
		return enums.COIN_REDEEM
	} else if strings.HasSuffix(url.String(), "/generate-coin-release-tx-dispute") {
		return enums.DISPUTE
	} else if strings.HasSuffix(url.String(), "/generate-coin-release-tx") {
		return enums.TRADE
	}
	log.Panic("Invalid TransactionType")
	return ""
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
