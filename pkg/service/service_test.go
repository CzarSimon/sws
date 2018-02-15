package service

import (
	"testing"
)

func TestVolumeSpec(t *testing.T) {
	expectedVolSpec := "source=example-service_volume,target=/var/lib/sws/data"
	volSpec := getTestService().volumeSpec()
	if volSpec != expectedVolSpec {
		t.Errorf(`service.volumeSpec wrong:
      Expected: %s
      Got: %s`, expectedVolSpec, volSpec)
	}
}

func TestDomains(t *testing.T) {
	s := getTestService()
	expectedDomains := "example.com www.example.com"
	if s.domains() != expectedDomains {
		t.Errorf("service.domain wrong: Expected: %s, Got: %s", expectedDomains, s.domains())
	}
	s.Domain = "www.example.com"
	if s.domains() != expectedDomains {
		t.Errorf("service.domain wrong: Expected: %s, Got: %s", expectedDomains, s.domains())
	}
}

func TestProxySpec(t *testing.T) {
	expectedProxySpec := "server {\n\tserver_name example.com www.example.com;\n\tlocation / {\n\t\tproxy_pass example-service:8080;\n\t}\n}"
	proxySpec := getTestService().ProxySpec()
	if proxySpec != expectedProxySpec {
		t.Errorf(`service.ProxySpec wrong:
      Expected: %s
      Got: %s`, expectedProxySpec, proxySpec)
	}
}

func TestEnvVars(t *testing.T) {
	expectedEnvVars := []string{
		"-e", "EXAMPLE_KEY=example-value",
	}
	s := getTestService()
	envVars := s.envVars()
	if len(envVars) != len(expectedEnvVars) {
		t.Fatalf("service.envVars wrong: Expected length: %d, Got: %d", len(expectedEnvVars), len(envVars))
	}
	for i, item := range expectedEnvVars {
		if envVars[i] != item {
			t.Errorf("%d - service.envVars wrong: Expected: %s, Got: %s",
				i, item, envVars[i])
		}
	}
	s.Env = make([]EnvVar, 0)
	envVars = s.envVars()
	if len(envVars) != 0 {
		t.Fatalf("service.envVars wrong: Expected length: %d, Got: %d", 0, len(envVars))
	}
}

func TestDockerRunCmd(t *testing.T) {
	expectedCmd := []string{
		"docker", "run", "-d", "--restart", "always", "--name", "example-service",
		"-e", "EXAMPLE_KEY=example-value",
		"--mount", "source=example-service_volume,target=/var/lib/sws/data",
		"czarsimon/sws/test-image:latest",
	}
	s := getTestService()
	runCmd := s.DockerRunCmd()
	if len(runCmd) != len(expectedCmd) {
		t.Fatalf("service.DockerRunCmd wrong: Expected length: %d, Got: %d",
			len(expectedCmd), len(runCmd))
	}
	for i, item := range expectedCmd {
		if runCmd[i] != item {
			t.Errorf("%d - service.DockerRunCmd wrong: Expected: %s, Got: %s",
				i, item, runCmd[i])
		}
	}
	s.VolumeMount = ""
	s.Env = make([]EnvVar, 0)
	runCmd = s.DockerRunCmd()
	expectedCmd = []string{
		"docker", "run", "-d", "--restart", "always", "--name", "example-service",
		"czarsimon/sws/test-image:latest",
	}
	if len(runCmd) != len(expectedCmd) {
		t.Fatalf("service.DockerRunCmd wrong: Expected length: %d, Got: %d",
			len(expectedCmd), len(runCmd))
	}
	for i, item := range expectedCmd {
		if runCmd[i] != item {
			t.Errorf("%d - service.DockerRunCmd wrong: Expected: %s, Got: %s",
				i, item, runCmd[i])
		}
	}
}

func getTestService() Service {
	return Service{
		Name:        "example-service",
		Port:        8080,
		Domain:      "example.com",
		Image:       "czarsimon/sws/test-image:latest",
		VolumeMount: "/var/lib/sws/data",
		Env: []EnvVar{
			EnvVar{
				Name:  "EXAMPLE_KEY",
				Value: "example-value",
			},
		},
	}
}
