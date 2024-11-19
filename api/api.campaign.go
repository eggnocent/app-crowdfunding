package api

import (
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	CampaignModule struct {
		db   *sqlx.DB
		name string
	}
)

type CampaignInput struct {
	UserID           uuid.UUID `json:"user_id" validate:"required"`
	Name             string    `json:"name" validate:"required"`
	ShortDescription string    `json:"short_description" validate:"required"`
	Description      string    `json:"description" validate:"required"`
	GoalAmount       int       `json:"goal_amount" validate:"required"`
	Perks            string    `json:"perks" validate:"required"`
	Slug             string    `json:"slug" validate:"required"`
}

type CampaignUpdateInput struct {
	Name             string `json:"name" validate:"required"`
	ShortDescription string `json:"short_description" validate:"required"`
	Description      string `json:"description" validate:"required"`
	GoalAmount       int    `json:"goal_amount" validate:"required"`
	Perks            string `json:"perks" validate:"required"`
}

func NewCampaignModule(db *sqlx.DB) *CampaignModule {
	return &CampaignModule{
		db:   db,
		name: "campaign-module",
	}
}

func (cm *CampaignModule) List(ctx context.Context) (interface{}, error) {
	campaigns, err := model.GetAllCampaign(ctx, cm.db)
	if err != nil {
		return nil, err
	}

	var campaignResponses []model.CampaignResponse
	for _, campaign := range campaigns {
		campaignResponses = append(campaignResponses, model.NewCampaignResponse(&campaign))
	}

	return campaignResponses, nil
}

func (cm *CampaignModule) DetailByID(ctx context.Context, id uuid.UUID) (model.CampaignResponse, error) {
	campaignModel, err := model.GetCampaignByID(ctx, cm.db, id)
	if err != nil {
		return model.CampaignResponse{}, err
	}

	campaignResponse := model.NewCampaignResponse(&campaignModel)
	return campaignResponse, nil
}

func (cm *CampaignModule) CreateCampaign(ctx context.Context, input CampaignInput) (*model.CampaignModel, error) {
	slug := strings.ToLower(strings.ReplaceAll(input.Name, " ", "-"))

	campaign := model.CampaignModel{
		ID:               uuid.New(),
		UserID:           input.UserID,
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		CurrentAmount:    0,
		Perks:            input.Perks,
		BackerCount:      0,
		Slug:             slug,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := campaign.CreateCampaign(ctx, cm.db)
	if err != nil {
		return nil, util.NewErrorWrap(err, cm.name, "create", ctx, "unable to create campaign", http.StatusInternalServerError)
	}

	return &campaign, nil
}

func (cm *CampaignModule) UpdateCampaign(ctx context.Context, id uuid.UUID, input CampaignUpdateInput) (model.CampaignResponse, error) {
	campaign, err := model.GetCampaignByID(ctx, cm.db, id)
	if err != nil {
		return model.CampaignResponse{}, util.NewErrorWrap(err, cm.name, "update", ctx, "unable to get campaign", http.StatusNotFound)
	}

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks

	campaign.ID = id
	err = campaign.UpdateCampaign(ctx, cm.db)
	if err != nil {
		return model.CampaignResponse{}, util.NewErrorWrap(err, cm.name, "update", ctx, "unable to update campaign", http.StatusInternalServerError)
	}

	return model.NewCampaignResponse(&campaign), nil
}
