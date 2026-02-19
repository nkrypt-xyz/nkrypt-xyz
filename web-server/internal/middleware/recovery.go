package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// Recovery recovers from panics, logs them, and returns a generic error
// response that matches the documented format.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.Error().
					Str("request_id", GetRequestID(r.Context())).
					Interface("panic", rvr).
					Bytes("stack", debug.Stack()).
					Msg("panic recovered")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"hasError": true,
					"error": map[string]interface{}{
						"code":    "GENERIC_SERVER_ERROR",
						"message": "We have encountered an unexpected server error.",
						"details": map[string]interface{}{},
					},
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

