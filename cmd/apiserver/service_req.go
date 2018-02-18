package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CzarSimon/httputil"
	"github.com/CzarSimon/sws/pkg/service"
	"github.com/CzarSimon/sws/pkg/user"
)

var (
	upsertServiceQuery      = getUpsertServiceQuery()
	upsertEnvVarQuery       = getUpsertEnvVarQuery()
	inactivateServiceQuery  = getInacivateServiceQuery()
	listServicesQuery       = getListServiesQuery()
	noRowsAffectedErr       = errors.New("No rows affected")
	couldNotListServicesErr = errors.New("Could not list services")
)

// HandleServiceRequest handles request releated to services.
func (env *Env) HandleServiceRequest(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		return env.addService(w, r)
	case http.MethodGet:
		return env.listServices(w, r)
	case http.MethodDelete:
		return env.deleteService(w, r)
	default:
		return httputil.MethodNotAllowed
	}
}

// addService adds a service with the requesting user as owner.
func (env *Env) addService(w http.ResponseWriter, r *http.Request) error {
	usr, err := getUser(r, env.DB)
	if err != nil {
		return err
	}
	var svc service.Service
	err = json.NewDecoder(r.Body).Decode(&svc)
	if err != nil {
		return httputil.BadRequest
	}
	err = insertService(svc, usr, env.DB)
	if err != nil {
		log.Println(err)
		return httputil.InternalServerError
	}
	httputil.SendString(w, fmt.Sprintf("Service \"%s\" added\n", svc.Name))
	return nil
}

// listServices list the current running services that the user has access to.
func (env *Env) listServices(w http.ResponseWriter, r *http.Request) error {
	usr, err := getUser(r, env.DB)
	if err != nil {
		return err
	}
	services, err := getServices(usr, env.DB)
	if err != nil {
		log.Println(err)
		return couldNotListServicesErr
	}
	httputil.SendJSON(w, services)
	return nil
}

// deleteService deletes a service given that its present and the user as access to it.
func (env *Env) deleteService(w http.ResponseWriter, r *http.Request) error {
	usr, err := getUser(r, env.DB)
	if err != nil {
		return err
	}
	var svc service.Service
	err = json.NewDecoder(r.Body).Decode(&svc)
	if err != nil {
		return httputil.BadRequest
	}
	err = scheduleServiceRemoval(svc, usr, env.DB)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not remove service: \"%s\"", svc.Name)
	}
	httputil.SendString(w, fmt.Sprintf("Deleted service \"%s\"\n", svc.Name))
	return nil
}

// insertService inserts a service and environment varibles with the supplied user as owner.
func insertService(svc service.Service, usr user.User, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = upsertServiceRecord(svc, usr, tx)
	if err != nil {
		rollbackTx(tx)
		return err
	}
	err = upsertEnvVars(svc, usr, tx)
	if err != nil {
		rollbackTx(tx)
		return err
	}
	return tx.Commit()
}

// upsertServiceRecord inserts service if new updates if already present.
func upsertServiceRecord(svc service.Service, usr user.User, tx *sql.Tx) error {
	stmt, err := tx.Prepare(upsertServiceQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		svc.Name, svc.Port, svc.Domain, svc.Image, svc.VolumeMount, getNow(), usr.Name)
	return err
}

