package proxy

import (
	"testing"

	"github.com/CzarSimon/sws/pkg/service"
)

func TestProxySpec(t *testing.T) {
	expectedProxySpec := "server proxy-1:18080;"
	spec := getTestProxy().ProxySpec()
	if expectedProxySpec != spec {
		t.Errorf("proxy.ProxySpec wrong. Expected=%s Got=%s", expectedProxySpec, spec)
	}
}

func TestRunCmd(t *testing.T) {
	expectedCmd := []string{
		"docker", "run", "-d", "--network", "sws-net", "--restart", "always",
		"--name", "proxy-1", "sws-proxy:fce5fe1ac0611f5ac0f9305abf77cfa7",
	}
	p := getTestProxy()
	runCmd := p.RunCmd("sws-net")
	if len(runCmd) != len(expectedCmd) {
		t.Fatalf("proxy.RunCmd wrong: Expected length: %d, Got: %d",
			len(expectedCmd), len(runCmd))
	}
	for i, item := range expectedCmd {
		if runCmd[i] != item {
			t.Errorf("%d - proxy.RunCmd wrong: Expected: %s, Got: %s",
				i, item, runCmd[i])
		}
	}
}

func TestPublicRunCmd(t *testing.T) {
	expectedCmd := []string{
		"docker", "run", "-d", "--network", "sws-net", "-p", "18080:18080",
		"--name", "proxy-1", "sws-proxy:fce5fe1ac0611f5ac0f9305abf77cfa7",
	}
	p := getTestProxy()
	runCmd := p.PublicRunCmd("sws-net")
	if len(runCmd) != len(expectedCmd) {
		t.Fatalf("proxy.RunCmd wrong: Expected length: %d, Got: %d",
			len(expectedCmd), len(runCmd))
	}
	for i, item := range expectedCmd {
		if runCmd[i] != item {
			t.Errorf("%d - proxy.PublicRunCmd wrong: Expected: %s, Got: %s",
				i, item, runCmd[i])
		}
	}
}

func TestAlignedPositive(t *testing.T) {
	pair := Pair{
		Primary:   New("proxy-1", DefaultPort, getTestServices()),
		Secondary: New("proxy-2", DefaultPort, getTestServices()),
	}
	err := pair.Alinged()
	if err != nil {
		t.Errorf("pair.Aligned wrong. Expeced no error, Got=%s", err)
	}
}

func TestAlignedNegative(t *testing.T) {
	pair := Pair{
		Primary:   New("proxy-1", DefaultPort, getTestServices()),
		Secondary: New("proxy-2", DefaultPort+1, getTestServices()),
	}
	err := pair.Alinged()
	if err == nil {
		t.Errorf("pair.Aligned wrong. Expeced expected error with different ports")
	}
	pair.Secondary = New("proxy-2", DefaultPort, []service.Service{getTestService()})
	err = pair.Alinged()
	if err == nil {
		t.Errorf("pair.Aligned wrong. Expeced expected error with different ID")
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
