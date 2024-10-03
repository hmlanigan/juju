// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelconfig

import (
	"github.com/juju/names/v5"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/state"
)

// Backend contains the state.State methods used in this package,
// allowing stubs to be created for testing.
type Backend interface {
	common.BlockGetter
	ControllerTag() names.ControllerTag
	ModelTag() names.ModelTag
	Sequences() (map[string]int, error)
	SetModelConstraints(value constraints.Value) error
	ModelConstraints() (constraints.Value, error)
}

type stateShim struct {
	*state.State
	model *state.Model
}

func (st stateShim) ModelTag() names.ModelTag {
	m, err := st.State.Model()
	if err != nil {
		return names.NewModelTag(st.State.ModelUUID())
	}

	return m.ModelTag()
}

// NewStateBackend creates a backend for the facade to use.
func NewStateBackend(m *state.Model) Backend {
	return stateShim{
		State: m.State(),
		model: m,
	}
}
