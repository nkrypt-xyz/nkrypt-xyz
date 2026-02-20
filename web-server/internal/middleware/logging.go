package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/redact"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.body != nil {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging logs request/response information using zerolog.
// At INFO level: logs basic request info (method, path, status, duration)
// At DEBUG/TRACE level: additionally logs headers and body (with sensitive data redacted)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := GetRequestID(r.Context())
		logLevel := zerolog.GlobalLevel()

		// Basic request logging (always at INFO level)
		log.Info().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Msg("request started")

		// Detailed request logging at DEBUG/TRACE level
		if logLevel <= zerolog.DebugLevel {
			// Log headers (redacted)
			headers := make(map[string][]string)
			for k, v := range r.Header {
				headers[k] = v
			}
			redactedHeaders := redact.RedactHeaders(headers)

			logEvent := log.Debug().
				Str("request_id", requestID).
				Interface("headers", redactedHeaders).
				Str("query", r.URL.RawQuery)

			// Log request body if present (redacted)
			if r.Body != nil && r.ContentLength > 0 {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					// Restore body for handlers to read
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					
					// Redact and log body
					if len(bodyBytes) > 0 {
						redactedBody := redact.RedactJSONString(string(bodyBytes))
						logEvent.Str("body", redactedBody)
					}
				}
			}

			logEvent.Msg("request details")
		}

		// Capture response
		var ww *responseWriter
		if logLevel <= zerolog.DebugLevel {
			// Capture response body at DEBUG/TRACE level
			ww = &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           &bytes.Buffer{},
			}
		} else {
			// Only capture status code at INFO level
			ww = &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           nil,
			}
		}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Basic response logging (always at INFO level)
		log.Info().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.statusCode).
			Dur("duration", duration).
			Msg("request completed")

		// Detailed response logging at DEBUG/TRACE level
		if logLevel <= zerolog.DebugLevel && ww.body != nil {
			responseBody := ww.body.String()
			if len(responseBody) > 0 {
				redactedResponse := redact.RedactJSONString(responseBody)
				log.Debug().
					Str("request_id", requestID).
					Str("response_body", redactedResponse).
					Msg("response details")
			}
		}
	})
}

