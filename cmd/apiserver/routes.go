package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/CzarSimon/httputil/auth"
	"github.com/CzarSimon/httputil/handler"
	"github.com/CzarSimon/sws/pkg/user"
)

// SetupRoutes sets up routes for apiserver
func SetupRoutes(env *Env) *http.ServeMux {
	mux := http.NewServeMux()
	check := auth.NewWrapper(env.validAccessKey)
	mux.Handle("/v1/service", check.Wrap(handler.New(env.HandleServiceRequest)))
	mux.Handle("/health", handler.HealthCheck)
	return mux
}

// validAccessKey checks if an access key is valid.
func (env *Env) validAccessKey(r *http.Request) bool {
	authKey := r.Header.Get("Authorization")
	accessKey, err := getAccessKey(authKey, env.DB)
	if err != nil {
		log.Println(err)
		return false
	}
	return accessKey.Valid()
}

// getAccessKey checks for access key in database and adds validty information.
func getAccessKey(authKey string, db *sql.DB) (user.AccessKey, error) {
	var accessKey user.AccessKey
	query := "SELECT KEY, VALID_TO FROM ACCESS_KEY WHERE KEY=$1"
	err := db.QueryRow(query, authKey).Scan(&accessKey.Key, &accessKey.ValidTo)
	return accessKey, err
}
