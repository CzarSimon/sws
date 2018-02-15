// Package service provides a interface and base type for starting and proxying to services.
package service

import (
	"fmt"
	"strings"
	"time"
)

const (
	wwwPrefix = "www."
)

// Interface service interface to specify service types.
type Interface interface {
	ProxySpec() string
	RunCmd() []string
}

// Service holds information to start, identify and proxy traffic to a containarized service.
type Service struct {
	Name        string   `yaml:"name"`
	Port        int      `yaml:"port"`
	Domain      string   `yaml:"domain"`
	Image       string   `yaml:"image"`
	VolumeMount string   `yaml:"volumeMount"`
	Env         []EnvVar `yaml:"env"`
}

// EnvVar key value pair to pass envionment values to a service.
type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// ServiceManifest struct matching the a service deployment manifest file.
type ServiceManifest struct {
	ApiVersion  string  `yaml:"apiVersion"`
	Spec        Interface `yaml:"spec"`
	DateChanged time.Time
}

// ProxySpec creates an server block in Nginx format in order to proxy traffic to a service.
func (s Service) ProxySpec() string {
	serverBlock := "server {\n\tserver_name %s;\n\tlocation / {\n\t\tproxy_pass http://%s:%d;\n\t}\n}"
	return fmt.Sprintf(serverBlock, s.domains(), s.Name, s.Port)
}

// DockerRunCmd creates a docker run command for a service based on its attributes.
func (s Service) RunCmd() []string {
	runCmd := []string{
		"docker", "run", "-d", "--restart", "always", "--name", s.Name,
	}
	runCmd = append(runCmd, s.envVars()...)
	if s.VolumeMount != "" {
		runCmd = append(runCmd, "--mount", s.volumeSpec())
	}
	return append(runCmd, s.Image)
}

// volumeSpec returns argument for hooking up a volume to a container.
func (s Service) volumeSpec() string {
	return fmt.Sprintf("source=%s,target=%s", s.Name+"_volume", s.VolumeMount)
}

// domains returns the domains which should route trafic to a service.
func (s Service) domains() string {
	if strings.HasPrefix(s.Domain, wwwPrefix) {
		shortDomain := strings.Replace(s.Domain, wwwPrefix, "", 1)
		return fmt.Sprintf("%s %s", shortDomain, s.Domain)
	}
	return fmt.Sprintf("%s %s", s.Domain, wwwPrefix+s.Domain)
}

// envVars returns a list of arguments to inject environment variables to a service.
func (s Service) envVars() []string {
	envVars := make([]string, 0, 2*len(s.Env))
	for _, envVar := range s.Env {
		envVars = append(envVars, "-e", envVar.Name+"="+envVar.Value)
	}
	return envVars
}