// upsertEnvVars inserts or updates environment varibles related to a service.
func upsertEnvVars(svc service.Service, usr user.User, tx *sql.Tx) error {
	stmt, err := tx.Prepare(upsertEnvVarQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, envVar := range svc.Env {
		_, err = stmt.Exec(envVar.Name, envVar.Value, svc.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// scheduleServiceRemoval schedules a service for removal.
func scheduleServiceRemoval(svc service.Service, usr user.User, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = removeServiceEnvVars(svc.Name, tx)
	if err != nil {
		rollbackTx(tx)
		return err
	}
	err = setServiceToInactive(svc.Name, usr.Name, tx)
	if err != nil {
		rollbackTx(tx)
		return err
	}
	return tx.Commit()
}

// setServiceToInactive sets service to inactive if user owns it.
func setServiceToInactive(svcName, username string, tx *sql.Tx) error {
	stmt, err := tx.Prepare(inactivateServiceQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(getNow(), svcName, username)
	if err != nil {
		return err
	}
	return ensureRowsAffected(res)
}

// removeServiceEnvVars removes a services environment variables.
func removeServiceEnvVars(svcName string, tx *sql.Tx) error {
	stmt, err := tx.Prepare("DELETE FROM ENV WHERE SERVICE = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(svcName)
	return err
}

// getUser gets the user which sent a request.
func getUser(r *http.Request, db *sql.DB) (user.User, error) {
	authKey := r.Header.Get("Authorization")
	if authKey == "" {
		return user.User{}, httputil.BadRequest
	}
	var usr user.User
	query := "SELECT USER_NAME, KEY, VALID_TO FROM ACCESS_KEY WHERE KEY = $1"
	err := db.QueryRow(query, authKey).Scan(&usr.Name, &usr.AccessKey.Key, &usr.AccessKey.ValidTo)
	if err != nil {
		log.Println(err)
		return user.User{}, httputil.InternalServerError
	}
	return usr, nil
}

// getServices gets all services belonging to a user.
func getServices(usr user.User, db *sql.DB) ([]service.Service, error) {
	rows, err := db.Query(listServicesQuery, usr.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return constructServiceList(rows)
}

// constructServiceList creates a list for services from a resuling list of rows.
func constructServiceList(rows *sql.Rows) ([]service.Service, error) {
	var svc service.Service
	services := make([]service.Service, 0)
	for rows.Next() {
		err := rows.Scan(&svc.Name, &svc.Port, &svc.Domain, &svc.Image, &svc.VolumeMount)
		if err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, nil
}

// -------- Implementation details -------- //

// getUpsertServiceQuery gets query to upsert service information.
func getUpsertServiceQuery() string {
	return `
		INSERT INTO SERVICE(
			NAME, PORT, DOMAIN, IMAGE, VOLUME_MOUNT, DATE_CHANGED, ACTIVE, USER_NAME)
			VALUES ($1, $2, $3, $4, $5, $6, 'TRUE', $7)
			ON CONFLICT(NAME) DO UPDATE SET
				PORT = $2,
				DOMAIN = $3,
				IMAGE = $4,
				VOLUME_MOUNT = $5,
				DATE_CHANGED = $6,
				ACTIVE = 'TRUE',
				USER_NAME = $7;`
}

// getUpsertEnvVarQuery gets query to upsert environment variable linked to a service.
func getUpsertEnvVarQuery() string {
	return `
		INSERT INTO ENV(NAME, VALUE, SERVICE)
			VALUES ($1, $2, $3)
			ON CONFLICT(SERVICE, NAME) DO UPDATE SET VALUE = $2`
}

// getInacivateServiceQuery gets query to set a service to inactive.
func getInacivateServiceQuery() string {
	return `
		UPDATE SERVICE SET ACTIVE = 'FALSE', DATE_CHANGED = $1
			WHERE NAME = $2 AND USER_NAME = $3`
}

func getListServiesQuery() string {
	return `
		SELECT NAME, PORT, DOMAIN, IMAGE, VOLUME_MOUNT FROM SERVICE
			WHERE USER_NAME = $1 AND ACTIVE = 'TRUE'`
}

// ensureRowsAffected check that query affected rows, returns an error if not.
func ensureRowsAffected(res sql.Result) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return noRowsAffectedErr
	}
	return nil
}

// rollbackTx attepts to rollback a transaction.
func rollbackTx(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Println(err)
	}
}

// getNow gets the current UTC timestamp.
func getNow() time.Time {
	return time.Now().UTC()
}
