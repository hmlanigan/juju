// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmigration

import (
	"context"

	"github.com/juju/clock"
	"github.com/juju/description/v9"

	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/modelmigration"
	"github.com/juju/juju/domain/relation/service"
	"github.com/juju/juju/domain/relation/state"
	"github.com/juju/juju/internal/errors"
)

// RegisterExport registers the export operations with the given coordinator.
func RegisterExport(
	coordinator Coordinator,
	clock clock.Clock,
	logger logger.Logger,
) {
	coordinator.Add(&exportOperation{
		clock:  clock,
		logger: logger,
	})
}

// ExportService provides a subset of the resource domain service methods
// needed for resource export.
type ExportService interface {
	// ExportResources returns the list of application and unit resources to
	// export for the given application.
	//
	// If the application exists but doesn't have any resources, no error are
	// returned, the result just contains an empty list.
	ExportRelations(ctx context.Context, name string) error //resource.ExportedResources,

}

// exportOperation describes a way to execute a migration for
// exporting applications.
type exportOperation struct {
	modelmigration.BaseOperation

	service ExportService

	clock  clock.Clock
	logger logger.Logger
}

// Name returns the name of this operation.
func (e *exportOperation) Name() string {
	return "export resources"
}

// Setup the export operation.
// This will create a new service instance.
func (e *exportOperation) Setup(scope modelmigration.Scope) error {
	e.service = service.NewService(
		state.NewState(scope.ModelDB(), e.clock, e.logger),
		e.logger,
	)
	return nil
}

// Execute the export. Go through each of the applications on the model and add
// their resources, and unit resources.
func (e *exportOperation) Execute(ctx context.Context, model description.Model) error {
	for _, app := range model.Applications() {
		exported, err := e.service.ExportRelations(ctx, app.Name())
		if err != nil {
			return errors.Errorf("getting resource of application %q: %w", app.Name(), err)
		}

	}
	return nil
}
