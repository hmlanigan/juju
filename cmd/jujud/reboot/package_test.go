// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package reboot

import (
	"testing"

	"github.com/juju/clock"
	gc "gopkg.in/check.v1"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package mocks -destination mocks/service_mock.go github.com/juju/juju/cmd/jujud/reboot AgentConfig,Manager,Model,RebootWaiter,Service
//go:generate go run go.uber.org/mock/mockgen -typed -package mocks -destination mocks/instance_mock.go github.com/juju/juju/environs/instances Instance
//go:generate go run go.uber.org/mock/mockgen -typed -package mocks -destination mocks/clock_mock.go github.com/juju/clock Clock

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

// NewRebootForTest returns a Reboot object to be used for testing.
func NewRebootForTest(acfg AgentConfig, reboot RebootWaiter, clock clock.Clock) *Reboot {
	return &Reboot{acfg: acfg, reboot: reboot, clock: clock}
}
