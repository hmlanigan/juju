// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"maps"
	"sort"
	"strings"

	"github.com/juju/collections/transform"

	corenetwork "github.com/juju/juju/core/network"
	"github.com/juju/juju/core/relation"
	"github.com/juju/juju/core/trace"
	"github.com/juju/juju/core/unit"
	"github.com/juju/juju/domain/unitstate"
	"github.com/juju/juju/domain/unitstate/internal"
	"github.com/juju/juju/internal/errors"
)

// CommitHookChanges persists a set of changes after a hook successfully
// completes and executes them in a single transaction.
func (s *LeadershipService) CommitHookChanges(ctx context.Context, arg unitstate.CommitHookChangesArg) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	hasChanges, err := arg.ValidateAndHasChanges()
	if err != nil {
		return errors.Capture(err)
	}
	if !hasChanges {
		return nil
	}

	unitUUID, err := s.st.GetUnitUUIDByName(ctx, arg.UnitName)
	if err != nil {
		return errors.Capture(err)
	}

	relationSettings, err := s.transformRelationSettings(ctx, arg.RelationSettings)
	if err != nil {
		return errors.Capture(err)
	}

	newArgs := internal.TransformCommitHookChangesArg(arg, unitUUID)
	newArgs.RelationSettings = relationSettings

	if arg.UpdateNetworkInfo {
		newArgs.RelationSettings, err = s.mergeRelationSettingsAndNetworkInfo(ctx, unitUUID, newArgs.RelationSettings)
		if err != nil {
			return errors.Capture(err)
		}
	}

	withCaveat, err := s.getManagementCaveat(arg)
	if err != nil {
		return err
	}
	return withCaveat(ctx, func(innerCtx context.Context) error {
		err := s.st.CommitHookChanges(innerCtx, newArgs)
		return errors.Capture(err)
	})
}

func (s *LeadershipService) getManagementCaveat(
	arg unitstate.CommitHookChangesArg,
) (func(context.Context, func(context.Context) error) error, error) {
	if arg.RequiresLeadership() {
		return func(ctx context.Context, fn func(context.Context) error) error {
			return s.leaderEnsurer.WithLeader(ctx, arg.UnitName.Application(), arg.UnitName.String(),
				func(ctx context.Context) error {
					return fn(ctx)
				},
			)
		}, nil
	}
	return func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	}, nil
}

func (s *Service) transformRelationSettings(
	ctx context.Context, in []unitstate.RelationSettings,
) ([]internal.RelationSettings, error) {
	return transform.SliceOrErr(in, func(in unitstate.RelationSettings) (internal.RelationSettings, error) {
		// TODO HEATHER - rework to only use relations included by unit.
		relationUUID, err := s.getRelationUUIDByKey(ctx, in.RelationKey)
		if err != nil {
			return internal.RelationSettings{}, errors.Capture(err)
		}
		settings := make(map[string]string, len(in.Settings))
		maps.Copy(settings, in.Settings)
		// The ingress address and egress subnets should be written only by juju.
		delete(settings, unitstate.IngressAddressKey)
		delete(settings, unitstate.EgressSubnetsKey)
		return internal.RelationSettings{
			RelationUUID:        relationUUID,
			Settings:            settings,
			ApplicationSettings: in.ApplicationSettings,
		}, nil
	})
}

// getRelationUUIDByKey returns a relation UUID for the given Key.
func (s *Service) getRelationUUIDByKey(ctx context.Context, relationKey relation.Key) (relation.UUID, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	eids := relationKey.EndpointIdentifiers()
	var uuid relation.UUID
	var err error
	switch len(eids) {
	case 1:
		uuid, err = s.st.GetPeerRelationUUIDByEndpointIdentifiers(
			ctx,
			eids[0],
		)
		if err != nil {
			return "", errors.Errorf("getting peer relation by key: %w", err)
		}
		return uuid, nil
	case 2:
		uuid, err = s.st.GetRegularRelationUUIDByEndpointIdentifiers(
			ctx,
			eids[0],
			eids[1],
		)
		if err != nil {
			return "", errors.Errorf("getting regular relation by key: %w", err)
		}
		return uuid, nil
	default:
		return "", errors.Errorf("internal error: unexpected number of endpoints %d", len(eids))
	}
}

