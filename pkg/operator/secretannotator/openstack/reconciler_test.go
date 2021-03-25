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
	"reflect"
	"testing"

	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/cloud-credential-operator/pkg/operator/constants"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileCloudCredSecret_Reconcile(t *testing.T) {
	type fields struct {
		Client client.Client
		Logger log.FieldLogger
	}
	type args struct {
		request reconcile.Request
	}

	infra := &configv1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Status: configv1.InfrastructureStatus{
			Platform:           configv1.GCPPlatformType,
			InfrastructureName: testInfraName,
		},
	}

	existing := append(test.existing, infra)

	tests := [...]struct {
		name             string
		fields           fields
		args             args
		wantReturnResult reconcile.Result
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			existing := append(test.existing, infra)
			fakeClient := fake.NewFakeClient(existing...)

			r := &ReconcileCloudCredSecret{
				Client: fakeClient,
				Logger: log.WithField("controller", "testController"),
			}

			gotReturnResult, err := r.Reconcile(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReconcileCloudCredSecret.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReturnResult, tt.wantReturnResult) {
				t.Errorf("ReconcileCloudCredSecret.Reconcile() = %v, want %v", gotReturnResult, tt.wantReturnResult)
			}
		})
	}
}

func testOperatorConfig(mode operatorv1.CloudCredentialsMode) *operatorv1.CloudCredential {
	return &operatorv1.CloudCredential{
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.CloudCredOperatorConfig,
		},
		Spec: operatorv1.CloudCredentialSpec{
			CredentialsMode: mode,
		},
	}
}
