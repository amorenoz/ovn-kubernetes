package util

import (
	"github.com/k8snetworkplumbingwg/govdpa/pkg/kvdpa"
)

type VdpaDevice interface {
	kvdpa.VdpaDevice
}

type VdpaOps interface {
	GetVdpaDeviceByPci(pciAddress string) (kvdpa.VdpaDevice, error)
	GetVduseVdpaDevice(device string) (kvdpa.VdpaDevice, error)
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
	// the PCI prefix is required by the govdpa library
	vdpaDevices, err := kvdpa.GetVdpaDevicesByPciAddress("pci/" + pciAddress)
	if len(vdpaDevices) > 0 {
		return vdpaDevices[0], nil
	}
	return nil, err
}

func (v *defaultVdpaOps) GetVduseVdpaDevice(name string) (kvdpa.VdpaDevice, error) {
	vduse, err := kvdpa.GetVduseDevice(name)
	if err != nil {
		return nil, err
	}
	return vduse.VdpaDevice()
}
