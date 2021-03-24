/*
Copyright 2021 The OpenShift Authors.

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

package openstack

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	yaml "gopkg.in/yaml.v3"
)

func TestOpenStackActuatorFixInvalidCACertFile(t *testing.T) {
	const noCACert = `
clouds:
  openstack:
    auth:
      auth_url: http://1.2.3.4:5000
      password: password
      project_domain_name: Default
      project_name: openshift
      user_domain_name: Default
      username: openshift
    identity_api_version: "3"
    region_name: regionOne
    verify: true
`

	const cloudsWithCACert = `
clouds:
  openstack:
    auth:
      auth_url: http://1.2.3.4:5000
      password: password
      project_domain_name: Default
      project_name: openshift
      user_domain_name: Default
      username: openshift
    cacert: %s
    identity_api_version: "3"
    region_name: regionOne
    verify: true
`

	const incorrectCACertFile = "/incorrect/path/to/ca-bundle.pem"

	parseClouds := func(clouds string) (map[string]interface{}, error) {
		parsed := make(map[string]interface{})
		err := yaml.Unmarshal([]byte(clouds), parsed)
		return parsed, err
	}

	yamlDiff := func(a, b string) ([]string, error) {
		aClouds, err := parseClouds(a)
		if err != nil {
			return []string{a}, err
		}
		bClouds, err := parseClouds(b)
		if err != nil {
			return []string{b}, err
		}

		return deep.Equal(aClouds, bClouds), nil
	}

	cloudsIncorrectCACert := fmt.Sprintf(cloudsWithCACert, incorrectCACertFile)
	cloudsCorrectCACert := fmt.Sprintf(cloudsWithCACert, caCertFile)

	// The diff go-test/deep should generate when replacing incorrectCACertFile with caCertFile
	expectedDiff := fmt.Sprintf("map[clouds].map[openstack].map[cacert]: %s != %s", incorrectCACertFile, caCertFile)

	tests := []struct {
		name         string
		arg          string
		expectedDiff []string
		wantErr      bool
	}{
		{"invalidYAML", "\"", nil, true},
		{"noCACert", noCACert, nil, false},
		{"incorrectCACert", cloudsIncorrectCACert, []string{expectedDiff}, false},
		{"correctCACert", cloudsCorrectCACert, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &OpenStackActuator{}
			got, err := a.fixInvalidCACertFile(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenStackActuator.fixInvalidCACertFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// We don't assert anything about the value if we returned an error
			if tt.wantErr {
				return
			}

			diff, err := yamlDiff(tt.arg, got)
			if err != nil {
				t.Errorf("OpenStackActuator.fixInvalidCACertFile() returned invalid YAML, returned = %v, error = %v", diff[0], err)
				return
			}
			if !reflect.DeepEqual(diff, tt.expectedDiff) {
				t.Errorf("OpenStackActuator.fixInvalidCACertFile() = %v, want %v", diff, tt.expectedDiff)
			}
		})
	}
}
