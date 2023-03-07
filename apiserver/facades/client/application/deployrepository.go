// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"fmt"
	"sync"

	"github.com/juju/charm/v10"
	jujuclock "github.com/juju/clock"
	"github.com/juju/collections/set"
	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/apiserver/facades/client/charms/services"
	"github.com/juju/juju/core/arch"
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/devices"
	"github.com/juju/juju/core/instance"
	coreseries "github.com/juju/juju/core/series"
	"github.com/juju/juju/environs/bootstrap"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
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
	AddCharmMetadata(info state.CharmInfo) (Charm, error)
	Machine(string) (Machine, error)
	ModelConstraints() (constraints.Value, error)

	services.StateBackend
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
	// From queueAsyncCharmDownload
	_, err := api.state.AddCharmMetadata(state.CharmInfo{
		Charm: dt.charm,
		ID:    dt.charmURL,
	})
	if err != nil {
		return nil, nil, []error{errors.Trace(err)}
	}

	// TODO:
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

func makeDeployFromRepositoryValidator(st DeployFromRepositoryState, m DeployFromRepositoryModel, charmhubHTTPClient facade.HTTPClient) DeployFromRepositoryValidator {
	v := &deployFromRepositoryValidator{
		charmhubHTTPClient: charmhubHTTPClient,
		model:              m,
		state:              st,
		newRepoFactory: func(cfg services.CharmRepoFactoryConfig) corecharm.RepositoryFactory {
			return services.NewCharmRepoFactory(cfg)
		},
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

type deployFromRepositoryValidator struct {
	model DeployFromRepositoryModel
	state DeployFromRepositoryState

	mu                 sync.Mutex
	repoFactory        corecharm.RepositoryFactory
	newRepoFactory     func(services.CharmRepoFactoryConfig) corecharm.RepositoryFactory
	charmhubHTTPClient facade.HTTPClient
}

// validateDeployFromRepositoryArgs does validation of all provided
// arguments. Returned is a deployTemplate which contains validate
// data necessary to deploy the application.
func (v *deployFromRepositoryValidator) validate(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	errs := make([]error, 0)
	// Are we deploying a charm? if not, fail fast here.
	// TODO: add a ErrorNotACharm or the like for the juju client.

	initialCurl, requestedOrigin, usedModelDefaultBase, err := v.createOrigin(arg)
	if err != nil {
		// HEATHER
	}
	// TODO:
	// The logic in resolveCharm and getCharm can be improved as there is some
	// duplication. We call ResolveCharmWithPreferredChannel, then pick a
	// series, then call GetEssentialMetadata, which again calls ResolveCharmWithPreferredChannel
	// then a refresh request.

	charmURL, resolvedOrigin, err := v.resolveCharm(initialCurl, requestedOrigin, arg.Force, usedModelDefaultBase)
	// TODO: determine base to use here.

	// get the charm data to validate against, either a previsouly deployed
	// charm or the essential metdata from a charm to be async downloaded.
	resolvedOrigin, resolvedCharm, err := v.getCharm(charmURL, resolvedOrigin)
	if err != nil {
		// HEATHER
	}

	// TODO
	//charmConfig, err := resolvedCharm.Config().ValidateSettings(arg.ConfigYAML)
	//if err != nil {
	//	// HEATHER
	//}

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

	//model, err := api.backend.Model()
	//if err != nil {
	//	// HEATHER
	//}
	//modelType := model.Type()
	//
	//appConfig, appSchema, charmSettings, appDefaults, err := parseCharmSettings(modelType, newCharm, params.AppName, params.ConfigSettingsStrings, params.ConfigSettingsYAML, environsconfig.NoDefaults)
	//if err != nil {
	//	// HEATHER
	//}
	//if err := appConfig.Validate(); err != nil {
	//	// HEATHER
	//}
	//var stateStorageConstraints map[string]state.StorageConstraints
	//if len(arg.Storage) > 0 {
	//	stateStorageConstraints = make(map[string]state.StorageConstraints)
	//	for name, cons := range arg.Storage {
	//		stateCons := state.StorageConstraints{Pool: cons.Pool}
	//		if cons.Size != nil {
	//			stateCons.Size = *cons.Size
	//		}
	//		if cons.Count != nil {
	//			stateCons.Count = *cons.Count
	//		}
	//		stateStorageConstraints[name] = stateCons
	//	}
	//}

	//if err := c.validateCharmFlags(); err != nil {
	//		return errors.Trace(err)
	//	}

	// Enforce "assumes" requirements if the feature flag is enabled.
	//if err := assertCharmAssumptions(resolvedCharm.Meta().Assumes, model, st.ControllerConfig); err != nil {
	//	if !errors.IsNotSupported(err) || !arg.Force {
	//		// HEATHER
	//	}
	//
	//	logger.Warningf("proceeding with deployment of application %q even though the charm feature requirements could not be met as --force was specified", args.ApplicationName)
	//}

	if corecharm.IsKubernetes(resolvedCharm) && charm.MetaFormat(resolvedCharm) == charm.FormatV1 {
		logger.Debugf("DEPRECATED: %q is a podspec charm, which will be removed in a future release", arg.CharmName)
	}

	// Validate the other args.
	return deployTemplate{
		applicationName: arg.ApplicationName,
		charm:           resolvedCharm,
		charmURL:        charmURL,
		origin:          resolvedOrigin,
		placement:       arg.Placement,
		//storage: stateStorageConstraints,
	}, errs
}

type caasDeployFromRepositoryValidator struct {
	validator *deployFromRepositoryValidator
}

func (v caasDeployFromRepositoryValidator) ValidateArg(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	// TODO: NumUnits
	// TODO: Storage
	// TODO: Warn on use of old kubernetes series in charms
	return v.validator.validate(arg)
}

type iaasDeployFromRepositoryValidator struct {
	validator *deployFromRepositoryValidator
}

func (v *iaasDeployFromRepositoryValidator) ValidateArg(arg params.DeployFromRepositoryArg) (deployTemplate, []error) {
	// TODO: NumUnits
	// TODO: Storage
	return v.validator.validate(arg)
}

func (v *deployFromRepositoryValidator) createOrigin(arg params.DeployFromRepositoryArg) (*charm.URL, corecharm.Origin, bool, error) {
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
	channel, err := charm.ParseChannelNormalize(arg.Channel)
	if err != nil {
		return nil, corecharm.Origin{}, false, err
	}
	if channel.Empty() {
		channel.Risk = corecharm.DefaultChannelString
	}

	plat, usedModelDefaultBase, err := v.deducePlatform(arg)
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

// platform is determined by the args: architecture constraint and provided base.
// Check placement to determine known machine platform. If diffs from other provided
// data return error.
// If no base provided, use model default base
// If model default base, will be determined later with help from Charmhub
// If no architecture provided, use model default.
func (v *deployFromRepositoryValidator) deducePlatform(arg params.DeployFromRepositoryArg) (corecharm.Platform, bool, error) {
	arch := arg.Cons.Arch
	base := arg.Base

	// Try base with provided arch and base first.
	platform := corecharm.Platform{}
	if arch != nil {
		platform.Architecture = *arch
	}
	var usedModelDefaultBase bool
	if base != nil {
		platform.OS = base.Name
		platform.Channel = base.Channel
	}
	_, err := corecharm.ParsePlatform(platform.String())
	if err != nil && !errors.Is(err, errors.BadRequest) {
		return corecharm.Platform{}, usedModelDefaultBase, err
	}
	argEmptyPlatform := errors.Is(err, errors.BadRequest)

	// Match against platforms from placement
	placementPlatform, placementsMatch, err := v.platformFromPlacement(arg.Placement)
	if err != nil && !errors.Is(err, errors.NotFound) {
		return corecharm.Platform{}, usedModelDefaultBase, err
	}
	if err == nil && !placementsMatch {
		return corecharm.Platform{}, usedModelDefaultBase, errors.BadRequestf("bases of existing placement machines do not match")
	}

	// No platform args, and one platform from placement, use that.
	if placementsMatch && argEmptyPlatform {
		return placementPlatform, usedModelDefaultBase, nil
	}

	// Fallback to defaults if set.
	if platform.Architecture == "" {
		mConst, err := v.state.ModelConstraints()
		if err != nil {
			return corecharm.Platform{}, usedModelDefaultBase, err
		}
		if mConst.Arch != nil {
			platform.Architecture = *mConst.Arch
		}
	}
	if platform.Channel == "" {
		mCfg, err := v.model.Config()
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
		}
	}
	return platform, usedModelDefaultBase, nil
}

func (v *deployFromRepositoryValidator) platformFromPlacement(placements []*instance.Placement) (corecharm.Platform, bool, error) {
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

func (v *deployFromRepositoryValidator) resolveCharm(curl *charm.URL, requestedOrigin corecharm.Origin, force, usedModelDefaultBase bool) (*charm.URL, corecharm.Origin, error) {
	repo, err := v.getCharmRepository(requestedOrigin.Source)
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
		cons, err := v.state.ModelConstraints()
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

	modelCfg, err := v.model.Config()
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
func (v *deployFromRepositoryValidator) getCharm(charmURL *charm.URL, resolvedOrigin corecharm.Origin) (corecharm.Origin, charm.Charm, error) {
	repo, err := v.getCharmRepository(corecharm.CharmHub)
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

func (v *deployFromRepositoryValidator) getCharmRepository(src corecharm.Source) (corecharm.Repository, error) {
	// The following is only required for testing, as we generate api new http
	// client here for production.
	v.mu.Lock()
	if v.repoFactory != nil {
		defer v.mu.Unlock()
		return v.repoFactory.GetCharmRepository(src)
	}
	v.mu.Unlock()

	repoFactory := v.newRepoFactory(services.CharmRepoFactoryConfig{
		Logger:             logger,
		CharmhubHTTPClient: v.charmhubHTTPClient,
		StateBackend:       v.state,
		ModelBackend:       v.model,
	})

	return repoFactory.GetCharmRepository(src)
}
