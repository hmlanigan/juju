// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"

	"github.com/canonical/sqlair"
	"github.com/juju/collections/set"
	"github.com/juju/errors"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/domain"
	"github.com/juju/juju/internal/uuid"
)

// State represents the persistence layer for opened ports.
type State struct {
	*domain.StateBase
}

// NewState returns a new state reference.
func NewState(factory database.TxnRunnerFactory) *State {
	return &State{
		StateBase: domain.NewStateBase(factory),
	}
}

// GetOpenedPorts returns the opened ports for a given unit uuid,
// grouped by endpoint.
func (st *State) GetOpenedPorts(ctx context.Context, unit string) (network.GroupedPortRanges, error) {
	db, err := st.DB()
	if err != nil {
		return nil, errors.Trace(err)
	}

	unitUUID := unitUUID{UUID: unit}

	query, err := st.Prepare(`
SELECT &endpointPortRange.*
FROM port_range
JOIN protocol ON port_range.protocol_id = protocol.id
JOIN unit_endpoint ON port_range.unit_endpoint_uuid = unit_endpoint.uuid
WHERE unit_endpoint.unit_uuid = $unitUUID.unit_uuid
`, endpointPortRange{}, unitUUID)
	if err != nil {
		return nil, errors.Annotate(err, "preparing get opened ports statement")
	}

	results := []endpointPortRange{}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, query, unitUUID).GetAll(&results)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return errors.Trace(err)
	})
	if err != nil {
		return nil, errors.Annotatef(err, "getting opened ports for unit %q", unit)
	}

	groupedPortRanges := network.GroupedPortRanges{}
	for _, endpointPortRange := range results {
		endpointName := endpointPortRange.Endpoint
		if _, ok := groupedPortRanges[endpointName]; !ok {
			groupedPortRanges[endpointPortRange.Endpoint] = []network.PortRange{}
		}
		groupedPortRanges[endpointName] = append(groupedPortRanges[endpointName], endpointPortRange.decode())
	}

	for _, portRanges := range groupedPortRanges {
		network.SortPortRanges(portRanges)
	}

	return groupedPortRanges, nil
}

// GetEndpointOpenedPorts returns the opened ports for a given endpoint of a
// given unit.
func (st *State) GetEndpointOpenedPorts(ctx domain.AtomicContext, unit string, endpoint string) ([]network.PortRange, error) {
	unitUUID := unitUUID{UUID: unit}
	endpointName := endpointName{Endpoint: endpoint}

	query, err := st.Prepare(`
SELECT &portRange.*
FROM port_range
JOIN protocol ON port_range.protocol_id = protocol.id
JOIN unit_endpoint ON port_range.unit_endpoint_uuid = unit_endpoint.uuid
WHERE unit_endpoint.unit_uuid = $unitUUID.unit_uuid
AND unit_endpoint.endpoint = $endpointName.endpoint
`, portRange{}, unitUUID, endpointName)
	if err != nil {
		return nil, errors.Annotate(err, "preparing get opened ports statement")
	}

	var portRanges []portRange
	err = domain.Run(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, query, unitUUID, endpointName).GetAll(&portRanges)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return errors.Trace(err)
	})
	if err != nil {
		return nil, errors.Annotatef(err, "getting opened ports for unit %q", unit)
	}

	decodedPortRanges := make([]network.PortRange, len(portRanges))
	for i, pr := range portRanges {
		decodedPortRanges[i] = pr.decode()
	}
	network.SortPortRanges(decodedPortRanges)
	return decodedPortRanges, nil
}

