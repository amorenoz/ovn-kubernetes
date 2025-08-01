package util

import (
	"fmt"
	"path/filepath"

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
	// Look up vdpa devices directly, which uses /sys, instead of looking up the VDUSE
	// device first which would require /dev to be mounted.
	vdpaDevices, err := kvdpa.GetVdpaDevicesByMgmtDev("", "vduse")
	if err != nil {
		return nil, err
	}

	for _, dev := range vdpaDevices {
		parent, err := dev.ParentDevicePath()
		if err != nil {
			return nil, err
		}
		if filepath.Base(parent) == name {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("vduse vdpa device %s not found", name)
}
