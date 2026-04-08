// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"testing"

	"github.com/juju/tc"

	corenetwork "github.com/juju/juju/core/network"
	"github.com/juju/juju/core/relation"
	"github.com/juju/juju/domain/life"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/internal/uuid"
)

type infoSuite struct {
	commitHookBaseSuite
	relationCount int
}

func TestInfoSuite(t *testing.T) {
	tc.Run(t, &infoSuite{})
}

func (s *infoSuite) TestGetUnitPublicAddressForEgress(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	subnetUUID := s.addSubnet(c, "198.51.100.0/24", spaceUUID)
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, "10.0.0.10/24",
		corenetwork.ScopeCloudLocal,
	)
	secondaryUUID := s.addIPAddressWithSubnetAndOrigin(
		c, deviceUUID, nodeUUID, subnetUUID, "198.51.100.1/24", 1,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, secondaryUUID)
	s.markIPAddressSecondary(c, secondaryUUID)
	publicUUID := s.addIPAddressWithSubnetAndOrigin(
		c, deviceUUID, nodeUUID, subnetUUID, "198.51.100.10/24", 1,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, publicUUID)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	address, err := s.state.GetUnitPublicAddressForEgress(
		c.Context(), unitUUID,
	)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(address, tc.Equals, "198.51.100.10/24")
}

func (s *infoSuite) TestGetUnitPublicAddressForEgressWithoutPublicAddress(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	subnetUUID := s.addSubnet(c, "10.0.0.0/24", spaceUUID)
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, "10.0.0.10/24",
		corenetwork.ScopeCloudLocal,
	)
	s.addIPAddressWithSubnetAndOrigin(
		c, deviceUUID, nodeUUID, subnetUUID, "198.51.100.10/24", 1,
	)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	address, err := s.state.GetUnitPublicAddressForEgress(
		c.Context(), unitUUID,
	)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(address, tc.Equals, "")
}

func (s *infoSuite) TestGetUnitPublicAddressForEgressDeadUnit(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	subnetUUID := s.addSubnet(c, "198.51.100.0/24", spaceUUID)
	publicUUID := s.addIPAddressWithSubnetAndOrigin(
		c, deviceUUID, nodeUUID, subnetUUID, "198.51.100.10/24", 1,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, publicUUID)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)
	s.query(c, `UPDATE unit SET life_id = ? WHERE uuid = ?`, life.Dead, unitUUID)

	address, err := s.state.GetUnitPublicAddressForEgress(
		c.Context(), unitUUID,
	)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(address, tc.Equals, "198.51.100.10/24")
}

func (s *infoSuite) TestGetModelEgressSubnets(c *tc.C) {
	s.query(c, `INSERT INTO model_config VALUES (?, ?)`,
		config.EgressSubnets, "10.0.1.0/24, 10.0.2.0/24")

	cidrs, err := s.state.GetModelEgressSubnets(c.Context())
	c.Assert(err, tc.ErrorIsNil)
	c.Check(cidrs, tc.DeepEquals, []string{"10.0.1.0/24", "10.0.2.0/24"})
}

func (s *infoSuite) TestGetModelEgressSubnetsEmpty(c *tc.C) {
	cidrs, err := s.state.GetModelEgressSubnets(c.Context())
	c.Assert(err, tc.ErrorIsNil)
	c.Check(cidrs, tc.HasLen, 0)
}

