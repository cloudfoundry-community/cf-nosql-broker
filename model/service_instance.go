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

// ProvisionBody represents the expected request body for provisioning.
type ProvisionBody struct {
	ServiceID      string `json:"service_id"`
	PlanID         string `json:"plan_id"`
	OrganizationID string `json:"organization_guid"`
	SpaceID        string `json:"space_guid"`
}

// ProvisionResponse could be populated with the URL of a web-based portal for the service instance management.
type ProvisionResponse struct {
	DashboardURL string `json:"dashboard_url, omitempty"`
	Operation    string `json:"operation, omitempty"`
}

// DeprovisionResponse expected {} but may return an identifier representing the operation.
type DeprovisionResponse struct {
	Operation string `json:"operation, omitempty"`
}

// ErrorResponse represents the error response during the provision and deprovision implementation.
type ErrorResponse struct {
	Description string `json:"description"`
}
