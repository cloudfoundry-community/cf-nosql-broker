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
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/cloudfoundry-community/cf-nosql-broker/model"
	"github.com/gorilla/mux"
)

const (
	provisionError        = "Error creating the database service."
	deprovisionError      = "Error deleting the database service."
	errorEmptyBodyRequest = "Please send a request body."
)

// GetCatalog returns the NoSQL database services offered.
func GetCatalog(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REQUEST] Getting catalog "+
		"{ Hostname: %s, URI: %s, Method: %s, Agent: %s } \n",
		r.RemoteAddr, r.RequestURI, r.Method, r.UserAgent())

	plans := []model.ServicePlan{
		{
			Name:        "Standard",
			ID:          "4f79aa95-b5ca-4030-a263-c58cb2c61dfc",
			Description: "MongoDB database",
			Metadata:    nil,
			Free:        true,
			Bindable:    true,
		},
	}

	services := []model.Service{
		{
			Name:            "MongoDB",
			ID:              "011ca270-ad21-44e2-95d6-60c70a840a80",
			Description:     "MongoDB database service based on Docker containers",
			Tags:            []string{"database", "no-sql", "container-based"},
			Requires:        []string{},
			Bindable:        true,
			Metadata:        nil,
			DashboardClient: nil,
			PlanUpdateable:  true,
			Plans:           plans,
		},
	}

	catalog := model.Catalog{
		Services: services,
	}

	log.Println("[RESPONSE] OK: Service catalog fetched.")
	writeResponse(w, http.StatusOK, catalog)
}