func (s *infoSuite) TestGetRelationsEgressSubnetsByUnitUUID(c *tc.C) {
	// Arrange: unit
	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, corenetwork.AlphaSpaceId.String())
	unitUUID := s.addUnitAndNetNode(c, "unit/7", appUUID, charmUUID)

	// Arrange: add endpoints
	endpointUUID := s.addApplicationEndpoint(c, appUUID, charmUUID, "endpoint1", corenetwork.AlphaSpaceId.String())
	endpoint2UUID := s.addApplicationEndpoint(c, appUUID, charmUUID, "endpoint8", corenetwork.AlphaSpaceId.String())

	// Arrange: add relation and join.
	relationUUID := s.addRelation(c)
	relationEndpointUUID := s.addRelationEndpoint(c, relationUUID.String(), endpointUUID)
	s.addRelationUnit(c, relationEndpointUUID, unitUUID.String())
	s.addRelationNetworkEgress(c, relationUUID.String(), "10.0.1.0/24")
	s.addRelationNetworkEgress(c, relationUUID.String(), "10.0.2.0/24")

	// Arrange: add second relation and join.
	relation2UUID := s.addRelation(c)
	relationEndpoint2UUID := s.addRelationEndpoint(c, relation2UUID.String(), endpoint2UUID)
	s.addRelationUnit(c, relationEndpoint2UUID, unitUUID.String())
	s.addRelationNetworkEgress(c, relation2UUID.String(), "10.0.4.0/24")

	// Act
	cidrs, err := s.state.GetRelationsEgressSubnetsByUnitUUID(c.Context(), unitUUID)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
	c.Check(cidrs, tc.DeepEquals, map[relation.UUID][]string{
		relationUUID:  {"10.0.1.0/24", "10.0.2.0/24"},
		relation2UUID: {"10.0.4.0/24"},
	})
}

func (s *infoSuite) TestGetRelationsEgressSubnetsByUnitUUIDEmpty(c *tc.C) {
	// Arrange: unit
	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, corenetwork.AlphaSpaceId.String())
	unitUUID := s.addUnitAndNetNode(c, "unit/3", appUUID, charmUUID)

	// Arrange: add endpoint
	endpointUUID := s.addApplicationEndpoint(
		c, appUUID, charmUUID, "endpoint1", corenetwork.AlphaSpaceId.String())

	// Arrange: add relation and join.
	relationUUID := s.addRelation(c)
	relationEndpointUUID := s.addRelationEndpoint(c, relationUUID.String(), endpointUUID)
	s.addRelationUnit(c, relationEndpointUUID, unitUUID.String())

	// Act
	cidrs, err := s.state.GetRelationsEgressSubnetsByUnitUUID(c.Context(), unitUUID)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
	c.Check(cidrs, tc.HasLen, 0)
}

func (s *infoSuite) TestGetUnitRelationsIngressAddresses(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	cidr := "10.0.0.0/24"
	subnetUUID := s.addSubnet(c, cidr, spaceUUID)
	expectedAddr := "10.0.0.1"
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, expectedAddr, corenetwork.ScopeCloudLocal,
	)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	endpointName := "endpoint1"
	endpointUUID := s.addApplicationEndpoint(c, appUUID, charmUUID, endpointName, "")

	relationUUID := s.addRelation(c)
	s.addRelationEndpoint(c, relationUUID.String(), endpointUUID)

	addrByRelation, err := s.state.GetUnitRelationsIngressAddresses(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(addrByRelation, tc.DeepEquals, map[relation.UUID]string{relationUUID: expectedAddr})
}

func (s *infoSuite) TestGetUnitRelationsIngressAddressesOrdersIngress(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	subnetUUID := s.addSubnet(c, "198.51.100.0/24", spaceUUID)
	s.addIPAddressWithSubnetAndOrigin(
		c, deviceUUID, nodeUUID, subnetUUID, "198.51.100.20", 0,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, "address-198.51.100.20-uuid")
	s.addIPAddressWithSubnetAndOrigin(
		c, deviceUUID, nodeUUID, subnetUUID, "198.51.100.10", 1,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, "address-198.51.100.10-uuid")

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	endpointName := "endpoint1"
	endpointUUID := s.addApplicationEndpoint(c, appUUID, charmUUID, endpointName, "")

	relationUUID := s.addRelation(c)
	s.addRelationEndpoint(c, relationUUID.String(), endpointUUID)

	addrByRelation, err := s.state.GetUnitRelationsIngressAddresses(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(
		addrByRelation,
		tc.DeepEquals,
		map[relation.UUID]string{relationUUID: "198.51.100.10"},
	)
}

func (s *infoSuite) TestGetUnitRelationsIngressAddressesPrioritisesPrimaryIngress(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	subnetUUID := s.addSubnet(c, "10.0.0.0/24", spaceUUID)
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, "10.0.0.20", corenetwork.ScopeCloudLocal,
	)
	secondaryUUID := s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, "10.0.0.10", corenetwork.ScopeCloudLocal,
	)
	s.markIPAddressSecondary(c, secondaryUUID)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	endpointName := "endpoint1"
	endpointUUID := s.addApplicationEndpoint(c, appUUID, charmUUID, endpointName, "")

	relationUUID := s.addRelation(c)
	s.addRelationEndpoint(c, relationUUID.String(), endpointUUID)

	addrByRelation, err := s.state.GetUnitRelationsIngressAddresses(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(
		addrByRelation,
		tc.DeepEquals,
		map[relation.UUID]string{relationUUID: "10.0.0.20"},
	)
}

