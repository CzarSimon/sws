package proxy

import (
	"testing"

	"github.com/CzarSimon/sws/pkg/service"
)

func TestProxySpec(t *testing.T) {
	expectedProxySpec := "server proxy-1:28080"
	spec := getTestProxy().ProxySpec()
	if expectedProxySpec != spec {
		t.Errorf("proxy.ProxySpec wrong. Expected=%s Got=%s", expectedProxySpec, spec)
	}
}

func TestRunCmd(t *testing.T) {
	expectedCmd := []string{
		"docker", "run", "-d", "--network", "sws-net",
		"--name", "proxy-1", "sws/proxy:fce5fe1ac0611f5ac0f9305abf77cfa7",
	}
	p := getTestProxy()
	runCmd := p.RunCmd("sws-net")
	if len(runCmd) != len(expectedCmd) {
		t.Fatalf("proxy.RunCmd wrong: Expected length: %d, Got: %d",
			len(expectedCmd), len(runCmd))
	}
	for i, item := range expectedCmd {
		if runCmd[i] != item {
			t.Errorf("%d - service.DockerRunCmd wrong: Expected: %s, Got: %s",
				i, item, runCmd[i])
		}
	}
}

func getTestProxy() Proxy {
	return New("proxy-1", DefaultPort, []service.Service{getTestService()})
}

func getTestService() service.Service {
	return service.Service{
		Name:        "example-service",
		Port:        8080,
		Domain:      "example.com",
		Image:       "czarsimon/sws/test-image:latest",
		VolumeMount: "/var/lib/sws/data",
		Env: []service.EnvVar{
			service.EnvVar{
				Name:  "EXAMPLE_KEY",
				Value: "example-value",
			},
		},
	}
}
