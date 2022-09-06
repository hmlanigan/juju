// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client

import (
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/base"
)

// SetServerAddress allows changing the URL to the internal API server
// that AddLocalCharm uses in order to test NotImplementedError.
func SetServerAddress(c *Client, scheme, addr string) {
	api.SetServerAddressForTesting(c.conn, scheme, addr)
}

func NewClientFromCaller(caller base.FacadeCaller) *Client {
	return &Client{
		facade: caller,
	}
}
