package v1

import (
	"app-crowdfunding/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAPIUser(r *mux.Router) {
	r.HandleFunc("/users", router.HandlerListUser).Methods(http.MethodGet)
	r.HandleFunc("/avatar", router.HandlerUploadAvatar).Methods(http.MethodPost)
}
