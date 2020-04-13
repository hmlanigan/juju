// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/juju/core/life"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/leadership"
	"github.com/juju/juju/worker/uniter/operation"
	"github.com/juju/juju/worker/uniter/relation"
	"github.com/juju/juju/worker/uniter/relation/mocks"
	"github.com/juju/juju/worker/uniter/remotestate"
	"github.com/juju/juju/worker/uniter/resolver"
	"github.com/juju/juju/worker/uniter/runner/context"
)

type newRelationResolverSuite struct {
	charmDir              string
	leadershipContextFunc relation.LeadershipContextFunc

	mockRelStTracker *mocks.MockRelationStateTracker
	mockSupDestroyer *mocks.MockSubordinateDestroyer
}

var _ = gc.Suite(&newRelationResolverSuite{})

func (s *newRelationResolverSuite) SetUpTest(c *gc.C) {
	s.leadershipContextFunc = func(accessor context.LeadershipSettingsAccessor, tracker leadership.Tracker, unitName string) context.LeadershipContext {
		return &stubLeadershipContext{isLeader: true}
	}
}

func (s *newRelationResolverSuite) TestNextOpNothing(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSyncScopesEmpty()

	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{}

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	_, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(errors.Cause(err), gc.Equals, resolver.ErrNoOperation)
}

func (s *newRelationResolverSuite) TestHookRelationJoined(c *gc.C) {
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life:      life.Alive,
				Suspended: false,
				Members: map[string]int64{
					"wordpress/0": 1,
				},
				ApplicationMembers: map[string]int64{
					"wordpress": 0,
				},
			},
		},
	}

	defer s.setupMocks(c).Finish()
	s.expectSyncScopes(remoteState)
	s.expectIsKnownFalse(1)
	s.expectIsImplicitFalse(1)
	s.expectStateUnknown(1)
	s.expectIsPeerRelationFalse(1)

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-joined on unit wordpress/0 with relation 1")
}

func (s *newRelationResolverSuite) TestHookRelationChangedApplication(c *gc.C) {
	//var numCalls int32
	//apiCalls := relationJoinedAPICalls()
	//r := s.assertHookRelationJoined(c, &numCalls, apiCalls...)
	//
	//// There will be an initial relation-changed regardless of
	//// members, due to the "changed pending" local persistent
	//// state.
	//s.assertHookRelationChanged(c, r, remotestate.RelationSnapshot{
	//	Life:      life.Alive,
	//	Suspended: false,
	//}, &numCalls)

	// wordpress app starts at 0, changing to 1 should trigger a
	// relation-changed hook for the app. We also leave wordpress/0 at 1 so that
	// it doesn't trigger relation-departed or relation-changed.
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life:      life.Alive,
				Suspended: false,
				Members: map[string]int64{
					"wordpress/0": 1,
				},
				ApplicationMembers: map[string]int64{
					"wordpress": 1,
				},
			},
		},
	}
	relationState := relation.State{
		RelationId: 1,
		Members: map[string]int64{
			"wordpress/0": 0,
		},
		ApplicationMembers: map[string]int64{
			"wordpress": 0,
		},
		ChangedPending: "",
	}
	defer s.setupMocks(c).Finish()
	s.expectSyncScopes(remoteState)
	s.expectIsKnownTrue(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectIsPeerRelationFalse(1)

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-changed on app wordpress with relation 1")
	c.Logf("%s", spew.Sdump(op.(*mockOperation).hookInfo))
}

func (s *newRelationResolverSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.mockRelStTracker = mocks.NewMockRelationStateTracker(ctrl)
	s.mockSupDestroyer = mocks.NewMockSubordinateDestroyer(ctrl)
	return ctrl
}

func (s *newRelationResolverSuite) expectSyncScopesEmpty() {
	exp := s.mockRelStTracker.EXPECT()
	exp.SynchronizeScopes(remotestate.Snapshot{}).Return(nil)
}

func (s *newRelationResolverSuite) expectSyncScopes(snapshot remotestate.Snapshot) {
	exp := s.mockRelStTracker.EXPECT()
	exp.SynchronizeScopes(snapshot).Return(nil)
}

func (s *newRelationResolverSuite) expectIsKnownTrue(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsKnown(id).Return(true)
}

func (s *newRelationResolverSuite) expectIsKnownFalse(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsKnown(id).Return(true)
}

func (s *newRelationResolverSuite) expectIsImplicitFalse(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsImplicit(id).Return(false, nil).AnyTimes()
}

func (s newRelationResolverSuite) expectStateUnknown(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.State(id).Return(nil, errors.Errorf("unknown relation: %d", id))
}

func (s newRelationResolverSuite) expectState(st relation.State) {
	exp := s.mockRelStTracker.EXPECT()
	exp.State(st.RelationId).Return(&st, nil)
}

func (s newRelationResolverSuite) expectIsPeerRelationFalse(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsPeerRelation(id).Return(false, nil)
}
