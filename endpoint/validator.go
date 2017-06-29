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
	"errors"
	"regexp"

	"github.com/cloudfoundry-community/cf-nosql-broker/model"
)

const (
	nullError    = "The field is requiered and cannot be null/empty"
	notValidUUID = "The input provided is not a valid UUID."
)

// validateProvisionInputs validates request fields for create a database.
func validateProvisionInputs(body *model.ProvisionBody, instanceID string) error {

	if isNull(body.ServiceID) || isNull(body.PlanID) || isNull(body.OrganizationID) || isNull(body.SpaceID) || isNull(instanceID) {
		return errors.New(nullError)
	}

	if !isUUID(body.ServiceID) || !isUUID(body.PlanID) || !isUUID(body.OrganizationID) || !isUUID(body.SpaceID) || !isUUID(instanceID) {
		return errors.New(notValidUUID)
	}

	return nil
}

// validateDeprovisionInputs validates request parameters for destroy a database.
func validateDeprovisionInputs(instanceID string, serviceID string, planID string) error {

	if isNull(instanceID) || isNull(serviceID) || isNull(planID) {
		return errors.New(nullError)
	}

	if !isUUID(instanceID) || !isUUID(serviceID) || !isUUID(planID) {
		return errors.New(notValidUUID)
	}

	return nil
}

// validateBindInputs validates request fields for binding database services.
func validateBindInputs(body *model.BindBody, instanceID string, bindingID string) error {

	if isNull(body.ServiceID) || isNull(body.PlanID) || isNull(body.Database.Name) || isNull(body.Database.UserName) || isNull(body.Database.Password) ||
		isNull(instanceID) || isNull(bindingID) {
		return errors.New(nullError)
	}

	if !isUUID(body.ServiceID) || !isUUID(body.PlanID) || !isUUID(instanceID) || !isUUID(bindingID) {
		return errors.New(notValidUUID)
	}

	return nil
}

// validateUnBindInputs validates request parameters for unbind a database services.
func validateUnBindInputs(instanceID string, bindingID string, serviceID string, planID string) error {

	if isNull(instanceID) || isNull(bindingID) || isNull(serviceID) || isNull(planID) {
		return errors.New(nullError)
	}

	if !isUUID(instanceID) || !isUUID(bindingID) || !isUUID(serviceID) || !isUUID(planID) {
		return errors.New(notValidUUID)
	}

	return nil
}

// isNull checks if a string is null.
func isNull(str string) bool {
	return len(str) == 0
}

// isUUID checks if the string is a UUID.
func isUUID(str string) bool {
	UUID := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
	return UUID.MatchString(str)
}
