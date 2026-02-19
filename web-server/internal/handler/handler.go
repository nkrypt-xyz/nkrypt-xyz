package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/apperror"
	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/pkg/validate"
)

// SendSuccess sends a 200 JSON response with hasError: false.
func SendSuccess(w http.ResponseWriter, data map[string]interface{}) {
	data["hasError"] = false
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

// SendErrorResponse sends an error JSON response according to the spec.
func SendErrorResponse(w http.ResponseWriter, err error) {
	statusCode := apperror.DetectHTTPStatusCode(err)
	serialized := apperror.SerializeError(err)

	response := map[string]interface{}{
		"hasError": true,
		"error":    serialized,
	}

	// Log non-user, non-validation errors.
	if _, ok := err.(*apperror.UserError); !ok {
		if _, ok := err.(*apperror.ValidationError); !ok {
			log.Error().Err(err).Msg("server error")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

// ParseAndValidateBody reads JSON body and validates with struct tags.
func ParseAndValidateBody(r *http.Request, dst interface{}) error {
	// 1. Limit body size to 100KB
	r.Body = http.MaxBytesReader(nil, r.Body, 100*1024)

	// 2. Decode JSON
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return &apperror.ValidationError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body: " + err.Error(),
		}
	}

	// 3. Validate struct
	if err := validate.Validator().Struct(dst); err != nil {
		return &apperror.ValidationError{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
			Details: formatValidationErrors(err),
		}
	}

	return nil
}

// formatValidationErrors converts validator errors into a serializable form.
func formatValidationErrors(err error) interface{} {
	// For now, just return the error string; can be enriched later.
	return map[string]interface{}{
		"error": err.Error(),
	}
}

