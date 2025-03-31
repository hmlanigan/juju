// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"fmt"

	"github.com/canonical/sqlair"
	"github.com/juju/clock"
	"github.com/juju/collections/transform"

	"github.com/juju/juju/core/application"
	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/logger"
	corerelation "github.com/juju/juju/core/relation"
	corestatus "github.com/juju/juju/core/status"
	"github.com/juju/juju/core/unit"
	"github.com/juju/juju/core/watcher/eventsource"
	"github.com/juju/juju/domain"
	"github.com/juju/juju/domain/relation"
	relationerrors "github.com/juju/juju/domain/relation/errors"
	"github.com/juju/juju/internal/charm"
	"github.com/juju/juju/internal/errors"
)

type State struct {
	*domain.StateBase
	clock  clock.Clock
	logger logger.Logger
}

// NewState returns a new state reference.
func NewState(factory database.TxnRunnerFactory, clock clock.Clock, logger logger.Logger) *State {
	return &State{
		StateBase: domain.NewStateBase(factory),
		clock:     clock,
		logger:    logger,
	}
}

// GetPrincipalApplicationID return the principal application's ID given a
// subordinate application's ID.
//
// The following error types can be expected to be returned:
//   - [relationerrors.ApplicationNotSubordinate] is returned if the application
//     RelationUUID is not found as a subordinate in the application_subordinate link table.
//   - [relationerrors.ApplicationNotFound] is returned if the application
//     RelationUUID is not found.
func (st *State) GetPrincipalApplicationID(ctx context.Context, id application.ID) (application.ID, error) {
	db, err := st.DB()
	if err != nil {
		return "", errors.Capture(err)
	}

	input := matchPrincipalSubordinateApps{
		Subordinate: id,
	}
	stmt, err := st.Prepare(`
SELECT principal_uuid AS &matchPrincipalSubordinateApps.principal_uuid
FROM   v_principal_subordinate AS ps
WHERE  ps.subordinate_uuid = $matchPrincipalSubordinateApps.subordinate_uuid
`, input)
	if err != nil {
		return "", errors.Capture(err)
	}

	output := matchPrincipalSubordinateApps{}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		exists, err := checkExistsByUUID(ctx, st, tx, "application", id.String())
		if err != nil {
			return errors.Errorf("checking application exists: %w", err)
		} else if !exists {
			return relationerrors.ApplicationNotFound
		}
		err = tx.Query(ctx, stmt, input).Get(&output)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.ApplicationNotSubordinate
		}
		return err
	})
	if err != nil {
		return "", errors.Capture(err)
	}

	return output.Principal, nil
}

