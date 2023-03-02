// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"fmt"

	"github.com/juju/charm/v10"
	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"github.com/juju/featureflag"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/devices"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/series"
	"github.com/juju/juju/feature"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/storage"
)

// DeployFromRepository is a one-stop deployment method for repository
// charms. Only a charm name is required to deploy. If argument validation
// fails, a list of all errors found in validation will be returned. If a
// local resource is provided, details required for uploading the validated
// resource will be returned.
func (api *APIBase) DeployFromRepository(args params.DeployFromRepositoryArgs) (params.DeployFromRepositoryResults, error) {
	if !featureflag.Enabled(feature.ServerSideCharmDeploy) {
		return params.DeployFromRepositoryResults{}, errors.NotImplementedf("this facade method is under develop")
	}

	if err := api.checkCanWrite(); err != nil {
		return params.DeployFromRepositoryResults{}, errors.Trace(err)
	}
	if err := api.check.RemoveAllowed(); err != nil {
		return params.DeployFromRepositoryResults{}, errors.Trace(err)
	}

	results := make([]params.DeployFromRepositoryResult, len(args.Args))
	for i, entity := range args.Args {
		info, pending, errs := api.deployOneFromRepository(entity)
		if len(errs) > 0 {
			results[i].Errors = apiservererrors.ServerErrors(errs)
			continue
		}
		results[i].Info = info
		results[i].PendingResourceUploads = pending
	}
	return params.DeployFromRepositoryResults{
		Results: results,
	}, nil
}

func (api *APIBase) deployOneFromRepository(arg params.DeployFromRepositoryArg) ([]string, []*params.PendingResourceUpload, []error) {
	// Validate the args.
	dt, errs := api.validateDeployFromRepositoryArgs(arg)
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

// validateDeployFromRepositoryArgs does validation of all provided
// arguments. Returned is a deployTemplate which contains validate
// data necessary to deploy the application.
func (api *APIBase) validateDeployFromRepositoryArgs(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	template := deployTemplate{}
	// Are we deploying a charm? if not, fail fast here.
	// TODO: add a ErrorNotACharm or the like for the juju client.

	origin, err := api.createOrigin(arg)
	if err != nil {

	}
	//resolvedOrigin, supportedBases, err := resolveCharm(origin)
	_, _, err = resolveCharm(origin)

	// Validate the other args.
	return template, nil
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

func (api *APIBase) createOrigin(arg params.DeployFromRepositoryArg) (corecharm.Origin, error) {
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

	plat, err := api.deducePlatform(arg)
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
func (api *APIBase) deducePlatform(arg params.DeployFromRepositoryArg) (corecharm.Platform, error) {
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
	placementPlatform, placementsMatch, err := api.platformFromPlacement(arg.Placement)
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
		mConst, err := api.backend.ModelConstraints()
		if err != nil {
			return corecharm.Platform{}, err
		}
		if mConst.Arch != nil {
			platform.Architecture = *mConst.Arch
		}
	}
	if platform.Channel == "" {
		mCfg, err := api.model.ModelConfig()
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

func (api *APIBase) platformFromPlacement(placements []*instance.Placement) (corecharm.Platform, bool, error) {
	if len(placements) == 0 {
		return corecharm.Platform{}, false, errors.NotFoundf("placements")
	}
	machines := make([]Machine, 0)
	// Find which machines in placement actually exist.
	for _, placement := range placements {
		m, err := api.backend.Machine(placement.Directive)
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
