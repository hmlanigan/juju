// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"net/url"
	"time"

	"github.com/juju/loggo"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/client/client"
	"github.com/juju/juju/api/common"
)

type streamDebugLogSuite struct {
}

var _ = gc.Suite(&streamDebugLogSuite{})

func (s *streamDebugLogSuite) TestWatchDebugLogConnected(c *gc.C) {
	// Use the no tail option so we don't try to start a tailing cursor
	// on the oplog when there is no oplog configured in mongo as the tests
	// don't set up mongo in replicaset mode.
	messages, err := client.WatchDebugLog(common.DebugLogParams{NoTail: true})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(messages, gc.NotNil)
}

func (s *streamDebugLogSuite) TestWatchDebugLogParamsEncoded(c *gc.C) {
	catcher := api.UrlCatcher{}
	s.PatchValue(&api.WebsocketDial, catcher.RecordLocation)

	params := common.DebugLogParams{
		IncludeEntity: []string{"a", "b"},
		IncludeModule: []string{"c", "d"},
		IncludeLabel:  []string{"e", "f"},
		ExcludeEntity: []string{"g", "h"},
		ExcludeModule: []string{"i", "j"},
		ExcludeLabel:  []string{"k", "l"},
		Limit:         100,
		Backlog:       200,
		Level:         loggo.ERROR,
		Replay:        true,
		NoTail:        true,
		StartTime:     time.Date(2016, 11, 30, 11, 48, 0, 100, time.UTC),
	}

	client := client.NewClient(s.APIState)
	_, err := client.WatchDebugLog(params)
	c.Assert(err, jc.ErrorIsNil)

	connectURL, err := url.Parse(catcher.Location())
	c.Assert(err, jc.ErrorIsNil)

	values := connectURL.Query()
	c.Assert(values, jc.DeepEquals, url.Values{
		"includeEntity": params.IncludeEntity,
		"includeModule": params.IncludeModule,
		"includeLabel":  params.IncludeLabel,
		"excludeEntity": params.ExcludeEntity,
		"excludeModule": params.ExcludeModule,
		"excludeLabel":  params.ExcludeLabel,
		"maxLines":      {"100"},
		"backlog":       {"200"},
		"level":         {"ERROR"},
		"replay":        {"true"},
		"noTail":        {"true"},
		"startTime":     {"2016-11-30T11:48:00.0000001Z"},
	})
}
