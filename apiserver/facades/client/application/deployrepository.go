// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"fmt"

	"github.com/juju/charm/v10"
	jujuclock "github.com/juju/clock"
	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"github.com/juju/featureflag"
	"github.com/kr/pretty"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facades/client/charms/services"
	"github.com/juju/juju/core/arch"
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/devices"
	"github.com/juju/juju/core/instance"
	coreseries "github.com/juju/juju/core/series"
	"github.com/juju/juju/environs/bootstrap"
	"github.com/juju/juju/feature"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
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

func (api *APIBase) deployOneFromRepository(arg params.DeployFromRepositoryArg) (params.DeployFromRepositoryInfo, []*params.PendingResourceUpload, []error) {
	logger.Tracef("deployOneFromRepository(%s)", pretty.Sprint(arg))
	// Validate the args.
	dt, errs := api.validateDeployFromRepositoryArgs(arg)
	if len(errs) > 0 {
		return params.DeployFromRepositoryInfo{}, nil, errs
	}

	logger.Tracef("%s", pretty.Sprint(dt))
	info := params.DeployFromRepositoryInfo{
		CharmURL:     dt.charmURL.String(),
		Risk:         string(dt.origin.Channel.Risk),
		Track:        nil,
		Branch:       nil,
		Architecture: dt.origin.Platform.Architecture,
		Base: params.Base{
			Name:    dt.origin.Platform.OS,
			Channel: dt.origin.Platform.Channel,
		},
		EffectiveChannel: nil,
	}
	if dt.dryRun {
		return info, nil, nil
	}
	// Queue async charm download.
	// AddCharmMetadata returns no error if the charm
	// has already been queue'd or downloaded.
	ch, err := api.backend.AddCharmMetadata(state.CharmInfo{
		Charm: dt.charm,
		ID:    dt.charmURL,
	})
	if err != nil {
		return params.DeployFromRepositoryInfo{}, nil, []error{errors.Trace(err)}
	}

	stOrigin, err := StateCharmOrigin(dt.origin)
	if err != nil {
		return params.DeployFromRepositoryInfo{}, nil, []error{errors.Trace(err)}
	}
	_, err = api.backend.AddApplication(state.AddApplicationArgs{
		Name:              dt.applicationName,
		Charm:             CharmToStateCharm(ch),
		CharmOrigin:       stOrigin,
		Storage:           nil,
		Devices:           nil,
		AttachStorage:     nil,
		EndpointBindings:  nil,
		ApplicationConfig: nil,
		CharmConfig:       nil,
		NumUnits:          dt.numUnits,
		Placement:         dt.placement,
		Constraints:       dt.constraints,
		Resources:         dt.resources,
	})
	if err != nil {
		return params.DeployFromRepositoryInfo{}, nil, []error{errors.Trace(err)}
	}

	// Last step, add pending resources.
	pendingResourceUploads, errs := addPendingResources()

	return info, pendingResourceUploads, errs
}

