// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm

import (
	"strings"

	"github.com/juju/errors"

	"github.com/juju/juju/core/series"
	"github.com/juju/juju/version"
)

const (
	msgUserRequestedBase = "with the user specified base %q"
	msgBundleBase        = "with the base %q defined by the bundle"
	msgLatestLTSBase     = "with the latest LTS base %q"
)

// BaseSelector is a helper type that determines what base the charm should
// be deployed to.
type BaseSelector struct {
	// BaseFlag is the base passed to the --base flag on the command line.
	BaseFlag series.Base
	// Conf is the configuration for the model we're deploying to.
	Conf modelConfig
	// SupportedBases is the list of base the charm supports.
	SupportedBases []series.Base
	// SupportedJujuBases is the list of base that juju supports.
	SupportedJujuBases []series.Base
	// Force indicates the user explicitly wants to deploy to a requested
	// base, regardless of whether the charm says it supports that base.
	Force bool
	// from bundle specifies the deploy request comes from a bundle spec.
	FromBundle bool
	Logger     SelectorLogger
	// UsingImageID is true when the user is using the image-id constraint
	// when deploying the charm. This is needed to validate that in that
	// case the user is also explicitly providing a base.
	UsingImageID bool
}

// TODO(nvinuesa): The Force flag is only valid if the BaseFlag is specified
// or to force the deploy of a LXD profile that doesn't pass validation, this
// should be added to these validation checks.
func (s BaseSelector) validate() error {
	// If the image-id constraint is provided then base must be explicitly
	// provided either by flag either by model-config default base.
	_, explicit := s.Conf.DefaultBase()
	if s.UsingImageID && s.BaseFlag.Empty() && !explicit {
		return errors.Forbiddenf("base must be explicitly provided when image-id constraint is used")
	}
	if len(s.SupportedBases) == 0 {
		return errors.NotValidf("charm does not define any bases,")
	}
	if len(s.SupportedJujuBases) == 0 {
		return errors.BadRequestf("programming error: no juju supported bases")
	}
	return nil
}

// CharmBase determines what base to use with a charm.
// Order of preference is:
//   - user requested with --base or defined by bundle when deploying
//   - model default, if set, acts like --base
//   - juju default ubuntu LTS from charm manifest
//   - first base listed in the charm manifest
//   - in the case of local charms with no manifest nor base in metadata,
//     base must be provided by the user.
func (s BaseSelector) CharmBase() (selectedBase series.Base, err error) {
	// TODO(sidecar): handle systems
	if err := s.validate(); err != nil {
		return series.Base{}, err
	}

	// User has requested a base with --base.
	if !s.BaseFlag.Empty() {
		return s.userRequested(s.BaseFlag)
	}

	// No base explicitly requested by the user.
	// Use model default base, if explicitly set and supported by the charm.
	if defaultBase, explicit := s.Conf.DefaultBase(); explicit {
		base, err := series.ParseBaseFromString(defaultBase)
		if err != nil {
			return series.Base{}, errors.Trace(err)
		}
		return s.userRequested(base)
	}

	// Next fall back to the charm's list of base, filtered to what's supported
	// by Juju.
	var supportedBase []series.Base
	for _, charmBase := range s.SupportedBases {
		for _, jujuCharmBase := range s.SupportedBases {
			if jujuCharmBase.IsCompatible(charmBase) {
				supportedBase = append(supportedBase, charmBase)
				s.Logger.Infof(msgUserRequestedBase, charmBase)
			}
		}
	}
	if len(supportedBase) == 0 {
		return series.Base{}, errors.NotSupportedf("the charm defined base(s) %q ", printBases(s.SupportedBases))
	}

	// Prefer latest Ubuntu LTS.
	preferredBase, err := BaseForCharm(series.LatestLTSBase(), supportedBase)
	if err == nil {
		s.Logger.Infof(msgLatestLTSBase, series.LatestLTSBase())
		return preferredBase, nil
	} else if IsMissingBaseError(err) {
		return series.Base{}, err
	}

	// Try juju's current default supported Ubuntu LTS
	jujuDefaultBase, err := BaseForCharm(version.DefaultSupportedLTSBase(), supportedBase)
	if err == nil {
		s.Logger.Infof(msgLatestLTSBase, version.DefaultSupportedLTSBase())
		return jujuDefaultBase, nil
	}

	// Last chance, the first base in the charm's manifest
	return BaseForCharm(series.Base{}, supportedBase)
}

// userRequested checks the base the user has requested, and returns it if it
// is supported, or if they used --Force.
func (s BaseSelector) userRequested(requestedBase series.Base) (series.Base, error) {
	// TODO(sidecar): handle computed base
	base, err := BaseForCharm(requestedBase, s.SupportedBases)
	if s.Force {
		base = requestedBase
	} else if err != nil {
		return series.Base{}, err
	}

	// validate the base we get from the charm
	if err := s.validateBase(base); err != nil {
		return series.Base{}, err
	}

	// either it's a supported base or the user used --Force, so just
	// give them what they asked for.
	if s.FromBundle {
		s.Logger.Infof(msgBundleBase, base)
		return base, nil
	}
	s.Logger.Infof(msgUserRequestedBase, base)
	return base, nil
}

func (s BaseSelector) validateBase(base series.Base) error {
	for _, jujuBase := range s.SupportedJujuBases {
		if jujuBase.IsCompatible(base) {
			return nil
		}
	}
	return errors.NotSupportedf("base: %q", base.DisplayString())
}

func printBases(bases []series.Base) string {
	baseStrings := make([]string, len(bases))
	for i, base := range bases {
		baseStrings[i] = base.DisplayString()
	}
	return strings.Join(baseStrings, ", ")
}
