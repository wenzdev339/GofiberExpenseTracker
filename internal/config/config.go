package config

import (
	"os"
)

type Config struct {
	Port            string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	CORSOrigins     string
	JWTSecret       string
	JWTExpireHours  int
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "3000"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "expense_tracker"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
		CORSOrigins:    getEnv("CORS_ORIGINS", "*"),
		JWTSecret:      getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpireHours: 24,
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
