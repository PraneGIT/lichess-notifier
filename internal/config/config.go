package config

import (
	"os"
	"strconv"
	"strings"
	"fmt"
    "github.com/joho/godotenv"
)

type EmailConfig struct {
	From     string
	To       []string
	Password string
	SMTPHost string
	SMTPPort int
}

type Config struct {
	LichessAPIBase string
	Users          []string
	Email          EmailConfig
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file: %w", err)
	}

	smtpPort, err := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	config := Config{
		LichessAPIBase: getEnvOrDefault("LICHESS_API_BASE", "https://lichess.org/api/games/user/"),
		Users:          strings.Split(getEnvOrDefault("LICHESS_USERS", ""), ","),
		Email: EmailConfig{
			From:     getEnvOrDefault("EMAIL_FROM", ""),
			To:       strings.Split(getEnvOrDefault("EMAIL_TO", ""), ","),
			Password: getEnvOrDefault("EMAIL_PASSWORD", ""),
			SMTPHost: getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: smtpPort,
		},
	}

	if err := validateConfig(config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// validateConfig 
func validateConfig(config Config) error {
	if config.Email.From == "" {
		return fmt.Errorf("EMAIL_FROM is required")
	}
	if len(config.Email.To) == 0 || (len(config.Email.To) == 1 && config.Email.To[0] == "") {
		return fmt.Errorf("EMAIL_TO is required")
	}
	if config.Email.Password == "" {
		return fmt.Errorf("EMAIL_PASSWORD is required")
	}
	if len(config.Users) == 0 || (len(config.Users) == 1 && config.Users[0] == "") {
		return fmt.Errorf("LICHESS_USERS is required")
	}
	return nil
}