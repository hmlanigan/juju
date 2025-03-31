// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package errors

import "github.com/juju/juju/internal/errors"

const (
	// ApplicationIDNotValid describes an error when the application ID is
	// not valid.
	ApplicationIDNotValid = errors.ConstError("application ID not valid")

	// ApplicationNotFound describes an error that occurs when the application
	// being operated on does not exist.
	ApplicationNotFound = errors.ConstError("application not found")

	// ApplicationNotFoundForRelation indicates that the application is not part of
	// the relation.
	ApplicationNotFoundForRelation = errors.ConstError("application not found for relation")

	// ApplicationNotSubordinate the application is either not a subordinate
	// or is not found.
	ApplicationNotSubordinate = errors.ConstError("application not a subordinate")

	// PotentialRelationUnitNotValid describes an error that occurs during
	// EnterScope pre-checks to ensure the created relation unit will be valid.
	//
	// This is not a fatal error. It replaces a boolean return value from a
	// prior implementation.
	PotentialRelationUnitNotValid = errors.ConstError("potential relation unit not valid")

	// RelationEndpointNotFound describes an error that occurs when the specified
	// relation endpoint does not exist.
	RelationEndpointNotFound = errors.ConstError("relation endpoint not found")

	// RelationNotFound describes an error that occurs when the specified
	// relation does not exist.
	RelationNotFound = errors.ConstError("relation not found")

	// RelationUUIDNotValid describes an error when the relation RelationUUID is
	// not valid.
	RelationUUIDNotValid = errors.ConstError("relation RelationUUID not valid")

	// RelationKeyNotValid describes an error when the relation key is
	// not valid.
	RelationKeyNotValid = errors.ConstError("relation key not valid")

	// UnitUUIDNotValid describes an error when the unit RelationUUID is
	// not valid.
	UnitUUIDNotValid = errors.ConstError("unit RelationUUID not valid")
)
