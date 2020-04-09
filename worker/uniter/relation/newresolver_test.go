// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation_test

import (
	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
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
	//s.charmDir = filepath.Join(c.MkDir(), "charm")
	//err := os.MkdirAll(s.charmDir, 0755)
	//c.Assert(err, jc.ErrorIsNil)
	//err = ioutil.WriteFile(filepath.Join(s.charmDir, "metadata.yaml"), []byte(minimalMetadata), 0755)
	//c.Assert(err, jc.ErrorIsNil)
	s.leadershipContextFunc = func(accessor context.LeadershipSettingsAccessor, tracker leadership.Tracker, unitName string) context.LeadershipContext {
		return &stubLeadershipContext{isLeader: true}
	}
}

func (s *newRelationResolverSuite) TestNextOpNothing(c *gc.C) {
	//unitTag := names.NewUnitTag("wordpress/0")
	//abort := make(chan struct{})
	//
	//var numCalls int32
	//unitEntity := params.Entities{Entities: []params.Entity{{Tag: "unit-wordpress-0"}}}
	//unitStateResults := params.UnitStateResults{Results: []params.UnitStateResult{{}}}
	//apiCaller := mockAPICaller(c, &numCalls,
	//	uniterAPICall("Refresh", unitEntity, params.UnitRefreshResults{Results: []params.UnitRefreshResult{{Life: life.Alive, Resolved: params.ResolvedNone}}}, nil),
	//	uniterAPICall("GetPrincipal", unitEntity, params.StringBoolResults{Results: []params.StringBoolResult{{Result: "", Ok: false}}}, nil),
	//	uniterAPICall("RelationsStatus", unitEntity, params.RelationUnitStatusResults{Results: []params.RelationUnitStatusResult{{RelationResults: []params.RelationUnitStatus{}}}}, nil),
	//	uniterAPICall("State", unitEntity, unitStateResults, nil),
	//)
	//st := uniter.NewState(apiCaller, unitTag)
	//u, err := st.Unit(unitTag)
	//c.Assert(err, jc.ErrorIsNil)
	//r, err := relation.NewRelationStateTracker(
	//	relation.RelationStateTrackerConfig{
	//		State:                st,
	//		Unit:                 u,
	//		CharmDir:             s.charmDir,
	//		NewLeadershipContext: s.leadershipContextFunc,
	//		Abort:                abort,
	//	})
	//c.Assert(err, jc.ErrorIsNil)
	//assertNumCalls(c, &numCalls, 4)
	defer s.setupMocks(c).Finish()

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
