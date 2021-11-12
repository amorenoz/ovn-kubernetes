package types

import (
	"github.com/containernetworking/cni/pkg/types"
	"net"
)

// NetConf is CNI NetConf with DeviceID
type NetConf struct {
	types.NetConf
	// PciAddrs in case of using sriov
	DeviceID string `json:"deviceID,omitempty"`
	// PCIAddress of the PF to create subfunctions on
	SubFunctionPF string `json:"subfunctionPF,omitempty"`
	// LogFile to log all the messages from cni shim binary to
	LogFile string `json:"logFile,omitempty"`
	// Level is the logging verbosity level
	LogLevel string `json:"logLevel,omitempty"`
	// LogFileMaxSize is the maximum size in bytes of the logfile
	// before it gets rolled.
	LogFileMaxSize int `json:"logfile-maxsize"`
	// LogFileMaxBackups represents the the maximum number of
	// old log files to retain
	LogFileMaxBackups int `json:"logfile-maxbackups"`
	// LogFileMaxAge represents the maximum number
	// of days to retain old log files
	LogFileMaxAge int `json:"logfile-maxage"`
	// RuntimeConfig
	RuntimeConfig struct {
		// HACK: Using DeviceID for runtime override of SubFunctionPF
		// TODO: Add a new well-known capability in
		// https://github.com/containernetworking/cni/blob/master/CONVENTIONS.md
		DeviceID string `json:"deviceID,omitempty"`
		// Path to the CNIDeviceInfoFile (TODO)
		CNIDeviceFile string `json:"CNIDeviceFile,omitempty"`
	} `json:"runtimeConfig,omitempty"`
}

// NetworkSelectionElement represents one element of the JSON format
// Network Attachment Selection Annotation as described in section 4.1.2
// of the CRD specification.
type NetworkSelectionElement struct {
	// Name contains the name of the Network object this element selects
	Name string `json:"name"`
	// Namespace contains the optional namespace that the network referenced
	// by Name exists in
	Namespace string `json:"namespace,omitempty"`
	// MacRequest contains an optional requested MAC address for this
	// network attachment
	MacRequest string `json:"mac,omitempty"`
	// GatewayRequest contains default route IP address for the pod
	GatewayRequest []net.IP `json:"default-route,omitempty"`
}
