package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sweatsparks/internal/commons/response"
	"sweatsparks/pkg/token"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			resp := response.UnauthorizedError("Missing Authorization header")
			w.WriteHeader(resp.StatusCode)
			json.NewEncoder(w).Encode(resp)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			resp := response.UnauthorizedError("Invalid Authorization header format")
			w.WriteHeader(resp.StatusCode)
			json.NewEncoder(w).Encode(resp)
			return
		}

		tokenStr := parts[1]
		token, err := token.ValidateToken(tokenStr)
		if err != nil {
			resp := response.UnauthorizedError("Invalid token")
			w.WriteHeader(resp.StatusCode)
			json.NewEncoder(w).Encode(resp)
			return
		}

		ctx := r.Context()
		ctx = contextWithUserID(ctx, int64(token.AuthId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func contextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value("userID").(int64)
	return userID, ok
}
