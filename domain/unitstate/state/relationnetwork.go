// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"
	"strings"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/core/relation"
	coreunit "github.com/juju/juju/core/unit"
	"github.com/juju/juju/domain/unitstate"
	"github.com/juju/juju/internal/errors"
)

// GetUnitIngressAddress retrieves ingress address for the specified unit
// when provider networking is not supported.
func (st *State) GetUnitIngressAddress(
	ctx context.Context,
	unitUUID coreunit.UUID,
) (string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return "", errors.Capture(err)
	}

	type InSpaceAddress spaceAddress
	type unitIngressAddressRow struct {
		IngressAddress sql.NullString `db:"ingress_address_value"`
	}

	var rows []unitIngressAddressRow
	ident := entityUUID{UUID: unitUUID.String()}
	stmt, err := st.Prepare(`
WITH candidate AS (
    SELECT urn.address_value,
           urn.device_uuid,
           urn.scope_rank,
           urn.origin_id,
           urn.is_secondary,
           urn.device_type_id
    FROM   v_unit_relation_network AS urn
    WHERE  urn.unit_uuid = $entityUUID.uuid
),
best_rank AS (
    SELECT MIN(scope_rank) AS scope_rank
    FROM   candidate
    WHERE  scope_rank IS NOT NULL
),
selected_ingress AS (
    SELECT c.address_value,
           c.device_uuid,
           c.origin_id,
           c.is_secondary,
           c.device_type_id
    FROM   candidate AS c
    JOIN   best_rank USING (scope_rank)
),
unit_address AS (
    SELECT urn.address_value,
           si.address_value AS ingress_address_value,
           si.origin_id AS ingress_origin_id,
           si.is_secondary AS ingress_is_secondary,
           si.device_type_id AS ingress_device_type_id
    FROM   v_unit_relation_network AS urn
    LEFT JOIN selected_ingress AS si
           ON si.device_uuid = urn.device_uuid
          AND si.address_value = urn.address_value
    WHERE  urn.unit_uuid = $entityUUID.uuid
)
SELECT &unitIngressAddressRow.*
FROM   unit_address
ORDER BY CASE WHEN ingress_address_value IS NULL THEN 1 ELSE 0 END,
         ingress_device_type_id,
         ingress_is_secondary, /* primary (0) before secondary (1) */
         ingress_origin_id DESC, /* provider (1) before machine (0) */
         address_value
`, unitIngressAddressRow{}, entityUUID{})

	if err != nil {
		return "", errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, ident).GetAll(&rows)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf("querying unit network: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", errors.Capture(err)
	}
	if len(rows) == 0 || !rows[0].IngressAddress.Valid {
		return "", nil
	}

	return rows[0].IngressAddress.String, nil
}

