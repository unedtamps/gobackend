package config

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/unedtamps/gobackend/pkg/logger"
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
	RateLimit RateLimitConfig `validate:"required" mapstructure:"rate_limit"`
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
	Primary   PostgresConfig `mapstructure:"primary"`
	Secondary MySQLConfig    `mapstructure:"secondary"`
	Backup    SQLiteConfig   `mapstructure:"backup"`
}

func (dc *DatabaseConfig) All() []Database {
	var dbs []Database

	if dc.Primary.RDBMS != "" {
		dc.Primary.Name = "primary"
		dbs = append(dbs, dc.Primary)
	}
	if dc.Secondary.RDBMS != "" {
		dc.Secondary.Name = "secondary"
		dbs = append(dbs, dc.Secondary)
	}
	if dc.Backup.RDBMS != "" {
		dc.Backup.Name = "backup"
		dbs = append(dbs, dc.Backup)
	}

	return dbs
}

func (dc *DatabaseConfig) GetByName(name string) (Database, bool) {
	for _, db := range dc.All() {
		if db.GetName() == name {
			return db, true
		}
	}
	return nil, false
}

func (dc *DatabaseConfig) ListNames() []string {
	var names []string
	for _, db := range dc.All() {
		names = append(names, db.GetName())
	}
	return names
}

type Database interface {
	GetName() string
	GetRDBMS() string
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetDBName() string
	GetPath() string
}

type MySQLConfig struct {
	Name     string `validate:"required"`
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,gt=0"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"             mapstructure:"db_name"`
	RDBMS    string `validate:"required,oneof=mysql" mapstructure:"rdbms"`
}

func (m MySQLConfig) GetName() string     { return m.Name }
func (m MySQLConfig) GetRDBMS() string    { return m.RDBMS }
func (m MySQLConfig) GetHost() string     { return m.Host }
func (m MySQLConfig) GetPort() int        { return m.Port }
func (m MySQLConfig) GetUser() string     { return m.User }
func (m MySQLConfig) GetPassword() string { return m.Password }
func (m MySQLConfig) GetDBName() string   { return m.DBName }
func (m MySQLConfig) GetPath() string     { return "" }

type PostgresConfig struct {
	Name     string `validate:"required"`
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,gt=0"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"                mapstructure:"db_name"`
	RDBMS    string `validate:"required,oneof=postgres" mapstructure:"rdbms"`
}

func (p PostgresConfig) GetName() string     { return p.Name }
func (p PostgresConfig) GetRDBMS() string    { return p.RDBMS }
func (p PostgresConfig) GetHost() string     { return p.Host }
func (p PostgresConfig) GetPort() int        { return p.Port }
func (p PostgresConfig) GetUser() string     { return p.User }
func (p PostgresConfig) GetPassword() string { return p.Password }
func (p PostgresConfig) GetDBName() string   { return p.DBName }
func (p PostgresConfig) GetPath() string     { return "" }

type SQLiteConfig struct {
	Name  string `validate:"required"`
	Path  string `validate:"required"`
	RDBMS string `validate:"required,oneof=sqlite" mapstructure:"rdbms"`
}

func (s SQLiteConfig) GetName() string     { return s.Name }
func (s SQLiteConfig) GetRDBMS() string    { return s.RDBMS }
func (s SQLiteConfig) GetHost() string     { return "" }
func (s SQLiteConfig) GetPort() int        { return 0 }
func (s SQLiteConfig) GetUser() string     { return "" }
func (s SQLiteConfig) GetPassword() string { return "" }
func (s SQLiteConfig) GetDBName() string   { return "" }
func (s SQLiteConfig) GetPath() string     { return s.Path }

type CacheConfig struct {
	Redis map[string]RedisConfig `validate:"required,dive"`
}

type RedisConfig struct {
	Host     string `validate:"required,hostname|ip"`
	Port     int    `validate:"required,gt=0"`
	Password string `validate:"-"`
	DB       int    `validate:"gte=0"`
}

func NewLogger(config *Config) *slog.Logger {
	return logger.New(config.Server.Environment)
}

func NewAppConfiguration(path string) (*Config, error) {
	config := &Config{}
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