// GetOtherRelatedEndpointApplicationData returns an OtherApplicationForWatcher struct
// for each Endpoint in a relation with the given application ID.
//
// The following error types can be expected to be returned:
//   - [relationerrors.ApplicationNotFound] is returned if application ID
//     is not used in any relations or if the other relation applications
//     are not found.
func (st *State) GetOtherRelatedEndpointApplicationData(
	ctx context.Context,
	applicationID application.ID,
) ([]relation.OtherApplicationForWatcher, error) {
	db, err := st.DB()
	if err != nil {
		return nil, errors.Capture(err)
	}

	otherApps := []otherApplicationsForWatcher{}
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		// Find all applications in a relation with the given application.
		getOtherApps, err := st.getOtherApplicationInRelations(ctx, tx, applicationID)
		if err != nil {
			return errors.Capture(err)
		}

		// For all applications, determine if it is a subordinate.
		otherApps, err = st.getApplicationSubordinate(ctx, tx, getOtherApps)
		if err != nil {
			return errors.Capture(err)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	results := make([]relation.OtherApplicationForWatcher, len(otherApps))
	for i, oneApp := range otherApps {
		results[i] = relation.OtherApplicationForWatcher{
			ApplicationID: oneApp.AppID,
			Subordinate:   oneApp.Subordinate,
		}
	}

	return results, nil
}

// getOtherApplicationInRelations returns a slice of applications ID
// in the given relations which are not the given application ID.
func (st *State) getOtherApplicationInRelations(
	ctx context.Context,
	tx *sqlair.TX,
	appID application.ID,
) ([]applicationID, error) {

	findOtherEndsStmt, err := st.Prepare(`
SELECT other_ae.application_uuid AS &applicationID.uuid
FROM   application_endpoint AS other_ae
JOIN   relation_endpoint AS other_re ON other_re.endpoint_uuid = other_ae.uuid
JOIN   relation_endpoint AS re ON re.relation_uuid = other_re.relation_uuid
JOIN   application_endpoint AS ae ON ae.uuid = re.endpoint_uuid
WHERE  ae.application_uuid = $applicationID.uuid
AND    other_ae.application_uuid != $applicationID.uuid
`, applicationID{})
	if err != nil {
		return nil, errors.Errorf("preparing other endpoint query: %w", err)
	}

	app := applicationID{ID: appID}
	otherApps := []applicationID{}

	err = tx.Query(ctx, findOtherEndsStmt, app).GetAll(&otherApps)
	if err != nil {
		return nil, errors.Capture(err)
	}
	return otherApps, nil
}

// getApplicationSubordinate returns a otherApplicationsForWatcher structure
// for each given application ID.
func (st *State) getApplicationSubordinate(
	ctx context.Context,
	tx *sqlair.TX,
	apps []applicationID,
) ([]otherApplicationsForWatcher, error) {

	appSubordinateStmt, err := st.Prepare(`
SELECT application_uuid AS &otherApplicationsForWatcher.application_uuid,
       subordinate AS &otherApplicationsForWatcher.subordinate
FROM   v_application_subordinate
WHERE  application_uuid = ($uuids[:])
`, otherApplicationsForWatcher{}, uuids{})
	if err != nil {
		return nil, errors.Errorf("preparing other application query: %w", err)
	}

	appIDs := uuids{}
	for _, uuid := range apps {
		appIDs = append(appIDs, uuid.ID.String())
	}
	otherApps := []otherApplicationsForWatcher{}
	err = tx.Query(ctx, appSubordinateStmt, appIDs).GetAll(&otherApps)
	if err != nil {
		return nil, errors.Capture(err)
	}

	return otherApps, nil
}

// GetRelationID returns the relation ID for the given relation RelationUUID.
//
// The following error types can be expected to be returned:
//   - [relationerrors.RelationNotFound] is returned if the relation RelationUUID
//     is not found.
func (st *State) GetRelationID(ctx context.Context, relationUUID corerelation.UUID) (int, error) {
	db, err := st.DB()
	if err != nil {
		return 0, errors.Capture(err)
	}

	id := relationIDAndUUID{
		UUID: relationUUID.String(),
	}
	stmt, err := st.Prepare(`
SELECT &relationIDAndUUID.relation_id
FROM   relation
WHERE  uuid = $relationIDAndUUID.uuid
`, id)
	if err != nil {
		return 0, errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, id).Get(&id)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.RelationNotFound
		}
		return err
	})
	if err != nil {
		return 0, errors.Capture(err)
	}

	return id.ID, nil
}

// GetRelationUUIDByID returns the relation RelationUUID based on the relation ID.
//
// The following error types can be expected to be returned:
//   - [relationerrors.RelationNotFound] is returned if the relation RelationUUID
//     relating to the relation ID cannot be found.
func (st *State) GetRelationUUIDByID(ctx context.Context, relationID int) (corerelation.UUID, error) {
	db, err := st.DB()
	if err != nil {
		return "", errors.Capture(err)
	}

	id := relationIDAndUUID{
		ID: relationID,
	}
	stmt, err := st.Prepare(`
SELECT &relationIDAndUUID.uuid
FROM   relation
WHERE  relation_id = $relationIDAndUUID.relation_id
`, id)
	if err != nil {
		return "", errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, id).Get(&id)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.RelationNotFound
		}
		return err
	})
	if err != nil {
		return "", errors.Capture(err)
	}

	return corerelation.UUID(id.UUID), nil
}

