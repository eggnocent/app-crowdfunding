package model

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionModel struct {
	ID         uuid.UUID `db:"id"`
	CampaignID uuid.UUID `db:"campaign_id"`
	UserID     uuid.UUID `db:"user_id"`
	Amount     int       `db:"amount"`
	Status     string    `db:"status"`
	PaymentURL string    `db:"payment_url"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type TransactionResponse struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CampaignID uuid.UUID `db:"campaign_id" json:"campaign_id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	Amount     int       `db:"amount" json:"amount"`
	Status     string    `db:"status" json:"status"`
	PaymentURL string    `db:"payment_url" json:"payment_url"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type MidtransNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
	GrossAmount       string `json:"gross_amount"`
}

func (t *TransactionModel) Response() *TransactionResponse {
	return &TransactionResponse{
		ID:         t.ID,
		CampaignID: t.CampaignID,
		UserID:     t.UserID,
		Amount:     t.Amount,
		Status:     t.Status,
		PaymentURL: t.PaymentURL,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func GetAllTransactions(ctx context.Context, db *sqlx.DB) ([]TransactionResponse, error) {
	query := `
		SELECT id, campaign_id, user_id, amount, status, 
		       COALESCE(payment_url, '') AS payment_url, 
		       created_at, updated_at
		FROM transactions
	`
	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []TransactionResponse
	for rows.Next() {
		var transaction TransactionResponse
		if err := rows.StructScan(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionByID(ctx context.Context, db *sqlx.DB, transactionID uuid.UUID) (TransactionModel, error) {
	query := `SELECT id, campaign_id, user_id, amount, status, payment_url, created_at, updated_at 
              FROM transactions 
              WHERE id = $1`

	var transaction TransactionModel
	err := db.GetContext(ctx, &transaction, query, transactionID)
	if err != nil {
		return TransactionModel{}, err
	}

	return transaction, nil
}

func GetTransactionsByCampaignID(ctx context.Context, db *sqlx.DB, campaignID uuid.UUID) ([]TransactionResponse, error) {
	query := `
		SELECT id, campaign_id, user_id, amount, status, 
		       COALESCE(payment_url, '') AS payment_url, 
		       created_at, updated_at
		FROM transactions 
		WHERE campaign_id = $1
	`
	rows, err := db.QueryxContext(ctx, query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []TransactionResponse
	for rows.Next() {
		var transaction TransactionResponse
		if err := rows.StructScan(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionsByUserID(ctx context.Context, db *sqlx.DB, userID uuid.UUID) ([]TransactionResponse, error) {
	query := `
		SELECT id, campaign_id, user_id, amount, status, 
		       COALESCE(payment_url, '') AS payment_url, 
		       created_at, updated_at
		FROM transactions 
		WHERE user_id = $1
	`
	rows, err := db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []TransactionResponse
	for rows.Next() {
		var transaction TransactionResponse
		if err := rows.StructScan(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (t *TransactionModel) CreateTransaction(ctx context.Context, db *sqlx.DB) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	query := `INSERT INTO transactions(id, campaign_id, user_id, amount, status, payment_url, created_at, updated_at) 
	          VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	log.Printf("Creating Transaction: %+v", t)

	_, err := db.ExecContext(ctx, query, t.ID, t.CampaignID, t.UserID, t.Amount, t.Status, t.PaymentURL, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting transaction: %v", err)
	}
	return err
}

func (t *TransactionModel) UpdatePaymentUrl(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE transactions SET payment_url = $1, updated_at = $2 WHERE id = $3`
	_, err := db.ExecContext(ctx, query, t.PaymentURL, time.Now(), t.ID)
	if err != nil {
		log.Printf("Error updating payment URL: %v", err)
		log.Printf("Transaction Data: %+v", t)
	}
	return err
}

func (t *TransactionModel) UpdateTransactionStatus(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE transactions 
			  SET status = $1, updated_at = $2 
			  WHERE id = $3`
	_, err := db.ExecContext(ctx, query, t.Status, time.Now(), t.ID)
	return err
}
