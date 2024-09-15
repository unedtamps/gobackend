package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Config ConfigType

type ConfigType struct {
	SERVER_HOST       string  `mapstructure:"SERVER_HOST"`
	SERVER_PORT       string  `mapstructure:"SERVER_PORT"`
	JWT_SECRET        string  `mapstructure:"JWT_SECRET"`
	LIMIT_RPS         float32 `mapstructure:"LIMIT_RPS"`
	LIMIT_BURST       uint    `mapstructure:"LIMIT_BURST"`
	LIMIT_ENABLE      bool    `mapstructure:"LIMIT_ENABLE"`
	DB_DRIVER         string  `mapstructure:"DB_DRIVER"`
	POSTGRES_USER     string  `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD string  `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_DB       string  `mapstructure:"POSTGRES_DB"`
	POSTGRES_HOST     string  `mapstructure:"POSTGRES_HOST"`
	POSTGRES_PORT     string  `mapstructure:"POSTGRES_PORT"`
	EMAIL_SENDER      string  `mapstructure:"EMAIL_SENDER"`
	SMTP_HOST         string  `mapstructure:"SMTP_HOST"`
	SMTP_PORT         uint    `mapstructure:"SMTP_PORT"`
	SMTP_USERNAME     string  `mapstructure:"SMTP_USERNAME"`
	SMTP_PASSWORD     string  `mapstructure:"SMTP_PASSWORD"`
}

func ConStr() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Config.POSTGRES_HOST,
		Config.POSTGRES_PORT,
		Config.POSTGRES_USER,
		Config.POSTGRES_PASSWORD,
		Config.POSTGRES_DB,
	)
}

func load_env() []string {
	var env []string
	for _, s := range os.Environ() {
		e := strings.Split(s, "=")
		env = append(env, e[0])
	}
	return env
}

func init() {
	v := viper.New()
	v.AutomaticEnv()
	env := load_env()
	for _, s := range env {
		v.BindEnv(s)
	}
	if err := v.Unmarshal(&Config); err != nil {
		log.Fatal(err)
	}
}