// UpdateUnitPorts opens and closes ports for the endpoints of a given unit.
// The service layer must ensure that opened and closed ports for the same
// endpoints must not conflict.
func (st *State) UpdateUnitPorts(
	ctx domain.AtomicContext, unit string, openPorts, closePorts network.GroupedPortRanges,
) error {
	endpointsUnderActionSet := set.NewStrings()
	for endpoint := range openPorts {
		endpointsUnderActionSet.Add(endpoint)
	}
	for endpoint := range closePorts {
		endpointsUnderActionSet.Add(endpoint)
	}
	endpointsUnderAction := endpoints(endpointsUnderActionSet.Values())

	unitUUID := unitUUID{UUID: unit}

	return domain.Run(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		currentOpenedPorts, err := st.getOpenedPorts(ctx, tx, unitUUID)
		if err != nil {
			return errors.Annotatef(err, "getting opened ports for unit %q", unit)
		}

		allInputPortRanges := append(openPorts.UniquePortRanges(), closePorts.UniquePortRanges()...)
		err = verifyNoPortRangeConflicts(currentOpenedPorts, allInputPortRanges)
		if err != nil {
			return errors.Annotatef(err, "port ranges conflict on unit %q", unit)
		}

		endpoints, err := st.ensureEndpoints(ctx, tx, unitUUID, endpointsUnderAction)
		if err != nil {
			return errors.Annotatef(err, "ensuring endpoints exist for unit %q", unit)
		}

		err = st.openPorts(ctx, tx, openPorts, currentOpenedPorts, endpoints)
		if err != nil {
			return errors.Annotatef(err, "opening ports for unit %q", unit)
		}

		err = st.closePorts(ctx, tx, closePorts, currentOpenedPorts)
		if err != nil {
			return errors.Annotatef(err, "closing ports for unit %q", unit)
		}

		return nil
	})
}

// ensureEndpoints ensures that the given endpoints are present in the database.
// Return all endpoints under action with their corresponding UUIDs.
//
// TODO(jack-w-shaw): Once it has been implemented, we should verify new endpoints
// are valid by checking the charm_relation table.
func (st *State) ensureEndpoints(
	ctx context.Context, tx *sqlair.TX, unitUUID unitUUID, endpointsUnderAction endpoints,
) ([]endpoint, error) {
	getUnitEndpoints, err := st.Prepare(`
SELECT &endpoint.*
FROM unit_endpoint
WHERE unit_uuid = $unitUUID.unit_uuid
AND endpoint IN ($endpoints[:])
`, endpoint{}, unitUUID, endpointsUnderAction)
	if err != nil {
		return nil, errors.Annotate(err, "preparing get unit endpoints statement")
	}

	insertUnitEndpoint, err := st.Prepare("INSERT INTO unit_endpoint (*) VALUES ($unitEndpoint.*)", unitEndpoint{})
	if err != nil {
		return nil, errors.Annotate(err, "preparing insert unit endpoint statement")
	}

	var endpoints []endpoint
	err = tx.Query(ctx, getUnitEndpoints, unitUUID, endpointsUnderAction).GetAll(&endpoints)
	if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
		return nil, errors.Trace(err)
	}

	foundEndpoints := set.NewStrings()
	for _, ep := range endpoints {
		foundEndpoints.Add(ep.Endpoint)
	}

	// Insert any new endpoints that are required.
	requiredEndpoints := set.NewStrings(endpointsUnderAction...).Difference(foundEndpoints)
	newUnitEndpoints := make([]unitEndpoint, requiredEndpoints.Size())
	for i, requiredEndpoint := range requiredEndpoints.Values() {
		uuid, err := uuid.NewUUID()
		if err != nil {
			return nil, errors.Annotatef(err, "generating UUID for unit endpoint")
		}
		newUnitEndpoints[i] = unitEndpoint{
			UUID:     uuid.String(),
			UnitUUID: unitUUID.UUID,
			Endpoint: requiredEndpoint,
		}
		endpoints = append(endpoints, endpoint{
			Endpoint: requiredEndpoint,
			UUID:     uuid.String(),
		})
	}

	if len(newUnitEndpoints) > 0 {
		err = tx.Query(ctx, insertUnitEndpoint, newUnitEndpoints).Run()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return endpoints, nil
}

