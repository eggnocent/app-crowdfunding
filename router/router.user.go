package router

import (
	"app-crowdfunding/api"
	"app-crowdfunding/helper"
	"app-crowdfunding/util"
	"encoding/json"
	"net/http"
)

func HandlerListUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := user.List(ctx)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	response := helper.APIResponse("Users fetched successfully", http.StatusOK, "success", users)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerCheckEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input api.CheckEmailInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	message, err := user.CheckEmail(ctx, input)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusConflict, "Failed to check email")
	}

	response := helper.APIResponse("email check success", http.StatusOK, "success", message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func HandlerUploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parsing file dari form-data
	file, header, err := r.FormFile("Avatar")
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Failed to retrieve avatar file")
		return
	}
	defer file.Close()

	input := api.UpdateAvatarInput{
		File:     file,
		Filename: header.Filename,
	}

	// Proses penyimpanan avatar, misalnya dengan user.UploadAvatar()
	uploadedFileName, err := user.UpdateAvatar(ctx, input)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update avatar")
		return
	}

	response := helper.APIResponse("Avatar updated successfully", http.StatusOK, "success", map[string]string{"avatar_file_name": uploadedFileName})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