// GetRelationEndpointScope returns the scope of the relation endpoint
// at the intersection of the relationUUID and applicationID.
//
// The following error types can be expected to be returned:
//   - [relationerrors.RelationNotFound] is returned if the relation RelationUUID
//     relating to the relation ID cannot be found.
func (st *State) GetRelationEndpointScope(
	ctx context.Context,
	relUUID corerelation.UUID,
	appID application.ID,
) (charm.RelationScope, error) {
	db, err := st.DB()
	if err != nil {
		return "", errors.Capture(err)
	}

	type scope struct {
		Scope charm.RelationScope `db:"name"`
	}

	stmt, err := st.Prepare(`
SELECT  scope AS &scope.name
FROM    v_relation_endpoint
WHERE   relation_uuid = $relationUUID.uuid
AND     application_uuid = $applicationID.uuid
`, scope{}, applicationID{}, relationUUID{})
	if err != nil {
		return "", errors.Capture(err)
	}

	rel := relationUUID{
		UUID: relUUID.String(),
	}
	app := applicationID{ID: appID}
	var output scope
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		// Check if the relation exists.
		relationFound, err := checkExistsByUUID(ctx, st, tx, "relation", relUUID.String())
		if err != nil {
			return errors.Capture(err)
		} else if !relationFound {
			return relationerrors.RelationNotFound
		}

		return tx.Query(ctx, stmt, rel, app).Get(&output)
	})
	if err != nil {
		return "", errors.Capture(err)
	}

	return output.Scope, nil
}

// GetRelationEndpointUUID retrieves the endpoint RelationUUID of a given relation
// for a specific application.
// It queries the database using the provided application ID and relation RelationUUID
// arguments.
//
// The following error types can be expected to be returned:
//   - [relationerrors.ApplicationNotFound] is returned if the application
//     is not found.
//   - [relationerrors.RelationEndpointNotFound] is returned if the relation
//     Endpoint is not found.
//   - [relationerrors.RelationNotFound] is returned if the relation RelationUUID
//     is not found.
func (st *State) GetRelationEndpointUUID(ctx context.Context, args relation.GetRelationEndpointUUIDArgs) (
	corerelation.EndpointUUID, error) {
	db, err := st.DB()
	if err != nil {
		return "", errors.Capture(err)
	}

	type relationEndpointUUID struct {
		UUID string `db:"uuid"`
	}
	type relationEndpointArgs struct {
		AppID        string `db:"application_uuid"`
		RelationUUID string `db:"relation_uuid"`
	}
	dbArgs := relationEndpointArgs{
		AppID:        args.ApplicationID.String(),
		RelationUUID: args.RelationUUID.String(),
	}
	stmt, err := st.Prepare(`
SELECT re.uuid AS &relationEndpointUUID.uuid
FROM   relation_endpoint re
JOIN   application_endpoint ae ON re.endpoint_uuid = ae.uuid
WHERE  ae.application_uuid = $relationEndpointArgs.application_uuid
AND    re.relation_uuid = $relationEndpointArgs.relation_uuid
`, relationEndpointUUID{}, dbArgs)
	if err != nil {
		return "", errors.Capture(err)
	}
	var relationEndpoint relationEndpointUUID
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt, dbArgs).Get(&relationEndpoint)
		if errors.Is(err, sqlair.ErrNoRows) {
			// Check if it is a missing application.
			appFound, err := checkExistsByUUID(ctx, st, tx, "application", args.ApplicationID.String())
			if err != nil {
				return errors.Capture(err)
			}
			// Check if the relation exists.
			relationFound, err := checkExistsByUUID(ctx, st, tx, "relation", args.RelationUUID.String())
			if err != nil {
				return errors.Capture(err)
			}
			var errs []error
			if !appFound {
				errs = append(errs, errors.Errorf("%w: %s", relationerrors.ApplicationNotFound, args.ApplicationID))
			}
			if !relationFound {
				errs = append(errs, errors.Errorf("%w: %s", relationerrors.RelationNotFound, args.RelationUUID))
			}
			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			return errors.Errorf("relationUUID %q with applicationID %q: %w",
				args.RelationUUID, args.ApplicationID, relationerrors.RelationEndpointNotFound)

		}
		return errors.Capture(err)
	})

	return corerelation.EndpointUUID(relationEndpoint.UUID), errors.Capture(err)
}