// validateDeployFromRepositoryArgs does validation of all provided
// arguments. Returned is a deployTemplate which contains validate
// data necessary to deploy the application.
func (api *APIBase) validateDeployFromRepositoryArgs(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	errs := make([]error, 0)
	// Are we deploying a charm? if not, fail fast here.
	// TODO: add a ErrorNotACharm or the like for the juju client.

	initialCurl, requestedOrigin, usedModelDefaultBase, err := api.createOrigin(arg)
	if err != nil {
		errs = append(errs, err)
		return deployTemplate{}, errs
	}
	logger.Criticalf("from createOrigin: %s, %s", initialCurl, pretty.Sprint(requestedOrigin))
	// TODO:
	// The logic in resolveCharm and getCharm can be improved as there is some
	// duplication. We call ResolveCharmWithPreferredChannel, then pick a
	// series, then call GetEssentialMetadata, which again calls ResolveCharmWithPreferredChannel
	// then a refresh request.

	charmURL, resolvedOrigin, err := api.resolveCharm(initialCurl, requestedOrigin, arg.Force, usedModelDefaultBase)
	if err != nil {
		errs = append(errs, err)
		return deployTemplate{}, errs
	}
	logger.Criticalf("from resolveCharm: %s, %s", charmURL, pretty.Sprint(resolvedOrigin))

	// get the charm data to validate against, either a previously deployed
	// charm or the essential metadata from a charm to be async downloaded.
	resolvedOrigin, resolvedCharm, err := api.getCharm(charmURL, resolvedOrigin)
	if err != nil {
		errs = append(errs, err)
		return deployTemplate{}, errs
	}
	logger.Criticalf("from getCharm: %s", charmURL, pretty.Sprint(resolvedOrigin))

	// TODO: validate config

	if resolvedCharm.Meta().Name == bootstrap.ControllerCharmName {
		errs = append(errs, errors.NotSupportedf("manual deploy of the controller charm"))
	}
	if resolvedCharm.Meta().Subordinate {
		if arg.NumUnits != nil && *arg.NumUnits != 0 {
			errs = append(errs, fmt.Errorf("subordinate application must be deployed without units"))
		}
		if !constraints.IsEmpty(&arg.Cons) {
			errs = append(errs, fmt.Errorf("subordinate application must be deployed without constraints"))
		}
	}

	appName := charmURL.Name
	if arg.ApplicationName != "" {
		appName = arg.ApplicationName
	}

	// Enforce "assumes" requirements if the feature flag is enabled.
	if err := assertCharmAssumptions(resolvedCharm.Meta().Assumes, api.model, api.backend.ControllerConfig); err != nil {
		if !errors.Is(err, errors.NotSupported) || !arg.Force {
			errs = append(errs, err)
		}
		logger.Warningf("proceeding with deployment of application %q even though the charm feature requirements could not be met as --force was specified", appName)
	}

	if corecharm.IsKubernetes(resolvedCharm) && charm.MetaFormat(resolvedCharm) == charm.FormatV1 {
		logger.Debugf("DEPRECATED: %q is a podspec charm, which will be removed in a future release", arg.CharmName)
	}

	var numUnits int
	if arg.NumUnits != nil {
		numUnits = *arg.NumUnits
	} else {
		numUnits = 1
	}

	// Validate the other args.
	dt := deployTemplate{
		applicationName: appName,
		charm:           resolvedCharm,
		charmURL:        charmURL,
		dryRun:          arg.DryRun,
		force:           arg.Force,
		numUnits:        numUnits,
		origin:          resolvedOrigin,
		placement:       arg.Placement,
		storage:         stateStorageConstraints(arg.Storage),
	}
	if !resolvedCharm.Meta().Subordinate {
		dt.constraints = arg.Cons
	}
	logger.Criticalf("validateDeployFromRepositoryArgs returning: %s", pretty.Sprint(dt))
	return dt, errs
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
	charm           charm.Charm
	charmURL        *charm.URL
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
	storage         map[string]state.StorageConstraints
	trust           bool
}

func (api *APIBase) createOrigin(arg params.DeployFromRepositoryArg) (*charm.URL, corecharm.Origin, bool, error) {
	path, err := charm.EnsureSchema(arg.CharmName, charm.CharmHub)
	if err != nil {
		return nil, corecharm.Origin{}, false, err
	}
	curl, err := charm.ParseURL(path)
	if err != nil {
		return nil, corecharm.Origin{}, false, err
	}
	if arg.Revision != nil {
		curl = curl.WithRevision(*arg.Revision)
	}
	if !charm.CharmHub.Matches(curl.Schema) {
		return nil, corecharm.Origin{}, false, errors.Errorf("unknown schema for charm URL %q", curl.String())
	}
	if arg.Channel == "" {
		arg.Channel = corecharm.DefaultChannelString
	}
	channel, err := charm.ParseChannelNormalize(arg.Channel)
	if err != nil {
		return nil, corecharm.Origin{}, false, err
	}

	plat, usedModelDefaultBase, err := api.deducePlatform(arg)
	if err != nil {
		return nil, corecharm.Origin{}, false, err
	}

	origin := corecharm.Origin{
		Channel:  &channel,
		Platform: plat,
		Revision: arg.Revision,
		Source:   corecharm.CharmHub,
	}
	return curl, origin, usedModelDefaultBase, nil
}

