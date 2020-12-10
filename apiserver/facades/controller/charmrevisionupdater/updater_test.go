// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmrevisionupdater_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/juju/charm/v8"
	"github.com/juju/charm/v8/resource"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/facades/controller/charmrevisionupdater"
	testmocks "github.com/juju/juju/apiserver/facades/controller/charmrevisionupdater/mocks"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/charmhub"
	"github.com/juju/juju/charmhub/transport"
	"github.com/juju/juju/cloud"
	statemocks "github.com/juju/juju/state/mocks"
)

type updaterSuite struct {
	state     *MockState
	chClient  *testmocks.MockCharmhubRefreshClient
	resources *statemocks.MockResources
}

var _ = gc.Suite(&updaterSuite{})

func (s *updaterSuite) TestNewAuthSuccess(c *gc.C) {
	authorizer := apiservertesting.FakeAuthorizer{Controller: true}
	facadeCtx := facadeContextShim{state: nil, authorizer: authorizer}
	updater, err := charmrevisionupdater.NewCharmRevisionUpdaterAPI(facadeCtx)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(updater, gc.NotNil)
}

func (s *updaterSuite) TestNewAuthFailure(c *gc.C) {
	authoriser := apiservertesting.FakeAuthorizer{Controller: false}
	facadeCtx := facadeContextShim{state: nil, authorizer: authoriser}
	updater, err := charmrevisionupdater.NewCharmRevisionUpdaterAPI(facadeCtx)
	c.Assert(updater, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "permission denied")
}

func (s *updaterSuite) TestCharmhubUpdate(c *gc.C) {
	ctrl := s.setupMocks(c)
	defer ctrl.Finish()
	s.expectSetCharmStoreResourcesZero()

	s.state.EXPECT().AllApplications().Return([]charmrevisionupdater.Application{
		makeApplication(ctrl, "ch", "mysql", "charm-1", "app-1", 22),
		makeApplication(ctrl, "ch", "postgresql", "charm-2", "app-2", 41),
	}, nil).AnyTimes()

	s.expectAddCharmPlaceholder("ch:mysql-23")
	s.expectAddCharmPlaceholder("ch:postgresql-42")
	s.expectRefresh(c, []string{"charm-1", "charm-2"})

	updater, err := charmrevisionupdater.NewCharmRevisionUpdaterAPIState(s.state, nil, s.newCharmHubClient)
	c.Assert(err, jc.ErrorIsNil)

	result, err := updater.UpdateLatestRevisions()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Error, gc.IsNil)
}

func (s *updaterSuite) TestCharmhubUpdateWithResources(c *gc.C) {
	ctrl := s.setupMocks(c)
	defer ctrl.Finish()
	s.expectSetCharmStoreResources(c)

	s.state.EXPECT().AllApplications().Return([]charmrevisionupdater.Application{
		makeApplication(ctrl, "ch", "resourcey", "charm-3", "app-1", 1),
	}, nil).AnyTimes()

	s.expectAddCharmPlaceholder("ch:resourcey-1")
	s.expectRefresh(c, []string{"charm-3"})

	updater, err := charmrevisionupdater.NewCharmRevisionUpdaterAPIState(s.state, nil, s.newCharmHubClient)
	c.Assert(err, jc.ErrorIsNil)

	result, err := updater.UpdateLatestRevisions()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Error, gc.IsNil)
}

func (s *updaterSuite) TestCharmhubNoUpdate(c *gc.C) {
	ctrl := s.setupMocks(c)
	defer ctrl.Finish()
	s.expectSetCharmStoreResourcesZero()

	s.state.EXPECT().AllApplications().Return([]charmrevisionupdater.Application{
		makeApplication(ctrl, "ch", "postgresql", "charm-2", "app-2", 42),
	}, nil).AnyTimes()

	s.expectAddCharmPlaceholder("ch:postgresql-42")
	s.expectRefresh(c, []string{"charm-2"})

	updater, err := charmrevisionupdater.NewCharmRevisionUpdaterAPIState(s.state, nil, s.newCharmHubClient)
	c.Assert(err, jc.ErrorIsNil)

	result, err := updater.UpdateLatestRevisions()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Error, gc.IsNil)
}

func (s *updaterSuite) newCharmHubClient(_ charmrevisionupdater.State, _ map[string]string) (charmrevisionupdater.CharmhubRefreshClient, error) {
	return s.chClient, nil
}

