// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provisioner

import (
	"fmt"
	"sync"

	"github.com/juju/errors"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/catacomb"

	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/watcher"
)

func NewContainerWorker(cs *ContainerSetup) (worker.Worker, error) {
	containerWatcher, err := cs.getContainerWatcherFunc()
	if err != nil {
		return nil, err
	}
	w := &ContainerWorker{
		catacomb:         catacomb.Catacomb{},
		containerWatcher: containerWatcher,
		logger:           cs.logger,
		cs:               cs,
	}

	if err := catacomb.Invoke(catacomb.Plan{
		Site: &w.catacomb,
		Work: w.loop,
		Init: []worker.Worker{w.containerWatcher},
	}); err != nil {
		return nil, errors.Trace(err)
	}
	return w, nil
}

type ContainerWorker struct {
	catacomb catacomb.Catacomb

	cs *ContainerSetup

	containerWatcher watcher.StringsWatcher
	containerType    instance.ContainerType
	logger           Logger
	provisioner      ContainerProvisioner

	// For introspection Report
	mu sync.Mutex
}

func (w *ContainerWorker) loop() error {
	// Wait for a container of w.ContainerType type to be
	// found.
	if err := w.containerWatcherLoop(); err != nil {
		return err
	}
	if err := w.checkDying(); err != nil {
		return err
	}

	// The container watcher is no longer needed
	if err := worker.Stop(w.containerWatcher); err != nil {
		return err
	}
	w.mu.Lock()
	w.containerWatcher = nil
	w.mu.Unlock()

	// Set up w.ContainerType provisioning dependencies
	// on this machine.
	if err := w.cs.initialiseContainers(w.catacomb.Dying()); err != nil {
		return err
	}
	if err := w.checkDying(); err != nil {
		return err
	}

	// Configure and Add the w.ContainerType Provisioner
	provisioner, err := w.cs.initialiseContainerProvisioner()
	if err != nil {
		return err
	}
	w.logger.Tracef("Starting %s provisioner for %q", w.containerType, w.cs.mTag)
	if err := w.checkDying(); err != nil {
		return err
	}
	if err := w.catacomb.Add(provisioner); err != nil {
		return err
	}

	// For introspection Report
	w.mu.Lock()
	w.provisioner = provisioner
	w.mu.Unlock()

	// Set the w.ContainerType provisioner to doing it's work.
	return w.provisioner.Loop()
}

func (w *ContainerWorker) checkDying() error {
	select {
	case <-w.catacomb.Dying():
		return w.catacomb.ErrDying()
	default:
		return nil
	}
}

func (w *ContainerWorker) containerWatcherLoop() error {
	for {
		select {
		case <-w.catacomb.Dying():
			return w.catacomb.ErrDying()
		case containerIds, ok := <-w.containerWatcher.Changes():
			if !ok {
				return errors.New("container watcher closed")
			}
			// Consume the initial watcher event.
			if len(containerIds) == 0 {
				continue
			}
			return nil
		}
	}
}

// Kill is part of the worker.Worker interface.
func (w *ContainerWorker) Kill() {
	w.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (w *ContainerWorker) Wait() error {
	return w.catacomb.Wait()
}

// Report provides information for the engine report.
func (w *ContainerWorker) Report() map[string]interface{} {
	w.mu.Lock()

	watcherName := fmt.Sprintf("%s-container-watcher", string(w.containerType))
	var watcherMsg string
	if w.containerWatcher == nil {
		watcherMsg = fmt.Sprintf("found containers, watcher stopped")
	} else {
		watcherMsg = fmt.Sprintf("waiting for containers")
	}
	provisionerName := fmt.Sprintf("%s-provisioner", string(w.containerType))
	var provisionerMsg string
	if w.provisioner == nil {
		provisionerMsg = fmt.Sprintf("not setup, nor running")
	} else {
		provisionerMsg = fmt.Sprintf("setup and running")
	}
	result := map[string]interface{}{
		watcherName:     watcherMsg,
		provisionerName: provisionerMsg,
	}
	w.mu.Unlock()
	return result
}
