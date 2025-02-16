package util

// Contains helper functions for OVN
// Eventually these should all be migrated to go-ovn bindings

import (
	"context"
	"fmt"
	"net"
	"time"

	ocpconfigapi "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	libovsdbclient "github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/config"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/libovsdbops"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/sbdb"
)

// CreateMACBinding Creates MAC binding in OVN SBDB
func CreateMACBinding(sbClient libovsdbclient.Client, logicalPort, datapathName string, portMAC net.HardwareAddr, nextHop net.IP) error {
	p := func(item *sbdb.DatapathBinding) bool {
		return item.ExternalIDs["name"] == datapathName
	}

	// It could take some time for the datapath to propagate to SBDB, wait for some time
	maxTimeout := 10 * time.Second
	var datapath *sbdb.DatapathBinding
	var err1 error
	err := wait.PollUntilContextTimeout(context.TODO(), 50*time.Millisecond, maxTimeout, true, func(ctx context.Context) (bool, error) {
		if datapath, err1 = libovsdbops.GetDatapathBindingWithPredicate(sbClient, p); err1 != nil {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("failed to find datpath: %s, after %s: %w, %v", datapathName, maxTimeout, err, err1)
	}

	// find Create mac_binding if needed
	mb := sbdb.MACBinding{
		LogicalPort: logicalPort,
		MAC:         portMAC.String(),
		Datapath:    datapath.UUID,
		IP:          nextHop.String(),
	}

	err = libovsdbops.CreateOrUpdateMacBinding(sbClient, &mb, &mb.Datapath, &mb.LogicalPort, &mb.IP, &mb.MAC)
	if err != nil {
		return fmt.Errorf("failed to create mac binding %+v: %v", mb, err)
	}

	return nil
}

func PlatformTypeIsEgressIPCloudProvider() bool {
	return config.Kubernetes.PlatformType == string(ocpconfigapi.AWSPlatformType) ||
		config.Kubernetes.PlatformType == string(ocpconfigapi.GCPPlatformType) ||
		config.Kubernetes.PlatformType == string(ocpconfigapi.AzurePlatformType) ||
		config.Kubernetes.PlatformType == string(ocpconfigapi.OpenStackPlatformType)
}
