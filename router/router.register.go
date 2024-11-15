package router

import (
	"app-crowdfunding/api"
	"app-crowdfunding/helper"
	"app-crowdfunding/util"
	"encoding/json"
	"net/http"
)

func HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input api.RegistrationUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Invalid request body")
		return
	}

	user, err := register.RegisterUser(ctx, input)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	response := helper.APIResponse("user registration success", http.StatusOK, "success", user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
