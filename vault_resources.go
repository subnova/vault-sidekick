/*
Copyright 2015 Home Office All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// VaultResources is a collection of type resource
type VaultResources struct {
	// an array of resource to retrieve
	items []*VaultResource
}

// Set is the implementation for the parser
// secret:test:file=filename.test,fmt=yaml
func (r *VaultResources) Set(value string) error {
	rn := defaultVaultResource()

	// step: split on the separator, default ':'
	sep := getEnv("VAULT_SIDEKICK_SEPARATOR", ":")
	items := strings.Split(os.ExpandEnv(value), sep)
	if len(items) < 2 {
		return fmt.Errorf("invalid resource, must have at least two sections TYPE:PATH")
	}
	if len(items) > 3 {
		return fmt.Errorf("invalid resource, can only has three sections, TYPE:PATH[:OPTIONS]")
	}
	if items[0] == "" || items[1] == "" {
		return fmt.Errorf("invalid resource, neither type or path can be empty")
	}

	// step: extract the elements
	rn.resource = items[0]
	rn.path = items[1]
	rn.options = make(map[string]string, 0)

	// step: extract any options
	if len(items) > 2 {
		for _, x := range strings.Split(items[2], ",") {
			kp := strings.Split(x, "=")
			if len(kp) != 2 {
				return fmt.Errorf("invalid resource option: %s, must be KEY=VALUE", x)
			}
			if kp[1] == "" {
				return fmt.Errorf("invalid resource option: %s, must have a value", x)
			}
			// step: set the name and value
			name := strings.TrimSpace(kp[0])
			value := strings.Replace(kp[1], "|", ",", -1)

			// step: extract the control options from the path resource parameters
			switch name {
			case optionMode:
				if !strings.HasPrefix(value, "0") {
					value = "0" + value
				}
				if len(value) != 4 {
					return errors.New("the file permission invalid, should be octal 0444 or alike")
				}
				v, err := strconv.ParseUint(value, 0, 32)
				if err != nil {
					return errors.New("invalid file permissions on resource")
				}
				rn.fileMode = os.FileMode(v)
			case optionFormat:
				if matched := resourceFormatRegex.MatchString(value); !matched {
					return fmt.Errorf("unsupported output format: %s", value)
				}
				rn.format = value
			case optionUpdate:
				duration, err := time.ParseDuration(value)
				if err != nil {
					return fmt.Errorf("update option: %s is not value, should be a duration format", value)
				}
				rn.update = duration
			case optionRevoke:
				choice, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("the revoke option: %s is invalid, should be a boolean", value)
				}
				rn.revoked = choice
			case optionsRevokeDelay:
				duration, err := time.ParseDuration(value)
				if err != nil {
					return fmt.Errorf("the revoke delay option: %s is not value, should be a duration format", value)
				}
				rn.revokeDelay = duration
			case optionRenewal:
				choice, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("the renewal option: %s is invalid, should be a boolean", value)
				}
				rn.renewable = choice
			case optionCreate:
				choice, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("the create option: %s is invalid, should be a boolean", value)
				}
				if rn.resource != "secret" {
					return fmt.Errorf("the create option is only supported for 'cn=secret' at this time")
				}
				rn.create = choice
			case optionSize:
				size, err := strconv.ParseInt(value, 10, 16)
				if err != nil {
					return fmt.Errorf("the size option: %s is invalid, should be an integer", value)
				}
				rn.size = size
			case optionExec:
				rn.execPath = value
			case optionFilename:
				rn.filename = value
			case optionTemplatePath:
				rn.templateFile = value
				rn.format = "tpl"
			case optionMaxRetries:
				maxRetries, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					return fmt.Errorf("the retries option: %s is invalid, should be an integer", value)
				}
				rn.maxRetries = int(maxRetries)
			case optionMaxJitter:
				maxJitter, err := time.ParseDuration(value)
				if err != nil {
					return fmt.Errorf("the jitter option: %s is invalid, should be in duration format", value)
				}
				rn.maxJitter = maxJitter
			default:
				rn.options[name] = value
			}
		}
	}
	// step: append to the list of resources
	r.items = append(r.items, rn)

	return nil
}

// String returns a string representation of the struct
func (r VaultResources) String() string {
	return ""
}
