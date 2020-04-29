// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation

import (
	"github.com/juju/juju/worker/uniter/runner/context"
)

type StateTrackerForTestConfig struct {
	St                StateTrackerState
	Unit              Unit
	LeadershipContext context.LeadershipContext
	Subordinate       bool
	PrincipalName     string
	CharmDir          string
}

func NewStateTrackerForTest(cfg StateTrackerForTestConfig) (RelationStateTracker, error) {
	rst := &relationStateTracker{
		st:              cfg.St,
		unit:            cfg.Unit,
		leaderCtx:       cfg.LeadershipContext,
		abort:           make(chan struct{}),
		subordinate:     cfg.Subordinate,
		principalName:   cfg.PrincipalName,
		charmDir:        cfg.CharmDir,
		relationers:     make(map[int]*Relationer),
		remoteAppName:   make(map[int]string),
		relationCreated: make(map[int]bool),
		isPeerRelation:  make(map[int]bool),
	}
	stateMgr, err := NewStateManager(cfg.Unit)
	if err != nil {
		return nil, err
	}
	rst.stateMgr = stateMgr

	return rst, nil
}
