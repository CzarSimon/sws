// Package loadbalancer implements types and functions to
// create and start a fronting loadbalancer
package loadbalancer

import "fmt"

const (
	HttpPort      = 80
	HttpsPort     = 443
	lbName        = "loadbalancer"
	baseImageName = "sws/loadbalancer"
)

// LoadBalancer information needed to create and start a proxy.
type LoadBalancer struct {
	Name       string
	Image      string
	Definition string
}

// RunCmd creates a docker run command for a proxy based on its attributes.
func (lb LoadBalancer) RunCmd() []string {
	portMap := fmt.Sprintf("%d:%d", HttpPort, HttpPort)
	return []string{
		"docker", "run", "-d", "-p", portMap, "--name", lb.Name, lb.Image,
	}
}

// ProxyCmd placeholder method in order to make LoadBalancer be compliant with service.Interface
func (lb LoadBalancer) ProxyCmd() string {
	return ""
}
