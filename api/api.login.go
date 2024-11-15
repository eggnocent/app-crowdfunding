package api

import (
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type LoginModule struct {
	db   *sqlx.DB
	name string
}

func NewLoginModule(db *sqlx.DB) *LoginModule {
	return &LoginModule{
		db:   db,
		name: "login-module",
	}
}

func (l *LoginModule) LoginUser(ctx context.Context, input LoginInput) (model.LoginResponse, error) {
	login, err := model.Login(ctx, l.db, input.Email)
	if err != nil {
		return model.LoginResponse{}, errors.New("invalid email or password")
	}

	if !util.CheckPasswordHash(input.Password, login.PasswordHash) {
		return model.LoginResponse{}, errors.New("invalid email or password")
	}
	return login.ToLoginResponse(), nil
}
