package router

import (
	"app-crowdfunding/api"
	"app-crowdfunding/helper"
	"app-crowdfunding/util"
	"encoding/json"
	"net/http"

	"github.com/iancoleman/orderedmap"
)

func HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input api.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := login.LoginUser(ctx, input)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid to login: invalid credentials")
		return
	}

	token, err := authService.GenerateToken(user.ID.String())
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	data := orderedmap.New()
	data.Set("user", user)
	data.Set("token", token)

	response := helper.APIResponse("User logged in successfully", http.StatusOK, "success", data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
