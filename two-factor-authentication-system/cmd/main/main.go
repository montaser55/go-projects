package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/montaser55/two-factor-authentication-service/pkg/routes"
)

func main() {

	r := mux.NewRouter()
	routes.RegisterRoutes(r)
	log.Fatal(http.ListenAndServe("localhost:10340", r))

}
