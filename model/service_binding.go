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

package model

// BindBody represents the expected request body for binding.
type BindBody struct {
	ServiceID string   `json:"service_id"`
	PlanID    string   `json:"plan_id"`
	Database  Database `json:"parameters"`
}

// Database contains the configuration options for database service binding.
type Database struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

// BindResponse contains the credentials that may be used by applications or
// users to access the database service.
type BindResponse struct {
	Credentials interface{} `json:"credentials, omitempty"`
}

// Credentials represents the set of information used by an application or
// a user to utilize the service instance.
type Credentials struct {
	ConnectionString string `json:"connection_string"`
	UserName         string `json:"username"`
	Password         string `json:"password"`
	Hostname         string `json:"hostname"`
	DatabaseName     string `json:"database_name"`
}
