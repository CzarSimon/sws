// Package proxy implements types and functions to create and start proxy services.
package proxy

import (
	"bytes"
	"crypto/md5"
	"fmt"

	"github.com/CzarSimon/sws/pkg/service"
)

const (
	DefaultPort   = 28080
	baseImageName = "sws/proxy"
)

// Proxy infromation needed to create and start a proxy.
type Proxy struct {
	Name       string `yaml:"name"`
	Port       int    `yaml:"port"`
	Image      string `yaml:"image"`
	Definition string `yaml:"definition"`
}

// RunCmd creates a docker run command for a proxy based on its attributes.
func (p Proxy) RunCmd() []string {
	portMap := fmt.Sprintf("%d:%d", p.Port, p.Port)
	return []string{
		"docker", "run", "-d", "-p", portMap, "--name", p.Name, p.Image,
	}
}

// ProxySpec creates an server block in Nginx format in order to proxy traffic to a proxy.
func (p Proxy) ProxySpec() string {
	return fmt.Sprintf("server %s:%d", p.Name, p.Port)
}

// ID returns a checksum of the proxy definition as an identifier.
func (p Proxy) ID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(p.Definition)))
}

// New builds a proxy struct based on definitions of backend services.
func New(name string, port int, services []service.Service) Proxy {
	proxy := Proxy{
		Name:       name,
		Port:       port,
		Definition: buildDefinition(services),
	}
	proxy.Image = baseImageName + ":" + proxy.ID()
	return proxy
}

// buildDefinition creats a specification of backend services to proxy to.
func buildDefinition(services []service.Service) string {
	var def bytes.Buffer
	for _, s := range services {
		def.WriteString(fmt.Sprintf("%s\n\n", s.ProxySpec()))
	}
	return def.String()
}