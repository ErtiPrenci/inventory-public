package response

import (
	"encoding/json"
	"net/http"
)

type Body struct {
	StatusCode  int         `json:"statusCode"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	Data        interface{} `json:"data,omitempty"` // omitempty: It doesn't send if nil
}

func generalResponse(w http.ResponseWriter, statusCode int, message, description string, data interface{}) {
	//Headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	payload := Body{
		StatusCode:  statusCode,
		Message:     message,
		Description: description,
		Data:        data,
	}

	// Convert to JSON and send
	json.NewEncoder(w).Encode(payload)
}

// Success (200)
func Success(w http.ResponseWriter, message, description string) {
	generalResponse(w, http.StatusOK, message, description, nil)
}

// SuccessData (200)
func SuccessData(w http.ResponseWriter, message, description string, data interface{}) {
	generalResponse(w, http.StatusOK, message, description, data)
}

// BadRequest (400)
func BadRequest(w http.ResponseWriter, message, description string) {
	generalResponse(w, http.StatusBadRequest, message, description, nil)
}

// NotFound (404)
func NotFound(w http.ResponseWriter, message, description string) {
	generalResponse(w, http.StatusNotFound, message, description, nil)
}

// Conflict (409)
func Conflict(w http.ResponseWriter, message, description string) {
	generalResponse(w, http.StatusConflict, message, description, nil)
}

// InternalServerError (500)
func InternalServerError(w http.ResponseWriter, message, description string) {
	generalResponse(w, http.StatusInternalServerError, message, description, nil)
}

func Unauthorized(w http.ResponseWriter, message, description string) {
	generalResponse(w, http.StatusUnauthorized, message, description, nil)
}
