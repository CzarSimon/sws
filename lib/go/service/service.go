package service

package service

import (
  "fmt"
  "time"
  "strings"
)

const (
  wwwPrefix = "www."
)

type Service struct {
  Name        string   `yaml:"name"`
  Port        int      `yaml:"port"`
  Domain      string   `yaml:"domain"`
  Image       string   `yaml:"image"`
  VolumeMount string   `yaml:"volumeMount"`
  Env         []EnvVar `yaml:"env"`
}

type EnvVar struct {
  Name  string `yaml:"name"`
  Value string `yaml:"value"`
}

type ServiceManifest struct {
  ApiVersion  string  `yaml:"apiVersion"`
  Spec        Service `yaml:"spec"`
  DateChanged time.Time
}

func (s Service) ProxySpec() string {
  return fmt.Sprintf(`
    server {
        server_name %s;
        location / {
          proxy_pass %s:%d;
        }
    }`, s.domains(), s.Name, s.Port);
}

func (s Service) DockerRunCmd() []string {
  runCmd := []string{
    "docker", "run", "-d", "--restart", "always", "--name", s.Name
  }
  runCmd = append(runCmd, ...s.envVars)
  if s.VolumeMount != "" {
    runCmd = append(runCmd, "--mount", s.volumeSpec)
  }
  return append(runCmd, s.Image)
}

func (s Service) volumeSpec() string {
  return fmt.Sprintf("source=%s,target=%s", s.Name+"_volume", s.VolumeMount)
}

func (s Service) domains() string {
  if strings.HasPrefix(s.Domain, wwwPrefix) {
    shortDomain := strings.Replace(s.Domain, wwwPrefix, "1")
    return fmt.Sprintf("%s %s", shortDomain, s.Domain)
  }
  return fmt.Sprintf("%s %s", s.Domain, wwwPrefix+s.Domain)
}

func (s Service) envVars() []string {
  envVars := make([]string, 0, 2 * len(s.Env))
  for _, envVar := range s.Env {
    envVars := append(envVars, "-e", envVar.Name+"="+envVar.Value)
  }
  return envVars
}
