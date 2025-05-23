// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package internal

import (
	corenetwork "github.com/juju/juju/core/network"
)

// ImportLinkLayerDevice represents a physical or virtual
// network interface and its IP addresses.
type ImportLinkLayerDevice struct {
	UUID             string
	IsAutoStart      bool
	IsEnabled        bool
	MTU              *int64
	MachineID        string
	MACAddress       *string
	NetNodeUUID      string
	Name             string
	ParentDeviceName string
	ProviderID       *string
	Type             corenetwork.LinkLayerDeviceType
	VirtualPortType  corenetwork.VirtualPortType
}

// ImportLinkLayerDevice represents a physical or virtual
// network interface and its IP addresses.
type ImportApplicationCloudService struct {
	Name       string
	ProviderID *string
	Addresses  []ImportCloudServiceAddress
}

// ImportLinkLayerDevice represents a physical or virtual
// network interface and its IP addresses.
type ImportCloudServiceAddress struct {
	// UUID is the ip address uuid
	UUID string
	// DeviceUUID is the link layer device uuid.
	DeviceUUID string
	Value      string
	Type       corenetwork.AddressType
	Scope      corenetwork.Scope
	Origin     corenetwork.Origin
	SpaceID    string
}
