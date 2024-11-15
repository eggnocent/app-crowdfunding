package v1

import (
	"app-crowdfunding/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRegistration(r *mux.Router) {
	r.HandleFunc("/checkemail", router.HandlerCheckEmail).Methods(http.MethodPost)
	r.HandleFunc("/register", router.HandlerRegisterUser).Methods(http.MethodPost)
}
