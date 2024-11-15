package api

import (
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RegistrationModule struct {
	db   *sqlx.DB
	name string
}

func NewRegistrationModule(db *sqlx.DB) *RegistrationModule {
	return &RegistrationModule{
		db:   db,
		name: "user-module",
	}
}

type RegistrationUserInput struct {
	Name           string `json:"name" validate:"required"`
	Occupation     string `json:"occupation" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	AvatarFileName string `json:"avatar_file_name"`
	Role           string `json:"role" validate:"required"`
}

func (r *RegistrationModule) RegisterUser(ctx context.Context, input RegistrationUserInput) (model.UserResponse, error) {
	user := model.UserModel{
		ID:             uuid.New(),
		Name:           input.Name,
		Occupation:     input.Occupation,
		Email:          input.Email,
		PasswordHash:   util.HashPassword(input.Password),
		AvatarFileName: input.AvatarFileName,
		Role:           input.Role,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := user.RegisterUser(ctx, r.db)
	if err != nil {
		return model.NewUserResponse(&user), util.NewErrorWrap(err, r.name, "register", ctx, "unable to register user", http.StatusInternalServerError)
	}
	return model.NewUserResponse(&user), nil
}
