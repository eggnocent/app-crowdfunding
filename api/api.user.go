package api

import (
	"app-crowdfunding/model"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserAPIModule struct {
	db   *sqlx.DB
	name string
}

func NewUserAPIModule(db *sqlx.DB) *UserAPIModule {
	return &UserAPIModule{
		db:   db,
		name: "user-module",
	}
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateAvatarInput struct {
	File     multipart.File
	Filename string
}

func (um *UserAPIModule) List(ctx context.Context) (interface{}, error) {
	users, err := model.GetAllUser(ctx, um.db)
	if err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, model.NewUserResponse(&user))
	}

	return userResponses, nil

}

func (um *UserAPIModule) GetByID(ctx context.Context, id uuid.UUID) (model.UserResponse, error) {
	userModel, err := model.GetUserByID(ctx, um.db, id)
	if err != nil {
		return model.UserResponse{}, err
	}

	userResponse := model.NewUserResponse(&userModel)
	return userResponse, nil
}

func (um *UserAPIModule) CheckEmail(ctx context.Context, input CheckEmailInput) (string, error) {
	emailExists, err := model.CheckEmail(ctx, um.db, input.Email)
	if err != nil {
		return "", err
	}

	if emailExists {
		return "Email already exists", errors.New("email already in use")
	}
	return "Email is available", nil
}

func (um *UserAPIModule) UpdateAvatar(ctx context.Context, input UpdateAvatarInput) (string, error) {
	path := fmt.Sprintf("./images/%s", input.Filename)
	outFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, input.File)
	if err != nil {
		return "", err
	}
	return path, nil
}
