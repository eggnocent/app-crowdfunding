package router

import (
	"app-crowdfunding/api"
	"app-crowdfunding/helper"
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerGetAllTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	transactions, err := transaction.GetAllTransactions(ctx)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch all transactions")
		return
	}

	response := helper.APIResponse("Transactions fetched successfully", http.StatusOK, "success", transactions)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerGetTransactionByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	transactionIDStr, ok := vars["transaction_id"]

	if !ok || transactionIDStr == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "transaction_id is required")
		return
	}

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid transaction ID format")
		return
	}

	txn, err := transaction.GetByID(ctx, transactionID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch transaction by ID")
		return
	}

	response := helper.APIResponse("Transaction fetched successfully", http.StatusOK, "success", txn)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerGetTransactionsByCampaignID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	campaignIDStr, ok := vars["campaign_id"]

	if !ok || campaignIDStr == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "campaign_id is required")
		return
	}

	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}

	transactions, err := transaction.GetByCampaignID(ctx, campaignID)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch transactions by campaign ID")
		return
	}

	response := helper.APIResponse("Transactions fetched successfully", http.StatusOK, "success", transactions)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerGetTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userIDStr, ok := vars["user_id"]

	if !ok || userIDStr == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "user_id is required")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Validasi otorisasi user
	loggedInUserID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok || loggedInUserID != userID {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "You are not authorized to view these transactions")
		return
	}

	// Fetch transactions with campaign details
	transactions, err := transaction.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch transactions by user ID")
		return
	}

	response := helper.APIResponse("Transactions fetched successfully", http.StatusOK, "success", transactions)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerCreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input api.TrasactionInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	input.UserID = userID
	txnResponse, err := transaction.CreateTransaction(ctx, input)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Format JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).SetIndent("", "  ") // Tambahkan indentasi
	_ = json.NewEncoder(w).Encode(helper.APIResponse("Transaction created successfully", http.StatusOK, "success", txnResponse))
}

func HandlerProcessMidtransNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var notification model.MidtransNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid notification payload")
		return
	}

	err := payment.ProcessNotification(ctx, notification)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to process notification")
		return
	}

	response := helper.APIResponse("Notification processed successfully", http.StatusOK, "success", nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)

}