// deducePlatform returns a platform for initial resolveCharm call.
// At minimum, it must contain an architecture.
// Platform is determined by the args: architecture constraint and
// provided base.
// - Check placement to determine known machine platform. If diffs from
// other provided data return error.
// - If no base provided, use model default base.
// - If no model default base, will be determined later.
// - If no architecture provided, use model default. Fallback
// to DefaultArchitecture.
func (api *APIBase) deducePlatform(arg params.DeployFromRepositoryArg) (corecharm.Platform, bool, error) {
	argArch := arg.Cons.Arch
	argBase := arg.Base

	// Try argBase with provided argArch and argBase first.
	platform := corecharm.Platform{}
	if argArch != nil {
		platform.Architecture = *argArch
	}
	// Fallback to model defaults if set. DefaultArchitecture otherwise.
	if platform.Architecture == "" {
		mConst, err := api.backend.ModelConstraints()
		if err != nil {
			return corecharm.Platform{}, false, err
		}
		if mConst.Arch != nil {
			logger.Criticalf("deducePlatform use model default argArch %s", *mConst.Arch)
			platform.Architecture = *mConst.Arch
		} else {
			platform.Architecture = arch.DefaultArchitecture
		}
	}
	var usedModelDefaultBase bool
	if argBase != nil {
		platform.OS = argBase.Name
		platform.Channel = argBase.Channel
	}

	// Initial validation of platform from known data.
	_, err := corecharm.ParsePlatform(platform.String())
	if err != nil && !errors.Is(err, errors.BadRequest) {
		logger.Criticalf("deducePlatform placements don't match %+v", err)
		return corecharm.Platform{}, usedModelDefaultBase, err
	}
	argEmptyPlatform := errors.Is(err, errors.BadRequest)

	// Match against platforms from placement
	placementPlatform, placementsMatch, err := api.platformFromPlacement(arg.Placement)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return corecharm.Platform{}, usedModelDefaultBase, err
	}
	if err == nil && !placementsMatch {
		return corecharm.Platform{}, usedModelDefaultBase, errors.BadRequestf("bases of existing placement machines do not match")
	}
	logger.Criticalf("deducePlatform placements don't match")

	// No platform args, and one platform from placement, use that.
	if placementsMatch && argEmptyPlatform {
		return placementPlatform, usedModelDefaultBase, nil
	}
	logger.Criticalf("deducePlatform no placements")
	if platform.Channel == "" {
		mCfg, err := api.model.ModelConfig()
		if err != nil {
			return corecharm.Platform{}, usedModelDefaultBase, err
		}
		if db, ok := mCfg.DefaultBase(); ok {
			defaultBase, err := coreseries.ParseBaseFromString(db)
			if err != nil {
				return corecharm.Platform{}, usedModelDefaultBase, err
			}
			platform.OS = defaultBase.OS
			platform.Channel = defaultBase.Channel.String()
			usedModelDefaultBase = true
			logger.Criticalf("deducePlatform use model default base %s", defaultBase)
		}
	}
	return platform, usedModelDefaultBase, nil
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

func (api *APIBase) resolveCharm(curl *charm.URL, requestedOrigin corecharm.Origin, force, usedModelDefaultBase bool) (*charm.URL, corecharm.Origin, error) {
	repo, err := api.getCharmRepository(requestedOrigin.Source)
	if err != nil {
		return nil, corecharm.Origin{}, errors.Trace(err)
	}

	resultURL, resolvedOrigin, supportedSeries, resolveErr := repo.ResolveWithPreferredChannel(curl, requestedOrigin)
	if charm.IsUnsupportedSeriesError(resolveErr) {
		if !force {
			msg := fmt.Sprintf("%v. Use --force to deploy the charm anyway.", resolveErr)
			if usedModelDefaultBase {
				msg += " Used the default-series."
			}
			return nil, corecharm.Origin{}, errors.Errorf(msg)
		}
	} else if resolveErr != nil {
		return nil, corecharm.Origin{}, errors.Trace(resolveErr)
	}
	// TODO: choose a base, if we weren't successful with NA.
	// look at logic above too for this.

	// The charmhub API can return "all" for architecture as it's not a real
	// arch we don't know how to correctly model it. "all " doesn't mean use the
	// default arch, it means use any arch which isn't quite the same. So if we
	// do get "all" we should see if there is a clean way to resolve it.
	if resolvedOrigin.Platform.Architecture == "all" {
		cons, err := api.backend.ModelConstraints()
		if err != nil {
			return nil, corecharm.Origin{}, errors.Trace(err)
		}
		resolvedOrigin.Platform.Architecture = arch.ConstraintArch(cons, nil)
	}

	var seriesFlag string
	if requestedOrigin.Platform.OS != "" {
		var err error
		seriesFlag, err = coreseries.GetSeriesFromChannel(requestedOrigin.Platform.OS, requestedOrigin.Platform.Channel)
		if err != nil {
			return nil, corecharm.Origin{}, errors.Trace(err)
		}
	}

	modelCfg, err := api.model.ModelConfig()
	if err != nil {
		return nil, corecharm.Origin{}, errors.Trace(err)
	}

	imageStream := modelCfg.ImageStream()

	workloadSeries, err := coreseries.WorkloadSeries(jujuclock.WallClock.Now(), seriesFlag, imageStream)
	if err != nil {
		return nil, corecharm.Origin{}, errors.Trace(err)
	}

	selector := corecharm.SeriesSelector{
		SeriesFlag:          seriesFlag,
		SupportedSeries:     supportedSeries,
		SupportedJujuSeries: workloadSeries,
		Force:               force,
		Conf:                modelCfg,
		FromBundle:          false,
		Logger:              logger,
	}

	// Get the series to use.
	series, err := selector.CharmSeries()
	logger.Tracef("Using series %q from %v to deploy %v", series, supportedSeries, curl)
	if charm.IsUnsupportedSeriesError(err) {
		msg := fmt.Sprintf("%v. Use --force to deploy the charm anyway.", err)
		if usedModelDefaultBase {
			msg += " Used the default-series."
		}
		return nil, corecharm.Origin{}, errors.Trace(err)
	}

	var base coreseries.Base
	if series == coreseries.Kubernetes.String() {
		base = coreseries.LegacyKubernetesBase()
	} else {
		base, err = coreseries.GetBaseFromSeries(series)
		if err != nil {
			return nil, corecharm.Origin{}, errors.Trace(err)
		}
	}
	resolvedOrigin.Platform.OS = base.OS
	resolvedOrigin.Platform.Channel = base.Channel.String()

	// handle actualSupportedSeries if possible here...
	return resultURL, resolvedOrigin, nil
}

