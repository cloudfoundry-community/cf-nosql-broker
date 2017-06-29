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

// Catalog contains the databases services and their description.
type Catalog struct {
	Services []Service `json:"services"`
}

// Service represents a databases service offering.
type Service struct {
	Name            string        `json:"name"`
	ID              string        `json:"id"`
	Description     string        `json:"description"`
	Tags            []string      `json:"tags, omitempty"`
	Requires        []string      `json:"requires, omitempty"`
	Bindable        bool          `json:"bindable"`
	Metadata        interface{}   `json:"metadata, omitempty"`
	DashboardClient interface{}   `json:"dashboard_client, omitempty"`
	PlanUpdateable  bool          `json:"plan_updateable, omitempty"`
	Plans           []ServicePlan `json:"plans"`
}

// ServicePlan represents the different plans available for a database service.
type ServicePlan struct {
	Name        string      `json:"name"`
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Metadata    interface{} `json:"metadata, omitempty"`
	Free        bool        `json:"free, omitempty"`
	Bindable    bool        `json:"bindable, omitempty"`
}
