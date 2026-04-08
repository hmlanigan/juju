// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"testing"

	"github.com/juju/tc"
	"go.uber.org/mock/gomock"

	coreerrors "github.com/juju/juju/core/errors"
	corerelation "github.com/juju/juju/core/relation"
	corerelationtesting "github.com/juju/juju/core/relation/testing"
	coreunit "github.com/juju/juju/core/unit"
	unittesting "github.com/juju/juju/core/unit/testing"
	relationerrors "github.com/juju/juju/domain/relation/errors"
	"github.com/juju/juju/domain/unitstate"
	"github.com/juju/juju/domain/unitstate/internal"
	"github.com/juju/juju/internal/errors"
	loggertesting "github.com/juju/juju/internal/logger/testing"
)

type commitHookSuite struct {
	svc *LeadershipService

	st                     *MockState
	leadershipEnsurer      *MockEnsurer
	networkProviderGetter  func(context.Context) (ProviderWithNetworking, error)
	providerWithNetworking *MockProviderWithNetworking
}

func TestCommitHookSuite(t *testing.T) {
	tc.Run(t, &commitHookSuite{})
}

func (s *commitHookSuite) TestCommitHookChangesNoChanges(c *tc.C) {
	defer s.setupMocks(c).Finish()

	// Arrange: args which no changes are needed
	arg := unitstate.CommitHookChangesArg{
		UnitName: unittesting.GenNewName(c, "test/0"),
	}

	// Act
	err := s.svc.CommitHookChanges(c.Context(), arg)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
}

func (s *commitHookSuite) TestCommitHookChangesNoLeadership(c *tc.C) {
	defer s.setupMocks(c).Finish()

	// Arrange: args which indicate leadership is not required
	key := corerelationtesting.GenNewKey(c, "app-1:fake-endpoint-name-1 app-2:fake-endpoint-name-2")
	eids := key.EndpointIdentifiers()
	expectedRelationUUID := tc.Must(c, corerelation.NewUUID)
	s.st.EXPECT().GetRegularRelationUUIDByEndpointIdentifiers(
		gomock.Any(), eids[0], eids[1],
	).Return(expectedRelationUUID, nil)

	unitName := unittesting.GenNewName(c, "test/0")
	unitUUID := tc.Must(c, coreunit.NewUUID)
	s.st.EXPECT().GetUnitUUIDByName(c.Context(), unitName).Return(unitUUID, nil)

	expected := internal.CommitHookChangesArg{
		UnitUUID: unitUUID,
		RelationSettings: []internal.RelationSettings{{
			RelationUUID: expectedRelationUUID,
			Settings:     map[string]string{"key": "value"},
		},
		},
	}
	s.st.EXPECT().CommitHookChanges(c.Context(), expected).Return(nil)

	arg := unitstate.CommitHookChangesArg{
		UnitName: unitName,
		RelationSettings: []unitstate.RelationSettings{{
			RelationKey: key,
			Settings:    map[string]string{"key": "value"},
		}},
	}

	// Act
	err := s.svc.CommitHookChanges(c.Context(), arg)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
}

func (s *commitHookSuite) TestCommitHookChangesLeadership(c *tc.C) {
	defer s.setupMocks(c).Finish()
	// Arrange: args which indicate leadership is required
	key := corerelationtesting.GenNewKey(c, "app-1:fake-endpoint-name-1 app-2:fake-endpoint-name-2")
	eids := key.EndpointIdentifiers()
	expectedRelationUUID := tc.Must(c, corerelation.NewUUID)
	s.st.EXPECT().GetRegularRelationUUIDByEndpointIdentifiers(
		gomock.Any(), eids[0], eids[1],
	).Return(expectedRelationUUID, nil)

	arg := unitstate.CommitHookChangesArg{
		UnitName: unittesting.GenNewName(c, "test/0"),
		RelationSettings: []unitstate.RelationSettings{{
			RelationKey:         key,
			ApplicationSettings: map[string]string{"key": "value"},
		}},
	}

	unitUUID := tc.Must(c, coreunit.NewUUID)
	s.st.EXPECT().GetUnitUUIDByName(c.Context(), arg.UnitName).Return(unitUUID, nil)
	s.leadershipEnsurer.EXPECT().WithLeader(c.Context(), "test", "test/0", gomock.Any()).Return(nil)

	// Act
	err := s.svc.CommitHookChanges(c.Context(), arg)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
}

