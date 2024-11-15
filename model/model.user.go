package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserModel struct {
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

type UserResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Occupation     string    `json:"occupation"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"password_hash"`
	AvatarFileName string    `json:"avatar_file_name"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Fungsi konversi UserModel ke UserResponse
func NewUserResponse(u *UserModel) UserResponse {
	return UserResponse{
		ID:             u.ID,
		Name:           u.Name,
		Occupation:     u.Occupation,
		Email:          u.Email,
		PasswordHash:   u.PasswordHash,
		AvatarFileName: u.AvatarFileName,
		Role:           u.Role,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func GetAllUser(ctx context.Context, db *sqlx.DB) ([]UserModel, error) {
	query := `SELECT id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at FROM users`

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []UserModel
	for rows.Next() {
		var user UserModel
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (UserModel, error) {
	var user UserModel
	query := `SELECT id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at FROM users WHERE id = $1`
	err := db.GetContext(ctx, &user, query, id)
	if err != nil {
		return UserModel{}, err
	}
	return user, nil
}

func (u *UserModel) RegisterUser(ctx context.Context, db *sqlx.DB) error {
	query := `
        INSERT INTO users (id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at
    `
	err := db.QueryRowxContext(ctx, query,
		u.ID,
		u.Name,
		u.Occupation,
		u.Email,
		u.PasswordHash,
		u.AvatarFileName,
		u.Role,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func CheckEmail(ctx context.Context, db *sqlx.DB, email string) (bool, error) {
	var count int
	query := `SELECT COUNT (*) FROM users WHERE email = $1`

	err := db.GetContext(ctx, &count, query, email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *UserModel) UpdateAvatar(ctx context.Context, db *sqlx.DB) error {
	query := `UPDATE users SET avatar_file_name = $1 WHERE id = $2 RETURNING updated_at`
	err := db.QueryRowxContext(ctx, query, u.AvatarFileName, u.ID).Scan(&u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
