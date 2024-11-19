package v1

import (
	"app-crowdfunding/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewCampaign(r *mux.Router) {
	r.HandleFunc("/campaigns", router.HandlerListCampaign).Methods(http.MethodGet)
	r.HandleFunc("/campaigns/{id}", router.HandlerDetailByIDCampaign).Methods(http.MethodGet)
	r.HandleFunc("/campaigns", router.HandlerCreateCampaign).Methods(http.MethodPost)
	r.HandleFunc("/campaigns/{id}", router.HandlerUpdateCampaign).Methods(http.MethodPut)
}
