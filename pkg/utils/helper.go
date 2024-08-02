package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func ResponseHandler(w http.ResponseWriter, message any, statusCode ...int) {

	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(message)
}

func Print(text ...string) {
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Println(green, text, reset)
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ErrorJSON sends a JSON error response with the given message and status code
func ErrorJSON(w http.ResponseWriter, message string, status ...int) {
	statusUsed := http.StatusBadRequest
	if len(status) > 0 {
		statusUsed = status[0]
	}
	response := ErrorResponse{
		Error:   http.StatusText(statusUsed),
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusUsed)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// If there's an error encoding the JSON, log it and send a plain text response
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred while processing the error response"))
	}
}

func PrintType(v interface{}) {
	fmt.Printf("Type: %s\n", reflect.TypeOf(v))
}