// Provision starts the database service creation using Docker commands.
func Provision(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REQUEST] Provisioning a database service "+
		"{ Hostname: %s, URI: %s, Method: %s, Agent: %s } \n",
		r.RemoteAddr, r.RequestURI, r.Method, r.UserAgent())

	instanceID := mux.Vars(r)["instance_id"]
	command := "docker"

	var body *model.ProvisionBody
	bodyErr := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close() // nolint: errcheck

	if bodyErr == io.EOF {
		log.Println("[RESPONSE] Error: " + bodyErr.Error())
		response := model.ErrorResponse{
			Description: errorEmptyBodyRequest,
		}
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	err := validateProvisionInputs(body, instanceID)
	if err != nil {
		log.Println("[RESPONSE] Error: " + err.Error())
		response := model.ErrorResponse{
			Description: err.Error(),
		}
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	port, err := nextAvailablePort()

	if err != nil {
		log.Println("[RESPONSE] Error: " + err.Error())
		response := model.ErrorResponse{
			Description: provisionError,
		}
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	containerPort := port + ":27017"
	containerName := "cf-mongo-" + instanceID
	_, err = exec.Command(command, "run", "-d",
		"--name", containerName,
		"-p", containerPort, "mongo").Output()

	if err != nil {
		log.Println("[RESPONSE] Error in " +
			"[" + command + "] run : " + err.Error())
		response := model.ErrorResponse{
			Description: provisionError,
		}
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := model.ProvisionResponse{
		DashboardURL: "https://dashboard.example.com",
		Operation:    "task_01",
	}

	log.Println("[RESPONSE] Created: The database service " + containerName +
		" has been created successfully.")
	writeResponse(w, http.StatusCreated, response)
}

// Bind associates the database service to a specific application.
func Bind(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REQUEST] Binding a service instance "+
		"{ Hostname: %s, URI: %s, Method: %s, Agent: %s } \n",
		r.RemoteAddr, r.RequestURI, r.Method, r.UserAgent())

	instanceID := mux.Vars(r)["instance_id"]
	bindingID := mux.Vars(r)["binding_id"]

	var body *model.BindBody
	bodyErr := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close() // nolint: errcheck

	if bodyErr == io.EOF {
		log.Println("[RESPONSE] Error: " + bodyErr.Error())
		response := model.ErrorResponse{
			Description: errorEmptyBodyRequest,
		}
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	err := validateBindInputs(body, instanceID, bindingID)
	if err != nil {
		log.Println("[RESPONSE] Error: " + err.Error())
		response := model.ErrorResponse{
			Description: err.Error(),
		}
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	//TODO: Implement bind functionality

	//TODO: Assign user credentials to the server instance specified
	//TODO: Validate if the server instance exists - returns 404
	//TODO: Get server instance Hostname

	credentials := model.Credentials{
		ConnectionString: "MONGO_URL=mongodb://" + body.Database.UserName + ":" +
			body.Database.Password + "@HOSTNAME:27017/" + body.Database.Name,
		UserName:     body.Database.UserName,
		Password:     body.Database.Password,
		Hostname:     "",
		DatabaseName: body.Database.Name,
	}

	response := model.BindResponse{
		Credentials: credentials,
	}

	log.Println("[RESPONSE] Bind: The credentials for the service " +
		instanceID + " has been generated successfully.")
	writeResponse(w, http.StatusCreated, response)
}

// UnBind deletes any resources associated with the binding.
func UnBind(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REQUEST] Unbinding a service instance "+
		"{ Hostname: %s, URI: %s, Method: %s, Agent: %s } \n",
		r.RemoteAddr, r.RequestURI, r.Method, r.UserAgent())

	instanceID := mux.Vars(r)["instance_id"]
	bindingID := mux.Vars(r)["binding_id"]
	serviceID := r.FormValue("service_id")
	planID := r.FormValue("plan_id")

	err := validateUnBindInputs(instanceID, bindingID, serviceID, planID)
	if err != nil {
		log.Println("[RESPONSE] Error: " + err.Error())
		response := model.ErrorResponse{
			Description: err.Error(),
		}
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	//TODO: Implement unbinding

	response := "{}"
	log.Println("[RESPONSE] Unbind: The resources associated to the service " +
		instanceID + " has been deleted successfully.")
	writeResponse(w, http.StatusOK, response)
}

// Deprovision destroy the container where the database service is running.
func Deprovision(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REQUEST] Deprovisioning a database service "+
		"{ Hostname: %s, URI: %s, Method: %s, Agent: %s } \n",
		r.RemoteAddr, r.RequestURI, r.Method, r.UserAgent())

	instanceID := mux.Vars(r)["instance_id"]
	serviceID := r.FormValue("service_id")
	planID := r.FormValue("plan_id")

	command := "docker"
	containerName := "cf-mongo-" + instanceID

	err := validateDeprovisionInputs(instanceID, serviceID, planID)
	if err != nil {
		log.Println("[RESPONSE] Error: " + err.Error())
		response := model.ErrorResponse{
			Description: err.Error(),
		}
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	//TODO: Validate that the server instance exists - returns 404

	_, err = exec.Command(command, "rm", "-f", containerName).Output()

	if err != nil {
		log.Println("[RESPONSE] Error in [" + command + "]: " + err.Error())
		response := model.ErrorResponse{
			Description: deprovisionError,
		}
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response := model.DeprovisionResponse{
		Operation: "task_01",
	}

	log.Println("[RESPONSE] Destroyed: The database service " + containerName +
		" has been deleted successfully.")
	writeResponse(w, http.StatusOK, response)
}

// writeResponse builds the response object and status code and sends it
// to Cloud Foundy in JSON format.
func writeResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data) // nolint: errcheck
}

// nextAvailablePort finds an available port to be used by Docker to run the
// container.
func nextAvailablePort() (string, error) {
	initHostPort := "59000"

	portsTaken, err := exec.Command("bash", "-c",
		"docker ps --format '{{.Ports}}' | grep -oE ':[^-]+' | cut -c2-").Output()

	if err != nil {
		return "", errors.New("[docker] ps error: " + err.Error())
	}

	ports := strings.Fields(string(portsTaken))

	if len(ports) == 0 {
		return initHostPort, nil
	}

	sort.Strings(ports)
	port, _ := strconv.Atoi(ports[len(ports)-1])
	port = port + 1

	return strconv.Itoa(port), nil
}