func (s *commitHookSuite) TestCommitHookChangesUpdateNetworkInfo(c *tc.C) {
	defer s.setupMocks(c).Finish()

	// Arrange: Unit
	unitName := unittesting.GenNewName(c, "test/0")
	unitUUID := tc.Must(c, coreunit.NewUUID)
	s.st.EXPECT().GetUnitUUIDByName(c.Context(), unitName).Return(unitUUID, nil)

	// Arrange: GetUnitRelationIngressAddress state call.
	relation1UUID := tc.Must(c, corerelation.NewUUID)
	relation2UUID := tc.Must(c, corerelation.NewUUID)
	relationIngressAddresses := map[corerelation.UUID]string{
		relation1UUID: "10.0.0.6",
		relation2UUID: "10.0.1.23",
	}
	s.st.EXPECT().GetUnitRelationsIngressAddress(c.Context(), unitUUID).Return(relationIngressAddresses, nil)

	// Arrange: GetUnitRelationEgressSubnetsByUnitUUID state call.
	relationEgressSubnets := map[corerelation.UUID][]string{
		relation1UUID: {"192.0.2.0/24", "192.51.100.0/24"},
		relation2UUID: {"203.0.113.0/24"},
	}
	relUUIDs := []corerelation.UUID{relation1UUID, relation2UUID}
	s.st.EXPECT().GetRelationUUIDsByUnitUUID(c.Context(), unitUUID).Return(relUUIDs, nil)
	s.st.EXPECT().GetRelationsEgressSubnetsByUnitUUID(c.Context(), unitUUID).Return(relationEgressSubnets, nil)

	// Arrange relation key
	key := corerelationtesting.GenNewKey(c, "app-1:fake-endpoint-name-1 app-2:fake-endpoint-name-2")
	eids := key.EndpointIdentifiers()
	s.st.EXPECT().GetRegularRelationUUIDByEndpointIdentifiers(
		gomock.Any(), eids[0], eids[1],
	).Return(relation1UUID, nil)

	// Arrange: CommitHookChanges state call
	expected := internal.CommitHookChangesArg{
		UnitUUID: unitUUID,
		RelationSettings: []internal.RelationSettings{
			{
				RelationUUID: relation1UUID,
				Settings: map[string]string{
					"key":                       "value",
					unitstate.IngressAddressKey: "10.0.0.6",
					unitstate.EgressSubnetsKey:  "192.0.2.0/24, 192.51.100.0/24",
				},
			}, {
				RelationUUID: relation2UUID,
				Settings: map[string]string{
					unitstate.IngressAddressKey: "10.0.1.23",
					unitstate.EgressSubnetsKey:  "203.0.113.0/24",
				},
			},
		},
	}
	s.st.EXPECT().CommitHookChanges(c.Context(), expected).Return(nil)

	// Arrange args for call
	arg := unitstate.CommitHookChangesArg{
		UnitName: unitName,
		RelationSettings: []unitstate.RelationSettings{{
			RelationKey: key,
			Settings:    map[string]string{"key": "value"},
		}},
		UpdateNetworkInfo: true,
	}

	// Act
	err := s.svc.CommitHookChanges(c.Context(), arg)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
}

func (s *commitHookSuite) TestGetUnitRelationNetworksNoNetworkingSupport(c *tc.C) {
	s.networkProviderGetter = func(context.Context) (ProviderWithNetworking, error) {
		return nil, errors.Errorf("provider %w", coreerrors.NotSupported)
	}
	s.setupMocks(c).Finish()

	// Arrange: Unit
	unitUUID := tc.Must(c, coreunit.NewUUID)

	// Arrange: GetUnitIngressAddress state call.
	relation1UUID := tc.Must(c, corerelation.NewUUID)
	relation2UUID := tc.Must(c, corerelation.NewUUID)
	relUUIDs := []corerelation.UUID{relation1UUID, relation2UUID}
	s.st.EXPECT().GetRelationUUIDsByUnitUUID(c.Context(), unitUUID).Return(relUUIDs, nil)
	s.st.EXPECT().GetUnitIngressAddress(c.Context(), unitUUID).Return("10.0.0.6", nil)

	// Arrange: GetUnitRelationEgressSubnetsByUnitUUID state call.
	relationEgressSubnets := map[corerelation.UUID][]string{
		relation1UUID: {"192.0.2.0/24", "192.51.100.0/24"},
		relation2UUID: {"203.0.113.0/24"},
	}
	s.st.EXPECT().GetRelationsEgressSubnetsByUnitUUID(c.Context(), unitUUID).Return(relationEgressSubnets, nil)

	// Arrange: CommitHookChanges state call
	expected := []internal.RelationNetworkInfo{
		{
			RelationUUID:   relation1UUID,
			IngressAddress: "10.0.0.6",
			EgressSubnets:  "192.0.2.0/24, 192.51.100.0/24",
		}, {
			RelationUUID:   relation2UUID,
			IngressAddress: "10.0.0.6",
			EgressSubnets:  "203.0.113.0/24",
		},
	}

	// Act:
	relationNetworkInfos, err := s.svc.getUnitRelationNetworks(c.Context(), unitUUID)

	// Assert:
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(relationNetworkInfos, tc.SameContents, expected)
}

