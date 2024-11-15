package v1

import (
	"app-crowdfunding/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewLogin(r *mux.Router) {
	r.HandleFunc("/login", router.HandlerLoginUser).Methods(http.MethodPost)
}
