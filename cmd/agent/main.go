package main

import (
	"database/sql"
	"log"

	"github.com/jasonlvhit/gocron"
	_ "github.com/lib/pq"
)

// Env database connection and configuration.
type Env struct {
	DB    *sql.DB
	Agent *AgentMetadata
}

// SetupEnv sets up environment based on config.
func SetupEnv(config Config) *Env {
	db, err := config.db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	return &Env{
		DB: db,
	}
}

func main() {
	config := GetConfig()
	env := SetupEnv(config)
	defer env.DB.Close()
	err := env.BootupAgent()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting %s with update frequency: %d s.\n",
		SERVICE_NAME, config.UpdateFreqSeconds)

	gocron.Every(config.UpdateFreqSeconds).Seconds().Do(env.ReconcileServices)
	<-gocron.Start()
}