// getOpenedPorts returns the opened ports for the given unit.
//
// NOTE: This differs from GetOpenedPorts in that it returns port ranges with
// their UUIDs, which are not needed by GetOpenedPorts.
func (st *State) getOpenedPorts(ctx context.Context, tx *sqlair.TX, unitUUID unitUUID) ([]endpointPortRangeUUID, error) {
	getOpenedPorts, err := st.Prepare(`
SELECT 
	port_range.uuid AS &endpointPortRangeUUID.uuid,
	protocol.protocol AS &endpointPortRangeUUID.protocol,
	port_range.from_port AS &endpointPortRangeUUID.from_port,
	port_range.to_port AS &endpointPortRangeUUID.to_port,
	unit_endpoint.endpoint AS &endpointPortRangeUUID.endpoint
FROM port_range
JOIN protocol ON port_range.protocol_id = protocol.id
INNER JOIN unit_endpoint ON port_range.unit_endpoint_uuid = unit_endpoint.uuid
WHERE unit_endpoint.unit_uuid = $unitUUID.unit_uuid
`, endpointPortRangeUUID{}, unitUUID)
	if err != nil {
		return nil, errors.Annotate(err, "preparing get opened ports statement")
	}

	openedPorts := []endpointPortRangeUUID{}
	err = tx.Query(ctx, getOpenedPorts, unitUUID).GetAll(&openedPorts)
	if errors.Is(err, sqlair.ErrNoRows) {
		return []endpointPortRangeUUID{}, nil
	}
	if err != nil {
		return nil, errors.Trace(err)
	}
	return openedPorts, nil
}

func verifyNoPortRangeConflicts(currentOpenedPorts []endpointPortRangeUUID, inputPortRanges []network.PortRange) error {
	for _, inputPortRange := range inputPortRanges {
		for _, openedPortRange := range currentOpenedPorts {
			if inputPortRange == openedPortRange.decode() {
				// We allow port ranges to conflict only if they are equal
				continue
			}
			if inputPortRange.ConflictsWith(openedPortRange.decode()) {
				return errors.Annotatef(ErrPortRangeConflict,
					"port range %q conflicts with existing port range %q",
					inputPortRange, openedPortRange.decode())
			}
		}
	}
	return nil
}

