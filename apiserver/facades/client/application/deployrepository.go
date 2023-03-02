// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"fmt"

	"github.com/juju/charm/v10"
	"github.com/juju/collections/set"
	"github.com/juju/errors"

	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/devices"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/series"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	"github.com/juju/juju/storage"
)

// DeployFromRepositoryValidator defines an deploy config validator.
type DeployFromRepositoryValidator interface {
	ValidateArg(params.DeployFromRepositoryArg) (deployTemplate, []error)
}

// DeployFromRepository defines an interface for deploying a charm
// from a repository.
type DeployFromRepository interface {
	DeployFromRepository(arg params.DeployFromRepositoryArg) ([]string, []*params.PendingResourceUpload, []error)
}

// DeployFromRepositoryState defines a common set of functions for retrieving state
// objects.
type DeployFromRepositoryState interface {
	Machine(string) (Machine, error)
	ModelConstraints() (constraints.Value, error)
}

// DeployFromRepositoryModel defines a common set of functions for retrieving model
// objects.
type DeployFromRepositoryModel interface {
	Config() (*config.Config, error)
	Type() state.ModelType
}

// DeployFromRepositoryAPI provides the deploy from repository
// API facade for any given version. It is expected that any API
// parameter changes should be performed before entering the API.
type DeployFromRepositoryAPI struct {
	state     DeployFromRepositoryState
	validator DeployFromRepositoryValidator
}

// NewDeployFromRepositoryAPI creates a new DeployFromRepositoryAPI.
func NewDeployFromRepositoryAPI(state DeployFromRepositoryState, validator DeployFromRepositoryValidator) DeployFromRepository {
	api := &DeployFromRepositoryAPI{
		state:     state,
		validator: validator,
	}
	return api
}

func (api *DeployFromRepositoryAPI) DeployFromRepository(arg params.DeployFromRepositoryArg) ([]string, []*params.PendingResourceUpload, []error) {
	// Validate the args.
	dt, errs := api.validator.ValidateArg(arg)

	if len(errs) > 0 {
		return nil, nil, errs
	}

	if dt.dryRun {

	}
	// TODO:
	// SetCharm equivalent method called here
	// AddApplication equivalent method called here.

	// Last step, add pending resources.
	pendingResourceUploads, errs := addPendingResources()

	return nil, pendingResourceUploads, errs
}

// addPendingResource adds a pending resource doc for all resources to be
// added when deploying the charm. PendingResourceUpload is only returned
// for local resources which will require the client to upload the
// resource once DeployFromRepository returns. All resources will be
// processed. Errors are not terminal.
// TODO: determine necessary args.
func addPendingResources() ([]*params.PendingResourceUpload, []error) {
	return nil, nil
}

func makeDeployFromRepositoryValidator(st DeployFromRepositoryState, m DeployFromRepositoryModel, client CharmhubClient) DeployFromRepositoryValidator {
	v := deployFromRepositoryValidator{
		client: client,
		model:  m,
		state:  st,
	}
	if m.Type() == state.ModelTypeCAAS {
		return &caasDeployFromRepositoryValidator{
			validator: v,
		}
	}
	return &iaasDeployFromRepositoryValidator{
		validator: v,
	}
}

type deployTemplate struct {
	applicationName string
	attachStorage   []string
	charmName       string
	configYaml      string
	constraints     constraints.Value
	devices         map[string]devices.Constraints
	endpoints       map[string]string
	dryRun          bool
	force           bool
	numUnits        int
	origin          corecharm.Origin
	placement       []*instance.Placement
	resources       map[string]string
	storage         map[string]storage.Constraints
	trust           bool
}

type deployFromRepositoryValidator struct {
	client CharmhubClient
	model  DeployFromRepositoryModel
	state  DeployFromRepositoryState
}

// validateDeployFromRepositoryArgs does validation of all provided
// arguments. Returned is a deployTemplate which contains validate
// data necessary to deploy the application.
func (v deployFromRepositoryValidator) validate(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	template := deployTemplate{}
	// Are we deploying a charm? if not, fail fast here.
	// TODO: add a ErrorNotACharm or the like for the juju client.

	origin, err := v.createOrigin(arg)
	if err != nil {

	}
	//resolvedOrigin, supportedBases, err := resolveCharm(origin)
	_, _, err = resolveCharm(origin)

	// Validate the other args.
	return template, nil
}

type caasDeployFromRepositoryValidator struct {
	validator deployFromRepositoryValidator
}

func (v caasDeployFromRepositoryValidator) ValidateArg(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	// TODO: NumUnits
	// TODO: Storage
	// TODO: Warn on use of old kubernetes series in charms
	return v.validator.validate(arg)
}

type iaasDeployFromRepositoryValidator struct {
	validator deployFromRepositoryValidator
}

func (v iaasDeployFromRepositoryValidator) ValidateArg(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	// TODO: NumUnits
	// TODO: Storage
	return v.validator.validate(arg)
}

func (v deployFromRepositoryValidator) createOrigin(arg params.DeployFromRepositoryArg) (corecharm.Origin, error) {
	curl, err := charm.ParseURL(arg.CharmName)
	if err != nil {
		return corecharm.Origin{}, err
	}
	if !charm.CharmHub.Matches(curl.Schema) {
		return corecharm.Origin{}, errors.Errorf("unknown schema for charm URL %q", curl.String())
	}
	channel, err := charm.ParseChannelNormalize(arg.Channel)
	if err != nil {
		return corecharm.Origin{}, err
	}
	if channel.Empty() {
		channel.Risk = corecharm.DefaultChannelString
	}

	plat, err := v.deducePlatform(arg)
	if err != nil {
		return corecharm.Origin{}, err
	}

	origin := corecharm.Origin{
		Channel:  &channel,
		Platform: plat,
		Revision: arg.Revision,
		Source:   corecharm.CharmHub,
	}
	return origin, nil
}