func (s *infoSuite) TestGetUnitRelationsIngressAddressesMultipleEndpointsSameSpace(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	cidr := "10.0.0.0/24"
	subnetUUID := s.addSubnet(c, cidr, spaceUUID)
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, "10.0.0.1", corenetwork.ScopeCloudLocal,
	)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	endpoint1UUID := s.addApplicationEndpoint(c, appUUID, charmUUID, "endpoint1", "")
	endpoint2UUID := s.addApplicationEndpoint(c, appUUID, charmUUID, "endpoint2", "")

	relation1UUID := s.addRelation(c)
	s.addRelationEndpoint(c, relation1UUID.String(), endpoint1UUID)
	relation2UUID := s.addRelation(c)
	s.addRelationEndpoint(c, relation2UUID.String(), endpoint2UUID)

	addrByRelation, err := s.state.GetUnitRelationsIngressAddresses(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(addrByRelation, tc.DeepEquals, map[relation.UUID]string{
		relation1UUID: "10.0.0.1",
		relation2UUID: "10.0.0.1",
	})
}

func (s *infoSuite) TestGetUnitRelationsIngressAddressesCaasUnit(c *tc.C) {
	podNodeUUID := s.addNetNode(c)
	svcNodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, podNodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := s.addSpace(c)
	cidr := "10.0.0.0/24"
	subnetUUID := s.addSubnet(c, cidr, spaceUUID)
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, podNodeUUID, subnetUUID, "10.0.0.1", corenetwork.ScopeMachineLocal,
	)

	svcDeviceUUID := s.addLinkLayerDevice(
		c, svcNodeUUID, "eth1", "00:11:22:33:44:66", corenetwork.EthernetDevice,
	)
	s.addIPAddressWithSubnetAndScope(
		c, svcDeviceUUID, svcNodeUUID, subnetUUID, "10.0.0.2", corenetwork.ScopeCloudLocal,
	)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, podNodeUUID)
	s.addK8sService(c, svcNodeUUID, appUUID)

	endpointName := "endpoint1"
	endpointUUID := s.addApplicationEndpoint(c, appUUID, charmUUID, endpointName, spaceUUID)

	relationUUID := s.addRelation(c)
	s.addRelationEndpoint(c, relationUUID.String(), endpointUUID)

	addrByRelation, err := s.state.GetUnitRelationsIngressAddresses(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(addrByRelation, tc.DeepEquals, map[relation.UUID]string{relationUUID: "10.0.0.2"})
}

func (s *infoSuite) TestGetUnitIngressAddress(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	ethDeviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	vethDeviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "veth0", "00:11:22:33:44:66", corenetwork.VirtualEthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	subnetUUID := s.addSubnet(c, "198.51.100.0/24", spaceUUID)

	s.addIPAddressWithSubnetAndOrigin(
		c, ethDeviceUUID, nodeUUID, subnetUUID, "198.51.100.20", 0,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, "address-198.51.100.20-uuid")

	s.addIPAddressWithSubnetAndOrigin(
		c, ethDeviceUUID, nodeUUID, subnetUUID, "198.51.100.10", 1,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, "address-198.51.100.10-uuid")

	s.addIPAddressWithSubnetAndOrigin(
		c, vethDeviceUUID, nodeUUID, subnetUUID, "198.51.100.30", 1,
	)
	s.query(c, `
UPDATE ip_address
SET scope_id = (SELECT id FROM ip_address_scope WHERE name = 'public')
WHERE uuid = ?
`, "address-198.51.100.30-uuid")

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	obtainedAddress, err := s.state.GetUnitIngressAddress(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(obtainedAddress, tc.Equals, "198.51.100.10")
}

func (s *infoSuite) TestGetUnitNetworkInfoPrioritisesPrimaryIngress(c *tc.C) {
	nodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, nodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	subnetUUID := s.addSubnet(c, "10.0.0.0/24", spaceUUID)
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, "10.0.0.20", corenetwork.ScopeCloudLocal,
	)
	secondaryUUID := s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, nodeUUID, subnetUUID, "10.0.0.10", corenetwork.ScopeCloudLocal,
	)
	s.markIPAddressSecondary(c, secondaryUUID)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, nodeUUID)

	obtainedAddress, err := s.state.GetUnitIngressAddress(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(obtainedAddress, tc.Equals, "10.0.0.20")
}

