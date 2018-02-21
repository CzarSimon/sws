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
	RunCmd(string) []string
	ID() string
}

// Service holds information to start, identify and proxy traffic to a containarized service.
type Service struct {
	Name        string   `yaml:"name" json:"name"`
	Port        int      `yaml:"port" json:"port"`
	Domain      string   `yaml:"domain" json:"domain"`
	Image       string   `yaml:"image" json:"image"`
	VolumeMount string   `yaml:"volumeMount" json:"volumeMount"`
	Env         []EnvVar `yaml:"env" json:"env"`
	Active      bool     `yaml:"active" json:"active"`
}

// EnvVar key value pair to pass envionment values to a service.
type EnvVar struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

// ServiceManifest struct matching the a service deployment manifest file.
type Manifest struct {
	ApiVersion  string  `yaml:"apiVersion" json:"apiVersion"`
	Spec        Service `yaml:"spec" json:"spec"`
	DateChanged time.Time
}

// ProxySpec creates an server block in Nginx format in order to proxy traffic to a service.
func (s Service) ProxySpec() string {
	serverBlock := "location /%s/ {\n\tproxy_pass http://%s:%d/;\n}"
	return fmt.Sprintf(serverBlock, s.Domain, s.Name, s.Port)
}

// DockerRunCmd creates a docker run command for a service based on its attributes.
func (s Service) RunCmd(network string) []string {
	runCmd := []string{
		"docker", "run", "-d", "--restart", "always", "--network", network, "--name", s.Name,
	}
	runCmd = append(runCmd, s.envVars()...)
	if s.VolumeMount != "" {
		runCmd = append(runCmd, "--mount", s.volumeSpec())
	}
	return append(runCmd, s.Image)
}

// ID returns and identifying string from a service.
func (s Service) ID() string {
	return s.Image
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