func (s *updaterSuite) TestCharmNotInStore(c *gc.C) {
	ctrl := s.setupMocks(c)
	defer ctrl.Finish()
	s.expectSetCharmStoreResourcesZero()

	s.state.EXPECT().AllApplications().Return([]charmrevisionupdater.Application{
		makeApplication(ctrl, "ch", "varnish", "charm-5", "app-1", 1),
		makeApplication(ctrl, "cs", "varnish", "charm-6", "app-2", 2),
	}, nil).AnyTimes()

	s.expectRefresh(c, []string{"charm-5"})

	updater, err := charmrevisionupdater.NewCharmRevisionUpdaterAPIState(s.state, newFakeCharmstoreClient, s.newCharmHubClient)
	c.Assert(err, jc.ErrorIsNil)

	result, err := updater.UpdateLatestRevisions()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Error, gc.IsNil)
}

func (s *updaterSuite) TestCharmstoreUpdate(c *gc.C) {
	ctrl := s.setupMocks(c)
	defer ctrl.Finish()
	s.expectSetCharmStoreResourcesZero()

	s.state.EXPECT().AllApplications().Return([]charmrevisionupdater.Application{
		makeApplication(ctrl, "cs", "mysql", "charm-1", "app-1", 22),
		makeApplication(ctrl, "cs", "wordpress", "charm-2", "app-2", 26),
		makeApplication(ctrl, "cs", "varnish", "charm-3", "app-3", 5), // doesn't exist in store
	}, nil)

	s.expectAddCharmPlaceholder("cs:mysql-23")
	s.expectAddCharmPlaceholder("cs:wordpress-26")

	updater, err := charmrevisionupdater.NewCharmRevisionUpdaterAPIState(s.state, newFakeCharmstoreClient, nil)
	c.Assert(err, jc.ErrorIsNil)

	result, err := updater.UpdateLatestRevisions()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Error, gc.IsNil)

	// Update mysql version and run update again.
	s.state.EXPECT().AllApplications().Return([]charmrevisionupdater.Application{
		makeApplication(ctrl, "cs", "mysql", "charm1", "app-1", 23),
	}, nil)

	s.expectAddCharmPlaceholder("cs:mysql-23")

	result, err = updater.UpdateLatestRevisions()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Error, gc.IsNil)
}

func (s *updaterSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.chClient = testmocks.NewMockCharmhubRefreshClient(ctrl)
	s.resources = statemocks.NewMockResources(ctrl)

	s.state = NewMockState(ctrl)
	s.state.EXPECT().Cloud(gomock.Any()).Return(cloud.Cloud{Type: "cloud"}, nil).AnyTimes()
	s.state.EXPECT().ControllerUUID().Return("controller-1").AnyTimes()
	s.state.EXPECT().Model().Return(makeModel(c, ctrl), nil).AnyTimes()
	s.state.EXPECT().Resources().Return(s.resources, nil).AnyTimes()

	return ctrl
}

func (s *updaterSuite) expectSetCharmStoreResourcesZero() {
	s.resources.EXPECT().SetCharmStoreResources(gomock.Any(), gomock.Len(0), gomock.Any()).Return(nil).AnyTimes()
}

func (s *updaterSuite) expectSetCharmStoreResources(c *gc.C) {
	expectedResources := []resource.Resource{
		makeResource(c, "reza", 7, 5, "59e1748777448c69de6b800d7a33bbfb9ff1b463e44354c3553bcdb9c666fa90125a3c79f90397bdf5f6a13de828684f"),
		makeResource(c, "rezb", 1, 6, "03130092073c5ac523ecb21f548b9ad6e1387d1cb05f3cb892fcc26029d01428afbe74025b6c567b6564a3168a47179a"),
	}
	s.resources.EXPECT().SetCharmStoreResources("app-1", expectedResources, gomock.Any()).Return(nil).AnyTimes()
}

func (s *updaterSuite) expectAddCharmPlaceholder(curl string) {
	s.state.EXPECT().AddCharmPlaceholder(charm.MustParseURL(curl)).Return(nil)
}

