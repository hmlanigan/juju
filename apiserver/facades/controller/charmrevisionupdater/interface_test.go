// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmrevisionupdater_test

import (
	"github.com/golang/mock/gomock"
	"github.com/juju/charm/v8"
	"github.com/juju/charm/v8/resource"
	"github.com/juju/charmrepo/v6/csclient"
	csparams "github.com/juju/charmrepo/v6/csclient/params"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/apiserver/facades/controller/charmrevisionupdater"
	"github.com/juju/juju/charmstore"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/state"
	statemocks "github.com/juju/juju/state/mocks"
	"github.com/juju/juju/testing"
)

func makeApplication(ctrl *gomock.Controller, schema, charmName, charmID, appID string, revision int) charmrevisionupdater.Application {
	source := "charm-hub"
	if schema == "cs" {
		source = "charm-store"
	}

	app := NewMockApplication(ctrl)
	app.EXPECT().CharmURL().Return(&charm.URL{
		Schema:   schema,
		Name:     charmName,
		Revision: revision,
	}, false).AnyTimes()
	app.EXPECT().CharmOrigin().Return(&state.CharmOrigin{
		Source:   source,
		Type:     "charm",
		ID:       charmID,
		Revision: &revision,
		Channel: &state.Channel{
			Track: "latest",
			Risk:  "stable",
		},
		Platform: &state.Platform{
			Architecture: "amd64",
			OS:           "ubuntu",
			Series:       "focal",
		},
	}).AnyTimes()
	app.EXPECT().Channel().Return(csparams.Channel("latest/stable")).AnyTimes()
	app.EXPECT().ApplicationTag().Return(names.ApplicationTag{Name: appID}).AnyTimes()

	return app
}

func makeModel(c *gc.C, ctrl *gomock.Controller) charmrevisionupdater.Model {
	model := NewMockModel(ctrl)
	model.EXPECT().CloudName().Return("testcloud").AnyTimes()
	model.EXPECT().CloudRegion().Return("juju-land").AnyTimes()
	uuid := testing.ModelTag.Id()
	cfg, err := config.New(true, map[string]interface{}{
		"charm-hub-url": "https://api.staging.snapcraft.io", // not actually used in tests
		"name":          "model",
		"type":          "type",
		"uuid":          uuid,
	})
	c.Assert(err, jc.ErrorIsNil)
	model.EXPECT().Config().Return(cfg, nil).AnyTimes()
	model.EXPECT().IsControllerModel().Return(false).AnyTimes()
	model.EXPECT().UUID().Return(uuid).AnyTimes()
	return model
}

func (s *updaterSuite) makeState(c *gc.C, ctrl *gomock.Controller, resources state.Resources) {
	if resources == nil {
		r := statemocks.NewMockResources(ctrl)
		r.EXPECT().SetCharmStoreResources(gomock.Any(), gomock.Len(0), gomock.Any()).Return(nil).AnyTimes()
		resources = r
	}

}

func makeResource(c *gc.C, name string, revision, size int, hexFingerprint string) resource.Resource {
	fingerprint, err := resource.ParseFingerprint(hexFingerprint)
	c.Assert(err, jc.ErrorIsNil)
	return resource.Resource{
		Meta: resource.Meta{
			Name: name,
			Type: resource.TypeFile,
		},
		Origin:      resource.OriginStore,
		Revision:    revision,
		Fingerprint: fingerprint,
		Size:        int64(size),
	}
}

func newFakeCharmstoreClient(st charmrevisionupdater.State) (charmstore.Client, error) {
	charms := map[string]charmstoreCharm{
		"mysql":     {name: "mysql", revision: 23},
		"wordpress": {name: "wordpress", revision: 26},
	}
	client := &fakeCharmstoreClient{charms: charms}
	return charmstore.NewCustomClient(client), nil
}

type charmstoreCharm struct {
	name     string
	revision int
}

type fakeCharmstoreClient struct {
	charms map[string]charmstoreCharm
}

func (c *fakeCharmstoreClient) Latest(_ csparams.Channel, ids []*charm.URL, _ map[string][]string) ([]csparams.CharmRevision, error) {
	revisions := make([]csparams.CharmRevision, len(ids))
	for i, id := range ids {
		charm, ok := c.charms[id.Name]
		if !ok {
			revisions[i] = csparams.CharmRevision{Err: errors.NotFoundf("charm %q", id.Name)}
			continue
		}
		revisions[i] = csparams.CharmRevision{Revision: charm.revision}
	}
	return revisions, nil
}

func (c *fakeCharmstoreClient) ListResources(_ csparams.Channel, _ *charm.URL) ([]csparams.Resource, error) {
	return nil, nil
}

func (c *fakeCharmstoreClient) GetResource(_ csparams.Channel, _ *charm.URL, _ string, revision int) (csclient.ResourceData, error) {
	panic("not implemented")
}

func (c *fakeCharmstoreClient) ResourceMeta(_ csparams.Channel, _ *charm.URL, _ string, revision int) (csparams.Resource, error) {
	panic("not implemented")
}

func (c *fakeCharmstoreClient) ServerURL() string {
	panic("not implemented")
}

type facadeContextShim struct {
	facade.Context // Make it fulfil the interface, but we only define a couple of methods
	state          *state.State
	authorizer     facade.Authorizer
}

func (s facadeContextShim) Auth() facade.Authorizer {
	return s.authorizer
}

func (s facadeContextShim) State() *state.State {
	return s.state
}
