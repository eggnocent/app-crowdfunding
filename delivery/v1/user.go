package v1

import (
	"app-crowdfunding/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIUser(r *mux.Router) {
	r.HandleFunc("/register", router.HandlerRegisterUser).Methods(http.MethodPost)
	r.HandleFunc("/checkemail", router.HandlerCheckEmail).Methods(http.MethodPost)

	r.HandleFunc("/users", router.HandlerListUser).Methods(http.MethodGet)
	r.HandleFunc("/avatar", router.HandlerUploadAvatar).Methods(http.MethodPost)
}
