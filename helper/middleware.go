// /helper/middleware.go

package helper

import (
	"app-crowdfunding/api"
	"app-crowdfunding/model"
	"app-crowdfunding/util"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type contextKey string

const (
	ContextUserKey contextKey = "currentUser"
)

func AuthMiddleware(authService *api.JWTService, userModule *api.UserAPIModule) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				util.WriteErrorResponse(w, http.StatusUnauthorized, "Missing Authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}

			tokenString := parts[1]
			token, err := authService.ValidateToken(tokenString)
			if err != nil || !token.Valid {
				util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Mengambil klaim dari token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			// Mengambil user_id dari klaim
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid user ID in token")
				return
			}

			// Mengonversi user_id ke UUID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				util.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid user ID format")
				return
			}

			// Mengambil pengguna dari API
			user, err := userModule.GetByID(r.Context(), userID)
			if err != nil {
				util.WriteErrorResponse(w, http.StatusUnauthorized, "User not found")
				return
			}

			// Menyimpan pengguna ke context
			ctx := context.WithValue(r.Context(), ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetCurrentUser mengambil pengguna yang tersimpan di context
func GetCurrentUser(ctx context.Context) (model.UserModel, bool) {
	user, ok := ctx.Value(ContextUserKey).(model.UserModel)
	return user, ok
}
