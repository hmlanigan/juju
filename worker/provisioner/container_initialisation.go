// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provisioner

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/juju/names/v4"

	"github.com/juju/juju/agent"
	apiprovisioner "github.com/juju/juju/api/agent/provisioner"
	"github.com/juju/juju/api/common"
	"github.com/juju/juju/container"
	"github.com/juju/juju/container/kvm"
	"github.com/juju/juju/container/lxd"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/machinelock"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/rpc/params"
	workercommon "github.com/juju/juju/worker/common"
)

// ContainerSetup is a StringsWatchHandler that is notified when containers
// are created on the given machine. It will set up the machine to be able
// to create containers and start a suitable provisioner.
type ContainerSetup struct {
	// abort used for init container dependencies
	logger        Logger
	containerType instance.ContainerType
	provisioner   *apiprovisioner.State
	machine       apiprovisioner.MachineProvisioner
	mTag          names.MachineTag
	config        agent.Config
	machineLock   machinelock.Lock
	managerConfig container.ManagerConfig

	// The number of provisioners started. Once all necessary provisioners have
	// been started, the container watcher can be stopped.
	numberProvisioners      int32
	credentialAPI           workercommon.CredentialAPI
	getNetConfig            func(network.ConfigSource) ([]params.NetworkConfig, error)
	getContainerWatcherFunc GetContainerWatcherFunc
}

// ContainerSetupParams are used to initialise a container setup handler.
type ContainerSetupParams struct {
	Logger                  Logger
	ContainerType           instance.ContainerType
	Machine                 apiprovisioner.MachineProvisioner
	MTag                    names.MachineTag
	Provisioner             *apiprovisioner.State
	Config                  agent.Config
	MachineLock             machinelock.Lock
	CredentialAPI           workercommon.CredentialAPI
	GetContainerWatcherFunc GetContainerWatcherFunc
}

// NewContainerSetup returns a ContainerSetup to start the container
// provisioner workers.
func NewContainerSetup(params ContainerSetupParams) *ContainerSetup {
	return &ContainerSetup{
		logger:                  params.Logger,
		machine:                 params.Machine,
		mTag:                    params.MTag,
		containerType:           params.ContainerType,
		provisioner:             params.Provisioner,
		config:                  params.Config,
		machineLock:             params.MachineLock,
		credentialAPI:           params.CredentialAPI,
		getNetConfig:            common.GetObservedNetworkConfig,
		getContainerWatcherFunc: params.GetContainerWatcherFunc,
	}
}

// initContainerDependencies ensures that the host machine is set-up to manage
// containers of the input type.
func (cs *ContainerSetup) initContainerDependencies(abort <-chan struct{}, managerCfg container.ManagerConfig) error {
	snapChannels := map[string]string{
		"lxd": managerCfg.PopValue(config.LXDSnapChannel),
	}
	initialiser := getContainerInitialiser(cs.containerType, snapChannels)

	releaser, err := cs.acquireLock(abort, fmt.Sprintf("%s container initialisation", cs.containerType))
	if err != nil {
		return errors.Annotate(err, "failed to acquire initialization lock")
	}
	defer releaser()

	if err := initialiser.Initialise(); err != nil {
		return errors.Trace(err)
	}

	// At this point, Initialiser likely has changed host network information,
	// so re-probe to have an accurate view.
	observedConfig, err := cs.observeNetwork()
	if err != nil {
		return errors.Annotate(err, "cannot discover observed network config")
	}
	if len(observedConfig) > 0 {
		machineTag := cs.machine.MachineTag()
		cs.logger.Tracef("updating observed network config for %q %s containers to %#v",
			machineTag, cs.containerType, observedConfig)
		if err := cs.provisioner.SetHostMachineNetworkConfig(machineTag, observedConfig); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (cs *ContainerSetup) observeNetwork() ([]params.NetworkConfig, error) {
	return cs.getNetConfig(network.DefaultConfigSource())
}

func (cs *ContainerSetup) acquireLock(abort <-chan struct{}, comment string) (func(), error) {
	spec := machinelock.Spec{
		Cancel:  abort,
		Worker:  "provisioner",
		Comment: comment,
	}
	return cs.machineLock.Acquire(spec)
}

// getContainerInitialiser exists to patch out in tests.
var getContainerInitialiser = func(ct instance.ContainerType, snapChannels map[string]string) container.Initialiser {
	if ct == instance.LXD {
		return lxd.NewContainerInitialiser(snapChannels["lxd"])
	}
	return kvm.NewContainerInitialiser()
}

func containerManagerConfig(
	containerType instance.ContainerType, configGetter ContainerManagerConfigGetter,
) (container.ManagerConfig, error) {
	// Ask the configGetter for the container manager configuration.
	managerConfigResult, err := configGetter.ContainerManagerConfig(
		params.ContainerManagerConfigParams{Type: containerType},
	)
	if err != nil {
		return nil, errors.Trace(err)
	}
	managerConfig := container.ManagerConfig(managerConfigResult.ManagerConfig)
	return managerConfig, nil
}

type ContainerManagerConfigGetter interface {
	ContainerManagerConfig(params.ContainerManagerConfigParams) (params.ContainerManagerConfig, error)
}
