// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation_test

import (
	"github.com/juju/charm/v7"
	"github.com/juju/charm/v7/hooks"
	"github.com/juju/juju/worker/uniter/hook"
	"github.com/juju/juju/worker/uniter/runner/context"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api/uniter"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/core/leadership"
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

func (s *relationResolverSuite) TestNewRelationsWithExistingRelationsLeader(c *gc.C) {
	s.assertNewRelationsWithExistingRelations(c, true)
}

func (s *relationResolverSuite) TestNewRelationsWithExistingRelationsNotLeader(c *gc.C) {
	s.assertNewRelationsWithExistingRelations(c, false)
}

func (s *relationResolverSuite) assertNewRelationsWithExistingRelations(c *gc.C, isLeader bool) {
	unitTag := names.NewUnitTag("wordpress/0")
	abort := make(chan struct{})
	s.leadershipContextFunc = func(accessor context.LeadershipSettingsAccessor, tracker leadership.Tracker, unitName string) context.LeadershipContext {
		return &stubLeadershipContext{isLeader: isLeader}
	}

	var numCalls int32
	unitEntity := params.Entities{Entities: []params.Entity{{Tag: "unit-wordpress-0"}}}
	relationUnits := params.RelationUnits{RelationUnits: []params.RelationUnit{
		{Relation: "relation-wordpress.db#mysql.db", Unit: "unit-wordpress-0"},
	}}
	relationResults := params.RelationResults{
		Results: []params.RelationResult{
			{
				Id:   1,
				Key:  "wordpress:db mysql:db",
				Life: life.Alive,
				Endpoint: params.Endpoint{
					ApplicationName: "wordpress",
					Relation:        params.CharmRelation{Name: "mysql", Role: string(charm.RoleProvider), Interface: "db"},
				}},
		},
	}
	relationStatus := params.RelationStatusArgs{Args: []params.RelationStatusArg{{
		UnitTag:    "unit-wordpress-0",
		RelationId: 1,
		Status:     params.Joined,
	}}}
	unitSetStateArgs := params.SetUnitStateArgs{
		Args: []params.SetUnitStateArg{{
			Tag:           "unit-wordpress-0",
			RelationState: &map[int]string{1: "id: 1\n"},
		},
		}}
	unitStateResults := params.UnitStateResults{Results: []params.UnitStateResult{{}}}

	apiCalls := []apiCall{
		uniterAPICall("Refresh", unitEntity, params.UnitRefreshResults{Results: []params.UnitRefreshResult{{Life: life.Alive, Resolved: params.ResolvedNone}}}, nil),
		uniterAPICall("GetPrincipal", unitEntity, params.StringBoolResults{Results: []params.StringBoolResult{{Result: "", Ok: false}}}, nil),
		uniterAPICall("RelationsStatus", unitEntity, params.RelationUnitStatusResults{Results: []params.RelationUnitStatusResult{
			{RelationResults: []params.RelationUnitStatus{{RelationTag: "relation-wordpress:db mysql:db", InScope: true}}}}}, nil),
		uniterAPICall("State", unitEntity, unitStateResults, nil),
		uniterAPICall("Relation", relationUnits, relationResults, nil),
		uniterAPICall("Relation", relationUnits, relationResults, nil),
		uniterAPICall("Watch", unitEntity, params.NotifyWatchResults{Results: []params.NotifyWatchResult{{NotifyWatcherId: "1"}}}, nil),
		uniterAPICall("SetState", unitSetStateArgs, noErrorResult, nil),
		uniterAPICall("EnterScope", relationUnits, params.ErrorResults{Results: []params.ErrorResult{{}}}, nil),
	}
	if isLeader {
		apiCalls = append(apiCalls,
			uniterAPICall("SetRelationStatus", relationStatus, noErrorResult, nil),
		)
	}
	apiCaller := mockAPICaller(c, &numCalls, apiCalls...)
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
	assertNumCalls(c, &numCalls, int32(len(apiCalls)))

	info := r.GetInfo()
	c.Assert(info, gc.HasLen, 1)
	oneInfo := info[1]
	c.Assert(oneInfo.RelationUnit.Relation().Tag(), gc.Equals, names.NewRelationTag("wordpress:db mysql:db"))
	c.Assert(oneInfo.RelationUnit.Endpoint(), jc.DeepEquals, uniter.Endpoint{
		Relation: charm.Relation{Name: "mysql", Role: "provider", Interface: "db", Optional: false, Limit: 0, Scope: ""},
	})
	c.Assert(oneInfo.MemberNames, gc.HasLen, 0)
}

func (s *relationResolverSuite) TestCommitHook(c *gc.C) {
	var numCalls int32
	apiCalls := relationJoinedAPICalls2SetState()
	relationUnits := params.RelationUnits{RelationUnits: []params.RelationUnit{
		{Relation: "relation-wordpress.db#mysql.db", Unit: "unit-wordpress-0"},
	}}
	unitSetStateArgs := params.SetUnitStateArgs{
		Args: []params.SetUnitStateArg{{
			Tag:           "unit-wordpress-0",
			RelationState: &map[int]string{1: "id: 1\nmembers:\n  wordpress/0: 2\n"},
		}}}
	unitSetStateArgs2 := params.SetUnitStateArgs{
		Args: []params.SetUnitStateArg{{
			Tag:           "unit-wordpress-0",
			RelationState: &map[int]string{1: "id: 1\n"},
		}}}
	// ops.Remove() via die()
	unitSetStateArgs3 := params.SetUnitStateArgs{
		Args: []params.SetUnitStateArg{{
			Tag:           "unit-wordpress-0",
			RelationState: &map[int]string{1: ""},
		}}}
	apiCalls = append(apiCalls,
		uniterAPICall("SetState", unitSetStateArgs, noErrorResult, nil),
		uniterAPICall("SetState", unitSetStateArgs2, noErrorResult, nil),
		uniterAPICall("LeaveScope", relationUnits, params.ErrorResults{Results: []params.ErrorResult{{}}}, nil),
		uniterAPICall("SetState", unitSetStateArgs3, noErrorResult, nil),
	)
	r := s.assertHookRelationJoined(c, &numCalls, apiCalls...)

	err := r.CommitHook(hook.Info{
		Kind:              hooks.RelationChanged,
		RemoteUnit:        "wordpress/0",
		RemoteApplication: "wordpress",
		RelationId:        1,
		ChangeVersion:     2,
	})
	c.Assert(err, jc.ErrorIsNil)

	err = r.CommitHook(hook.Info{
		Kind:              hooks.RelationDeparted,
		RemoteUnit:        "wordpress/0",
		RemoteApplication: "wordpress",
		RelationId:        1,
	})
	c.Assert(err, jc.ErrorIsNil)
}
