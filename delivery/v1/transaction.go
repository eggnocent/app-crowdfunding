package v1

import (
	"app-crowdfunding/router"
	"net/http"

	"github.com/gorilla/mux"
)

func NewTransaction(r *mux.Router) {
	r.HandleFunc("/transactions", router.HandlerCreateTransaction).Methods(http.MethodPost)
	r.HandleFunc("/transactions", router.HandlerGetAllTransactions).Methods(http.MethodGet)
	r.HandleFunc("/campaigns/{campaign_id}/transactions", router.HandlerGetTransactionsByCampaignID).Methods(http.MethodGet)
	r.HandleFunc("/users/{user_id}/transactions", router.HandlerGetTransactionsByUserID).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{transaction_id}", router.HandlerGetTransactionByID).Methods(http.MethodGet)
	r.HandleFunc("/transactions/midtrans-notification", router.HandlerProcessMidtransNotification).Methods(http.MethodPost)
}