func (s *commitHookSuite) TestGetUnitRelationNetworksUseModelEgressSubnets(c *tc.C) {
	s.setupMocks(c).Finish()

	// Arrange: Unit
	unitUUID := tc.Must(c, coreunit.NewUUID)

	// Arrange: GetUnitRelationIngressAddress state call.
	relation1UUID := tc.Must(c, corerelation.NewUUID)
	relation2UUID := tc.Must(c, corerelation.NewUUID)
	relationIngressAddresses := map[corerelation.UUID]string{
		relation1UUID: "10.0.0.6",
		relation2UUID: "10.0.1.23",
	}
	relUUIDs := []corerelation.UUID{relation1UUID, relation2UUID}
	s.st.EXPECT().GetRelationUUIDsByUnitUUID(c.Context(), unitUUID).Return(relUUIDs, nil)
	s.st.EXPECT().GetUnitRelationsIngressAddress(c.Context(), unitUUID).Return(relationIngressAddresses, nil)

	// Arrange: setup egress addresses from relation and model
	relationEgressSubnets := map[corerelation.UUID][]string{
		relation2UUID: {"203.0.113.0/24"},
	}
	s.st.EXPECT().GetRelationsEgressSubnetsByUnitUUID(c.Context(), unitUUID).Return(relationEgressSubnets, nil)
	s.st.EXPECT().GetModelEgressSubnets(c.Context()).Return([]string{"192.0.2.0/24", "192.51.100.0/24"}, nil)

	// Arrange: CommitHookChanges state call
	expected := []internal.RelationNetworkInfo{
		{
			RelationUUID:   relation1UUID,
			IngressAddress: "10.0.0.6",
			EgressSubnets:  "192.0.2.0/24, 192.51.100.0/24",
		}, {
			RelationUUID:   relation2UUID,
			IngressAddress: "10.0.1.23",
			EgressSubnets:  "203.0.113.0/24",
		},
	}

	// Act:
	relationNetworkInfos, err := s.svc.getUnitRelationNetworks(c.Context(), unitUUID)

	// Assert:
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(relationNetworkInfos, tc.SameContents, expected)
}

func (s *commitHookSuite) TestGetUnitRelationNetworksMixEgressSubnets(c *tc.C) {
	s.setupMocks(c).Finish()

	// Arrange: Unit
	unitUUID := tc.Must(c, coreunit.NewUUID)

	// Arrange: GetUnitRelationsIngressAddress state call.
	relation1UUID := tc.Must(c, corerelation.NewUUID)
	relation2UUID := tc.Must(c, corerelation.NewUUID)
	relationIngressAddresses := map[corerelation.UUID]string{
		relation1UUID: "10.0.0.6",
		relation2UUID: "10.0.1.23",
	}
	relUUIDs := []corerelation.UUID{relation1UUID, relation2UUID}
	s.st.EXPECT().GetRelationUUIDsByUnitUUID(c.Context(), unitUUID).Return(relUUIDs, nil)
	s.st.EXPECT().GetUnitRelationsIngressAddress(c.Context(), unitUUID).Return(relationIngressAddresses, nil)

	// Arrange: setup egress addresses from model
	s.st.EXPECT().GetRelationsEgressSubnetsByUnitUUID(c.Context(), unitUUID).Return(nil, nil)
	s.st.EXPECT().GetModelEgressSubnets(c.Context()).Return([]string{"192.0.2.0/24", "192.51.100.0/24"}, nil)

	// Arrange: CommitHookChanges state call
	expected := []internal.RelationNetworkInfo{
		{
			RelationUUID:   relation1UUID,
			IngressAddress: "10.0.0.6",
			EgressSubnets:  "192.0.2.0/24, 192.51.100.0/24",
		}, {
			RelationUUID:   relation2UUID,
			IngressAddress: "10.0.1.23",
			EgressSubnets:  "192.0.2.0/24, 192.51.100.0/24",
		},
	}

	// Act:
	relationNetworkInfos, err := s.svc.getUnitRelationNetworks(c.Context(), unitUUID)

	// Assert:
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(relationNetworkInfos, tc.SameContents, expected)
}

