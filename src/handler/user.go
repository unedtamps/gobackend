package handler

import (
	"context"
	"net/http"

	"github.com/unedtamps/gobackend/src/dto"
	"github.com/unedtamps/gobackend/src/middleware"
	"github.com/unedtamps/gobackend/utils"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	request := r.Context().Value("req").(dto.UserRegister)
	id, err := h.User.RegisterUser(context.Background(), request)
	if err != nil {
		utils.ResponseError(w, err.Code, err.Error)
		return
	}
	utils.ResponseSuccess(w, id, 201, "berhasil register")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	request := r.Context().Value("req").(dto.UserLogin)
	id, err := h.User.LoginUser(context.Background(), request)
	if err != nil {
		utils.ResponseError(w, err.Code, err.Error)
		return
	}
	_, token, _ := middleware.TokenAuth.Encode(
		map[string]interface{}{"email": request.Email, "id": id},
	)
	utils.ResponseSuccess(w, token, 200, "berhasil login")
}
