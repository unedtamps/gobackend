package util

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(w http.ResponseWriter, data interface{}, status int, mes string) {
	w.Header().Set("Content-Type", "application/json")
	res := Response{
		Success: true,
		Message: mes,
		Data:    data,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&res)
}

func ResponseError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	res := Response{
		Error: err.Error(),
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&res)
}
