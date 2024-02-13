/*
Copyright 2023 The Kubernetes Authors.


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

package metrics

import (
	"testing"

	computealpha "google.golang.org/api/compute/v0.alpha"
	computebeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	gce "sigs.k8s.io/gcp-compute-persistent-disk-csi-driver/pkg/gce-cloud-provider/compute"
)

const (
	hyperdiskBalanced = "hyperdisk-balanced"
)

func CreateDiskWithConfidentialCompute(betaVersion bool, confidentialCompute bool, diskType string) *gce.CloudDisk {
	if betaVersion {
		return gce.CloudDiskFromBeta(&computebeta.Disk{
			EnableConfidentialCompute: confidentialCompute,
			Type:                      diskType,
		})
	}
	return gce.CloudDiskFromV1(&compute.Disk{
		Type: diskType,
	})
}

func CreateDiskWithStoragePool(storagePool string, diskType string) *gce.CloudDisk {
	return gce.CloudDiskFromAlpha(&computealpha.Disk{
		StoragePool: storagePool,
		Type:        diskType,
	})
}

func TestGetMetricParameters(t *testing.T) {
	testCases := []struct {
		name                              string
		disk                              *gce.CloudDisk
		expectedEnableConfidentialCompute string
		expectedDiskType                  string
		expectedEnableStoragePools        string
	}{
		{
			name:                              "test betaDisk with enableConfidentialCompute=false",
			disk:                              CreateDiskWithConfidentialCompute(true, false, hyperdiskBalanced),
			expectedEnableConfidentialCompute: "false",
			expectedDiskType:                  hyperdiskBalanced,
			expectedEnableStoragePools:        "false",
		},
		{
			name:                              "test betaDisk with enableConfidentialCompute=true",
			disk:                              CreateDiskWithConfidentialCompute(true, true, hyperdiskBalanced),
			expectedEnableConfidentialCompute: "true",
			expectedDiskType:                  hyperdiskBalanced,
			expectedEnableStoragePools:        "false",
		},
		{
			name:                              "test disk without enableConfidentialCompute",
			disk:                              CreateDiskWithConfidentialCompute(false, false, hyperdiskBalanced),
			expectedEnableConfidentialCompute: "false",
			expectedDiskType:                  hyperdiskBalanced,
			expectedEnableStoragePools:        "false",
		},
		{
			name:                              "test alphaDisk with storage pool projects/my-project/zone/us-central1-a/storagePools/sp1",
			disk:                              CreateDiskWithStoragePool("projects/my-project/zone/us-central1-a/storagePools/sp1", hyperdiskBalanced),
			expectedEnableConfidentialCompute: "false",
			expectedDiskType:                  hyperdiskBalanced,
			expectedEnableStoragePools:        "true",
		},
		{
			name:                              "test alphaDisk with no storage pool",
			disk:                              CreateDiskWithStoragePool("", hyperdiskBalanced),
			expectedEnableConfidentialCompute: "false",
			expectedDiskType:                  hyperdiskBalanced,
			expectedEnableStoragePools:        "false",
		},
		{
			name:                              "test nil disk",
			disk:                              nil,
			expectedEnableConfidentialCompute: DefaultEnableConfidentialCompute,
			expectedDiskType:                  DefaultDiskTypeForMetric,
			expectedEnableStoragePools:        DefaultEnableStoragePools,
		},
	}

	for _, tc := range testCases {
		t.Logf("Running test: %v", tc.name)
		diskType, confidentialCompute, enableStoragePools := GetMetricParameters(tc.disk)
		if confidentialCompute != tc.expectedEnableConfidentialCompute {
			t.Fatalf("Got confidentialCompute value %q expected %q", confidentialCompute, tc.expectedEnableConfidentialCompute)
		}
		if diskType != tc.expectedDiskType {
			t.Fatalf("Got diskType value %q expected %q", diskType, tc.expectedDiskType)
		}
		if enableStoragePools != tc.expectedEnableStoragePools {
			t.Fatalf("Got enableStoragePools value %q expected %q", enableStoragePools, tc.expectedEnableStoragePools)
		}
	}
}
