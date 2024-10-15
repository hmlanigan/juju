// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasfirewaller

import (
	"testing"

	gc "gopkg.in/check.v1"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package caasfirewaller_test -destination common_service_mocks_test.go github.com/juju/juju/apiserver/common/charms CharmService
//go:generate go run go.uber.org/mock/mockgen -typed -package caasfirewaller_test -destination service_mocks_test.go github.com/juju/juju/apiserver/facades/controller/caasfirewaller ApplicationService,PortService

func TestAll(t *testing.T) {
	gc.TestingT(t)
}
