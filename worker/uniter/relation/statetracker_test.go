// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/names.v3"

	"github.com/juju/juju/api/uniter"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/worker/uniter/relation"
)

func (s *newRelationResolverSuite) TestNewRelationsNoRelations(c *gc.C) {
	r := s.setupRelations(c)
	//No relations created.
	c.Assert(r.GetInfo(), gc.HasLen, 0)
}

func (s *newRelationResolverSuite) setupRelations(c *gc.C) relation.RelationStateTracker {
	unitTag := names.NewUnitTag("wordpress/0")
	abort := make(chan struct{})

	var numCalls int32
	unitEntity := params.Entities{Entities: []params.Entity{{Tag: "unit-wordpress-0"}}}
	unitStateResults := params.UnitStateResults{Results: []params.UnitStateResult{{}}}
	apiCaller := mockAPICaller(c, &numCalls,
		uniterAPICall("Refresh", unitEntity, params.UnitRefreshResults{Results: []params.UnitRefreshResult{{Life: life.Alive, Resolved: params.ResolvedNone}}}, nil),
		uniterAPICall("GetPrincipal", unitEntity, params.StringBoolResults{Results: []params.StringBoolResult{{Result: "", Ok: false}}}, nil),
		uniterAPICall("RelationsStatus", unitEntity, params.RelationUnitStatusResults{Results: []params.RelationUnitStatusResult{{RelationResults: []params.RelationUnitStatus{}}}}, nil),
		uniterAPICall("State", unitEntity, unitStateResults, nil),
	)
	st := uniter.NewState(apiCaller, unitTag)
	u, err := st.Unit(unitTag)
	c.Assert(err, jc.ErrorIsNil)
	r, err := relation.NewRelationStateTracker(
		relation.RelationStateTrackerConfig{
			State:                st,
			Unit:                 u,
			CharmDir:             s.charmDir,
			NewLeadershipContext: s.leadershipContextFunc,
			Abort:                abort,
		})
	c.Assert(err, jc.ErrorIsNil)
	assertNumCalls(c, &numCalls, 4)
	return r
}
