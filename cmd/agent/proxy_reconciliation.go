package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/CzarSimon/sws/pkg/proxy"
	"github.com/CzarSimon/sws/pkg/service"
	"github.com/CzarSimon/sws/pkg/swsutil"
)

// 11. Remove old image

// Proxy names
const (
	primaryProxy     = "sws-proxy-1"
	secondaryProxy   = "sws-proxy-2"
	initalProxyImage = "czarsimon/sws-proxy:init"
)

// Queries to handle proxy information in the database.
var (
	selectProxyQuery          = getSelectProxyQuery()
	selectActiveServicesQuery = getSelectAllServicesQuery()
	updateProxyQuery          = getUpdateProxyQuery()
)

// updateProxies updates the proxy configuration and redeployes both service proxies.
func updateProxies(tx *sql.Tx) error {
	proxies, err := getProxyConfigs(tx)
	if err != nil {
		return err
	}
	oldImage := proxies.Primary.Image
	defer removeImage(oldImage)
	newProxy, err := createUpdatedProxy(tx)
	if err != nil {
		return err
	}
	err = deployUpdatedProxy(&proxies, newProxy)
	if err != nil {
		return err
	}
	return updateProxyPairInDB(proxies, tx)
}

// removeImage removes a replaced proxy container image.
func removeImage(image string) error {
	if image == initalProxyImage {
		return nil
	}
	_, err := swsutil.RunShellCommand("docker", "rmi", image)
	return err
}

// deployUpdatedProxy performs a roling update of the primary and secondary proxy.
func deployUpdatedProxy(pair *proxy.Pair, newProxy proxy.Proxy) error {
	err := toggleProxy(pair.Primary, newProxy)
	if err != nil {
		return err
	}
	primaryName := pair.Primary.Name
	pair.Primary = newProxy
	pair.Primary.Name = primaryName
	time.Sleep(1 * time.Second)
	err = toggleProxy(pair.Secondary, newProxy)
	if err != nil {
		return err
	}
	secondaryName := pair.Secondary.Name
	pair.Secondary = newProxy
	pair.Secondary.Name = secondaryName
	return nil
}

// toggleProxy replaces an an old proxy with a new.
func toggleProxy(old, new proxy.Proxy) error {
	msg, err := stopAndRemoveContainer(old.Name)
	if err != nil {
		log.Println(msg)
		return err
	}
	new.Name = old.Name
	runCmd := new.RunCmd(NetworkName)
	msg, err = swsutil.RunShellCommand(runCmd[0], runCmd[1:]...)
	if err != nil {
		log.Println(msg)
	}
	return err
}

// createUpdatedProxy creates a proxy config form all updates services.
func createUpdatedProxy(tx *sql.Tx) (proxy.Proxy, error) {
	services, err := getActiveServices(tx)
	if err != nil {
		return proxy.Proxy{}, err
	}
	candidate := proxy.New("sws-proxy-candidate", proxy.DefaultPort, services)
	err = testNewProxy(candidate)
	_, undeployErr := stopAndRemoveContainer(candidate.Name)
	if err != nil {
		return proxy.Proxy{}, err
	}
	return candidate, undeployErr
}

// getProxyConfig gets the proxy configuration of the current primary and secondary proxies
func getProxyConfigs(tx *sql.Tx) (proxy.Pair, error) {
	pair := proxy.Pair{}
	p, err := getProxyConfig(primaryProxy, tx)
	if err != nil {
		return pair, err
	}
	pair.Primary = p
	p, err = getProxyConfig(secondaryProxy, tx)
	if err != nil {
		return pair, err
	}
	pair.Secondary = p
	return pair, nil
}

// getProxyConfig gets the stored configuration of a specified proxy.
func getProxyConfig(name string, tx *sql.Tx) (proxy.Proxy, error) {
	var p proxy.Proxy
	err := tx.QueryRow(selectProxyQuery, name).Scan(&p.Name, &p.Image)
	if err != nil {
		return p, err
	}
	return p, nil
}

// getActiveServices queries the database for all active serices.
func getActiveServices(tx *sql.Tx) ([]service.Service, error) {
	rows, err := tx.Query(selectActiveServicesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var svc service.Service
	services := make([]service.Service, 0)
	for rows.Next() {
		err = rows.Scan(&svc.Name, &svc.Domain, &svc.Port)
		if err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, nil
}

// testNewProxy deploys a new proxy candidate in test mode and verifies it.
func testNewProxy(candidate proxy.Proxy) error {
	err := buildAndDeployNewProxy(candidate)
	if err != nil {
		return err
	}
	startPause := 1 * time.Second
	time.Sleep(startPause)
	return healthCheckProxy(candidate)
}

// buildAndDeployNewProxy builds and deploys the candidate proxy
// as a container accessable on the local network.
func buildAndDeployNewProxy(candidate proxy.Proxy) error {
	conf := candidate.CreateConf(proxy.HealthRoute)
	err := proxy.WriteConf(conf, filepath.Join(ProxyFolder, "nginx.conf"))
	if err != nil {
		return err
	}
	msg, err := swsutil.RunShellCommand("docker", "build", "-t", candidate.Image, ProxyFolder)
	if err != nil {
		log.Println(msg)
		return err
	}
	runCmd := candidate.PublicRunCmd(NetworkName)
	msg, err = swsutil.RunShellCommand(runCmd[0], runCmd[1:]...)
	log.Println(msg)
	return err
}

// healthCheckProxy pings a proxy to check its health.
func healthCheckProxy(candidate proxy.Proxy) error {
	healthURL := fmt.Sprintf("http://localhost:%d%s", candidate.Port, proxy.HealthRoute)
	res, err := getHttpClient(200).Get(healthURL)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Health check error 2")
		return fmt.Errorf("Proxy candidate health check failed. StatusCode: %d", res.StatusCode)
	}
	log.Printf("Proxy candidte healthcheck successfull. Code: %d", res.StatusCode)
	return res.Body.Close()
}

func updateProxyPairInDB(pair proxy.Pair, tx *sql.Tx) error {
	fmt.Printf("Primary: %+v\n", pair.Primary)
	err := updateProxyInDB(pair.Primary, tx)
	if err != nil {
		return err
	}
	fmt.Printf("Secondary: %+v\n", pair.Secondary)
	return updateProxyInDB(pair.Secondary, tx)
}

func updateProxyInDB(p proxy.Proxy, tx *sql.Tx) error {
	stmt, err := tx.Prepare(updateProxyQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.Port, p.Image, proxy.HealthRoute, p.Name)
	return err
}

// -------- Implementation details -------- //

func getHttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout * time.Millisecond,
	}
}

// getSelectProxyQuery gets the query to select proxy info.
func getSelectProxyQuery() string {
	return `SELECT NAME, IMAGE FROM INGRESS_RESOURCE
            WHERE ROLE = 'PROXY' AND NAME = $1`
}

// getSelectAllServicesQuery gets the query to select all active serices.
func getSelectAllServicesQuery() string {
	return "SELECT NAME, DOMAIN, PORT FROM SERVICE WHERE ACTIVE = 'TRUE'"
}

// getUpdateProxyQuery gets the query to update proxies.
func getUpdateProxyQuery() string {
	return `UPDATE INGRESS_RESOURCE
            SET
              PORT = $1,
              IMAGE = $2 ,
              HEALTH_ROUTE = $3
            WHERE
              ROLE = 'PROXY' AND NAME = $4`
}
