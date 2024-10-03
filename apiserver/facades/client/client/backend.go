// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/mgo/v3"
	"github.com/juju/names/v5"
	"github.com/juju/replicaset/v3"
	"github.com/juju/version/v2"

	"github.com/juju/juju/apiserver/common/storagecommon"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/status"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/internal/charm"
	"github.com/juju/juju/state"
)

// Backend contains the state.State methods used in this package,
// allowing stubs to be created for testing.
type Backend interface {
	AddRelation(...state.Endpoint) (*state.Relation, error)
	AllApplications() ([]*state.Application, error)
	AllApplicationOffers() ([]*crossmodel.ApplicationOffer, error)
	AllRemoteApplications() ([]*state.RemoteApplication, error)
	AllMachines() ([]*state.Machine, error)
	AllIPAddresses() ([]*state.Address, error)
	AllLinkLayerDevices() ([]*state.LinkLayerDevice, error)
	AllRelations() ([]*state.Relation, error)
	Application(string) (Application, error)
	ControllerNodes() ([]state.ControllerNode, error)
	ControllerTag() names.ControllerTag
	ControllerTimestamp() (*time.Time, error)
	HAPrimaryMachine() (names.MachineTag, error)
	LatestPlaceholderCharm(*charm.URL) (*state.Charm, error)
	Machine(string) (*state.Machine, error)
	Model() (Model, error)
	ModelTag() names.ModelTag
	ModelUUID() string
	RemoteApplication(string) (*state.RemoteApplication, error)
	RemoteConnectionStatus(string) (*state.RemoteConnectionStatus, error)
	Unit(string) (Unit, error)
}

// MongoSession provides a way to get the status for the mongo replicaset.
type MongoSession interface {
	CurrentStatus() (*replicaset.Status, error)
}

// Model contains the state.Model methods used in this package.
type Model interface {
	Name() string
	Type() state.ModelType
	UUID() string
	Life() state.Life
	CloudName() string
	CloudRegion() string
	CloudCredentialTag() (names.CloudCredentialTag, bool)
	Config() (*config.Config, error)
	Owner() names.UserTag
	StatusHistory(status.StatusHistoryFilter) ([]status.StatusInfo, error)
	LatestToolsVersion() version.Number
	Status() (status.StatusInfo, error)
}

// Application represents a state.Application.
type Application interface {
	StatusHistory(status.StatusHistoryFilter) ([]status.StatusInfo, error)
}

// Unit represents a state.Unit.
type Unit interface {
	status.StatusHistoryGetter
	Life() state.Life
	IsPrincipal() bool
	PublicAddress() (network.SpaceAddress, error)
	AgentHistory() status.StatusHistoryGetter
}

// TODO - CAAS(ericclaudejones): This should contain state alone, model will be
// removed once all relevant methods are moved from state to model.
type stateShim struct {
	*state.State
	model   *state.Model
	session MongoSession
}

func (s *stateShim) Application(name string) (Application, error) {
	a, err := s.State.Application(name)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *stateShim) Unit(name string) (Unit, error) {
	u, err := s.State.Unit(name)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *stateShim) AllApplicationOffers() ([]*crossmodel.ApplicationOffer, error) {
	offers := state.NewApplicationOffers(s.State)
	return offers.AllApplicationOffers()
}

func (s stateShim) ModelTag() names.ModelTag {
	return names.NewModelTag(s.State.ModelUUID())
}

func (s stateShim) Model() (Model, error) {
	m, err := s.State.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &modelShim{Model: m}, nil
}

func (s stateShim) ControllerNodes() ([]state.ControllerNode, error) {
	nodes, err := s.State.ControllerNodes()
	if err != nil {
		return nil, errors.Trace(err)
	}
	result := make([]state.ControllerNode, len(nodes))
	for i, n := range nodes {
		result[i] = n
	}
	return result, nil
}

func (s stateShim) MongoSession() MongoSession {
	if s.session != nil {
		return s.session
	}
	return MongoSessionShim{s.State.MongoSession()}
}

type modelShim struct {
	*state.Model
}

// MongoSessionShim wraps a *mgo.Session to conform to the
// MongoSession interface.
type MongoSessionShim struct {
	*mgo.Session
}

// CurrentStatus returns the current status of the replicaset.
func (s MongoSessionShim) CurrentStatus() (*replicaset.Status, error) {
	return replicaset.CurrentStatus(s.Session)
}

type StorageInterface interface {
	storagecommon.StorageAccess
	storagecommon.VolumeAccess
	storagecommon.FilesystemAccess

	AllStorageInstances() ([]state.StorageInstance, error)
	AllFilesystems() ([]state.Filesystem, error)
	AllVolumes() ([]state.Volume, error)

	StorageAttachments(names.StorageTag) ([]state.StorageAttachment, error)
	FilesystemAttachments(names.FilesystemTag) ([]state.FilesystemAttachment, error)
	VolumeAttachments(names.VolumeTag) ([]state.VolumeAttachment, error)
}

var getStorageState = func(st *state.State) (StorageInterface, error) {
	sb, err := state.NewStorageBackend(st)
	if err != nil {
		return nil, err
	}
	return sb, nil
}
