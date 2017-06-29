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

package endpoint

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Start enables the service broker endpoints and specified their handlers.
func Start(port string, tlsConfig *tls.Config) {
	router := mux.NewRouter()
	router.HandleFunc("/v2/catalog", GetCatalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", Provision).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", Bind).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", UnBind).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{instance_id}", Deprovision).Methods("DELETE")

	http.Handle("/", router)

	server := http.Server{
		Addr:      ":" + port,
		TLSConfig: tlsConfig,
	}

	log.Println("Server started, listening on port " + port)
	err := server.ListenAndServeTLS("", "")
	if err != nil {
		log.Printf("[FATAL ERROR] %s \n", err)
	}
}
