// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package firewaller

import (
	"context"
	"reflect"

	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/common/cloudspec"
	"github.com/juju/juju/apiserver/common/firewall"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/domain/application/service"
	"github.com/juju/juju/internal/storage"
)

// Register is called to expose a package of facades onto a given registry.
func Register(registry facade.FacadeRegistry) {
	registry.MustRegister("Firewaller", 7, func(stdCtx context.Context, ctx facade.ModelContext) (facade.Facade, error) {
		return newFirewallerAPIV7(ctx)
	}, reflect.TypeOf((*FirewallerAPI)(nil)))
}

// newFirewallerAPIV7 creates a new server-side FirewallerAPIv7 facade.
func newFirewallerAPIV7(ctx facade.ModelContext) (*FirewallerAPI, error) {
	st := ctx.State()
	m, err := st.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	serviceFactory := ctx.ServiceFactory()
	cloudSpecAPI := cloudspec.NewCloudSpecV2(
		ctx.Resources(),
		cloudspec.MakeCloudSpecGetterForModel(st, serviceFactory.Cloud(), serviceFactory.Credential(), serviceFactory.Config()),
		cloudspec.MakeCloudSpecWatcherForModel(st, serviceFactory.Cloud()),
		cloudspec.MakeCloudSpecCredentialWatcherForModel(st),
		cloudspec.MakeCloudSpecCredentialContentWatcherForModel(st, serviceFactory.Credential()),
		common.AuthFuncForTag(m.ModelTag()),
	)
	controllerConfigAPI := common.NewControllerConfigAPI(
		st,
		serviceFactory.ControllerConfig(),
		serviceFactory.ExternalController(),
	)

	stShim := stateShim{st: st, State: firewall.StateShim(st, m)}
	return NewStateFirewallerAPI(
		stShim,
		serviceFactory.Network(),
		ctx.Resources(),
		ctx.WatcherRegistry(),
		ctx.Auth(),
		cloudSpecAPI,
		controllerConfigAPI,
		serviceFactory.ControllerConfig(),
		serviceFactory.Config(),
		serviceFactory.Application(service.ApplicationServiceParams{
			StorageRegistry: storage.NotImplementedProviderRegistry{},
			Secrets:         service.NotImplementedSecretService{},
		}),
		ctx.Logger().Child("firewaller"),
	)
}