func (s *infoSuite) TestGetUnitNetworkInfoCaasUnit(c *tc.C) {
	podNodeUUID := s.addNetNode(c)
	svcNodeUUID := s.addNetNode(c)
	deviceUUID := s.addLinkLayerDevice(
		c, podNodeUUID, "eth0", "00:11:22:33:44:55", corenetwork.EthernetDevice,
	)
	spaceUUID := corenetwork.AlphaSpaceId.String()
	cidr := "10.0.0.0/24"
	subnetUUID := s.addSubnet(c, cidr, spaceUUID)
	s.addIPAddressWithSubnetAndScope(
		c, deviceUUID, podNodeUUID, subnetUUID, "10.0.0.1", corenetwork.ScopeMachineLocal,
	)

	svcVethUUID := s.addLinkLayerDevice(
		c, svcNodeUUID, "veth0", "00:11:22:33:44:66", corenetwork.VirtualEthernetDevice,
	)
	s.addIPAddressWithSubnetAndScope(
		c, svcVethUUID, svcNodeUUID, subnetUUID, "10.0.0.3", corenetwork.ScopeCloudLocal,
	)

	svcEthUUID := s.addLinkLayerDevice(
		c, svcNodeUUID, "eth1", "00:11:22:33:44:77", corenetwork.EthernetDevice,
	)
	s.addIPAddressWithSubnetAndScope(
		c, svcEthUUID, svcNodeUUID, subnetUUID, "10.0.0.2", corenetwork.ScopeCloudLocal,
	)

	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, spaceUUID)
	unitUUID := s.addUnit(c, appUUID, charmUUID, podNodeUUID)
	s.addK8sService(c, svcNodeUUID, appUUID)

	obtainedAddress, err := s.state.GetUnitIngressAddress(c.Context(), unitUUID)

	c.Assert(err, tc.ErrorIsNil)
	c.Check(obtainedAddress, tc.Equals, "10.0.0.2")
}

func (s *infoSuite) TestGetRelationUUIDsByUnitUUID(c *tc.C) {
	// Arrange: unit
	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, corenetwork.AlphaSpaceId.String())
	unitUUID := s.addUnitAndNetNode(c, "unit/7", appUUID, charmUUID)

	// Arrange: add endpoints
	endpoint1UUID := s.addApplicationEndpoint(c, appUUID, charmUUID, "endpoint1", corenetwork.AlphaSpaceId.String())
	endpoint2UUID := s.addApplicationEndpoint(c, appUUID, charmUUID, "endpoint8", corenetwork.AlphaSpaceId.String())

	// Arrange: add relation and join.
	relationUUID := s.addRelation(c)
	relationEndpointUUID := s.addRelationEndpoint(c, relationUUID.String(), endpoint1UUID)
	s.addRelationUnit(c, relationEndpointUUID, unitUUID.String())

	// Arrange: add second relation and join.
	relation2UUID := s.addRelation(c)
	relationEndpoint2UUID := s.addRelationEndpoint(c, relation2UUID.String(), endpoint2UUID)
	s.addRelationUnit(c, relationEndpoint2UUID, unitUUID.String())

	// Act
	obtained, err := s.state.GetRelationUUIDsByUnitUUID(c.Context(), unitUUID)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
	c.Check(obtained, tc.DeepEquals, []relation.UUID{relationUUID, relation2UUID})
}

