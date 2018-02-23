package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CzarSimon/sws/pkg/service"
	"github.com/CzarSimon/sws/pkg/swsutil"
)

var (
	errNoUpdatedServices        = errors.New("No updated services")
	listUpdatedServiceQuery     = getUpdatedServiesQuery()
	checkIfServicesUpdatedQuery = getCheckIfServicesUpdatedQuery()
)

// RunStateReconsciliation triggers reconciliation of services and proxy state.
func (env *Env) RunStateReconsciliation() {
	if env.Agent.Locked {
		return
	}
	env.Agent.Lock()
	err := env.reconcileState()
	if err != nil && err != errNoUpdatedServices {
		log.Println(err)
		env.Agent.Unlock()
		return
	}
	env.upsertAgentMetadata(true)
}

// reconcileState performs service reconciliation.
func (env *Env) reconcileState() error {
	err := env.updatedServicesExist()
	if err != nil {
		return err
	}
	tx, err := env.DB.Begin()
	if err != nil {
		return err
	}
	err = reconcileServices(tx, env.Agent.LastUpdated)
	if err != nil {
		swsutil.RollbackTx(tx)
		return err
	}
	err = updateProxies(tx)
	if err != nil {
		swsutil.RollbackTx(tx)
		return err
	}
	return tx.Commit()
}

// updatedServicesExist checks for any updates services since last update.
func (env *Env) updatedServicesExist() error {
	var count int
	err := env.DB.QueryRow(
		checkIfServicesUpdatedQuery, env.Agent.LastUpdated).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		log.Println(errNoUpdatedServices)
		return errNoUpdatedServices
	}
	return nil
}

// reconcileServices reconciles the state of running services to state configuration.
func reconcileServices(tx *sql.Tx, fromTime time.Time) error {
	services, err := getUpdatedServies(fromTime, tx)
	if err != nil {
		return err
	}
	for i, svc := range services {
		fmt.Printf("%d. - %v\n", i, svc)
		err = stopService(svc, tx)
		if err != nil {
			log.Println(err)
		}
		err = startService(svc)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

// stopService stops service and removes inactive from the database.
func stopService(svc service.Service, tx *sql.Tx) error {
	msg, err := stopAndRemoveContainer(svc.Name)
	if err != nil {
		log.Println(msg)
		return err
	}
	if !svc.Active {
		err = removeServiceFromDB(svc.Name, tx)
		if err != nil {
			return err
		}
	}
	log.Printf("Stopped service \"%s\"\n", svc.Name)
	return nil
}

func stopAndRemoveContainer(name string) (string, error) {
	msg, err := swsutil.RunShellCommand("docker", "stop", name)
	if err != nil && !noSuchContainer(msg) {
		return msg, err
	}
	msg, err = swsutil.RunShellCommand("docker", "rm", name)
	if err != nil && !noSuchContainer(msg) {
		return msg, err
	}
	return fmt.Sprintf("Removed container: %s", name), nil
}

// removeServiceFromDB removes records of a service from the database.
func removeServiceFromDB(svcName string, tx *sql.Tx) error {
	envStmt, err := tx.Prepare("DELETE FROM ENV WHERE SERVICE = $1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer envStmt.Close()
	_, err = envStmt.Exec(svcName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	svcStmt, err := tx.Prepare("DELETE FROM SERVICE WHERE NAME = $1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer svcStmt.Close()
	_, err = svcStmt.Exec(svcName)
	return err
}

// startService starts the supplied service if active.
func startService(svc service.Service) error {
	if !svc.Active {
		log.Printf("Inactive service: \"%s\" not restarted\n", svc.Name)
		return nil
	}
	runCmd := svc.RunCmd(NetworkName)
	_, err := swsutil.RunShellCommand(runCmd[0], runCmd[1:]...)
	if err != nil {
		return err
	}
	log.Printf("Started service \"%s\"\n", svc.Name)
	return nil
}

// getUpdatedServies gets the list of updated services since last reconciliation.
func getUpdatedServies(fromTime time.Time, tx *sql.Tx) ([]service.Service, error) {
	rows, err := tx.Query(listUpdatedServiceQuery, fromTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	services, err := constructServiceList(rows)
	if err != nil {
		return nil, err
	}
	for i, svc := range services {
		services[i].Env, err = getServiceEnvVars(svc.Name, tx)
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
func getServiceEnvVars(serviceName string, tx *sql.Tx) ([]service.EnvVar, error) {
	rows, err := tx.Query("SELECT NAME, VALUE FROM ENV WHERE SERVICE = $1", serviceName)
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

// getCheckIfServicesUpdatedQuery gets the query to check if serices has been
// updated since a specified time.
func getCheckIfServicesUpdatedQuery() string {
	return "SELECT COUNT(*) FROM SERVICE WHERE DATE_CHANGED > $1;"
}

// noSuchContainer checks if an error message contains "no such container".
func noSuchContainer(errorMsg string) bool {
	return strings.Contains(errorMsg, "No such container")
}
