# NoSQL Service Broker for the CLOUD FOUNDRY* Platform.
This service broker allows application developers to setup a NoSQL database service based on Docker containers as a managed service for their Cloud Foundry environment.

#### Supported NoSQL databases
* MongoDB

#### Features
* Advertising database services and plans offered (catalog)
* Provisioning of database instances (create)
* Creation of credentials (bind) - In Progress
* Removal of credentials (unbind) - In Progress
* Deprovisioning of database instances (delete)

## Usage
### Installing dependencies
Assuming you have a valid [Golang](https://golang.org/doc/install) and [Docker](https://docs.docker.com/engine/installation/linux/ubuntu/) environment installed on your Linux system.

Install Gorilla Mux:
```
$ go get github.com/gorilla/mux
```

### Building from source
Get the latest version of the broker:
```
$ go get github.com/cloudfoundry-community/cf-nosql-broker
$ cd $GOPATH/src/github.com/cloudfoundry-community/cf-nosql-broker
```
- To generate a secure binary (PIE-Position Independent Executables), specify this build mode as an argument `-buildmode=pie`.
- Pass the `'-w'` flag to the linker to omit the debug information.

```
$ go build -buildmode=pie -ldflags "-w" -o $GOPATH/bin/nosql-broker
$ $GOPATH/bin/nosql-broker
```

You can specify a port where the broker will run by setting `$CF_NOSQL_BROKER_PORT` as environment variable.

#### Enabling TLS to use HTTPS
In order to establish a secure connection (HTTPS) between Cloud Foundry and the service broker, a x509 encoded RSA certificate will be required to run the broker.

Security requirements for cryptographic:
* RSA as digital signature for the private key certificate.
* Private key certificate must be at least 2048 bit length for RSA.
* Hash functions: SHA256, SHA384 or SHA512.

Using RSA:2048 with SHA384 is the most recommended configuration.

##### Generating certificates
Use OpenSSL tool on Linux to generate certificates. E.g. RSA:2048 signed certificate with SHA384 hash function:

```
$ openssl req -x509 -sha384 -new -nodes -newkey rsa:2048 -keyout key.pem -out cert.pem
```
##### Specifying the location
To send the location of the PEM files to the broker, set their paths as environment variable.
```
$ export CF_NOSQL_BROKER_KEY="/path/to/key.pem"
$ export CF_NOSQL_BROKER_CERT="/path/to/cert.pem"
```

## Broker registration in Cloud Foundry
Login in your Cloud Foundry environment:
```
$ cf login --skip-ssl-validation -a https://api.bosh-lite.com -u admin -p admin
```

Register the service broker:
```
$ cf create-service-broker nosql-broker <USER> <PASSWORD> <http://BROKER-SERVER:BROKER-PORT>
```

Validate the service broker installation:
```
$ cf service-brokers
```

See the database services offered:
```
$ cf service-access
```

Enable the NoSQL database services offered by the broker:
```
$ cf enable-service-access <DATABASE-SERVICE>
```

See the NoSQL database services listed in the Cloud Foundry marketplace:
```
$ cf marketplace
```

For more details, take a look in the [Managing Service Brokers](https://docs.cloudfoundry.org/services/managing-service-brokers.html) documentation.

## API documentation
This project implements the Cloud Foundry [Service Broker API v2.11](https://docs.cloudfoundry.org/services/api.html) specification. For details about its architecture, requests, parameters and responses review the official documentation.

## License
This project is under Apache License 2.0. See LICENSE for details.