// GetUnitRelationsIngressAddresses retrieves an ingress address for all the
// relations where the given unit is in scope.
func (st *State) GetUnitRelationsIngressAddresses(
	ctx context.Context,
	unitUUID coreunit.UUID,
) (map[relation.UUID]string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}

	type ingressAddressRelationRow struct {
		RelationUUID   string         `db:"relation_uuid"`
		IngressAddress sql.NullString `db:"ingress_address_value"`
	}

	ident := entityUUID{UUID: unitUUID.String()}
	var rows []ingressAddressRelationRow
	stmt, err := st.Prepare(`
WITH endpoint_binding AS (
    SELECT ae.application_uuid,
           cr.name,
           ae.space_uuid,
           re.relation_uuid
    FROM   application_endpoint AS ae
    JOIN   charm_relation AS cr ON ae.charm_relation_uuid = cr.uuid
    JOIN   relation_endpoint AS re ON ae.uuid = re.endpoint_uuid
),
endpoint_space AS (
    SELECT eb.name AS endpoint_name,
           eb.relation_uuid,
           IFNULL(eb.space_uuid, a.space_uuid) AS space_uuid
    FROM   endpoint_binding AS eb
    JOIN   unit AS u ON eb.application_uuid = u.application_uuid
    JOIN   application AS a ON u.application_uuid = a.uuid
    WHERE  u.uuid = $entityUUID.uuid
),
ingress_candidate AS (
    SELECT es.endpoint_name,
           es.relation_uuid,
           urn.address_value,
           urn.device_uuid,
           urn.scope_rank,
           urn.origin_id,
           urn.is_secondary,
           urn.device_type_id
    FROM   endpoint_space AS es
    JOIN   v_unit_relation_network AS urn
           ON urn.unit_uuid = $entityUUID.uuid
          AND urn.space_uuid = es.space_uuid
),
best_rank AS (
    SELECT endpoint_name, MIN(scope_rank) AS scope_rank
    FROM   ingress_candidate
    WHERE  scope_rank IS NOT NULL
    GROUP BY endpoint_name
),
selected_ingress AS (
    SELECT ic.endpoint_name,
           ic.relation_uuid,
           ic.address_value,
           ic.device_uuid,
           ic.origin_id,
           ic.is_secondary,
           ic.device_type_id
    FROM   ingress_candidate AS ic
    JOIN   best_rank USING (endpoint_name, scope_rank)
),
lld AS (
    SELECT lld.uuid,
           lld.name
    FROM   link_layer_device AS lld
    JOIN   link_layer_device_type AS lldt ON lld.device_type_id = lldt.id
),
endpoint_address AS (
    SELECT es.endpoint_name,
           es.relation_uuid,
           urn.address_value,
           lld.name,
           si.address_value AS ingress_address_value,
           si.origin_id AS ingress_origin_id,
           si.is_secondary AS ingress_is_secondary,
           si.device_type_id AS ingress_device_type_id
    FROM   endpoint_space AS es
    JOIN   v_unit_relation_network AS urn
           ON urn.unit_uuid = $entityUUID.uuid
          AND urn.space_uuid = es.space_uuid
    JOIN   lld ON urn.device_uuid = lld.uuid
    LEFT JOIN selected_ingress AS si
           ON si.endpoint_name = es.endpoint_name
          AND si.device_uuid = urn.device_uuid
          AND si.address_value = urn.address_value
)
SELECT &ingressAddressRelationRow.*
FROM   endpoint_address
ORDER BY endpoint_name,
         CASE WHEN ingress_address_value IS NULL THEN 1 ELSE 0 END,
         ingress_device_type_id,
         ingress_is_secondary, /* primary (0) before secondary (1) */
         ingress_origin_id DESC, /* provider (1) before machine (0) */
         address_value,
         name
`, ingressAddressRelationRow{}, entityUUID{})
	if err != nil {
		return nil, errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, ident).GetAll(&rows)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf(
				"querying endpoint networks: %w", err,
			)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	// Only keep the first valid ingress address per relation UUID.
	addrByRelationUUIDs := make(map[relation.UUID]string)
	for _, row := range rows {
		if !row.IngressAddress.Valid {
			continue
		}
		_, ok := addrByRelationUUIDs[relation.UUID(row.RelationUUID)]
		if !ok {
			addrByRelationUUIDs[relation.UUID(row.RelationUUID)] = row.IngressAddress.String
		}
	}

	return addrByRelationUUIDs, nil
}

