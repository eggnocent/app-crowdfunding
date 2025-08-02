package api

import (
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionModule struct {
	db   *sqlx.DB
	name string
}

func NewTransactionModule(db *sqlx.DB) *TransactionModule {
	return &TransactionModule{
		db:   db,
		name: "transaction-module",
	}
}

type TrasactionInput struct {
	CampaignID uuid.UUID `json:"campaign_id" validate:"required"`
	UserID     uuid.UUID `json:"-"`
	Amount     int       `json:"amount" validate:"required"`
	Status     string    `json:"-"`
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

func (tm *TransactionModule) GetAllTransactions(ctx context.Context) ([]model.TransactionResponse, error) {
	transactions, err := model.GetAllTransactions(ctx, tm.db)
	if err != nil {
		return nil, util.NewErrorWrap(err, tm.name, "fetch all", ctx, "Failed to fetch all transactions", http.StatusInternalServerError)
	}
	return transactions, nil
}

func (tm *TransactionModule) GetByID(ctx context.Context, id uuid.UUID) (model.TransactionResponse, error) {
	transaction, err := model.GetTransactionByID(ctx, tm.db, id)
	if err != nil {
		return model.TransactionResponse{}, util.NewErrorWrap(err, tm.name, "fetch by id", ctx, "Failed to fetch transaction by ID", http.StatusInternalServerError)
	}
	response := transaction.Response()
	return *response, nil
}

func (tm *TransactionModule) GetByCampaignID(ctx context.Context, id uuid.UUID) ([]model.TransactionResponse, error) {
	transactions, err := model.GetTransactionsByCampaignID(ctx, tm.db, id)
	if err != nil {
		return nil, util.NewErrorWrap(err, tm.name, "fetch by campaign", ctx, "Failed to fetch transaction", http.StatusInternalServerError)
	}
	return transactions, nil
}

func (tm *TransactionModule) GetByUserID(ctx context.Context, id uuid.UUID) ([]model.TransactionResponse, error) {
	transactions, err := model.GetTransactionsByUserID(ctx, tm.db, id)
	if err != nil {
		return nil, util.NewErrorWrap(err, tm.name, "fetch by user", ctx, "Failed to fetch transaction", http.StatusInternalServerError)
	}
	return transactions, nil
}

func (tm *TransactionModule) CreateTransaction(ctx context.Context, input TrasactionInput) (*model.TransactionResponse, error) {
	transaction := model.TransactionModel{
		ID:         uuid.New(),
		CampaignID: input.CampaignID,
		UserID:     input.UserID,
		Amount:     input.Amount,
		Status:     "pending",
	}

	// Save transaction to database
	err := transaction.CreateTransaction(ctx, tm.db)
	if err != nil {
		return nil, util.NewErrorWrap(err, tm.name, "create-transaction", ctx, "Failed to create transaction", http.StatusInternalServerError)
	}

	// Get user details
	user, err := model.GetUserByID(ctx, tm.db, input.UserID)
	if err != nil {
		return nil, util.NewErrorWrap(err, tm.name, "fetch-user", ctx, "Failed to fetch user", http.StatusInternalServerError)
	}

	// Generate payment URL
	paymentModule := NewPaymentModule(tm.db)
	paymentURL, err := paymentModule.GetPaymentURL(ctx, transaction, user)
	if err != nil {
		return nil, util.NewErrorWrap(err, tm.name, "create-transaction", ctx, "Failed to generate payment URL", http.StatusInternalServerError)
	}

	transaction.PaymentURL = paymentURL
	err = transaction.UpdatePaymentUrl(ctx, tm.db)
	if err != nil {
		return nil, util.NewErrorWrap(err, tm.name, "create-transaction", ctx, "Failed to update payment URL", http.StatusInternalServerError)
	}

	return transaction.Response(), nil
}
