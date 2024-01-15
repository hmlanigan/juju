// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"context"
	"fmt"

	"github.com/juju/names/v5"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v4/workertest"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	statetesting "github.com/juju/juju/state/testing"
)

type agentEntityWatcherSuite struct{}

var _ = gc.Suite(&agentEntityWatcherSuite{})

type fakeAgentEntityWatcher struct {
	state.Entity
	fetchError
}

func (a *fakeAgentEntityWatcher) Watch() state.NotifyWatcher {
	return apiservertesting.NewFakeNotifyWatcher()
}

func (*agentEntityWatcherSuite) TestWatch(c *gc.C) {
	st := &fakeState{
		entities: map[names.Tag]entityWithError{
			u("x/0"): &fakeAgentEntityWatcher{fetchError: "x0 fails"},
			u("x/1"): &fakeAgentEntityWatcher{},
			u("x/2"): &fakeAgentEntityWatcher{},
		},
	}
	getCanWatch := func() (common.AuthFunc, error) {
		x0 := u("x/0")
		x1 := u("x/1")
		return func(tag names.Tag) bool {
			return tag == x0 || tag == x1
		}, nil
	}
	resources := common.NewResources()
	a := common.NewAgentEntityWatcher(st, resources, getCanWatch)
	entities := params.Entities{Entities: []params.Entity{
		{Tag: "unit-x-0"}, {Tag: "unit-x-1"}, {Tag: "unit-x-2"}, {Tag: "unit-x-3"},
	}}
	result, err := a.Watch(context.Background(), entities)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, gc.DeepEquals, params.NotifyWatchResults{
		Results: []params.NotifyWatchResult{
			{Error: &params.Error{Message: "x0 fails"}},
			{NotifyWatcherId: "1", Error: nil},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
		},
	})
}

func (*agentEntityWatcherSuite) TestWatchError(c *gc.C) {
	getCanWatch := func() (common.AuthFunc, error) {
		return nil, fmt.Errorf("pow")
	}
	resources := common.NewResources()
	a := common.NewAgentEntityWatcher(
		&fakeState{},
		resources,
		getCanWatch,
	)
	_, err := a.Watch(context.Background(), params.Entities{Entities: []params.Entity{{Tag: "x0"}}})
	c.Assert(err, gc.ErrorMatches, "pow")
}

func (*agentEntityWatcherSuite) TestWatchNoArgsNoError(c *gc.C) {
	getCanWatch := func() (common.AuthFunc, error) {
		return nil, fmt.Errorf("pow")
	}
	resources := common.NewResources()
	a := common.NewAgentEntityWatcher(
		&fakeState{},
		resources,
		getCanWatch,
	)
	result, err := a.Watch(context.Background(), params.Entities{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 0)
}

type multiNotifyWatcherSuite struct{}

var _ = gc.Suite(&multiNotifyWatcherSuite{})

func (*multiNotifyWatcherSuite) TestMultiNotifyWatcher(c *gc.C) {
	w0 := apiservertesting.NewFakeNotifyWatcher()
	w1 := apiservertesting.NewFakeNotifyWatcher()

	mw := common.NewMultiNotifyWatcher(w0, w1)
	defer workertest.CleanKill(c, mw)

	wc := statetesting.NewNotifyWatcherC(c, mw)
	wc.AssertOneChange()

	w0.C <- struct{}{}
	wc.AssertOneChange()
	w1.C <- struct{}{}
	wc.AssertOneChange()

	w0.C <- struct{}{}
	w1.C <- struct{}{}
	wc.AssertOneChange()
}

func (*multiNotifyWatcherSuite) TestMultiNotifyWatcherStop(c *gc.C) {
	w0 := apiservertesting.NewFakeNotifyWatcher()
	w1 := apiservertesting.NewFakeNotifyWatcher()

	mw := common.NewMultiNotifyWatcher(w0, w1)
	wc := statetesting.NewNotifyWatcherC(c, mw)
	wc.AssertOneChange()
	statetesting.AssertCanStopWhenSending(c, mw)
	wc.AssertClosed()
}