// GetRelationsStatusForUnit returns RelationUnitStatus for any relation the
// unit is part of.
func (st *State) GetRelationsStatusForUnit(
	ctx context.Context,
	unitUUID unit.UUID,
) ([]relation.RelationUnitStatusResult, error) {
	db, err := st.DB()
	if err != nil {
		return nil, errors.Capture(err)
	}

	type relationUnitStatus struct {
		RelationUUID string `db:"relation_uuid"`
		InScope      bool   `db:"in_scope"`
		Status       string `db:"status"`
	}
	type unitUUIDArg struct {
		UUID string `db:"unit_uuid"`
	}

	uuid := unitUUIDArg{
		UUID: unitUUID.String(),
	}

	stmt, err := st.Prepare(`
SELECT (ru.relation_uuid, ru.in_scope, vrs.status) AS (&relationUnitStatus.*)
FROM   relation_unit ru
JOIN   v_relation_status vrs ON ru.relation_uuid = vrs.relation_uuid
WHERE  ru.unit_uuid = $unitUUIDArg.unit_uuid
`, uuid, relationUnitStatus{})
	if err != nil {
		return nil, errors.Capture(err)
	}

	var relationUnitStatuses []relation.RelationUnitStatusResult
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		var statuses []relationUnitStatus
		err := tx.Query(ctx, stmt, uuid).GetAll(&statuses)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Capture(err)
		}

		for _, status := range statuses {
			endpoints, err := st.getEndpoints(ctx, tx, corerelation.UUID(status.RelationUUID))
			if err != nil {
				return errors.Errorf("getting endpoints of relation %q: %w", status.RelationUUID, err)
			}

			relationUnitStatuses = append(relationUnitStatuses, relation.RelationUnitStatusResult{
				Endpoints: endpoints,
				InScope:   status.InScope,
				Suspended: status.Status == corestatus.Suspended.String(),
			})
		}

		return nil
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	return relationUnitStatuses, nil
}

// GetRelationEndpoints retrieves the endpoints of a given relation specified via its RelationUUID.
//
// The following error types can be expected to be returned:
//   - [relationerrors.RelationNotFound] is returned if the relation RelationUUID is not
//     found.
func (st *State) GetRelationEndpoints(ctx context.Context, uuid corerelation.UUID) ([]relation.Endpoint, error) {
	db, err := st.DB()
	if err != nil {
		return nil, errors.Capture(err)
	}

	var endpoints []relation.Endpoint
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		endpoints, err = st.getEndpoints(ctx, tx, uuid)
		return err
	})
	if err != nil {
		return nil, errors.Capture(err)
	}

	return endpoints, nil
}

