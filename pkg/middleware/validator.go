package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/utils"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterCustomTypeFunc(validateUUID, uuid.UUID{})
}

func Validate[T any](next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var jsonData T

		if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
			utils.ResponseError(w, 400, err)
			return
		}

		if err := validate.Struct(jsonData); err != nil {
			utils.ResponseError(w, 400, err)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "req", jsonData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateUUID(field reflect.Value) interface{} {
	if value, ok := field.Interface().(uuid.UUID); ok {
		return value.String()
	}
	return nil
}