// getUnitEndpointNetworks retrieves network relation information for a given
// unit UUID and its in scope relations. This data will be used to update the
// networks in the relation unit settings.
func (s *LeadershipService) getUnitRelationNetworks(
	ctx context.Context,
	unitUUID unit.UUID,
) ([]internal.RelationNetworkInfo, error) {
	relationUUIDs, err := s.st.GetRelationUUIDsByUnitUUID(ctx, unitUUID)
	if err != nil {
		return nil, errors.Errorf("getting relation UUIDs for unit %q: %w", unitUUID, err)
	}
	if len(relationUUIDs) == 0 {
		return nil, nil
	}

	supportsNetworking, err := s.supportsNetworking(ctx)
	if err != nil {
		return nil, err
	}

	var ingressAddrByRelation map[relation.UUID]string
	if supportsNetworking {
		ingressAddrByRelation, err = s.getUnitRelationsIngressAddress(ctx, unitUUID, relationUUIDs)
	} else {
		ingressAddrByRelation, err = s.getUnitIngressAddress(ctx, unitUUID, relationUUIDs)
	}
	if err != nil {
		return nil, errors.Errorf("getting unit's relations ingress addresses: %w", err)
	}

	egressSubnets, err := s.st.GetRelationsEgressSubnetsByUnitUUID(ctx, unitUUID)
	if err != nil {
		return nil, errors.Errorf(
			"getting egress subnets for relations of unit %q: %w", unitUUID, err,
		)
	}

	var fallbackEgressSubnets []string
	if len(ingressAddrByRelation) != len(egressSubnets) {
		fallbackEgressSubnets, err = s.getFallbackEgressSubnets(ctx, unitUUID)
		if err != nil {
			return nil, errors.Errorf("getting fallback egress subnets: %w", err)
		}
	}

	return transform.MapToSlice(ingressAddrByRelation, func(key relation.UUID, ingressAddr string,
	) []internal.RelationNetworkInfo {
		egressSubnets, ok := egressSubnets[key]
		if !ok {
			// fallback egress subnets are unit, not relation specific.
			egressSubnets = fallbackEgressSubnets
		}
		sort.Strings(egressSubnets)
		return []internal.RelationNetworkInfo{{
			RelationUUID:   key,
			IngressAddress: ingressAddr,
			EgressSubnets:  strings.Join(egressSubnets, ", "),
		}}
	}), nil
}

func (s *LeadershipService) getUnitRelationsIngressAddress(
	ctx context.Context, unitUUID unit.UUID, relationUUIDs []relation.UUID,
) (map[relation.UUID]string, error) {
	ingressAddrByRelation, err := s.st.GetUnitRelationsIngressAddress(ctx, unitUUID)
	if err != nil {
		return nil, errors.Capture(err)
	}
	// ensure every relation has an ingress address, even if it's an empty string.
	for _, relationUUID := range relationUUIDs {
		_, ok := ingressAddrByRelation[relationUUID]
		if !ok {
			ingressAddrByRelation[relationUUID] = ""
		}
	}
	return ingressAddrByRelation, nil
}

func (s *LeadershipService) getUnitIngressAddress(
	ctx context.Context, unitUUID unit.UUID, relationUUIDs []relation.UUID,
) (map[relation.UUID]string, error) {
	unitIngressAddress, err := s.st.GetUnitIngressAddress(ctx, unitUUID)
	if err != nil {
		return nil, errors.Capture(err)
	}
	// ensure every relation has an ingress address.
	return transform.SliceToMap(relationUUIDs, func(relationUUID relation.UUID) (relation.UUID, string) {
		return relationUUID, unitIngressAddress
	}), nil
}

func (s *LeadershipService) getFallbackEgressSubnets(
	ctx context.Context,
	unitUUID unit.UUID,
) ([]string, error) {
	modelEgressSubnets, err := s.st.GetModelEgressSubnets(ctx)
	if err != nil {
		return nil, errors.Errorf("getting model egress subnets: %w", err)
	}
	if len(modelEgressSubnets) > 0 {
		return modelEgressSubnets, nil
	}

	publicEgressSubnets, err := s.getUnitPublicEgressSubnets(ctx, unitUUID)
	if err != nil {
		return nil, errors.Errorf(
			"getting fallback egress subnet for unit %q: %w", unitUUID, err,
		)
	}
	return publicEgressSubnets, nil
}

func (s *LeadershipService) getUnitPublicEgressSubnets(
	ctx context.Context,
	unitUUID unit.UUID,
) ([]string, error) {
	address, err := s.st.GetUnitPublicAddressForEgress(ctx, unitUUID)
	if err != nil {
		s.logger.Warningf(
			ctx,
			"getting unit public address for egress fallback for unit %q: %v",
			unitUUID,
			err,
		)
		return []string{}, nil
	}
	if address == "" {
		return []string{}, nil
	}
	return corenetwork.SubnetsForAddresses([]string{normaliseAddress(address)}), nil
}

func (s *LeadershipService) mergeRelationSettingsAndNetworkInfo(
	ctx context.Context,
	unitUUID unit.UUID,
	argSettings []internal.RelationSettings,
) ([]internal.RelationSettings, error) {
	updates, err := s.getUnitRelationNetworks(ctx, unitUUID)
	if err != nil {
		return nil, err
	}

	if len(updates) == 0 {
		return argSettings, nil
	}

	// Iterate on updates, contains data for all relations.
	// argSettings cannot guarantee that.
	argSettingsMap := transform.SliceToMap(argSettings,
		func(info internal.RelationSettings) (relation.UUID, internal.RelationSettings) {
			return info.RelationUUID, info
		},
	)

	result := make([]internal.RelationSettings, len(updates))
	for i, update := range updates {
		set, ok := argSettingsMap[update.RelationUUID]
		if ok {
			if set.Settings == nil {
				set.Settings = make(map[string]string, 2)
			}
			if update.IngressAddress != "" {
				set.Settings[unitstate.IngressAddressKey] = update.IngressAddress
			}
			if update.RelationUUID != "" {
				set.Settings[unitstate.EgressSubnetsKey] = update.EgressSubnets
			}
			result[i] = set
			continue
		}
		result[i] = internal.RelationSettings{
			RelationUUID: update.RelationUUID,
			Settings: map[string]string{
				unitstate.IngressAddressKey: update.IngressAddress,
				unitstate.EgressSubnetsKey:  update.EgressSubnets,
			},
		}
	}

	return result, nil
}

func normaliseAddress(address string) string {
	before, _, _ := strings.Cut(address, "/")
	return before
}
