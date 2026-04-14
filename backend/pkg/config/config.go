package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	JWTRefreshSecret string
	WebhookSecret    string
	Port             string
	MigrationPath    string
	BaseURL          string
	SMTPHost         string
	SMTPPort         string
	SMTPUser         string
	SMTPPass         string
	EmailFrom        string
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
	KeycloakSecret   string
}

func Load() *Config {
	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://itsm:itsm@localhost:5432/itsm?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "change-me-in-production"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "change-me-refresh-in-production"),
		WebhookSecret:    getEnv("WEBHOOK_SECRET", "change-me-webhook-in-production"),
		Port:             getEnv("PORT", "8080"),
		MigrationPath:    getEnv("MIGRATION_PATH", "file://migrations"),
		BaseURL:          getEnv("BASE_URL", "http://localhost:3000"),
		SMTPHost:         getEnv("APP_EMAIL_SMTP_HOST", ""),
		SMTPPort:         getEnv("APP_EMAIL_SMTP_PORT", "587"),
		SMTPUser:         getEnv("APP_EMAIL_SMTP_USER", ""),
		SMTPPass:         getEnv("APP_EMAIL_SMTP_PASS", ""),
		EmailFrom:        getEnv("APP_EMAIL_SENDER_EMAIL", "noreply@pcsindonesia.com"),
		KeycloakURL:      getEnv("KEYCLOAK_URL", "https://jupyter.pcsindonesia.com"),
		KeycloakRealm:    getEnv("KEYCLOAK_REALM", "sso-internal"),
		KeycloakClientID: getEnv("KEYCLOAK_CLIENT_ID", "itsm-app"),
		KeycloakSecret:   getEnv("KEYCLOAK_CLIENT_SECRET", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