// getEndpoints retrieves the endpoints of the specified relation.
func (st *State) getEndpoints(
	ctx context.Context,
	tx *sqlair.TX,
	uuid corerelation.UUID,
) ([]relation.Endpoint, error) {
	id := relationUUID{
		UUID: uuid.String(),
	}
	stmt, err := st.Prepare(`
SELECT &endpoint.*
FROM   v_relation_endpoint
WHERE  relation_uuid = $relationUUID.uuid
`, id, endpoint{})
	if err != nil {
		return nil, errors.Capture(err)
	}

	var endpoints []endpoint
	err = tx.Query(ctx, stmt, id).GetAll(&endpoints)
	if errors.Is(err, sqlair.ErrNoRows) {
		return nil, relationerrors.RelationNotFound
	} else if err != nil {
		return nil, errors.Capture(err)
	}

	if length := len(endpoints); length > 2 {
		return nil, errors.Errorf("internal error: expected 1 or 2 endpoints in relation, got %d", length)
	}

	var relationEndpoints []relation.Endpoint
	for _, e := range endpoints {
		relationEndpoints = append(relationEndpoints, e.toRelationEndpoint())
	}

	return relationEndpoints, nil
}

// GetRegularRelationUUIDByEndpointIdentifiers gets the RelationUUID of a regular
// relation specified by two endpoint identifiers.
//
// The following error types can be expected to be returned:
//   - [relationerrors.RelationNotFound] is returned if endpoints cannot be
//     found.
func (st *State) GetRegularRelationUUIDByEndpointIdentifiers(
	ctx context.Context,
	endpoint1, endpoint2 relation.EndpointIdentifier,
) (corerelation.UUID, error) {
	db, err := st.DB()
	if err != nil {
		return "", errors.Capture(err)
	}

	type endpointIdentifier1 endpointIdentifier
	type endpointIdentifier2 endpointIdentifier
	e1 := endpointIdentifier1{
		ApplicationName: endpoint1.ApplicationName,
		EndpointName:    endpoint1.EndpointName,
	}
	e2 := endpointIdentifier2{
		ApplicationName: endpoint2.ApplicationName,
		EndpointName:    endpoint2.EndpointName,
	}

	stmt, err := st.Prepare(`
SELECT &relationUUID.*
FROM   relation r
JOIN   v_relation_endpoint_identifier e1 ON r.uuid = e1.relation_uuid
JOIN   v_relation_endpoint_identifier e2 ON r.uuid = e2.relation_uuid
WHERE  e1.application_name = $endpointIdentifier1.application_name 
AND    e1.endpoint_name    = $endpointIdentifier1.endpoint_name
AND    e2.application_name = $endpointIdentifier2.application_name 
AND    e2.endpoint_name    = $endpointIdentifier2.endpoint_name
`, relationUUID{}, e1, e2)
	if err != nil {
		return "", errors.Capture(err)
	}

	var uuid []relationUUID
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err = tx.Query(ctx, stmt, e1, e2).GetAll(&uuid)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.RelationNotFound
		}
		return errors.Capture(err)
	})
	if err != nil {
		return "", errors.Capture(err)
	}

	if len(uuid) > 1 {
		return "", errors.Errorf("found multiple relations for endpoint pair")
	}

	return corerelation.UUID(uuid[0].UUID), nil
}

// GetPeerRelationUUIDByEndpointIdentifiers gets the RelationUUID of a peer
// relation specified by a single endpoint identifier.
//
// The following error types can be expected to be returned:
//   - [relationerrors.RelationNotFound] is returned if endpoint cannot be
//     found.
func (st *State) GetPeerRelationUUIDByEndpointIdentifiers(
	ctx context.Context,
	endpoint relation.EndpointIdentifier,
) (corerelation.UUID, error) {
	db, err := st.DB()
	if err != nil {
		return "", errors.Capture(err)
	}

	e := endpointIdentifier{
		ApplicationName: endpoint.ApplicationName,
		EndpointName:    endpoint.EndpointName,
	}

	stmt, err := st.Prepare(`
SELECT &relationUUIDAndRole.*
FROM   relation r
JOIN   v_relation_endpoint e ON r.uuid = e.relation_uuid
WHERE  e.application_name = $endpointIdentifier.application_name 
AND    e.endpoint_name    = $endpointIdentifier.endpoint_name
`, relationUUIDAndRole{}, e)
	if err != nil {
		return "", errors.Capture(err)
	}

	var uuidAndRole []relationUUIDAndRole
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err = tx.Query(ctx, stmt, e).GetAll(&uuidAndRole)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.RelationNotFound
		}
		return errors.Capture(err)
	})
	if err != nil {
		return "", errors.Capture(err)
	}

	if len(uuidAndRole) > 1 {
		return "", errors.Errorf("found multiple relations for peer application endpoint combination")
	}

	// Verify that the role is peer. Endpoint names are unique per charm, so if
	// the role is not peer the application does not have a peer relation with
	// the specified endpoint name, so return RelationNotFound.
	if uuidAndRole[0].Role != string(charm.RolePeer) {
		return "", relationerrors.RelationNotFound
	}

	return corerelation.UUID(uuidAndRole[0].UUID), nil
}

