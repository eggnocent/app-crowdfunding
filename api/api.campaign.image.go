package api

import (
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CampaignImageInput struct {
	CampaignID uuid.UUID `json:"campaign_id" validate:"required"`
	IsPrimary  bool      `json:"is_primary"`
}

func (cm *CampaignModule) UploadImage(ctx context.Context, input CampaignImageInput, fileName string) (*model.CampaignImageResponse, error) {
	// Validasi campaign_id
	campaign, err := model.GetCampaignByID(ctx, cm.db, input.CampaignID)
	if err != nil || campaign.ID == uuid.Nil {
		return nil, util.NewErrorWrap(err, cm.name, "upload", ctx, "invalid campaign_id", http.StatusNotFound)
	}

	// Inisialisasi data gambar
	image := model.CampaignImage{
		CampaignID: input.CampaignID,
		FileName:   fileName,
		IsPrimary:  input.IsPrimary,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Simpan gambar ke database
	err = image.Save(ctx, cm.db)
	if err != nil {
		return nil, util.NewErrorWrap(err, cm.name, "upload", ctx, "failed to save image", http.StatusInternalServerError)
	}

	// Return response
	return image.ToResponse(), nil
}
