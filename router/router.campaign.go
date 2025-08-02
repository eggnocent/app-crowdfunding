package router

import (
	"app-crowdfunding/api"
	"app-crowdfunding/helper"
	"app-crowdfunding/util"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HandlerListCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	campaigns, err := campaign.List(ctx)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch campaigns")
		return
	}

	response := helper.APIResponse("Campaigns fetched successfully", http.StatusOK, "success", campaigns)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerDetailByIDCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr, ok := vars["id"]

	if !ok {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}

	campaign, err := campaign.DetailByID(ctx, id)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch campaign")
		return
	}

	response := helper.APIResponse("Campaign detail fetched successfully", http.StatusOK, "success", campaign)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerCreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input api.CampaignInput

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
	

	campaign, err := campaign.CreateCampaign(ctx, input)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create campaign")
		return
	}

	response := helper.APIResponse("Campaign created successfully", http.StatusOK, "success", campaign)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerUpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr, ok := vars["id"]

	if !ok {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid campaign ID format")
		return
	}

	var input api.CampaignUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	campaign, err := campaign.UpdateCampaign(ctx, id, input)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update campaign")
		return
	}

	response := helper.APIResponse("Campaign updated successfully", http.StatusOK, "success", campaign)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
