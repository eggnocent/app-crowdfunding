package router

import (
	"app-crowdfunding/api"
	"app-crowdfunding/helper"
	"app-crowdfunding/util"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func HandlerUploadCampaignImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Ambil file dari form-data
	file, handler, err := r.FormFile("image")
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Failed to retrieve image file. Please ensure 'image' key exists in form-data.")
		return
	}
	defer file.Close()

	// Simpan file ke folder
	fileName, err := util.SaveAvatarCampaign(file, handler.Filename)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to save image. Please check server permissions.")
		return
	}

	// Ambil input campaign_id dan is_primary dari form-data
	campaignIdStr := r.FormValue("campaign_id")
	isPrimaryStr := r.FormValue("is_primary")

	// Validasi apakah campaign_id kosong
	if campaignIdStr == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Campaign ID is required.")
		return
	}

	// Validasi campaign_id apakah UUID valid
	campaignID, err := uuid.Parse(campaignIdStr)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID format. Please provide a valid UUID.")
		return
	}

	// Konversi is_primary ke boolean
	isPrimary := isPrimaryStr == "true"

	// Panggil API UploadImage
	input := api.CampaignImageInput{
		CampaignID: campaignID,
		IsPrimary:  isPrimary,
	}

	imageResponse, err := campaign.UploadImage(ctx, input, fileName)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to upload image. Please try again later.")
		return
	}

	// Return response
	response := helper.APIResponse("Image uploaded successfully", http.StatusOK, "success", imageResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
