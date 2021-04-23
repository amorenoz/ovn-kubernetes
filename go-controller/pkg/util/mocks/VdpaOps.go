// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	kvdpa "github.com/redhat-virtio-net/govdpa/pkg/kvdpa"
	mock "github.com/stretchr/testify/mock"
)

// VdpaOps is an autogenerated mock type for the VdpaOps type
type VdpaOps struct {
	mock.Mock
}

// GetVdpaDeviceByPci provides a mock function with given fields: pciAddress
func (_m *VdpaOps) GetVdpaDeviceByPci(pciAddress string) (kvdpa.VdpaDevice, error) {
	ret := _m.Called(pciAddress)

	var r0 kvdpa.VdpaDevice
	if rf, ok := ret.Get(0).(func(string) kvdpa.VdpaDevice); ok {
		r0 = rf(pciAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(kvdpa.VdpaDevice)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(pciAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