func (s *updaterSuite) expectRefresh(c *gc.C, ids []string) {
	data := getRefreshResponse()
	s.chClient.EXPECT().Refresh(gomock.Any(), refreshConfigMatcher{c: c, ids: ids}).DoAndReturn(
		func(_ context.Context, config charmhub.RefreshConfig) ([]transport.RefreshResponse, error) {
			resp := make([]transport.RefreshResponse, len(ids))
			for i, v := range ids {
				r, ok := data[v]
				if !ok {
					r = transport.RefreshResponse{
						Error: &transport.APIError{
							Code:    "not-found",
							Message: fmt.Sprintf("charm ID %q not found", v),
						},
					}
				}
				r.Entity.CreatedAt = time.Now()
				resp[i] = r
			}
			return resp, nil
		},
	)
}

type refreshConfigMatcher struct {
	c   *gc.C
	ids []string
}

func (m refreshConfigMatcher) Matches(x interface{}) bool {
	config, ok := x.(charmhub.RefreshConfig)
	if !ok {
		return false
	}
	// This is an ugly way to check that the Config contains what is expected.
	for _, v := range m.ids {
		m.c.Assert(
			strings.Replace(config.String(), "\n", "", -1),
			gc.Matches,
			fmt.Sprintf(".* %s .*", v),
			gc.Commentf("RefreshConfig did not contain id %q", v),
		)
	}

	return true
}

func (m refreshConfigMatcher) String() string {
	return fmt.Sprintf("RefreshConfig contains ids %s", strings.Join(m.ids, ", "))
}

func getRefreshResponse() map[string]transport.RefreshResponse {
	return map[string]transport.RefreshResponse{
		"charm-1": {
			Entity: transport.RefreshEntity{
				Download: transport.Download{HashSHA256: "c97e1efc5367d2fdcfdf29f4a2243b13765cc9cbdfad19627a29ac903c01ae63", Size: 5487460, URL: "https://api.staging.charmhub.io/api/v1/charms/download/jmeJLrjWpJX9OglKSeUHCwgyaCNuoQjD_208.charm"},
				ID:       "charm-1",
				Name:     "mysql",
				Summary:  "some test",
				Revision: 23,
			},
			EffectiveChannel: "latest/stable",
			Error:            (*transport.APIError)(nil),
			Name:             "mysql",
			Result:           "refresh",
		},
		"charm-2": {
			Entity: transport.RefreshEntity{
				Download: transport.Download{HashSHA256: "c97e1efc5367d2fdcfdf29f4a2243b13765cc9cbdfad19627a29ac903c01ae63", Size: 5487460, URL: "https://api.staging.charmhub.io/api/v1/charms/download/jmeJLrjWpJX9OglKSeUHCwgyaCNuoQjD_208.charm"},
				ID:       "charm-2",
				Name:     "postgresql",
				Summary:  "some test",
				Revision: 42,
			},
			EffectiveChannel: "latest/stable",
			Error:            (*transport.APIError)(nil),
			Name:             "postgresql",
			Result:           "refresh",
		},
		"charm-3": {
			Entity: transport.RefreshEntity{
				Download: transport.Download{HashSHA256: "c97e1efc5367d2fdcfdf29f4a2243b13765cc9cbdfad19627a29ac903c01ae63", Size: 5487460, URL: "https://api.staging.charmhub.io/api/v1/charms/download/jmeJLrjWpJX9OglKSeUHCwgyaCNuoQjD_208.charm"},
				ID:       "charm-3",
				Name:     "resourcey",
				Resources: []transport.ResourceRevision{
					{
						Download: transport.ResourceDownload{
							HashSHA384: "59e1748777448c69de6b800d7a33bbfb9ff1b463e44354c3553bcdb9c666fa90125a3c79f90397bdf5f6a13de828684f",
							Size:       5,
						},
						Name:     "reza",
						Revision: 7,
						Type:     "file",
					},
					{
						Download: transport.ResourceDownload{
							HashSHA384: "03130092073c5ac523ecb21f548b9ad6e1387d1cb05f3cb892fcc26029d01428afbe74025b6c567b6564a3168a47179a",
							Size:       6,
						},
						Name:     "rezb",
						Revision: 1,
						Type:     "file",
					},
				},
				Summary:  "some test",
				Revision: 1,
			},
			EffectiveChannel: "latest/stable",
			Error:            (*transport.APIError)(nil),
			Name:             "resourcey",
			Result:           "refresh",
		},
	}
}
