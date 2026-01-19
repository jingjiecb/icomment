package main

import (
	"flag"
	"os"
)

type Config struct {
	DBPath        string
	Port          string
	AdminPort     string
	BarkDeviceKey string
}

func LoadConfig() *Config {
	cfg := &Config{}

	// Helper to get env with default
	getEnv := func(key, fallback string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return fallback
	}

	defaultDB := getEnv("ICOMMENT_DB", "./comments.db")
	defaultPort := getEnv("ICOMMENT_PORT", "7001")
	defaultAdminPort := getEnv("ICOMMENT_ADMIN_PORT", "7002")
	defaultBark := getEnv("ICOMMENT_BARK", "")

	flag.StringVar(&cfg.DBPath, "db", defaultDB, "Path to SQLite database file")
	flag.StringVar(&cfg.Port, "port", defaultPort, "Public API port")
	flag.StringVar(&cfg.AdminPort, "admin-port", defaultAdminPort, "Admin panel port")
	flag.StringVar(&cfg.BarkDeviceKey, "bark", defaultBark, "Bark device key for notifications (optional)")
	flag.Parse()
	return cfg
}
