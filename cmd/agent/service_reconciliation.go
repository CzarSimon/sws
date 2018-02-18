package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/CzarSimon/sws/pkg/service"
)

var (
	listUpdatedServiceQuery = getUpdatedServiesQuery()
)

// RunStateReconsciliation triggers reconciliation of services and proxy state.
func (env *Env) RunStateReconsciliation() {
	if env.Agent.Locked {
		return
	}
	log.Println("Reconciling services")
	env.Agent.Lock()
	err := env.reconcileState()
	if err != nil {
		log.Println(err)
		env.Agent.Unlock()
		return
	}
	env.upsertAgentMetadata(true)
}

// reconcileState performs service reconciliation.
func (env *Env) reconcileState() error {
	return env.reconcileServices()
}

// reconcileServices reconciles the state of running services to state configuration.
func (env *Env) reconcileServices() error {
	services, err := getUpdatedServies(env.Agent.LastUpdated, env.DB)
	if err != nil {
		return err
	}
	for i, svc := range services {
		fmt.Printf("%d. - %v\n", i, svc)
	}
	return nil
}

// getUpdatedServies gets the list of updated services since last reconciliation.
func getUpdatedServies(fromTime time.Time, db *sql.DB) ([]service.Service, error) {
	rows, err := db.Query(listUpdatedServiceQuery, fromTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	services, err := constructServiceList(rows)
	if err != nil {
		return nil, err
	}
	for i, svc := range services {
		services[i].Env, err = getServiceEnvVars(svc.Name, db)
		if err != nil {
			return nil, err
		}
	}
	return services, nil
}

// constructServiceList creates a list for services from a resuling list of rows.
func constructServiceList(rows *sql.Rows) ([]service.Service, error) {
	var svc service.Service
	services := make([]service.Service, 0)
	for rows.Next() {
		err := rows.Scan(
			&svc.Name, &svc.Port, &svc.Domain, &svc.Image, &svc.VolumeMount, &svc.Active)
		if err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, nil
}

// getServiceEnvVars gets the environment varables of a service.
func getServiceEnvVars(serviceName string, db *sql.DB) ([]service.EnvVar, error) {
	rows, err := db.Query("SELECT NAME, VALUE FROM ENV WHERE SERVICE = $1", serviceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return constructEnvVarList(rows)
}

// constructEnvVarList creates a list of env vars from a resulting list of rows.
func constructEnvVarList(rows *sql.Rows) ([]service.EnvVar, error) {
	var env service.EnvVar
	envVars := make([]service.EnvVar, 0)
	for rows.Next() {
		err := rows.Scan(&env.Name, &env.Value)
		if err != nil {
			return nil, err
		}
		envVars = append(envVars, env)
	}
	return envVars, nil
}

// -------- Implementation details -------- //

// getUpdatedAgentMetadataQuery gets query to fetch upated services.
func getUpdatedServiesQuery() string {
	return `
		SELECT NAME, PORT, DOMAIN, IMAGE, VOLUME_MOUNT, ACTIVE FROM SERVICE
			WHERE DATE_CHANGED > $1`
}
