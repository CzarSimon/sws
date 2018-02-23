package proxy

import (
	"fmt"
	"io/ioutil"
)

const (
	baseConf = "worker_processes 1;\n\npid /var/run/nginx.pid;\n\nevents {\n\tworker_connections 1024;\n}\n\nhttp {\n\tinclude mime.types;\n\tdefault_type application/octet-stream;\n\tsendfile on;\n\n\tserver {\n\t\tlisten %d;\n\n%s\t\tlocation = %s {\n\t\t\treturn 200;\n\t\t\taccess_log off;\n\t\t}\n\t}\n}"
	fileMode = 0666
)

// CreateConf creates proxy configuration.
func (p Proxy) CreateConf(healthRoute string) string {
	return fmt.Sprintf(baseConf, p.Port, p.Definition, healthRoute)
}

// WriteConf creates a configuration file with the proxy configuration.
func WriteConf(conf, filename string) error {
	return ioutil.WriteFile(filename, []byte(conf), fileMode)
}
