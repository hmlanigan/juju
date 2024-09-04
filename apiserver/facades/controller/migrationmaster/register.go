// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migrationmaster

import (
	"context"
	"reflect"

	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/common/cloudspec"
	"github.com/juju/juju/apiserver/facade"
	coremodel "github.com/juju/juju/core/model"
	"github.com/juju/juju/domain/application/service"
	"github.com/juju/juju/internal/migration"
	"github.com/juju/juju/internal/storage"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("MigrationMaster", 3, func(stdCtx context.Context, ctx facade.ModelContext) (facade.Facade, error) {
		return newMigrationMasterFacade(ctx) // Adds MinionReportTimeout.
	}, reflect.TypeOf((*API)(nil)))
}

// newMigrationMasterFacade exists to provide the required signature for API
// registration, converting st to backend.
func newMigrationMasterFacade(ctx facade.ModelContext) (*API, error) {
	pool := ctx.StatePool()
	modelState := ctx.State()

	controllerState, err := pool.SystemState()
	if err != nil {
		return nil, errors.Trace(err)
	}

	preCheckBackend, err := migration.PrecheckShim(modelState, controllerState)
	if err != nil {
		return nil, errors.Annotate(err, "creating precheck backend")
	}

	leadership, err := ctx.LeadershipReader()
	if err != nil {
		return nil, errors.Trace(err)
	}

	backend := newBacked(modelState)

	serviceFactory := ctx.ServiceFactory()
	credentialService := serviceFactory.Credential()

	modelConfigServiceGetter := func(modelID coremodel.UUID) common.ModelConfigService {
		return serviceFactory.Config()
	}

	return NewAPI(
		controllerState,
		backend,
		ctx.ModelExporter(ctx.ModelUUID(), backend),
		ctx.ObjectStore(),
		preCheckBackend,
		migration.PoolShim(pool),
		ctx.Resources(),
		ctx.Auth(),
		ctx.Presence(),
		cloudspec.MakeCloudSpecGetter(pool, serviceFactory.Cloud(), credentialService, modelConfigServiceGetter),
		leadership,
		credentialService,
		serviceFactory.ControllerConfig(),
		serviceFactory.Config(),
		serviceFactory.ModelInfo(),
		serviceFactory.Model(),
		serviceFactory.Application(service.ApplicationServiceParams{
			StorageRegistry: storage.NotImplementedProviderRegistry{},
			Secrets:         service.NotImplementedSecretService{},
		}),
		serviceFactory.Upgrade(),
	)
}
