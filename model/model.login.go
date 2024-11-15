package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LoginModel struct {
	ID             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	Occupation     string    `db:"occupation"`
	Email          string    `db:"email"`
	PasswordHash   string    `db:"password_hash"`
	AvatarFileName string    `db:"avatar_file_name"`
	Role           string    `db:"role"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type LoginResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Occupation     string    `json:"occupation"`
	Email          string    `json:"email"`
	AvatarFileName string    `json:"avatar_file_name"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (l *LoginModel) ToLoginResponse() LoginResponse {
	return LoginResponse{
		ID:             l.ID,
		Name:           l.Name,
		Occupation:     l.Occupation,
		Email:          l.Email,
		AvatarFileName: l.AvatarFileName,
		Role:           l.Role,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
	}
}

func Login(ctx context.Context, db *sqlx.DB, email string) (LoginModel, error) {
	var user LoginModel
	query := `SELECT id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at FROM users WHERE email = $1`
	err := db.GetContext(ctx, &user, query, email)
	if err != nil {
		return LoginModel{}, err
	}
	return user, nil
}
