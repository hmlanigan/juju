// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client

import (
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/base"
)

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/connection_mock.go github.com/juju/juju/api Connection
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/httprequest_mock.go gopkg.in/httprequest.v1 Doer

func TestAll(t *testing.T) {
	gc.TestingT(t)
}

func NewClientFromCaller(caller base.FacadeCaller) *Client {
	return &Client{
		facade: caller,
	}
}

func NewClientFromCallerAndConnection(caller base.FacadeCaller, conn api.Connection) *Client {
	return &Client{
		facade: caller,
		conn:   conn,
	}
}
