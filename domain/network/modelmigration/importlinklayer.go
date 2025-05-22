// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmigration

import (
	"context"

	"github.com/juju/description/v9"

	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/modelmigration"
	corenetwork "github.com/juju/juju/core/network"
	"github.com/juju/juju/domain/network/internal"
	"github.com/juju/juju/domain/network/service"
	"github.com/juju/juju/domain/network/state"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/uuid"
)

// RegisterLinkLayerImport registers the import operations with the
// given coordinator.
func RegisterLinkLayerImport(coordinator Coordinator, logger logger.Logger) {
	coordinator.Add(&importOperation{
		logger: logger,
	})
}

// MigrationService defines methods needed to import and export
// link layer devices as part of model migration.
type MigrationService interface {
	// ImportLinkLayerDevices imports the given link layer device data into
	// the model.
	ImportLinkLayerDevices(ctx context.Context, data []internal.ImportLinkLayerDevice) error

	// DeleteImportedLinkLayerDevices removes all link layer device data
	// imported via the ImportLinkLayerDevices method.
	DeleteImportedLinkLayerDevices(ctx context.Context) error
}

type importLinkLayerOperation struct {
	modelmigration.BaseOperation

	migrationService MigrationService
	logger           logger.Logger
}

// Name returns the name of this operation.
func (i *importLinkLayerOperation) Name() string {
	return "import link layer devices"
}

// Setup implements Operation.
func (i *importLinkLayerOperation) Setup(scope modelmigration.Scope) error {
	i.migrationService = service.NewMigrationService(
		state.NewState(scope.ModelDB(), i.logger),
		i.logger,
	)
	return nil
}

// Execute the import of the link layer devices, an application's cloud service
// addresses and ip addresses contained in the model.
func (i *importLinkLayerOperation) Execute(ctx context.Context, model description.Model) error {
	if err := i.importLinkLayerDevices(ctx, model.LinkLayerDevices()); err != nil {
		return errors.Capture(err)
	}
	return nil
}

// Rollback the resource import operation by deleting all data imported
// within the network domain.
func (i *importLinkLayerOperation) Rollback(ctx context.Context, model description.Model) error {
	if len(model.LinkLayerDevices()) == 0 {
		return nil
	}
	err := i.migrationService.DeleteImportedLinkLayerDevices(ctx)
	if err != nil {
		return errors.Errorf("link layer device import rollback failed: %w", err)
	}
	return nil
}

func (i *importLinkLayerOperation) importLinkLayerDevices(ctx context.Context, modelLLD []description.LinkLayerDevice) error {
	if len(modelLLD) == 0 {
		return nil
	}
	data, err := i.transformLinkLayerDevices(modelLLD)
	if err != nil {
		return err
	}
	if err := i.migrationService.ImportLinkLayerDevices(ctx, data); err != nil {
		return errors.Errorf("importing link layer devices: %w", err)
	}
	return nil
}

func (i *importLinkLayerOperation) transformLinkLayerDevices(modelLLD []description.LinkLayerDevice) ([]internal.ImportLinkLayerDevice, error) {
	data := make([]internal.ImportLinkLayerDevice, len(modelLLD))

	for i, lld := range modelLLD {
		var (
			mac        *string
			mtu        *int64
			providerID *string
		)
		if lld.ProviderID() != "" {
			providerID = ptr(lld.ProviderID())
		}
		if lld.MTU() > 0 {
			mtu = ptr(int64(lld.MTU()))
		}
		if lld.MACAddress() != "" {
			mac = ptr(lld.MACAddress())
		}
		lldUUID, err := uuid.NewUUID()
		if err != nil {
			return nil, errors.Errorf("creating UUID for link layer device %q", lld.Name())
		}
		data[i] = internal.ImportLinkLayerDevice{
			UUID:             lldUUID.String(),
			Name:             lld.Name(),
			MachineID:        lld.MachineID(),
			MTU:              mtu,
			MACAddress:       mac,
			ProviderID:       providerID,
			Type:             corenetwork.LinkLayerDeviceType(lld.Type()),
			VirtualPortType:  corenetwork.VirtualPortType(lld.VirtualPortType()),
			IsAutoStart:      lld.IsAutoStart(),
			IsEnabled:        lld.IsUp(),
			ParentDeviceName: lld.ParentName(),
		}
	}

	return data, nil
}

// ptr returns a reference to a copied value of type T.
func ptr[T any](i T) *T {
	return &i
}
