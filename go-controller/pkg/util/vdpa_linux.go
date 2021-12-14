package util

import (
	"github.com/redhat-virtio-net/govdpa/pkg/kvdpa"
)

type VdpaDevice interface {
	kvdpa.VdpaDevice
}

type VdpaOps interface {
	GetVdpaDeviceByPci(pciAddress string) (kvdpa.VdpaDevice, error)
}

type defaultVdpaOps struct {
}

var vdpaOps VdpaOps = &defaultVdpaOps{}

// SetVdpaOpsInst method should be used by unit tests in
func SetVdpaOpsInst(mockInst VdpaOps) {
	vdpaOps = mockInst
}

// GetVdpaOps will be invoked by functions in other packages that would need access to the govdpa library methods.
func GetVdpaOps() VdpaOps {
	return vdpaOps
}

func (v *defaultVdpaOps) GetVdpaDeviceByPci(pciAddress string) (kvdpa.VdpaDevice, error) {
	return kvdpa.GetVdpaDeviceByPci(pciAddress)
}
