package main

import "flag"

type Config struct {
	DBPath    string
	Port      string
	AdminPort string
}

func LoadConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.DBPath, "db", "./comments.db", "Path to SQLite database file")
	flag.StringVar(&cfg.Port, "port", "7001", "Public API port")
	flag.StringVar(&cfg.AdminPort, "admin-port", "7002", "Admin panel port")
	flag.Parse()
	return cfg
}
