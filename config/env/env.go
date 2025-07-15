package env

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/iamolegga/enviper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	env           *Env
	viperInstance *viper.Viper
	once          sync.Once
)

type Env struct {
	AppEnv             string `mapstructure:"APP_ENV"`
	Port               string `mapstructure:"PORT"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPassword         string `mapstructure:"DB_PASS"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             string `mapstructure:"DB_PORT"`
	DBName             string `mapstructure:"DB_NAME"`
	JwtExpiredTime     time.Duration
	JwtSecretKey       []byte
	SmtpHost           string `mapstructure:"SMTP_HOST"`
	SmtpPort           int    `mapstructure:"SMTP_PORT"`
	SmtpEmail          string `mapstructure:"SMTP_EMAIL"`
	SmtpPassword       string `mapstructure:"SMTP_PASSWORD"`
	GoogleRedirectUrl  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleTokenUrl     string `mapstructure:"GOOGLE_TOKEN_URL"`
	GeminiApiKey       string `mapstructure:"GEMINI_API_KEY"`
	GeminiApiModel     string `mapstructure:"GEMINI_API_MODEL"`
	MidtransServerKey  string `mapstructure:"MIDTRANS_SERVER_KEY"`
	BcryptCost         int    `mapstructure:"BCRYPT_COST"`
	ImageModelUrl      string `mapstructure:"IMAGE_MODEL_URL"`
	AudioModelUrl      string `mapstructure:"AUDIO_MODEL_URL"`
	GRPCPort           string `mapstructure:"GRPC_PORT"`
	GRPCHost           string `mapstructure:"GRPC_HOST"`
	HttpHeader         string `mapstructure:"HTTP_HEADER"`
}

func NewEnv() *Env {
	once.Do(func() {
		viperInstance = viper.New()

		env = &Env{}

		viperInstance.AutomaticEnv()

		if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
			log.Infof("[ENV] Using %s environment variables", appEnv)

			// Unmarshal configuration with enviper due to issue with viper
			if err := enviper.New(viperInstance).Unmarshal(env); err != nil {
				log.Fatalf("[ENV] failed to unmarshal configuration: %s", err.Error())
				return
			}
		} else {
			// If APP_ENV not found in environment, try .env file
			if _, err := os.Stat(".env"); err != nil {
				log.Fatal("[ENV] APP_ENV is not set in environment variables")
				return
			}

			viperInstance.SetConfigFile(".env")
			if err := viperInstance.ReadInConfig(); err != nil {
				log.Fatalf("[ENV] Failed to read .env file")
				return
			}

			log.Infof("[ENV] Using .env file")

			// Unmarshal configuration
			if err := viperInstance.Unmarshal(env); err != nil {
				log.Fatalf("[ENV] failed to unmarshal configuration: %s", err.Error())
				return
			}
		}

		env.JwtSecretKey = []byte(viperInstance.GetString("JWT_SECRET_KEY"))

		if err := parseDurations(env); err != nil {
			log.Fatal("failed to parse jwt duration", err)
		}
	})

	return env
}

func GetEnv() *Env {
	return env
}

func parseDurations(env *Env) error {
	var err error

	env.JwtExpiredTime, err = time.ParseDuration(viperInstance.GetString("JWT_EXPIRED_TIME"))
	if err != nil {
		return fmt.Errorf("invalid JWT_ACCESS_EXPIRE_DURATION: %w", err)
	}

	return nil
}
