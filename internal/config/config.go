package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func Validate(cfg *Config) error {
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				return fmt.Errorf(
					"config error: field=%s rule=%s value=%v",
					e.Namespace(),
					e.Tag(),
					e.Value(),
				)
			}
		}
		return err
	}
	return nil
}

type Config struct {
	Server    ServerConfig    `validate:"required"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit" validate:"required"`
	Email     EmailConfig     `validate:"required"`
	Databases DatabaseConfig  `validate:"required"`
	Cache     CacheConfig     `validate:"required"`
}

type ServerConfig struct {
	Host        string `validate:"required,hostname|ip"`
	Port        int    `validate:"required,gt=0,lt=65536"`
	Environment string `validate:"required,oneof=dev staging prod"`
	JWTSecret   string `validate:"required"`
}

type RateLimitConfig struct {
	Enable bool    `validate:"-"`
	RPS    float64 `validate:"gte=0"`
	Burst  int     `validate:"gte=0"`
}

type EmailConfig struct {
	Sender string     `validate:"required,email"`
	SMTP   SMTPConfig `validate:"required"`
}

type SMTPConfig struct {
	Host     string `validate:"required,hostname"`
	Port     int    `validate:"required,gt=0"`
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type DatabaseConfig struct {
	Postgres map[string]PostgresConfig `validate:"required,dive"`
	Mysql    map[string]MysqlConfig    `validate:"required,dive"`
}

type PostgresConfig struct {
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,gt=0"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DB       string `validate:"required"`
}

type MysqlConfig struct {
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,gt=0"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DB       string `validate:"required"`
}

type CacheConfig struct {
	Redis map[string]RedisConfig `validate:"required,dive"`
}

type RedisConfig struct {
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,gt=0"`
	Password string `validate:"-"`
	DB       int    `validate:"gte=0"`
}

func NewAppConfiguration(path string) (*Config, error) {
	var config *Config
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	if err := Validate(config); err != nil {
		return nil, err
	}

	return config, nil
}