// openPorts inserts the given port ranges into the database, unless they're already open.
func (st *State) openPorts(
	ctx context.Context, tx *sqlair.TX,
	openPorts network.GroupedPortRanges, currentOpenedPorts []endpointPortRangeUUID, endpoints []endpoint,
) error {
	insertPortRange, err := st.Prepare("INSERT INTO port_range (*) VALUES ($unitPortRange.*)", unitPortRange{})
	if err != nil {
		return errors.Annotate(err, "preparing insert port range statement")
	}

	protocolMap, err := st.getProtocolMap(ctx, tx)
	if err != nil {
		return errors.Annotate(err, "getting protocol map")
	}

	// index the current opened ports by endpoint and port range
	currentOpenedPortRangeExistenceIndex := make(map[string]map[network.PortRange]bool)
	for _, openedPortRange := range currentOpenedPorts {
		if _, ok := currentOpenedPortRangeExistenceIndex[openedPortRange.Endpoint]; !ok {
			currentOpenedPortRangeExistenceIndex[openedPortRange.Endpoint] = make(map[network.PortRange]bool)
		}
		currentOpenedPortRangeExistenceIndex[openedPortRange.Endpoint][openedPortRange.decode()] = true
	}

	// Construct the new port ranges to open
	var openPortRanges []unitPortRange

	for _, ep := range endpoints {
		ports, ok := openPorts[ep.Endpoint]
		if !ok {
			continue
		}

		for _, portRange := range ports {
			// skip port range if it's already open on this endpoint
			if _, ok := currentOpenedPortRangeExistenceIndex[ep.Endpoint][portRange]; ok {
				continue
			}

			uuid, err := uuid.NewUUID()
			if err != nil {
				return errors.Annotatef(err, "generating UUID for unit endpoint")
			}
			openPortRanges = append(openPortRanges, unitPortRange{
				UUID:             uuid.String(),
				ProtocolID:       protocolMap[portRange.Protocol],
				FromPort:         portRange.FromPort,
				ToPort:           portRange.ToPort,
				UnitEndpointUUID: ep.UUID,
			})
		}
	}

	if len(openPortRanges) > 0 {
		err = tx.Query(ctx, insertPortRange, openPortRanges).Run()
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// getProtocolMap returns a map of protocol names to their IDs in DQLite.
func (st *State) getProtocolMap(ctx context.Context, tx *sqlair.TX) (map[string]int, error) {
	getProtocols, err := st.Prepare("SELECT &protocol.* FROM protocol", protocol{})
	if err != nil {
		return nil, errors.Annotate(err, "preparing get protocol ID statement")
	}

	protocols := []protocol{}
	err = tx.Query(ctx, getProtocols).GetAll(&protocols)
	if err != nil {
		return nil, errors.Trace(err)
	}

	protocolMap := map[string]int{}
	for _, protocol := range protocols {
		protocolMap[protocol.Name] = protocol.ID
	}

	return protocolMap, nil
}

// closePorts removes the given port ranges from the database, if they exist.
func (st *State) closePorts(
	ctx context.Context, tx *sqlair.TX, closePorts network.GroupedPortRanges, currentOpenedPorts []endpointPortRangeUUID,
) error {
	closePortRanges, err := st.Prepare(`
DELETE FROM port_range
WHERE uuid IN ($portRangeUUIDs[:])
`, portRangeUUIDs{})
	if err != nil {
		return errors.Annotate(err, "preparing close port range statement")
	}

	// index the uuids of current opened ports by endpoint and port range
	openedPortRangeUUIDIndex := make(map[string]map[network.PortRange]string)
	for _, openedPortRange := range currentOpenedPorts {
		if _, ok := openedPortRangeUUIDIndex[openedPortRange.Endpoint]; !ok {
			openedPortRangeUUIDIndex[openedPortRange.Endpoint] = make(map[network.PortRange]string)
		}
		openedPortRangeUUIDIndex[openedPortRange.Endpoint][openedPortRange.decode()] = openedPortRange.UUID
	}

	// Find the uuids of port ranges to close
	var closePortRangeUUIDs portRangeUUIDs
	for endpoint, portRanges := range closePorts {
		index, ok := openedPortRangeUUIDIndex[endpoint]
		if !ok {
			continue
		}

		for _, closePortRange := range portRanges {
			openedRangeUUID, ok := index[closePortRange]
			if !ok {
				continue
			}
			closePortRangeUUIDs = append(closePortRangeUUIDs, openedRangeUUID)
		}
	}

	if len(closePortRangeUUIDs) > 0 {
		err = tx.Query(ctx, closePortRanges, closePortRangeUUIDs).Run()
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// GetEndpoints returns all the endpoints for the given unit
//
// TODO(jack-w-shaw): At the moment, we calculate this by checking the unit_endpoints
// table. However, this will not always return a complete list, as it only includes
// endpoints that have had ports opened on them at some point.
//
// Once it has been implemented, we should check the charm_relation table to get a
// complete list of endpoints instead.
func (st *State) GetEndpoints(ctx domain.AtomicContext, unit string) ([]string, error) {
	unitUUID := unitUUID{UUID: unit}

	getEndpoints, err := st.Prepare(`
SELECT &endpointName.*
FROM unit_endpoint
WHERE unit_endpoint.unit_uuid = $unitUUID.unit_uuid
`, endpointName{}, unitUUID)
	if err != nil {
		return nil, errors.Annotate(err, "preparing get endpoints statement")
	}

	var endpointNames []endpointName
	err = domain.Run(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, getEndpoints, unitUUID).GetAll(&endpointNames)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return errors.Trace(err)
	})
	if err != nil {
		return nil, errors.Annotatef(err, "getting endpoints for unit %q", unit)
	}

	endpoints := make([]string, len(endpointNames))
	for i, ep := range endpointNames {
		endpoints[i] = ep.Endpoint
	}
	return endpoints, nil
}