// getCharm returns the charm being deployed. Either it already has been
// used once and we get the data from state. Or we get the essential metadata.
func (api *APIBase) getCharm(charmURL *charm.URL, resolvedOrigin corecharm.Origin) (corecharm.Origin, charm.Charm, error) {
	repo, err := api.getCharmRepository(corecharm.CharmHub)
	if err != nil {
		return resolvedOrigin, nil, err
	}

	// Check if a charm doc already exists for this charm URL. If so, the
	// charm has already been queued for download so this is a no-op. We
	// still need to resolve and return back a suitable origin as charmhub
	// may refer to the same blob using the same revision in different
	// channels.
	//
	// We need to use GetDownloadURL instead of ResolveWithPreferredChannel
	// to ensure that the resolved origin has the ID/Hash fields correctly
	// populated.
	// TODO: HEATHER need resolved charmurl here.
	// TODO: Handle already deployed charm.
	//deployedCharm, err := api.backend.Charm(charmURL)
	//if err == nil {
	//	// Heather
	//	_, resolvedOrigin, err = repo.GetDownloadURL(charmURL, resolvedOrigin)
	//	if err != nil {
	//		//HEATHER
	//	}
	//} else if !errors.Is(err, errors.NotFound) {
	//	return resolvedOrigin, nil, err
	//}

	// Fetch the essential metadata that we require to deploy the charm
	// without downloading the full archive. The remaining metadata will
	// be populated once the charm gets downloaded.
	essentialMeta, err := repo.GetEssentialMetadata(corecharm.MetadataRequest{
		CharmURL: charmURL,
		Origin:   resolvedOrigin,
	})
	if err != nil {
		return resolvedOrigin, nil, errors.Annotatef(err, "retrieving essential metadata for charm %q", charmURL)
	}
	metaRes := essentialMeta[0]
	resolvedCharm := corecharm.NewCharmInfoAdapter(metaRes)
	return resolvedOrigin, resolvedCharm, nil
}

func (api *APIBase) getCharmRepository(src corecharm.Source) (corecharm.Repository, error) {
	// The following is only required for testing, as we generate api new http
	// client here for production.
	api.mu.Lock()
	if api.repoFactory != nil {
		defer api.mu.Unlock()
		return api.repoFactory.GetCharmRepository(src)
	}
	api.mu.Unlock()

	repoFactory := api.newRepoFactory(services.CharmRepoFactoryConfig{
		Logger:             logger,
		CharmhubHTTPClient: api.charmhubHTTPClient,
		StateBackend:       api.backend,
		ModelBackend:       api.model,
	})

	return repoFactory.GetCharmRepository(src)
}
