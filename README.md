# SWS #

Tool for running multiple outward facing services in containers on a single host and proxy to them based on domain.
Supports declaritve description of services, updating a service without interupting others.

### swsctl ##
For interacting with the system on the remote host, the command line tool swsctl is provided. It supports applying service definitons, listing and deleting services.

`$ swsctl apply <service-definitions-file>  `
Instructs the remote host to run a service an set up proxying and load balancing to it. An applied change may take up to minute to become active.

`$ swsctl delete service-name  `
Deletes a running service, and removes the capability to proxy to it. A deletion instruction may take up to a minute to apply.

`$ swsctl ls  `
List all of a users running services.

`$ swsctl configure  `
Promts the user for a configuration information, such as remote host and access token.

## Service declaration ##
Services are declared in YAML of JSON files and contain an api version reference and a service specification.

Example of a service declaration in YAML:
```yaml
apiVersion: v1
spec:
  name: service-name
  port: 8081
  domain: example.com
  image: my-repo/my-image:somve-version
  volumeMount: /foo/bar
  env:
    - name: KEY_1
      value: foo
    - name: KEY_2
      bar: bar
```

A specification consists of 6 parts, 4 of which are required.

**name**: Name of the service, must be unique. (required) 

**port**: Port on which the service can be addressed. (required)

**domain**: Domain name which will be used to proxy to the service. (required)

**image**: Docker image that will be used to start the service container. (required)

**volumeMount**: Path where a volume would be mounted, if set a peristant volume will be creaede for the servie named <service-name>-volume. (optional)

**env**: List of environment variable to be passed to the service when it starts. (optional)

## Architecture ##
![architecture - page 1-2](https://user-images.githubusercontent.com/9406331/36811484-47e2e3fe-1cce-11e8-80c9-e0ea7a9204eb.png)

