package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unedtamps/gobackend/config"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestHealthCheck(t *testing.T) {

	resp, err := http.Get(
		fmt.Sprintf("http://%s:%s/",
			config.Config.SERVER_HOST, config.Config.SERVER_PORT))

	assert.NoError(t, err, "Response not Error")
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateUser(t *testing.T) {
	user := User{
		Email:    "user1@gmail.com",
		Password: "password",
	}
	data, _ := json.Marshal(user)

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%s/user/register",
			config.Config.SERVER_HOST,
			config.Config.SERVER_PORT), "application/json", bytes.NewBuffer(data))

	assert.NoError(t, err, "Response not Error")
	assert.Equal(t, 201, resp.StatusCode)
	// check user email

	data_res, err := repo.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err, "No error")
	assert.Equal(t, user.Email, data_res.Email)
}
