package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Env database connection and configuration.
type Env struct {
	DB     *sql.DB
	config Config
}

// SetupEnv sets up environment based on config.
func SetupEnv(config Config) *Env {
	db, err := config.db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	return &Env{
		DB:     db,
		config: config,
	}
}

// SetupServer sets up server on specified port and with creaeted route hanlder.
func SetupServer(env *Env) *http.Server {
	return &http.Server{
		Addr:    ":" + env.config.server.Port,
		Handler: SetupRoutes(env),
	}
}

func main() {
	config := GetConfig()
	env := SetupEnv(config)

	server := SetupServer(env)

	log.Printf("Starting %s on port: %s\n", SERVER_NAME, config.server.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
