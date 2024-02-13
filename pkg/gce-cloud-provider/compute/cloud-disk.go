/*
Copyright 2018 The Kubernetes Authors.

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

package gcecloudprovider

import (
	"strings"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	computealpha "google.golang.org/api/compute/v0.alpha"
	computebeta "google.golang.org/api/compute/v0.beta"
	computev1 "google.golang.org/api/compute/v1"
)

type CloudDisk struct {
	disk      *computev1.Disk
	betaDisk  *computebeta.Disk
	alphaDisk *computealpha.Disk
}

type CloudDiskType string

func CloudDiskFromV1(disk *computev1.Disk) *CloudDisk {
	return &CloudDisk{
		disk: disk,
	}
}

func CloudDiskFromBeta(disk *computebeta.Disk) *CloudDisk {
	return &CloudDisk{
		betaDisk: disk,
	}
}

func CloudDiskFromAlpha(disk *computealpha.Disk) *CloudDisk {
	return &CloudDisk{
		alphaDisk: disk,
	}
}

func (d *CloudDisk) LocationType() meta.KeyType {
	var zone, region string
	switch {
	case d.disk != nil:
		zone = d.disk.Zone
		region = d.disk.Region
	case d.betaDisk != nil:
		zone = d.betaDisk.Zone
		region = d.betaDisk.Region
	case d.alphaDisk != nil:
		zone = d.alphaDisk.Zone
		region = d.alphaDisk.Region
	}
	switch {
	case zone != "":
		return meta.Zonal
	case region != "":
		return meta.Regional
	default:
		return meta.Global
	}
}

func (d *CloudDisk) GetUsers() []string {
	switch {
	case d.disk != nil:
		return d.disk.Users
	case d.betaDisk != nil:
		return d.betaDisk.Users
	case d.alphaDisk != nil:
		return d.alphaDisk.Users
	default:
		return nil
	}
}

func (d *CloudDisk) GetName() string {
	switch {
	case d.disk != nil:
		return d.disk.Name
	case d.betaDisk != nil:
		return d.betaDisk.Name
	case d.alphaDisk != nil:
		return d.alphaDisk.Name
	default:
		return ""
	}
}

func (d *CloudDisk) GetKind() string {
	switch {
	case d.disk != nil:
		return d.disk.Kind
	case d.betaDisk != nil:
		return d.betaDisk.Kind
	case d.alphaDisk != nil:
		return d.alphaDisk.Kind
	default:
		return ""
	}
}

func (d *CloudDisk) GetStatus() string {
	switch {
	case d.disk != nil:
		return d.disk.Status
	case d.betaDisk != nil:
		return d.betaDisk.Status
	case d.alphaDisk != nil:
		return d.alphaDisk.Status
	default:
		return "Unknown"
	}
}

// GetPDType returns the type of the PD, which is stored as a url like
// projects/project/zones/zone/diskTypes/pd-standard. The returned type is not
// validated, it is just passed verbatium from GCP.
func (d *CloudDisk) GetPDType() string {
	var pdType string
	switch {
	case d.disk != nil:
		pdType = d.disk.Type
	case d.betaDisk != nil:
		pdType = d.betaDisk.Type
	case d.alphaDisk != nil:
		pdType = d.alphaDisk.Type
	default:
		return ""
	}
	respType := strings.Split(pdType, "/")
	return strings.TrimSpace(respType[len(respType)-1])
}

func (d *CloudDisk) GetSelfLink() string {
	switch {
	case d.disk != nil:
		return d.disk.SelfLink
	case d.betaDisk != nil:
		return d.betaDisk.SelfLink
	case d.alphaDisk != nil:
		return d.alphaDisk.SelfLink
	default:
		return ""
	}
}

func (d *CloudDisk) GetSizeGb() int64 {
	switch {
	case d.disk != nil:
		return d.disk.SizeGb
	case d.betaDisk != nil:
		return d.betaDisk.SizeGb
	case d.alphaDisk != nil:
		return d.alphaDisk.SizeGb
	default:
		return -1
	}
}

// setSizeGb sets the size of the disk used ONLY
// for testing purposes.
func (d *CloudDisk) setSizeGb(size int64) {
	switch {
	case d.disk != nil:
		d.disk.SizeGb = size
	case d.betaDisk != nil:
		d.betaDisk.SizeGb = size
	case d.alphaDisk != nil:
		d.alphaDisk.SizeGb = size
	}
}

func (d *CloudDisk) GetZone() string {
	switch {
	case d.disk != nil:
		return d.disk.Zone
	case d.betaDisk != nil:
		return d.betaDisk.Zone
	case d.alphaDisk != nil:
		return d.alphaDisk.Zone
	default:
		return ""
	}
}

func (d *CloudDisk) GetSnapshotId() string {
	switch {
	case d.disk != nil:
		return d.disk.SourceSnapshotId
	case d.betaDisk != nil:
		return d.betaDisk.SourceSnapshotId
	case d.alphaDisk != nil:
		return d.alphaDisk.SourceSnapshotId
	default:
		return ""
	}
}

func (d *CloudDisk) GetSourceDiskId() string {
	switch {
	case d.disk != nil:
		return d.disk.SourceDiskId
	case d.betaDisk != nil:
		return d.betaDisk.SourceDiskId
	case d.alphaDisk != nil:
		return d.alphaDisk.SourceDiskId
	default:
		return ""
	}
}

func (d *CloudDisk) GetImageId() string {
	switch {
	case d.disk != nil:
		return d.disk.SourceImageId
	case d.betaDisk != nil:
		return d.betaDisk.SourceImageId
	case d.alphaDisk != nil:
		return d.alphaDisk.SourceImageId
	default:
		return ""
	}
}

func (d *CloudDisk) GetKMSKeyName() string {
	switch {
	case d.disk != nil:
		if dek := d.disk.DiskEncryptionKey; dek != nil {
			return dek.KmsKeyName
		}
	case d.betaDisk != nil:
		if dek := d.betaDisk.DiskEncryptionKey; dek != nil {
			return dek.KmsKeyName
		}
	case d.alphaDisk != nil:
		if dek := d.alphaDisk.DiskEncryptionKey; dek != nil {
			return dek.KmsKeyName
		}
	}
	return ""
}

func (d *CloudDisk) GetMultiWriter() bool {
	switch {
	case d.disk != nil:
		return false
	case d.betaDisk != nil:
		return d.betaDisk.MultiWriter
	case d.alphaDisk != nil:
		return d.alphaDisk.MultiWriter
	default:
		return false
	}
}

func (d *CloudDisk) GetEnableConfidentialCompute() bool {
	switch {
	case d.disk != nil:
		return false
	case d.betaDisk != nil:
		return d.betaDisk.EnableConfidentialCompute
	case d.alphaDisk != nil:
		return d.alphaDisk.EnableConfidentialCompute
	default:
		return false
	}
}

func (d *CloudDisk) GetEnableStoragePools() bool {
	switch {
	case d.disk != nil:
		return false
	case d.betaDisk != nil:
		return false
	case d.alphaDisk != nil:
		return d.alphaDisk.StoragePool != ""
	default:
		return false
	}
}
