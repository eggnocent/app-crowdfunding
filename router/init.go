package router

import (
	"app-crowdfunding/api"
	"app-crowdfunding/helper"
	"net/http"

	"github.com/jmoiron/sqlx"
)

var (
	authService *api.JWTService
	user        *api.UserAPIModule
	login       *api.LoginModule
	register    *api.RegistrationModule
	campaign    *api.CampaignModule
)

func Init(db *sqlx.DB) {
	authService = api.NewJWTService()
	user = api.NewUserAPIModule(db)
	login = api.NewLoginModule(db)
	register = api.NewRegistrationModule(db)
	campaign = api.NewCampaignModule(db)
}

func GetAuthMiddleware() func(http.Handler) http.Handler {
	return helper.AuthMiddleware(authService)
}
