package nginx

import "fmt"

const baseConf = `
worker_processes 1;

pid /var/run/nginx.pid;

events {
  worker_connections 1024;
}

http {
  include       mime.types;
  default_type  application/octet-stream;

  sendfile     on;

  %s
}`

const baseImage = `
FROM nginx:1.13.8-alpine

WORKDIR /etc/nginx
COPY %s nginx.conf
RUN nginx -t
`

// MakeConf creates nginx configuration.
func MakeConf(serverBlock string) string {
	return fmt.Sprintf(baseConf, serverBlock)
}

// MakeDockerfile creates dockerfile specification.
func MakeDockerfile(confSource string) string {
	return fmt.Sprintf(baseImage, confSource)
}
