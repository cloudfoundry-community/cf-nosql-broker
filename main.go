/*
// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
*/

package main

import (
	"crypto/tls"
	"log"
	"os"

	server "github.com/cloudfoundry-community/cf-nosql-broker/endpoint"
	"github.com/cloudfoundry-community/cf-nosql-broker/security"
)

func main() {
	log.Println("NoSQL Service Broker for the CLOUD FOUNDRY* Platform.")
	port := os.Getenv("CF_NOSQL_BROKER_PORT")

	if port == "" {
		port = "8080"
		log.Println("[WARNING] Requires $CF_NOSQL_BROKER_PORT environment variable, defaulting to:" + port)
	}

	// Validate expected key pair files before start their cryptographic validation
	keyFile := os.Getenv("CF_NOSQL_BROKER_KEY")
	certFile := os.Getenv("CF_NOSQL_BROKER_CERT")

	if keyFile == "" || certFile == "" {
		log.Println("[ERROR] Requires $CF_NOSQL_BROKER_KEY and $CF_NOSQL_BROKER_CERT environment variables to start the server.")
		return
	}

	if _, err := os.Stat(keyFile); err != nil {
		log.Println("[ERROR] The key file " + keyFile + " does not exists in the filesystem.")
		return
	}

	if _, err := os.Stat(certFile); err != nil {
		log.Println("[ERROR] The certificate file " + certFile + " does not exists in the filesystem.")
		return
	}

	// Generate the certificate chain to enable TLS
	// Before generation, the public/private key pair will be validated to see if they meet cryptographic requirements
	cert, err := security.GetCertificateChain(certFile, keyFile)
	if err != nil {
		log.Println("[ERROR] " + err.Error())
		return
	}

	// Set TLS configurations
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	security.SetCipherSuites(&tlsConfig)

	// Start the HTTPS server using TLS
	server.Start(port, &tlsConfig)
}
