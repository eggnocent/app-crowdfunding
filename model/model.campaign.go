package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CampaignModel struct {
	ID               uuid.UUID `db:"id"`
	UserID           uuid.UUID `db:"user_id"`
	Name             string    `db:"name"`
	ShortDescription string    `db:"short_description"`
	Description      string    `db:"description"`
	GoalAmount       int       `db:"goal_amount"`
	CurrentAmount    int       `db:"current_amount"`
	Perks            string    `db:"perks"`
	BackerCount      int       `db:"backer_count"`
	Slug             string    `db:"slug"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type CampaignResponse struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	GoalAmount       int       `json:"goal_amount"`
	CurrentAmount    int       `json:"current_amount"`
	Perks            string    `json:"perks"`
	BackerCount      int       `json:"backer_count"`
	Slug             string    `json:"slug"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func NewCampaignResponse(c *CampaignModel) CampaignResponse {
	return CampaignResponse{
		ID:               c.ID,
		UserID:           c.UserID,
		Name:             c.Name,
		ShortDescription: c.ShortDescription,
		Description:      c.Description,
		GoalAmount:       c.GoalAmount,
		CurrentAmount:    c.CurrentAmount,
		Perks:            c.Perks,
		BackerCount:      c.BackerCount,
		Slug:             c.Slug,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func GetAllCampaign(ctx context.Context, db *sqlx.DB) ([]CampaignModel, error) {
	query := `SELECT id, user_id, name, short_description, description, goal_amount, current_amount, perks, backer_count, slug, created_at, updated_at FROM campaigns`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var campaigns []CampaignModel
	for rows.Next() {
		var campaign CampaignModel
		if err := rows.StructScan(&campaign); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}
	return campaigns, nil
}

func GetCampaignByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (CampaignModel, error) {
	var campaign CampaignModel
	query := `SELECT id, user_id, name, short_description, description, goal_amount, current_amount, perks, backer_count, slug, created_at, updated_at FROM campaigns WHERE id = $1`
	err := db.GetContext(ctx, &campaign, query, id)
	if err != nil {
		return CampaignModel{}, err
	}
	return campaign, nil
}

func (c *CampaignModel) CreateCampaign(ctx context.Context, db *sqlx.DB) error {
	query := `
		INSERT INTO campaigns (id, user_id, name, short_description, description, goal_amount, current_amount, perks, backer_count, slug, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id, created_at, updated_at
    `

	err := db.QueryRowxContext(ctx, query,
		c.ID,
		c.UserID,
		c.Name,
		c.ShortDescription,
		c.Description,
		c.GoalAmount,
		c.CurrentAmount,
		c.Perks,
		c.BackerCount,
		c.Slug,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(
		&c.ID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (c *CampaignModel) UpdateCampaign(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE campaigns 
		SET name = $1, short_description = $2, description = $3, goal_amount = $4, perks = $5, updated_at = $6
		WHERE id = $7
		RETURNING updated_at
    `

	err := db.QueryRowxContext(ctx, query,
		c.Name,
		c.ShortDescription,
		c.Description,
		c.GoalAmount,
		c.Perks,
		time.Now(),
		c.ID,
	).Scan(
		&c.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}