// GetModelEgressSubnets retrieves the egress-subnets configuration from model
// config.
func (st *State) GetModelEgressSubnets(ctx context.Context) ([]string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}

	type modelConfigEntry struct {
		Key   string `db:"key"`
		Value string `db:"value"`
	}

	egressSubnetsConfig := modelConfigEntry{Key: unitstate.EgressSubnetsKey}
	stmt, err := st.Prepare(`
SELECT &modelConfigEntry.value
FROM   model_config
WHERE  key = $modelConfigEntry.key
`, modelConfigEntry{})
	if err != nil {
		return nil, errors.Errorf("preparing model egress subnets statement: %w", err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, egressSubnetsConfig).Get(&egressSubnetsConfig)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		return errors.Capture(err)
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	if egressSubnetsConfig.Value == "" {
		return nil, nil
	}

	cidrs := strings.Split(egressSubnetsConfig.Value, ",")
	result := make([]string, 0, len(cidrs))
	for _, cidr := range cidrs {
		trimmed := strings.TrimSpace(cidr)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result, nil
}

// GetUnitPublicAddressForEgress retrieves the best unit address to use when
// deriving fallback egress subnets.
func (st *State) GetUnitPublicAddressForEgress(
	ctx context.Context,
	unitUUID coreunit.UUID,
) (string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return "", errors.Capture(err)
	}

	type addressValue struct {
		Value string `db:"address_value"`
	}

	ident := entityUUID{UUID: unitUUID.String()}
	stmt, err := st.Prepare(`
WITH unit_net_node AS (
    SELECT
        s.net_node_uuid
    FROM unit AS u
    JOIN application AS a ON u.application_uuid = a.uuid
    JOIN k8s_service AS s ON a.uuid = s.application_uuid
    WHERE u.uuid = $entityUUID.uuid
    UNION ALL
    SELECT
        net_node_uuid
    FROM unit
    WHERE uuid = $entityUUID.uuid
)
SELECT address_value AS &addressValue.address_value
FROM unit_net_node AS unn
JOIN ip_address AS ipa ON ipa.net_node_uuid = unn.net_node_uuid
WHERE ipa.scope_id = 1 /* public */
ORDER BY ipa.type_id, /* ipv4 (0) before ipv6 (1) */
         ipa.is_secondary, /* primary (0) before secondary (1) */
         ipa.origin_id DESC, /* provider (1) before machine (0) */
         address_value
LIMIT 1
`, addressValue{}, entityUUID{})
	if err != nil {
		return "", errors.Errorf(
			"preparing select unit public address statement: %w", err,
		)
	}

	var address addressValue
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, ident).Get(&address)
		if errors.Is(err, sqlair.ErrNoRows) {
			return nil
		}
		if err != nil {
			return errors.Errorf("querying unit public address: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", errors.Capture(err)
	}

	return address.Value, nil
}

// GetRelationsEgressSubnetsByUnitUUID retrieves the egress subnets for all
// relations the specific unit is in scope for, grouped by relation UUID.
func (st *State) GetRelationsEgressSubnetsByUnitUUID(
	ctx context.Context, unitUUID coreunit.UUID,
) (map[relation.UUID][]string, error) {
	db, err := st.DB(ctx)
	if err != nil {
		return nil, errors.Capture(err)
	}

	ident := entityUUID{UUID: unitUUID.String()}
	stmt, err := st.Prepare(`
SELECT DISTINCT rne.* AS &egressCIDRAndRelationUUID.*
FROM   relation_network_egress AS rne
JOIN   relation AS r ON rne.relation_uuid = r.uuid
JOIN   relation_endpoint AS re ON r.uuid = re.relation_uuid
JOIN   relation_unit AS ru ON re.uuid = ru.relation_endpoint_uuid
WHERE  ru.unit_uuid = $entityUUID.uuid
`, egressCIDRAndRelationUUID{}, entityUUID{})
	if err != nil {
		return nil, errors.Errorf("preparing select unit's relations egress subnets statement: %w", err)
	}

	var cidrs []egressCIDRAndRelationUUID
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, ident).GetAll(&cidrs)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf("querying unit's relations egress subnets: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	if len(cidrs) == 0 {
		return nil, nil
	}

	result := make(map[relation.UUID][]string, 0)
	for _, c := range cidrs {
		result[relation.UUID(c.RelationUUID)] = append(result[relation.UUID(c.RelationUUID)], c.CIDR)
	}

	return result, nil
}
