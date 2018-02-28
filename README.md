# SWS #

Tool for running multiple outward facing services in containers on a single host and proxy to them based on domain.
Supports declaritve description of services, updating a service without interupting others.

## swsctl #
For interacting with the system on the remote host, the command line tool swsctl is provided. It supports applying service definitons, listing and deleting services.

`$ swsctl apply <service-definition-file>  `
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

The architecture consists of the following components

#### sws-loadbalancer ####
Front facing loadbalancer between the two active proxy servers, used to allow updating of proxy configuration without disrupting current traffic. Uses least connections as a loadbalancing algorithm. Also acts as an edge proxy for the sws-apiserver (described below).

#### sws-proxy-[1/2] ####
Service proxies for the running services. Routes traffic to a specific service based soley on the domain of the incomming requests. When new services are deployed or old ones removed the proxy configuration will be updated automatically, while always keeping one proxy active during the update.

#### sws-confdb ####
Configuration store for the services running on the host, information about the current proxy and loadbalancer configuration as well as user information and other system metadata.

#### sws-apiserver ####
Apiserver that exposes the capability to schedule services for deployment, upgrade and deletion. The apiserver only alters the configuration in sws-confdb, but does not apply the changes on the host.

#### sws-agent ####
Periodically runnig job that check for changed in the conifguration store and applies them to the host. This component is responsible for starting / stopping / updating services as well as performing updates of the proxy configuration.

#### sws-net ####
All services exept for the sws-agent run on a subnet set up by docker (network driver = bridge). This allows communiction between the components within the subnet while not allowing outside traffic in. Services with ports exposed to the host is the sws-loadbalancer, which at present exposes port 80 for ingress traffic over http and the sws-confdb which exposes port 5432 for the sws-agent to connect to the configuration store from outside the subnet. It is strongly advised to prevent outside traffic other than the sws-agent to the configuraion store as this holds all service configuration.
