package apperror

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode string                 `json:"error_code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *AppError
	var ok bool
	if appErr, ok = err.(*AppError); !ok {
		appErr = ErrInternal
	}

	w.WriteHeader(appErr.HTTPStatus)
	json.NewEncoder(w).Encode(ErrorResponse{
		ErrorCode: appErr.Code,
		Message:   appErr.Message,
		Details:   appErr.Details,
	})
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
