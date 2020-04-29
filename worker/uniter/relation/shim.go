// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation

import (
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/names/v4"

	"github.com/juju/juju/api/uniter"
)

type stateTrackerStateShim struct {
	*uniter.State
}

func (s *stateTrackerStateShim) Relation(tag names.RelationTag) (StateTrackerRelation, error) {
	rel, err := s.State.Relation(tag)
	if err != nil {
		return nil, err
	}
	return &stateTrackerRelationShim{rel}, nil
}

// RelationById returns the existing relation with the given id.
func (s *stateTrackerStateShim) RelationById(id int) (StateTrackerRelation, error) {
	rel, err := s.State.RelationById(id)
	if err != nil {
		return nil, err
	}
	return &stateTrackerRelationShim{rel}, nil
}

type stateTrackerRelationShim struct {
	*uniter.Relation
}

func (s *stateTrackerRelationShim) Unit(stUnit StateTrackerUnit) (RelationUnit, error) {
	u, err := s.Relation.Unit(stUnit.unit())
	if err != nil {
		return nil, err
	}
	return &relationUnitShim{rel: s, ru: u}, nil
}

type stateTrackerUnitShim struct {
	u *uniter.Unit
}

func (s *stateTrackerUnitShim) Application() (StateTrackerApplication, error) {
	app, err := s.u.Application()
	if err != nil {
		return nil, err
	}
	return &stateTrackerApplicationShim{app}, nil
}

func (s *stateTrackerUnitShim) RelationsStatus() ([]uniter.RelationStatus, error) {
	return s.u.RelationsStatus()
}

func (s *stateTrackerUnitShim) Watch() (watcher.NotifyWatcher, error) {
	return s.u.Watch()
}

func (s *stateTrackerUnitShim) Destroy() error {
	return s.u.Destroy()
}

func (s *stateTrackerUnitShim) Name() string {
	return s.u.Name()
}

func (s *stateTrackerUnitShim) Refresh() error {
	return s.u.Refresh()
}

func (s *stateTrackerUnitShim) Life() life.Value {
	return s.u.Life()
}

func (s *stateTrackerUnitShim) State() (params.UnitStateResult, error) {
	return s.u.State()
}

func (s *stateTrackerUnitShim) SetState(unitState params.SetUnitStateArg) error {
	return s.u.SetState(unitState)
}

func (s *stateTrackerUnitShim) unit() *uniter.Unit {
	return s.u
}

type stateTrackerApplicationShim struct {
	*uniter.Application
}

type relationUnitShim struct {
	rel StateTrackerRelation
	ru  *uniter.RelationUnit
}

func (r *relationUnitShim) Relation() StateTrackerRelation {
	return r.rel
}

func (r *relationUnitShim) Endpoint() uniter.Endpoint {
	return r.ru.Endpoint()
}

func (r *relationUnitShim) EnterScope() error {
	return r.ru.EnterScope()
}

func (r *relationUnitShim) LeaveScope() error {
	return r.ru.LeaveScope()
}

func (r *relationUnitShim) unit() *uniter.RelationUnit {
	return r.ru
}
