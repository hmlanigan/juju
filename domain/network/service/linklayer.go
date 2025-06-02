// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"

	"github.com/juju/juju/core/trace"
	"github.com/juju/juju/domain/network/internal"
	"github.com/juju/juju/internal/errors"
)

// DeleteImportedLinkLayerDevices is part of the [modelmigration.MigrationService]
// interface.
func (s *MigrationService) DeleteImportedLinkLayerDevices(ctx context.Context) error {
	return s.st.DeleteImportedLinkLayerDevices(ctx)
}

// ImportCloudServicesForApplications is part of the [modelmigration.MigrationService]
// interface.
func (s *MigrationService) ImportCloudServicesForApplications(
	ctx context.Context,
	arg internal.ImportApplicationCloudService,
) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	netNodeUUID, err := s.st.GetNetNodeUUIDByApplicationName(ctx, arg.Name)
	if err != nil {
		return errors.Capture(err)
	}

	// Get application uuid
	// Get net node uuid for application.

	// ip_address
	// link_layer_device - UUID at import.
	// provider_ip_address???? in theory, the provider id can have multiple addresses.
	// k8s service row - contains the provider_id (part of application)?

	// Q: why does ip_address have a device_uuid AND a net_node_uuid? - each net node can have multipl
	// lld, so which lld/device is this address linked to also. but if the device_uuid why net_node_uuid?
	// Q: what to do with the space id?
	// Q: where does the config_type come from? - set to unknown during migration
	return s.st.ImportCloudServices(ctx, arg, netNodeUUID)
}

// ImportLinkLayerDevices is part of the [modelmigration.MigrationService]
// interface.
func (s *MigrationService) ImportLinkLayerDevices(ctx context.Context, data []internal.ImportLinkLayerDevice) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if len(data) == 0 {
		return nil
	}

	namesToUUIDs, err := s.st.AllMachinesAndNetNodes(ctx)
	if err != nil {
		return errors.Capture(err)
	}
	useData := data

	// Net node uuids were created when machines were imported.
	for i, device := range data {
		netNodeUUID, ok := namesToUUIDs[device.MachineID]
		if !ok {
			return errors.Errorf("no net node found for machineID %q", device.MachineID)
		}
		useData[i].NetNodeUUID = netNodeUUID
	}
	return s.st.ImportLinkLayerDevices(ctx, useData)
}
