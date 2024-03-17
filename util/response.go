package util

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(w http.ResponseWriter, data interface{}, mes string) {
	w.Header().Set("Content-Type", "application/json")
	res := Response{
		Status:  200,
		Message: mes,
		Data:    data,
	}
	json.NewEncoder(w).Encode(&res)
}

func ResponseError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	res := Response{
		Status: status,
		Error:  err.Error(),
	}
	json.NewEncoder(w).Encode(&res)
}
