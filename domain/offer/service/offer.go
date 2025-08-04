// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"

	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/internal/uuid"
)

// State describes retrieval and persistence methods for offers.
type State interface{}

// Service provides the API for working with offers.
type Service struct {
	st     State
	logger logger.Logger
}

// NewService returns a new service reference wrapping the input state.
func NewService(
	st State,
	logger logger.Logger,
) *Service {
	return &Service{
		st:     st,
		logger: logger,
	}
}

func (s *Service) GetOfferUUID(ctx context.Context, offerURL *crossmodel.OfferURL) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
