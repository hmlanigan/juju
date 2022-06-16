// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm

// NewDeployer returns a manifest deployer. It is a var so that it can be
// patched for uniter tests.
var NewDeployer = newDeployer

// NewDeployerFunc returns a func used to create a deployer.
type NewDeployerFunc func(charmPath, dataPath string, charmReader CharmReader, logger Logger) (Deployer, error)

func newDeployer(charmPath, dataPath string, charmReader CharmReader, logger Logger) (Deployer, error) {
	return NewManifestDeployer(charmPath, dataPath, charmReader, logger), nil
}
