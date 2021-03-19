// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apicaller

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/juju/juju/cmd/jujud/agent/engine"
	"github.com/juju/worker/v2"
	"github.com/juju/worker/v2/dependency"

	"github.com/juju/juju/agent"
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/base"
)

// Logger represents the methods used by the worker to log details.
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

// ConnectFunc is responsible for making and validating an API connection
// on behalf of an agent.
type ConnectFunc func(agent.Agent, api.OpenFunc, Logger) (api.Connection, error)

// ManifoldConfig defines a Manifold's dependencies.
type ManifoldConfig struct {

	// AgentName is the name of the Agent resource that supplies
	// connection information.
	AgentName string

	// APIConfigWatcherName identifies a resource that will be
	// invalidated when api configuration changes. It's not really
	// fundamental, because it's not used directly, except to create
	// Inputs; it would be perfectly reasonable to wrap a Manifold
	// to report an extra Input instead.
	APIConfigWatcherName string

	// APIOpen is passed into NewConnection, and should be used to
	// create an API connection. You should probably just set it to
	// the local APIOpen func.
	APIOpen api.OpenFunc

	// NewConnection is responsible for getting a connection from an
	// agent, and may be responsible for other things that need to be
	// done before anyone else gets to see the connection.
	//
	// You should probably set it to ScaryConnect when running a
	// machine agent, and to OnlyConnect when running a model agent
	// (which doesn't have its own process). Unit agents should use
	// ScaryConnect at the moment; and probably switch to OnlyConnect
	// when they move into machine agent processes.
	NewConnection ConnectFunc

	// Filter is used to specialize responses to connection errors
	// made on behalf of different kinds of agent.
	Filter dependency.FilterFunc

	// Logger is used to write logging statements for the worker.
	Logger Logger

	APIServerFlagName    string
	IsControllerFlagName string
}

// Manifold returns a manifold whose worker wraps an API connection
// made as configured.
func Manifold(config ManifoldConfig) dependency.Manifold {
	inputs := []string{config.AgentName}
	if config.APIConfigWatcherName != "" {
		// config.APIConfigWatcherName is only applicable to unit
		// and machine scoped manifold.
		// It will be empty for model manifolds.
		inputs = append(inputs, config.APIConfigWatcherName)
	}
	if config.APIServerFlagName != "" {
		inputs = append(inputs, config.APIServerFlagName)
	}
	// Confirm we're running in a state server by asking the
	// stateconfigwatcher manifold.
	//var haveStateConfig bool
	//if err := context.Get(config.StateConfigWatcherName, &haveStateConfig); err != nil {
	//	config.Logger.Tracef("state config watcher not available")
	//	return nil, err
	//}
	//if !haveStateConfig {
	//	config.Logger.Tracef("not a state server, not needed")
	//	config.Recorder.Disable()
	//	return nil, dependency.ErrMissing
	//}
	return dependency.Manifold{
		Inputs: inputs,
		Output: outputFunc,
		Start:  config.startFunc(),
		Filter: config.Filter,
	}
}

// startFunc returns a StartFunc that creates a connection based on the
// supplied manifold config and wraps it in a worker.
func (config ManifoldConfig) startFunc() dependency.StartFunc {
	return func(context dependency.Context) (worker.Worker, error) {

		//var isController bool
		//if err := context.Get("state", nil); err == nil {
		//	isController = true
		//}
		//
		//var waiter gate.Waiter
		//if isController && config.APIServerFlagName != "" {
		//	if err := context.Get(config.APIServerFlagName, &waiter); err != nil {
		//		return nil, err

		//	}
		//	if !waiter.IsUnlocked() {
		//		config.Logger.Debugf("controller, but api-server not running")
		//		return nil, dependency.ErrMissing
		//	}
		//}

		var isControllerFlag engine.Flag
		if err := context.Get(config.IsControllerFlagName, &isControllerFlag); err != nil {
			//return nil, err
			config.Logger.Debugf("get IsControllerFlagName err %s", err)
		}

		if isControllerFlag != nil && !isControllerFlag.Check() {
			config.Logger.Debugf("IsNotController")
		} else if isControllerFlag != nil && isControllerFlag.Check() {
			config.Logger.Debugf("IsController")
		}

		var isAPIServerInitializedFlag engine.Flag
		var apiErr error
		if apiErr = context.Get(config.APIServerFlagName, &isAPIServerInitializedFlag); apiErr != nil {
			//return nil, err
			config.Logger.Debugf("get APIServerFlagName err %s", apiErr)
		}
		if isAPIServerInitializedFlag != nil && !isAPIServerInitializedFlag.Check() {
			config.Logger.Debugf("%q not set", config.APIServerFlagName)
			//return nil, errors.Annotatef(dependency.ErrMissing, "isController and %q not set", config.APIServerFlagName)
		} else if isAPIServerInitializedFlag != nil && isAPIServerInitializedFlag.Check() {
			config.Logger.Debugf("%q set", config.APIServerFlagName)
		}

		//var isControllerFlag Flag
		//if err := context.Get(name, &isControllerFlag); err != nil {
		//	return nil, errors.Trace(err)
		//}
		//if !isControllerFlag.Check() {
		//	return nil, errors.Annotatef(dependency.ErrMissing, "%q not set", name)
		//}

		var agent agent.Agent
		if err := context.Get(config.AgentName, &agent); err != nil {
			return nil, err
		}

		agentConfig := agent.CurrentConfig()
		info, ok := agentConfig.APIInfo()
		if !ok {
			return nil, errors.New("API info not available")
		}
		config.Logger.Debugf("addresses are: %v ", info.Addrs)
		config.Logger.Debugf("ports are: %v ", info.Ports())

		one := fmt.Sprintf("localhost:%d", info.Ports()[0])
		if info.Addrs[0] == one {
			if apiErr != nil || (isAPIServerInitializedFlag != nil && !isAPIServerInitializedFlag.Check()) {
				//config.Logger.Debugf("should dependency.ErrMissing: want connect to %s but %q not set", one, config.APIServerFlagName)
				return nil, errors.Annotatef(dependency.ErrMissing, "want connect to %s but %q not set", one, config.APIServerFlagName)
			}
		}

		conn, err := config.NewConnection(agent, config.APIOpen, config.Logger)
		if errors.Cause(err) == ErrChangedPassword {
			return nil, dependency.ErrBounce
		} else if err != nil {
			cfg := agent.CurrentConfig()
			return nil, errors.Annotatef(err, "[%s] %q cannot open api",
				shortModelUUID(cfg.Model()), cfg.Tag().String())
		}
		return newAPIConnWorker(conn), nil
	}
}

// outputFunc extracts an API connection from a *apiConnWorker.
func outputFunc(in worker.Worker, out interface{}) error {
	inWorker, _ := in.(*apiConnWorker)
	if inWorker == nil {
		return errors.Errorf("in should be a %T; got %T", inWorker, in)
	}

	switch outPointer := out.(type) {
	case *base.APICaller:
		*outPointer = inWorker.conn
	case *api.Connection:
		// Using api.Connection is strongly discouraged as consumers
		// of this API connection should not be able to close it. This
		// option is only available to support legacy upgrade steps.
		*outPointer = inWorker.conn
	default:
		return errors.Errorf("out should be *base.APICaller or *api.Connection; got %T", out)
	}
	return nil
}
