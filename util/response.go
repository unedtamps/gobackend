package util

import (
	"encoding/json"
	"net/http"
)

type Success struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ResponseSuccess(w http.ResponseWriter, data interface{}, status int, mes string) {
	w.Header().Set("Content-Type", "application/json")
	res := Success{
		Message: mes,
		Data:    data,
		Code:    status,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&res)
}

func ResponseError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	res := Error{
		Message: err.Error(),
		Code:    status,
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&res)
}