func (s *commitHookSuite) TestGetUnitRelationNetworksUseUnitPublicEgressSubnets(c *tc.C) {
	s.setupMocks(c).Finish()

	// Arrange: Unit
	unitUUID := tc.Must(c, coreunit.NewUUID)

	// Arrange: GetUnitRelationsIngressAddress state call.
	relation1UUID := tc.Must(c, corerelation.NewUUID)
	relation2UUID := tc.Must(c, corerelation.NewUUID)
	relationIngressAddresses := map[corerelation.UUID]string{
		relation1UUID: "10.0.0.6",
		relation2UUID: "10.0.1.23",
	}
	relUUIDs := []corerelation.UUID{relation1UUID, relation2UUID}
	s.st.EXPECT().GetRelationUUIDsByUnitUUID(c.Context(), unitUUID).Return(relUUIDs, nil)
	s.st.EXPECT().GetUnitRelationsIngressAddress(c.Context(), unitUUID).Return(relationIngressAddresses, nil)

	// Arrange: setup egress addresses from unit public egress address
	s.st.EXPECT().GetRelationsEgressSubnetsByUnitUUID(c.Context(), unitUUID).Return(nil, nil)
	s.st.EXPECT().GetModelEgressSubnets(c.Context()).Return(nil, nil)
	s.st.EXPECT().GetUnitPublicAddressForEgress(c.Context(), unitUUID).Return("192.51.100.0/24", nil)

	// Arrange: CommitHookChanges state call
	expected := []internal.RelationNetworkInfo{
		{
			RelationUUID:   relation1UUID,
			IngressAddress: "10.0.0.6",
			EgressSubnets:  "192.51.100.0/32",
		}, {
			RelationUUID:   relation2UUID,
			IngressAddress: "10.0.1.23",
			EgressSubnets:  "192.51.100.0/32",
		},
	}

	// Act:
	relationNetworkInfos, err := s.svc.getUnitRelationNetworks(c.Context(), unitUUID)

	// Assert:
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(relationNetworkInfos, tc.SameContents, expected)

}

func (s *commitHookSuite) TestGetRelationUUIDByKeyPeer(c *tc.C) {
	defer s.setupMocks(c).Finish()

	// Arrange:
	key := corerelationtesting.GenNewKey(c, "app-1:fake-endpoint-name-1")

	expectedRelationUUID := corerelationtesting.GenRelationUUID(c)

	s.st.EXPECT().GetPeerRelationUUIDByEndpointIdentifiers(
		gomock.Any(), key.EndpointIdentifiers()[0],
	).Return(expectedRelationUUID, nil)

	// Act:
	uuid, err := s.svc.getRelationUUIDByKey(c.Context(), key)

	// Assert:
	c.Assert(err, tc.ErrorIsNil)
	c.Check(uuid, tc.Equals, expectedRelationUUID)
}

func (s *commitHookSuite) TestGetRelationUUIDByKeyRegular(c *tc.C) {
	defer s.setupMocks(c).Finish()

	// Arrange:
	key := corerelationtesting.GenNewKey(c, "app-1:fake-endpoint-name-1 app-2:fake-endpoint-name-2")
	eids := key.EndpointIdentifiers()

	expectedRelationUUID := corerelationtesting.GenRelationUUID(c)
	s.st.EXPECT().GetRegularRelationUUIDByEndpointIdentifiers(
		gomock.Any(), eids[0], eids[1],
	).Return(expectedRelationUUID, nil)

	// Act:
	uuid, err := s.svc.getRelationUUIDByKey(c.Context(), key)

	// Assert:
	c.Assert(err, tc.ErrorIsNil)
	c.Check(uuid, tc.Equals, expectedRelationUUID)
}

func (s *commitHookSuite) TestGetRelationUUIDByKeyRelationNotFound(c *tc.C) {
	defer s.setupMocks(c).Finish()

	// Arrange:
	s.st.EXPECT().GetRegularRelationUUIDByEndpointIdentifiers(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return("", relationerrors.RelationNotFound)

	// Act:
	_, err := s.svc.getRelationUUIDByKey(
		c.Context(),
		corerelationtesting.GenNewKey(c, "app-1:fake-endpoint-name-1 app-2:fake-endpoint-name-2"),
	)

	// Assert:
	c.Assert(err, tc.ErrorIs, relationerrors.RelationNotFound)
}

func (s *commitHookSuite) setupMocks(c *tc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.st = NewMockState(ctrl)
	s.leadershipEnsurer = NewMockEnsurer(ctrl)

	if s.networkProviderGetter == nil {
		s.providerWithNetworking = NewMockProviderWithNetworking(ctrl)
		s.networkProviderGetter = func(context.Context) (ProviderWithNetworking, error) {
			return s.providerWithNetworking, nil
		}
	}

	s.svc = NewLeadershipService(
		s.st,
		s.leadershipEnsurer,
		s.networkProviderGetter,
		loggertesting.WrapCheckLog(c),
	)

	c.Cleanup(func() {
		s.svc = nil
		s.st = nil
		s.leadershipEnsurer = nil
		s.providerWithNetworking = nil
		s.networkProviderGetter = nil
	})

	return ctrl
}
