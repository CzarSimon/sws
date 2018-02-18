package main

import endpoint "github.com/CzarSimon/go-endpoint"

// Name constatns
const (
	SERVER_NAME = "SWS_API_SERVER"
	DB_NAME     = "SWS_CONFDB"
)

// Config holds configuration values.
type Config struct {
	db     endpoint.SQLConfig
	server endpoint.ServerAddr
}

// GetConfig gets configuration values from the environment.
func GetConfig() Config {
	return Config{
		db:     endpoint.NewPGConfig(DB_NAME),
		server: endpoint.NewServerAddr(SERVER_NAME),
	}
}
