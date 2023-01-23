package routes

import (
	"github.com/gorilla/mux"
	"github.com/montaser55/two-factor-authentication-service/pkg/controllers"
)

var RegisterRoutes = func(router *mux.Router) {
	contextPath := "/two-factor-authentication-service"
	tfaStatus := "/api/tfa-status"
	router.HandleFunc(contextPath+tfaStatus+"/check/{userId}", controllers.CheckTfa).Methods("GET")
	router.HandleFunc(contextPath+tfaStatus+"/post", controllers.PostTfa).Methods("POST")
	router.HandleFunc(contextPath+tfaStatus+"/generate-secret", controllers.CheckTfa).Methods("POST")
	router.HandleFunc(contextPath+tfaStatus+"/enable", controllers.CheckTfa).Methods("POST")
	router.HandleFunc(contextPath+tfaStatus+"/init-disable", controllers.CheckTfa).Methods("POST")
	router.HandleFunc(contextPath+tfaStatus+"/disable", controllers.CheckTfa).Methods("POST")

}
