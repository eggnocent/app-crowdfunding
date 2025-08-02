package api

import (
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/veritrans/go-midtrans"
)

type PaymentModule struct {
	db   *sqlx.DB
	name string
}

func NewPaymentModule(db *sqlx.DB) *PaymentModule {
	return &PaymentModule{
		db:   db,
		name: "payment-module",
	}
}

func (pm *PaymentModule) GetPaymentURL(ctx context.Context, transaction model.TransactionModel, user model.UserModel) (string, error) {
	midClient := midtrans.NewClient()
	midClient.ServerKey = "SB-Mid-server-dnT0rb4AQ4G1Yrss64G6OBpB"
	midClient.ClientKey = "SB-Mid-client-GRodqP_bd2veZmIa"
	midClient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midClient,
	}

	orderID := fmt.Sprintf("%s-%d", transaction.ID.String(), time.Now().Unix())

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", util.NewErrorWrap(err, pm.name, "get payment-url", ctx, "Failed to generate payment URL", http.StatusInternalServerError)
	}

	transaction.PaymentURL = snapTokenResp.RedirectURL
	err = transaction.UpdatePaymentUrl(ctx, pm.db)
	if err != nil {
		return "", util.NewErrorWrap(err, pm.name, "update payment-url", ctx, "Failed to save payment URL", http.StatusInternalServerError)
	}

	return snapTokenResp.RedirectURL, nil
}

func (pm *PaymentModule) ProcessNotification(ctx context.Context, notification model.MidtransNotification) error {
	log.Printf("Notification payload: %+v", notification)

	// Parse transaction ID
	transactionID, err := uuid.Parse(notification.OrderID)
	if err != nil {
		log.Printf("Error parsing transaction ID: %v", err)
		return util.NewErrorWrap(err, pm.name, "process-notification", ctx, "Invalid transaction ID format", http.StatusBadRequest)
	}

	// Retrieve transaction
	transaction, err := model.GetTransactionByID(ctx, pm.db, transactionID)
	if err != nil {
		log.Printf("Error retrieving transaction: %v", err)
		return util.NewErrorWrap(err, pm.name, "process-notification", ctx, "Transaction not found", http.StatusNotFound)
	}

	log.Printf("Transaction found: %+v", transaction)

	// Update transaction status based on notification
	switch notification.TransactionStatus {
	case "capture", "settlement":
		transaction.Status = "paid"
	case "deny", "cancel", "expire":
		transaction.Status = "failed"
	default:
		transaction.Status = "pending"
	}

	log.Printf("Before updating transaction: %+v", transaction)
	err = transaction.UpdateTransactionStatus(ctx, pm.db)
	if err != nil {
		log.Printf("Error updating transaction status: %v", err)
		return util.NewErrorWrap(err, pm.name, "update-transaction", ctx, "Failed to update transaction status", http.StatusInternalServerError)
	}
	log.Printf("After updating transaction: %+v", transaction)

	// Update campaign details if transaction is successful
	if transaction.Status == "paid" {
		campaign, err := model.GetCampaignByID(ctx, pm.db, transaction.CampaignID)
		if err != nil {
			log.Printf("Error fetching campaign: %v", err)
			return util.NewErrorWrap(err, pm.name, "fetch-campaign", ctx, "Campaign not found", http.StatusNotFound)
		}

		// Update campaign's current amount and backer count
		campaign.CurrentAmount += transaction.Amount
		campaign.BackerCount++
		err = campaign.UpdateCampaign(ctx, pm.db)
		if err != nil {
			log.Printf("Error updating campaign: %v", err)
			return util.NewErrorWrap(err, pm.name, "update-campaign", ctx, "Failed to update campaign", http.StatusInternalServerError)
		}
	}

	return nil
}
