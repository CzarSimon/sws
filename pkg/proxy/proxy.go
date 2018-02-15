// Package proxy implements types and functions to create and start proxy services.
package proxy

import (
  "fmt"
)

// Proxy infromation needed to create and start a proxy.
type Proxy struct {
  Name string       `yaml:"name"`
  Port int          `yaml:"port"`
  Image string      `yaml:"image"`
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
