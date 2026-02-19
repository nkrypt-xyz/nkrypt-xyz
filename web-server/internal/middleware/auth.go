package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/model"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/service"
)

type authContextKey string

const authDataKey authContextKey = "auth_data"

// Auth is the authentication middleware that verifies the Authorization header
// and populates auth data in the request context.
func Auth(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			authData, err := authSvc.Authenticate(r.Context(), authHeader)
			if err != nil {
				sendAuthErrorResponse(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), authDataKey, authData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAuthData extracts authentication data from the context.
func GetAuthData(ctx context.Context) *model.AuthData {
	if data, ok := ctx.Value(authDataKey).(*model.AuthData); ok {
		return data
	}
	return nil
}

// sendAuthErrorResponse is a local copy of the handler error serialization logic
// to avoid an import cycle between middleware and handler packages.
func sendAuthErrorResponse(w http.ResponseWriter, err error) {
	statusCode := apperror.DetectHTTPStatusCode(err)
	serialized := apperror.SerializeError(err)

	response := map[string]interface{}{
		"hasError": true,
		"error":    serialized,
	}

	if _, ok := err.(*apperror.UserError); !ok {
		if _, ok := err.(*apperror.ValidationError); !ok {
			log.Error().Err(err).Msg("server error")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}