func (s *infoSuite) TestGetRelationUUIDsByUnitUUIDEmpty(c *tc.C) {
	// Arrange: unit
	charmUUID := s.addCharm(c)
	appUUID := s.addApplication(c, charmUUID, corenetwork.AlphaSpaceId.String())
	unitUUID := s.addUnitAndNetNode(c, "unit/3", appUUID, charmUUID)

	// Act
	obtained, err := s.state.GetRelationUUIDsByUnitUUID(c.Context(), unitUUID)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
	c.Check(obtained, tc.HasLen, 0)
}

// Helper methods

// addApplicationEndpoint creates a charm relation and an application endpoint
// in the database, returning its UUID.
func (s *infoSuite) addApplicationEndpoint(c *tc.C, appUUID, charmUUID, endpointName, spaceUUID string) string {
	// Arrange: add charm relation
	relationUUID := tc.Must(c, relation.NewUUID).String()
	s.query(c, `INSERT INTO charm_relation (uuid, charm_uuid, name, role_id, scope_id) VALUES (?, ?, ?, 0, 0)`,
		relationUUID, charmUUID, endpointName)

	// Arrange: add application endpoint
	var spacePtr *string
	if spaceUUID != "" {
		spacePtr = &spaceUUID
	}
	appEndpointUUID := uuid.MustNewUUID().String()
	s.query(c, `
INSERT INTO application_endpoint (uuid, application_uuid, charm_relation_uuid, space_uuid) 
VALUES (?, ?,?, ?)`,
		appEndpointUUID, appUUID, relationUUID, spacePtr)
	return appEndpointUUID
}

// addRelation creates a relation in the database, returning its UUID.
func (s *infoSuite) addRelation(c *tc.C) relation.UUID {
	relationUUID := tc.Must(c, relation.NewUUID)
	s.relationCount++
	s.query(c, `INSERT INTO relation (uuid, life_id, relation_id, scope_id) VALUES (?, 0, ?, 0)`,
		relationUUID, s.relationCount)
	return relationUUID
}

// addRelationEndpoint creates a relation_endpoint linking a relation to an application endpoint.
func (s *infoSuite) addRelationEndpoint(c *tc.C, relationUUID, endpointUUID string) string {
	relationEndpointUUID := tc.Must(c, relation.NewUUID).String()
	s.query(c, `INSERT INTO relation_endpoint (uuid, relation_uuid, endpoint_uuid) VALUES (?, ?, ?)`,
		relationEndpointUUID, relationUUID, endpointUUID)
	return relationEndpointUUID
}

// addRelationUnit creates a relation_unit linking a relation endpoint to a unit.
func (s *infoSuite) addRelationUnit(c *tc.C, relationEndpointUUID, unitUUID string) {
	relationUnitUUID := tc.Must(c, relation.NewUUID).String()
	s.query(c, `INSERT INTO relation_unit (uuid, relation_endpoint_uuid, unit_uuid) VALUES (?, ?, ?)`,
		relationUnitUUID, relationEndpointUUID, unitUUID)
}

// addRelationNetworkEgress adds an egress CIDR to a relation.
func (s *infoSuite) addRelationNetworkEgress(c *tc.C, relationUUID, cidr string) {
	s.query(c, `INSERT INTO relation_network_egress (relation_uuid, cidr) VALUES (?, ?)`,
		relationUUID, cidr)
}

// addIPAddressWithSubnet adds an IP address to the database and returns its UUID.
func (s *infoSuite) addIPAddressWithSubnetAndOrigin(c *tc.C, deviceUUID, netNodeUUID,
	subnetUUID, addressValue string, origin int) string {

	addressUUID := "address-" + addressValue + "-uuid"

	s.query(c, `
		INSERT INTO ip_address (uuid, device_uuid, address_value, net_node_uuid, subnet_uuid, type_id, config_type_id, origin_id, scope_id, is_secondary, is_shadow)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, addressUUID, deviceUUID, addressValue, netNodeUUID, subnetUUID, 0, 4, origin, 0,
		false, false)

	return addressUUID
}

func (s *infoSuite) markIPAddressSecondary(c *tc.C, addressUUID string) {
	s.query(c, `
UPDATE ip_address
SET is_secondary = true
WHERE uuid = ?
`, addressUUID)
}
