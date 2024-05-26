package handler

import (
	"context"
	"net/http"

	"github.com/unedtamps/gobackend/pkg/dto"
	"github.com/unedtamps/gobackend/pkg/middleware"
	"github.com/unedtamps/gobackend/util"
)

// Register User
//
//		@Summary		Register User
//		@Description	User this api to Register your account
//		@Tags			accounts
//		@Accept			json
//		@Produce		json
//		@Param			 request  body  dto.UserRegister  true  "register using data request"
//	 @Success      200  {object}  util.Success
//	 @Failure      404  {object}  util.Error
//	 @Failure      400  {object}  util.Error
//	 @Failure      500  {object}  util.Error
//		@Router			/user/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	request := r.Context().Value("req").(dto.UserRegister)
	id, err := h.User.RegisterUser(context.Background(), request)
	if err != nil {
		util.ResponseError(w, err.Code, err.Error)
		return
	}
	util.ResponseSuccess(w, id, 201, "berhasil register")
}

// Login User
//
//		@Summary		Login User
//		@Description	User this api to Login your account
//		@Tags			accounts
//		@Accept			json
//		@Produce		json
//		@Param			 request  body  dto.UserLogin  true  "login using data request"
//	 @Success      200  {object}  util.Success
//	 @Failure      404  {object}  util.Error
//	 @Failure      400  {object}  util.Error
//	 @Failure      500  {object}  util.Error
//		@Router			/user/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	request := r.Context().Value("req").(dto.UserLogin)
	id, err := h.User.LoginUser(context.Background(), request)
	if err != nil {
		util.ResponseError(w, err.Code, err.Error)
		return
	}
	_, token, _ := middleware.TokenAuth.Encode(
		map[string]interface{}{"email": request.Email, "id": id},
	)
	util.ResponseSuccess(w, token, 200, "berhasil login")
}
