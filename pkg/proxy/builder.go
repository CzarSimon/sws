package proxy

import (
	"fmt"
)

const baseConf = "worker_processes 1;\n\npid /var/run/nginx.pid;\n\nevents {\n\tworker_connections 1024;\n}\n\nhttp {\n\tinclude mime.types;\n\tdefault_type application/octet-stream;\n\tsendfile on;\n\n\tserver {\n\t\tlisten 18080;\n\n%s\t\tlocation = %s {\n\t\t\treturn 200;\n\t\t\taccess_log off;\n\t\t}\n\t}\n}"

// CreateConf creates proxy configuration.
func (p Proxy) CreateConf(healthRoute string) string {
	return fmt.Sprintf(baseConf, p.Definition, healthRoute)
}
