package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError représente une erreur de l'API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

// Constantes pour les codes d'erreur
const (
	ErrBadRequest          = 400
	ErrNotFound            = 404
	ErrInternalServerError = 500
)

// Constantes pour les messages d'erreur communs
const (
	MsgInvalidInput        = "Invalid input provided"
	MsgResourceNotFound    = "Requested resource not found"
	MsgInternalServerError = "An internal server error occurred"
)

// NewAPIError crée une nouvelle instance de APIError
func NewAPIError(code int, message string) APIError {
	return APIError{Code: code, Message: message}
}

// WriteJSONError écrit une réponse d'erreur JSON
func WriteJSONError(w http.ResponseWriter, err APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