// platform is determined by the args: architecture constraint and provided base.
// Check placement to determine known machine platform. If diffs from other provided
// data return error.
// If no base provided, use model default base
// If model default base, will be determined later with help from Charmhub
// If no architecture provided, use model default.
func (v deployFromRepositoryValidator) deducePlatform(arg params.DeployFromRepositoryArg) (corecharm.Platform, error) {
	arch := arg.Cons.Arch
	base := arg.Base

	// Try base with provided arch and base first.
	platform := corecharm.Platform{}
	if arch != nil {
		platform.Architecture = *arch
	}
	if base != nil {
		platform.OS = base.Name
		platform.Channel = base.Channel
	}
	_, err := corecharm.ParsePlatform(platform.String())
	if err != nil && !errors.Is(err, errors.BadRequest) {
		return corecharm.Platform{}, err
	}
	argEmptyPlatform := errors.Is(err, errors.BadRequest)

	// Match against platforms from placement
	placementPlatform, placementsMatch, err := v.platformFromPlacement(arg.Placement)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return corecharm.Platform{}, err
	}
	if err == nil && !placementsMatch {
		return corecharm.Platform{}, errors.BadRequestf("bases of existing placement machines do not match")
	}

	// No platform args, and one platform from placement, use that.
	if placementsMatch && argEmptyPlatform {
		return placementPlatform, nil
	}

	// Fallback to defaults if set.
	if platform.Architecture == "" {
		mConst, err := v.state.ModelConstraints()
		if err != nil {
			return corecharm.Platform{}, err
		}
		if mConst.Arch != nil {
			platform.Architecture = *mConst.Arch
		}
	}
	if platform.Channel == "" {
		mCfg, err := v.model.Config()
		if err != nil {
			return corecharm.Platform{}, err
		}
		if db, ok := mCfg.DefaultBase(); ok {
			defaultBase, err := series.ParseBaseFromString(db)
			if err != nil {
				return corecharm.Platform{}, err
			}
			platform.OS = defaultBase.OS
			platform.Channel = defaultBase.Channel.String()
		}
	}
	return platform, nil
}

func (v deployFromRepositoryValidator) platformFromPlacement(placements []*instance.Placement) (corecharm.Platform, bool, error) {
	if len(placements) == 0 {
		return corecharm.Platform{}, false, errors.NotFoundf("placements")
	}
	machines := make([]Machine, 0)
	// Find which machines in placement actually exist.
	for _, placement := range placements {
		m, err := v.state.Machine(placement.Directive)
		if errors.Is(err, errors.NotFound) {
			continue
		}
		if err != nil {
			return corecharm.Platform{}, false, err
		}
		machines = append(machines, m)
	}
	if len(machines) == 0 {
		return corecharm.Platform{}, false, errors.NotFoundf("machines in placements")
	}

	// Gather platforms for existing machines
	var platform corecharm.Platform
	platStrings := set.NewStrings()
	for _, machine := range machines {
		b := machine.Base()
		a, err := machine.HardwareCharacteristics()
		if err != nil {
			return corecharm.Platform{}, false, err
		}
		platString := fmt.Sprintf("%s/%s/%s", *a.Arch, b.OS, b.Channel)
		p, err := corecharm.ParsePlatformNormalize(platString)
		if err != nil {
			return corecharm.Platform{}, false, err
		}
		platform = p
		platStrings.Add(p.String())
	}

	return platform, platStrings.Size() == 1, nil
}

func resolveCharm(origin corecharm.Origin) (corecharm.Origin, []string, error) {

	//repo, err := a.getCharmRepository(corecharm.Source(arg.Origin.Source))
	//if err != nil {
	//	result.Error = apiservererrors.ServerError(err)
	//	return result
	//}
	//
	//resultURL, origin, supportedSeries, err := repo.ResolveWithPreferredChannel(curl, requestedOrigin)
	//if err != nil {
	//	result.Error = apiservererrors.ServerError(err)
	//	return result
	//}
	//result.URL = resultURL.String()
	//
	//apiOrigin, err := convertOrigin(origin)
	//if err != nil {
	//	result.Error = apiservererrors.ServerError(err)
	//	return result
	//}
	//
	//// The charmhub API can return "all" for architecture as it's not a real
	//// arch we don't know how to correctly model it. "all " doesn't mean use the
	//// default arch, it means use any arch which isn't quite the same. So if we
	//// do get "all" we should see if there is a clean way to resolve it.
	//archOrigin := apiOrigin
	//if apiOrigin.Architecture == "all" {
	//	cons, err := a.backendState.ModelConstraints()
	//	if err != nil {
	//		result.Error = apiservererrors.ServerError(err)
	//		return result
	//	}
	//	archOrigin.Architecture = arch.ConstraintArch(cons, nil)
	//}
	//
	//result.Origin = archOrigin
	//
	//switch {
	//case resultURL.Series != "" && len(supportedSeries) == 0:
	//	result.SupportedSeries = []string{resultURL.Series}
	//default:
	//	result.SupportedSeries = supportedSeries
	//}
	//
	//return result
	return corecharm.Origin{}, nil, nil
}
