// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
)

var (
	ParseSettingsCompatible = parseSettingsCompatible
	GetStorageState         = getStorageState
)

func GetState(st *state.State) Backend {
	return stateShim{st}
}

func GetModel(m *state.Model) Model {
	return modelShim{m}
}

func DeducePlatformForTest(api *APIBase, arg params.DeployFromRepositoryArg) (corecharm.Platform, bool, error) {
	return api.deducePlatform(arg)
}
