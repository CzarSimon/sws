package main

import (
	"log"
	"os"
	"strconv"

	endpoint "github.com/CzarSimon/go-endpoint"
)

// Name constatns
const (
	SERVICE_NAME            = "SWS_AGENT"
	DB_NAME                 = "SWS_CONFDB"
	DEFAULT_UPDATE_FREQ     = 60
	DEFAULT_UPDATE_FREQ_STR = "60"
	UPDATE_FREQ_KEY         = "SWS_AGENT_UPDATE_FREQ"
)

// Config holds configuration values.
type Config struct {
	db                endpoint.SQLConfig
	UpdateFreqSeconds uint64
}

// GetConfig gets configuration values from the environment.
func GetConfig() Config {
	updateFreqStr := os.Getenv(UPDATE_FREQ_KEY)
	if updateFreqStr == "" {
		updateFreqStr = DEFAULT_UPDATE_FREQ_STR
	}
	updateFreq, err := strconv.ParseUint(updateFreqStr, 10, 64)
	if err != nil {
		log.Println(err)
		updateFreq = DEFAULT_UPDATE_FREQ
	}
	return Config{
		db:                endpoint.NewPGConfig(DB_NAME),
		UpdateFreqSeconds: updateFreq,
	}
}
