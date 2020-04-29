// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation_test

import (
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
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectStateUnknown(1)
	s.expectIsPeerRelationFalse(1)

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-joined on unit wordpress/0 with relation 1")
}

func (s *newRelationResolverSuite) TestHookRelationChangedApplication(c *gc.C) {
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
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectIsPeerRelationFalse(1)

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-changed on app wordpress with relation 1")
}

func (s *newRelationResolverSuite) TestHookRelationChangedSuspended(c *gc.C) {
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life:      life.Alive,
				Suspended: true,
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
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectLocalUnitAndApplicationLife()

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-departed on unit wordpress/0 with relation 1")
}

func (s *newRelationResolverSuite) TestHookRelationDeparted(c *gc.C) {
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life:      life.Alive,
				Suspended: true,
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
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectLocalUnitAndApplicationLife()

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-departed on unit wordpress/0 with relation 1")
}

func (s *newRelationResolverSuite) TestHookRelationBroken(c *gc.C) {
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life: life.Dying,
			},
		},
	}
	relationState := relation.State{
		RelationId:         1,
		Members:            map[string]int64{},
		ApplicationMembers: map[string]int64{},
		ChangedPending:     "",
	}
	defer s.setupMocks(c).Finish()
	s.expectSyncScopes(remoteState)
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectIsPeerRelationFalse(1)
	s.expectStateFound(1)
	s.expectRemoteApplication(1, "")

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-broken with relation 1")
}

func (s *newRelationResolverSuite) TestHookRelationBrokenWhenSuspended(c *gc.C) {
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life:      life.Alive,
				Suspended: true,
			},
		},
	}
	relationState := relation.State{
		RelationId:         1,
		Members:            map[string]int64{},
		ApplicationMembers: map[string]int64{},
		ChangedPending:     "",
	}
	defer s.setupMocks(c).Finish()
	s.expectSyncScopes(remoteState)
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectIsPeerRelationFalse(1)
	s.expectStateFound(1)
	s.expectRemoteApplication(1, "")

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-broken with relation 1")
}

func (s *newRelationResolverSuite) TestHookRelationBrokenOnlyOnce(c *gc.C) {
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life: life.Dying,
			},
		},
	}
	relationState := relation.State{
		RelationId:         1,
		Members:            map[string]int64{},
		ApplicationMembers: map[string]int64{},
		ChangedPending:     "",
	}
	defer s.setupMocks(c).Finish()
	s.expectSyncScopes(remoteState)
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectIsPeerRelationFalse(1)
	s.expectStateFoundFalse(1)

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	_, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(errors.Cause(err), gc.Equals, resolver.ErrNoOperation)
}

func (s *newRelationResolverSuite) TestImplicitRelationNoHooks(c *gc.C) {
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life: life.Alive,
				Members: map[string]int64{
					"wordpress": 1,
				},
			},
		},
	}
	defer s.setupMocks(c).Finish()
	s.expectSyncScopes(remoteState)
	s.expectIsKnown(1)
	s.expectIsImplicit(1)

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, s.mockSupDestroyer)
	_, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(errors.Cause(err), gc.Equals, resolver.ErrNoOperation)
}

func (s *newRelationResolverSuite) TestPrincipalDyingDestroysSubordinates(c *gc.C) {
	// So now we have a relation between a principal (wordpress) and a
	// subordinate (nrpe). If the wordpress unit is being destroyed,
	// the subordinate must be also queued for destruction.
	localState := resolver.LocalState{
		State: operation.State{
			Kind: operation.Continue,
		},
	}
	remoteState := remotestate.Snapshot{
		Life: life.Dying,
		Relations: map[int]remotestate.RelationSnapshot{
			1: {
				Life: life.Alive,
				Members: map[string]int64{
					"nrpe/0": 1,
				},
			},
		},
	}
	relationState := relation.State{
		RelationId:         1,
		Members:            map[string]int64{},
		ApplicationMembers: map[string]int64{},
		ChangedPending:     "",
	}
	ctrl := s.setupMocks(c)
	defer ctrl.Finish()
	s.expectSyncScopes(remoteState)
	s.expectIsKnown(1)
	s.expectIsImplicitFalse(1)
	s.expectState(relationState)
	s.expectIsPeerRelationFalse(1)
	s.expectHasContainerScope(1)
	s.expectStateFound(1)
	s.expectRemoteApplication(1, "")
	destroyer := mocks.NewMockSubordinateDestroyer(ctrl)
	destroyer.EXPECT().DestroyAllSubordinates().Return(nil)

	relationsResolver := relation.NewRelationResolver(s.mockRelStTracker, destroyer)
	op, err := relationsResolver.NextOp(localState, remoteState, &mockOperations{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(op.String(), gc.Equals, "run hook relation-broken with relation 1")
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

func (s *newRelationResolverSuite) expectIsKnown(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsKnown(id).Return(true).AnyTimes()
}

func (s *newRelationResolverSuite) expectIsImplicit(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsImplicit(id).Return(true, nil).AnyTimes()
}

func (s *newRelationResolverSuite) expectIsImplicitFalse(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsImplicit(id).Return(false, nil).AnyTimes()
}

func (s *newRelationResolverSuite) expectStateUnknown(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.State(id).Return(nil, errors.Errorf("unknown relation: %d", id))
}

func (s *newRelationResolverSuite) expectState(st relation.State) {
	exp := s.mockRelStTracker.EXPECT()
	exp.State(st.RelationId).Return(&st, nil)
}

func (s *newRelationResolverSuite) expectStateMaybe(st relation.State) {
	exp := s.mockRelStTracker.EXPECT()
	exp.State(st.RelationId).Return(&st, nil).AnyTimes()
}

func (s *newRelationResolverSuite) expectIsPeerRelationFalse(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.IsPeerRelation(id).Return(false, nil)
}

func (s *newRelationResolverSuite) expectLocalUnitAndApplicationLife() {
	exp := s.mockRelStTracker.EXPECT()
	exp.LocalUnitAndApplicationLife().Return(life.Alive, life.Alive, nil)
}

func (s *newRelationResolverSuite) expectStateFound(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.StateFound(id).Return(true)
}

func (s *newRelationResolverSuite) expectStateFoundFalse(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.StateFound(id).Return(false)
}

func (s *newRelationResolverSuite) expectRemoteApplication(id int, app string) {
	exp := s.mockRelStTracker.EXPECT()
	exp.RemoteApplication(id).Return(app)
}

func (s *newRelationResolverSuite) expectHasContainerScope(id int) {
	exp := s.mockRelStTracker.EXPECT()
	exp.HasContainerScope(id).Return(true, nil)
}