// GetRelationDetails returns relation details for the given relationID.
//
// The following error types can be expected to be returned:
//   - [relationerrors.RelationNotFound] is returned if the relation RelationUUID
//     is not found.
func (st *State) GetRelationDetails(ctx context.Context, relationID int) (relation.RelationDetailsResult, error) {
	db, err := st.DB()
	if err != nil {
		return relation.RelationDetailsResult{}, errors.Capture(err)
	}

	type getRelation struct {
		UUID corerelation.UUID `db:"uuid"`
		ID   int               `db:"relation_id"`
		Life life.Value        `db:"value"`
	}
	rel := getRelation{
		ID: relationID,
	}
	stmt, err := st.Prepare(`
SELECT (r.uuid, r.relation_id, l.value) AS (&getRelation.*)
FROM   relation r
JOIN   life l ON r.life_id = l.id
WHERE  relation_id = $getRelation.relation_id
`, rel)
	if err != nil {
		return relation.RelationDetailsResult{}, errors.Capture(err)
	}

	var endpoints []relation.Endpoint
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err = tx.Query(ctx, stmt, rel).Get(&rel)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.RelationNotFound
		} else if err != nil {
			return errors.Capture(err)
		}

		endpoints, err = st.getEndpoints(ctx, tx, rel.UUID)
		if err != nil {
			return errors.Errorf("getting relation endpoints: %w", err)
		}
		return errors.Capture(err)
	})
	if err != nil {
		return relation.RelationDetailsResult{}, errors.Capture(err)
	}

	return relation.RelationDetailsResult{
		Life:      rel.Life,
		UUID:      rel.UUID,
		ID:        rel.ID,
		Endpoints: endpoints,
	}, nil
}

// InitialWatchLifeSuspendedStatus returns the two tables to watch for
// a relation's Life and Suspended status when the relation contains
// the provided application and the initial namespace query.
func (st *State) InitialWatchLifeSuspendedStatus(id application.ID) (string, string, eventsource.NamespaceQuery) {
	queryFunc := func(ctx context.Context, runner database.TxnRunner) ([]string, error) {
		stmt, err := st.Prepare(`
SELECT  re.relation_uuid AS &relationUUID.uuid
FROM    relation_endpoint re
JOIN    application_endpoint ae ON ae.uuid = re.endpoint_uuid
WHERE   ae.application_uuid = $applicationID.ID
`, applicationID{})
		if err != nil {
			return nil, errors.Capture(err)
		}

		appID := applicationID{ID: id}

		var results []relationUUID
		err = runner.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
			err := tx.Query(ctx, stmt, appID).GetAll(&results)
			if errors.Is(err, sqlair.ErrNoRows) {
				return nil
			}
			return errors.Capture(err)
		})
		if err != nil {
			return nil, errors.Errorf("querying requested applications that have pending charms: %w", err)
		}

		return transform.Slice(results, func(r relationUUID) string { return r.UUID }), nil
	}

	return "relation", "relation_status", queryFunc
}

