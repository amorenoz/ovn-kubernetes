// +build linux

package util

import (
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"

	"github.com/Mellanox/sriovnet"
	"github.com/redhat-virtio-net/govdpa/pkg/kvdpa"
)

type SriovnetOps interface {
	GetNetDevicesFromPci(pciAddress string) ([]string, error)
	GetUplinkRepresentor(vfPciAddress string) (string, error)
	GetVfIndexByPciAddress(vfPciAddress string) (int, error)
	GetVfRepresentor(uplink string, vfIndex int) (string, error)
	GetPfPciFromVfPci(vfPciAddress string) (string, error)
	GetVfRepresentorSmartNIC(pfID, vfIndex string) (string, error)
	GetRepresentorPeerMacAddress(netdev string) (net.HardwareAddr, error)
	GetRepresentorPortFlavour(netdev string) (sriovnet.PortFlavour, error)
}

type defaultSriovnetOps struct {
}

var sriovnetOps SriovnetOps = &defaultSriovnetOps{}

// SetSriovnetOpsInst method would be used by unit tests in other packages
func SetSriovnetOpsInst(mockInst SriovnetOps) {
	sriovnetOps = mockInst
}

// GetSriovnetOps will be invoked by functions in other packages that would need access to the sriovnet library methods.
func GetSriovnetOps() SriovnetOps {
	return sriovnetOps
}

func (defaultSriovnetOps) GetNetDevicesFromPci(pciAddress string) ([]string, error) {
	vdpaDev, err := kvdpa.GetVdpaDeviceByPci(pciAddress)
	if err == nil && vdpaDev.GetDriver() == kvdpa.VirtioVdpaDriver {
		virtioNetDir := filepath.Join(vdpaDev.GetPath(), "net")
		netDeviceFiles, err := ioutil.ReadDir(virtioNetDir)
		if err != nil || len(netDeviceFiles) != 1 {
			return nil, fmt.Errorf("failed to get network device name from vdpa device in %v %v", vdpaDev.GetParent(), err)
		}

		return []string{strings.TrimSpace(netDeviceFiles[0].Name())}, nil
	}
	return sriovnet.GetNetDevicesFromPci(pciAddress)
}

func (defaultSriovnetOps) GetUplinkRepresentor(vfPciAddress string) (string, error) {
	return sriovnet.GetUplinkRepresentor(vfPciAddress)
}

func (defaultSriovnetOps) GetVfIndexByPciAddress(vfPciAddress string) (int, error) {
	return sriovnet.GetVfIndexByPciAddress(vfPciAddress)
}

func (defaultSriovnetOps) GetVfRepresentor(uplink string, vfIndex int) (string, error) {
	return sriovnet.GetVfRepresentor(uplink, vfIndex)
}

func (defaultSriovnetOps) GetPfPciFromVfPci(vfPciAddress string) (string, error) {
	return sriovnet.GetPfPciFromVfPci(vfPciAddress)
}

func (defaultSriovnetOps) GetVfRepresentorSmartNIC(pfID, vfIndex string) (string, error) {
	return sriovnet.GetVfRepresentorSmartNIC(pfID, vfIndex)
}

func (defaultSriovnetOps) GetRepresentorPeerMacAddress(netdev string) (net.HardwareAddr, error) {
	return sriovnet.GetRepresentorPeerMacAddress(netdev)
}

func (defaultSriovnetOps) GetRepresentorPortFlavour(netdev string) (sriovnet.PortFlavour, error) {
	return sriovnet.GetRepresentorPortFlavour(netdev)
}
