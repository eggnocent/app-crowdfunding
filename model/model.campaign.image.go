package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CampaignImage struct {
	ID         uuid.UUID `db:"id"`
	CampaignID uuid.UUID `db:"campaign_id"`
	FileName   string    `db:"file_name"`
	IsPrimary  bool      `db:"is_primary"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type CampaignImageResponse struct {
	ID         uuid.UUID `json:"id"`
	CampaignID uuid.UUID `json:"campaign_id"`
	FileName   string    `json:"file_name"`
	IsPrimary  bool      `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ToResponse: Konversi struct ke JSON response
func (ci *CampaignImage) ToResponse() *CampaignImageResponse {
	return &CampaignImageResponse{
		ID:         ci.ID,
		CampaignID: ci.CampaignID,
		FileName:   ci.FileName,
		IsPrimary:  ci.IsPrimary,
		CreatedAt:  ci.CreatedAt,
		UpdatedAt:  ci.UpdatedAt,
	}
}

// Save: Simpan gambar ke database
func (ci *CampaignImage) Save(ctx context.Context, db *sqlx.DB) error {
	query := `INSERT INTO campaigns_images (id, campaign_id, file_name, is_primary, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	ci.ID = uuid.New()
	ci.CreatedAt = time.Now()
	ci.UpdatedAt = time.Now()

	_, err := db.ExecContext(ctx, query, ci.ID, ci.CampaignID, ci.FileName, ci.IsPrimary, ci.CreatedAt, ci.UpdatedAt)
	return err
}