// WatcherApplicationSettingsNamespace returns the namespace string used for
// tracking application settings in the database.
func (st *State) WatcherApplicationSettingsNamespace() string {
	return "relation_application_setting"
}

// WatchLifeSuspendedStatusMapperData returns data needed to evaluate a relation
// uuid as part of WatchLifeSuspendedStatus eventmapper.
//
// The following error types can be expected to be returned:
//   - [relationerrors.ApplicationNotFoundForRelation] is returned if the
//     application is not part of the relation.
//   - [relationerrors.RelationNotFound] is returned if the relation RelationUUID
//     is not found.
func (st *State) WatchLifeSuspendedStatusMapperData(
	ctx context.Context,
	relUUID corerelation.UUID,
	appID application.ID,
) (relation.RelationLifeSuspendedData, error) {
	db, err := st.DB()
	if err != nil {
		return relation.RelationLifeSuspendedData{}, errors.Capture(err)
	}

	data := watcherMapperData{
		RelationUUID: relUUID.String(),
		AppUUID:      appID.String(),
	}

	relAppStmt, err := st.Prepare(`
SELECT  re.relation_uuid AS &watcherMapperData.uuid
FROM    relation_endpoint re
JOIN    application_endpoint ae ON ae.uuid = re.endpoint_uuid
WHERE   ae.application_uuid = $watcherMapperData.application_uuid
AND     re.relation_uuid = $watcherMapperData.uuid
`, watcherMapperData{})
	if err != nil {
		return relation.RelationLifeSuspendedData{}, errors.Capture(err)
	}

	lifeStatusStmt, err := st.Prepare(`
SELECT (rst.name, l.value) AS (&watcherMapperData.*)
FROM   relation r
JOIN   life l ON r.life_id = l.id
JOIN   relation_status rs ON rs.relation_uuid = r.uuid
JOIN   relation_status_type rst ON rst.id = rs.relation_status_type_id
WHERE  r.uuid = $watcherMapperData.uuid
`, watcherMapperData{})
	if err != nil {
		return relation.RelationLifeSuspendedData{}, errors.Capture(err)
	}

	var endpoints []relation.Endpoint
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err = tx.Query(ctx, relAppStmt, data).Get(&data)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.ApplicationNotFoundForRelation
		} else if err != nil {
			return errors.Errorf("verifying relation application intersection: %w", err)
		}

		err = tx.Query(ctx, lifeStatusStmt, data).Get(&data)
		if errors.Is(err, sqlair.ErrNoRows) {
			return relationerrors.RelationNotFound
		} else if err != nil {
			return errors.Errorf("getting relation life and status: %w", err)
		}

		endpoints, err = st.getEndpoints(ctx, tx, corerelation.UUID(data.RelationUUID))
		if err != nil {
			return errors.Errorf("getting relation endpoints: %w", err)
		}
		return errors.Capture(err)
	})
	if err != nil {
		return relation.RelationLifeSuspendedData{}, errors.Capture(err)
	}

	return relation.RelationLifeSuspendedData{
		Life:      life.Value(data.Life),
		Suspended: data.Suspended == corestatus.Suspended.String(),
		Endpoints: endpoints,
	}, nil
}

// checkExistsByUUID checks if a record with the specified RelationUUID exists in the given
// table using a transaction and context.
func checkExistsByUUID(ctx context.Context, st *State, tx *sqlair.TX, table string, uuid string) (bool,
	error) {
	type search struct {
		UUID string `db:"uuid"`
	}

	searched := search{UUID: uuid}
	query := fmt.Sprintf(`
SELECT &search.* 
FROM   %s 
WHERE  uuid = $search.uuid
`, table)
	checkStmt, err := st.Prepare(query, searched)
	if err != nil {
		return false, errors.Capture(err)
	}

	err = tx.Query(ctx, checkStmt, searched).Get(&searched)
	if errors.Is(err, sqlair.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, errors.Errorf("query %q: %w", query, err)
	}
	return true, nil
}
