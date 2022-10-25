// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provisioner

import (
	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"

	"github.com/juju/juju/agent"
	apiprovisioner "github.com/juju/juju/api/agent/provisioner"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/container/broker"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/machinelock"
	workercommon "github.com/juju/juju/worker/common"
)

// ManifoldConfig defines an environment provisioner's dependencies. It's not
// currently clear whether it'll be easier to extend this type to include all
// provisioners, or to create separate (Environ|Container)Manifold[Config]s;
// for now we dodge the question because we don't need container provisioners
// in dependency engines. Yet.
type ContainerManifoldConfig struct {
	AgentName                    string
	APICallerName                string
	Logger                       Logger
	MachineLock                  machinelock.Lock
	NewCredentialValidatorFacade func(base.APICaller) (workercommon.CredentialAPI, error)
	ContainerType                instance.ContainerType
}

// Validate is called by start to check for bad configuration.
func (cfg ContainerManifoldConfig) Validate() error {
	if cfg.AgentName == "" {
		return errors.NotValidf("empty AgentName")
	}
	if cfg.APICallerName == "" {
		return errors.NotValidf("empty APICallerName")
	}
	if cfg.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	if cfg.MachineLock == nil {
		return errors.NotValidf("missing MachineLock")
	}
	if cfg.NewCredentialValidatorFacade == nil {
		return errors.NotValidf("missing NewCredentialValidatorFacade")
	}
	if cfg.ContainerType == "" {
		return errors.NotValidf("missing Container Type")
	}
	return nil
}

func (cfg ContainerManifoldConfig) start(context dependency.Context) (worker.Worker, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Trace(err)
	}

	var a agent.Agent
	if err := context.Get(cfg.AgentName, &a); err != nil {
		return nil, errors.Trace(err)
	}

	agentConfig := a.CurrentConfig()
	tag := agentConfig.Tag()
	mTag, ok := tag.(names.MachineTag)
	if !ok {
		return nil, errors.NotValidf("%q machine tag", a)
	}

	var apiCaller base.APICaller
	if err := context.Get(cfg.APICallerName, &apiCaller); err != nil {
		return nil, errors.Trace(err)
	}
	pr := apiprovisioner.NewState(apiCaller)
	result, err := pr.Machines(mTag)
	if err != nil {
		return nil, errors.Annotatef(err, "cannot load machine %s from state", tag)
	}
	if len(result) != 1 {
		return nil, errors.Errorf("expected 1 result, got %d", len(result))
	}
	if errors.IsNotFound(result[0].Err) || (result[0].Err == nil && result[0].Machine.Life() == life.Dead) {
		return nil, dependency.ErrUninstall
	}

	m := result[0].Machine
	types, known, err := m.SupportedContainers()
	if err != nil {
		return nil, errors.Annotatef(err, "retrieving supported container types")
	}
	if known == false {
		return nil, errors.Annotatef(dependency.ErrMissing, "no container types determined")
	}
	if len(types) == 0 {
		return nil, errors.Annotatef(dependency.ErrUninstall, "no supported containers on %q", mTag)
	}

	cfg.Logger.Debugf("%s supported containers types set as %q", mTag, types)

	typeSet := set.NewStrings()
	for _, v := range types {
		typeSet.Add(string(v))
	}
	if !typeSet.Contains(string(cfg.ContainerType)) {
		return nil, errors.Annotatef(dependency.ErrUninstall, "%s does not support %s containers", mTag, string(cfg.ContainerType))
	}

	credentialAPI, err := workercommon.NewCredentialInvalidatorFacade(apiCaller)
	if err != nil {
		return nil, errors.Annotatef(err, "cannot get credential invalidator facade")
	}

	cs := NewContainerSetup(ContainerSetupParams{
		Abort:         context.Abort(),
		Logger:        cfg.Logger,
		ContainerType: cfg.ContainerType,
		Machine:       m,
		Provisioner:   pr,
		Config:        agentConfig,
		MachineLock:   cfg.MachineLock,
		CredentialAPI: credentialAPI,
	})

	return cs.initialiseAndStartProvisioner(agentConfig, mTag, pr, m, credentialAPI)
}

func (cs *ContainerSetup) initialiseAndStartProvisioner(
	agentCfg agent.Config,
	mTag names.MachineTag,
	provisioner *apiprovisioner.State,
	machineZone broker.AvailabilityZoner,
	credentialAPI workercommon.CredentialAPI,
) (worker.Worker, error) {
	cs.logger.Debugf("setup and start provisioner for %s containers", cs.containerType)
	managerConfig, err := containerManagerConfig(cs.containerType, provisioner)
	if err != nil {
		return nil, errors.Annotate(err, "generating container manager config")
	}
	managerConfigWithZones, err := broker.ConfigureAvailabilityZone(managerConfig, machineZone)
	if err != nil {
		return nil, errors.Annotate(err, "configuring availability zones")
	}

	if err := cs.initContainerDependencies(cs.containerType, managerConfig); err != nil {
		return nil, errors.Annotate(err, "setting up container dependencies on host machine")
	}

	instanceBroker, err := broker.New(broker.Config{
		Name:          "provisioner",
		ContainerType: cs.containerType,
		ManagerConfig: managerConfigWithZones,
		APICaller:     provisioner,
		AgentConfig:   cs.config,
		MachineTag:    mTag,
		MachineLock:   cs.machineLock,
		GetNetConfig:  cs.getNetConfig,
	})
	if err != nil {
		return nil, errors.Annotate(err, "initialising container infrastructure on host machine")
	}

	toolsFinder := getToolsFinder(provisioner)
	w, err := NewContainerProvisioner(cs.containerType,
		provisioner,
		cs.logger.Child(string(cs.containerType)),
		agentCfg,
		instanceBroker,
		toolsFinder,
		getDistributionGroupFinder(provisioner),
		credentialAPI,
	)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return w, nil
}

// Manifold creates a manifold that runs an environment provisioner. See the
// ManifoldConfig type for discussion about how this can/should evolve.
func ContainerManifold(config ContainerManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			config.AgentName,
			config.APICallerName,
		},
		Start: config.start,
	}
}
