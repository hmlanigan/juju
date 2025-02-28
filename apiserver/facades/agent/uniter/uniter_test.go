// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/names/v6"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/model"
	corerelation "github.com/juju/juju/core/relation"
	relationtesting "github.com/juju/juju/core/relation/testing"
	domaincharm "github.com/juju/juju/domain/application/charm"
	applicationerrors "github.com/juju/juju/domain/application/errors"
	"github.com/juju/juju/domain/relation"
	"github.com/juju/juju/internal/charm"
	coretesting "github.com/juju/juju/internal/testing"
	"github.com/juju/juju/rpc/params"
)

type uniterSuite struct {
	testing.IsolationSuite

	applicationService *MockApplicationService

	uniter *UniterAPI
}

var _ = gc.Suite(&uniterSuite{})

func (s *uniterSuite) TestCharmArchiveSha256Local(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.applicationService.EXPECT().GetAvailableCharmArchiveSHA256(gomock.Any(), domaincharm.CharmLocator{
		Source:   domaincharm.LocalSource,
		Name:     "foo",
		Revision: 1,
	}).Return("sha256:foo", nil)

	results, err := s.uniter.CharmArchiveSha256(context.Background(), params.CharmURLs{
		URLs: []params.CharmURL{
			{URL: "local:foo-1"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.StringResults{
		Results: []params.StringResult{{
			Result: "sha256:foo",
		}},
	})
}

func (s *uniterSuite) TestCharmArchiveSha256Charmhub(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.applicationService.EXPECT().GetAvailableCharmArchiveSHA256(gomock.Any(), domaincharm.CharmLocator{
		Source:   domaincharm.CharmHubSource,
		Name:     "foo",
		Revision: 1,
	}).Return("sha256:foo", nil)

	results, err := s.uniter.CharmArchiveSha256(context.Background(), params.CharmURLs{
		URLs: []params.CharmURL{
			{URL: "foo-1"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.StringResults{
		Results: []params.StringResult{{
			Result: "sha256:foo",
		}},
	})
}

func (s *uniterSuite) TestCharmArchiveSha256Errors(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.applicationService.EXPECT().GetAvailableCharmArchiveSHA256(gomock.Any(), domaincharm.CharmLocator{
		Source:   domaincharm.CharmHubSource,
		Name:     "foo",
		Revision: 1,
	}).Return("", applicationerrors.CharmNotFound)
	s.applicationService.EXPECT().GetAvailableCharmArchiveSHA256(gomock.Any(), domaincharm.CharmLocator{
		Source:   domaincharm.CharmHubSource,
		Name:     "foo",
		Revision: 2,
	}).Return("", applicationerrors.CharmNotFound)
	s.applicationService.EXPECT().GetAvailableCharmArchiveSHA256(gomock.Any(), domaincharm.CharmLocator{
		Source:   domaincharm.CharmHubSource,
		Name:     "foo",
		Revision: 3,
	}).Return("", applicationerrors.CharmNotResolved)

	results, err := s.uniter.CharmArchiveSha256(context.Background(), params.CharmURLs{
		URLs: []params.CharmURL{
			{URL: "foo-1"},
			{URL: "ch:foo-2"},
			{URL: "ch:foo-3"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.StringResults{
		Results: []params.StringResult{
			{Error: &params.Error{Message: `charm "foo-1" not found`, Code: params.CodeNotFound}},
			{Error: &params.Error{Message: `charm "ch:foo-2" not found`, Code: params.CodeNotFound}},
			{Error: &params.Error{Message: `charm "ch:foo-3" not available`, Code: params.CodeNotYetAvailable}},
		},
	})
}

func (s *uniterSuite) TestLeadershipSettings(c *gc.C) {
	defer s.setupMocks(c).Finish()

	s.uniter.Merge(context.Background(), struct{}{}, struct{}{})
	s.uniter.Read(context.Background(), struct{}{}, struct{}{})
	s.uniter.WatchLeadershipSettings(context.Background(), struct{}{}, struct{}{})
}

func (s *uniterSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.applicationService = NewMockApplicationService(ctrl)

	s.uniter = &UniterAPI{
		applicationService: s.applicationService,
	}

	return ctrl
}

type leadershipSettings interface {
	// Merge merges in the provided leadership settings. Only leaders for
	// the given service may perform this operation.
	Merge(ctx context.Context, bulkArgs params.MergeLeadershipSettingsBulkParams) (params.ErrorResults, error)

	// Read reads leadership settings for the provided service ID. Any
	// unit of the service may perform this operation.
	Read(ctx context.Context, bulkArgs params.Entities) (params.GetLeadershipSettingsBulkResults, error)

	// WatchLeadershipSettings will block the caller until leadership settings
	// for the given service ID change.
	WatchLeadershipSettings(ctx context.Context, bulkArgs params.Entities) (params.NotifyWatchResults, error)
}

type leadershipUniterSuite struct {
	testing.IsolationSuite

	watcherRegistry *MockWatcherRegistry

	uniter leadershipSettings

	setupMocks func(c *gc.C) *gomock.Controller
}

func (s *leadershipUniterSuite) TestLeadershipSettingsMerge(c *gc.C) {
	defer s.setupMocks(c).Finish()

	results, err := s.uniter.Merge(context.Background(), params.MergeLeadershipSettingsBulkParams{
		Params: []params.MergeLeadershipSettingsParam{
			{
				ApplicationTag: "app1",
				Settings: params.Settings{
					"key1": "value1",
				},
			},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.ErrorResults{
		Results: []params.ErrorResult{{}},
	})
}

func (s *leadershipUniterSuite) TestLeadershipSettingsRead(c *gc.C) {
	defer s.setupMocks(c).Finish()

	results, err := s.uniter.Read(context.Background(), params.Entities{
		Entities: []params.Entity{
			{
				Tag: "app1",
			},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.GetLeadershipSettingsBulkResults{
		Results: []params.GetLeadershipSettingsResult{{}},
	})
}

func (s *leadershipUniterSuite) TestLeadershipSettingsWatchLeadershipSettings(c *gc.C) {
	defer s.setupMocks(c).Finish()

	results, err := s.uniter.WatchLeadershipSettings(context.Background(), params.Entities{
		Entities: []params.Entity{
			{
				Tag: "app1",
			},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.NotifyWatchResults{
		Results: []params.NotifyWatchResult{{
			NotifyWatcherId: "watcher1",
		}},
	})
}

type uniterv19Suite struct {
	leadershipUniterSuite
}

var _ = gc.Suite(&uniterv19Suite{})

func (s *uniterv19Suite) SetUpTest(c *gc.C) {
	s.setupMocks = func(c *gc.C) *gomock.Controller {
		ctrl := gomock.NewController(c)

		s.watcherRegistry = NewMockWatcherRegistry(ctrl)
		s.watcherRegistry.EXPECT().Register(gomock.Any()).Return("watcher1", nil).AnyTimes()

		s.uniter = &UniterAPIv19{
			UniterAPIv20: &UniterAPIv20{
				UniterAPI: &UniterAPI{
					watcherRegistry: s.watcherRegistry,
				},
			},
		}

		return ctrl
	}
}

type uniterv20Suite struct {
	leadershipUniterSuite
}

var _ = gc.Suite(&uniterv20Suite{})

func (s *uniterv20Suite) SetUpTest(c *gc.C) {
	s.setupMocks = func(c *gc.C) *gomock.Controller {
		ctrl := gomock.NewController(c)

		s.watcherRegistry = NewMockWatcherRegistry(ctrl)
		s.watcherRegistry.EXPECT().Register(gomock.Any()).Return("watcher1", nil).AnyTimes()

		s.uniter = &UniterAPIv20{
			UniterAPI: &UniterAPI{
				watcherRegistry: s.watcherRegistry,
			},
		}

		return ctrl
	}
}

type uniterRelationSuite struct {
	testing.IsolationSuite

	//leadershipChecker *fakeLeadershipChecker
	unitTag1 names.UnitTag

	applicationService *MockApplicationService
	modelInfoService   *MockModelInfoService
	relationService    *MockRelationService

	uniter *UniterAPI
}

var _ = gc.Suite(&uniterRelationSuite{})

func (s *uniterRelationSuite) SetUpTest(c *gc.C) {
	s.unitTag1 = names.NewUnitTag("wordpress/0")
}

func (s *uniterRelationSuite) TestRelation(c *gc.C) {
	defer s.setupMocks(c).Finish()
	relTag := names.NewRelationTag("mysql:database wordpress:mysql")
	relUUID := relationtesting.GenRelationUUID(c)
	relID := 42

	s.expectModelUUID(c)
	s.expectGetRelationUUIDFromKey(corerelation.Key(relTag.Id()), relUUID)
	s.expectGetRelationDetailsForUnit1(relUUID, relID, relTag)

	args := params.RelationUnits{RelationUnits: []params.RelationUnit{
		{Relation: "relation-42", Unit: "unit-foo-0"},
		{Relation: relTag.String(), Unit: "unit-wordpress-0"},
		{Relation: relTag.String(), Unit: "unit-mysql-0"},
		{Relation: relTag.String(), Unit: "unit-foo-0"},
		{Relation: "relation-blah", Unit: "unit-wordpress-0"},
		{Relation: "application-foo", Unit: "user-foo"},
		{Relation: "foo", Unit: "bar"},
		{Relation: "unit-wordpress-0", Unit: relTag.String()},
	}}
	result, err := s.uniter.Relation(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(result, gc.DeepEquals, params.RelationResultsV2{
		Results: []params.RelationResultV2{
			{Error: apiservertesting.ErrUnauthorized},
			{
				Id:   relID,
				Key:  relTag.Id(),
				Life: life.Alive,
				Endpoint: params.Endpoint{
					ApplicationName: "wordpress",
					Relation: params.CharmRelation{
						Name:      "database",
						Role:      string(charm.RoleRequirer),
						Interface: "mysql",
						Optional:  false,
						Limit:     0,
						Scope:     string(charm.ScopeGlobal),
					},
				},
				OtherApplication: params.RelatedApplicationDetails{
					ApplicationName: "mysql",
					ModelUUID:       coretesting.ModelTag.Id(),
				},
			},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
			{Error: apiservertesting.ErrUnauthorized},
		},
	})
}

func (s *uniterRelationSuite) TestRelationById(c *gc.C) {
	defer s.setupMocks(c).Finish()
	relTag := names.NewRelationTag("mysql:database wordpress:mysql")
	relUUID := relationtesting.GenRelationUUID(c)
	relID := 37
	relIDNotFound := -1
	relIDUnexpectedAppName := 42

	s.expectModelUUID(c)
	s.expectGetRelationDetails(relUUID, relID, relTag)
	s.expectGetRelationDetailsNotFound(relIDNotFound)
	s.expectGetRelationDetailsUnexpectedAppName(c, relIDUnexpectedAppName)
	args := params.RelationIds{
		RelationIds: []int{relIDNotFound, relID, relIDUnexpectedAppName},
	}
	result, err := s.uniter.RelationById(context.Background(), args)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, gc.DeepEquals, params.RelationResultsV2{
		Results: []params.RelationResultV2{
			// The relation ID does not exist: ErrUnauthorized.
			{Error: apiservertesting.ErrUnauthorized},
			// Successful result
			{
				Id:   relID,
				Key:  relTag.Id(),
				Life: life.Alive,
				Endpoint: params.Endpoint{
					ApplicationName: "wordpress",
					Relation: params.CharmRelation{
						Name:      "database",
						Role:      string(charm.RoleRequirer),
						Interface: "mysql",
						Optional:  false,
						Limit:     0,
						Scope:     string(charm.ScopeGlobal),
					},
				},
				OtherApplication: params.RelatedApplicationDetails{
					ApplicationName: "mysql",
					ModelUUID:       coretesting.ModelTag.Id(),
				},
			},
			// The auth application is not part of the relation: ErrUnauthorized.
			{Error: apiservertesting.ErrUnauthorized},
		},
	})
}

func (s *uniterRelationSuite) TestReadSettings(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadSettingsForApplicationWhenNotLeader(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadSettingsForApplicationInPeerRelation(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadLocalApplicationSettingsWhenNotLeader(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadLocalApplicationSettingsInPeerRelation(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadRemoteSettings(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadRemoteSettingsForApplication(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadRemoteApplicationSettingsForPeerRelation(c *gc.C) {

}

func (s *uniterRelationSuite) TestReadRemoteSettingsForCAASApplicationInPeerRelationSidecar(c *gc.C) {
}

func (s *uniterRelationSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.applicationService = NewMockApplicationService(ctrl)
	s.modelInfoService = NewMockModelInfoService(ctrl)
	s.relationService = NewMockRelationService(ctrl)

	unitAuthFunc := func() (common.AuthFunc, error) {
		return func(tag names.Tag) bool {
			return tag.Id() == s.unitTag1.Id()
		}, nil
	}

	authorizer := &apiservertesting.FakeAuthorizer{
		Tag:        s.unitTag1,
		Controller: true,
	}

	s.uniter = &UniterAPI{
		accessUnit: unitAuthFunc,
		auth:       authorizer,

		applicationService: s.applicationService,
		modelInfoService:   s.modelInfoService,
		relationService:    s.relationService,
	}

	return ctrl
}

func (s *uniterRelationSuite) expectModelUUID(c *gc.C) {
	modelInfo := model.ModelInfo{
		UUID: model.UUID(coretesting.ModelTag.Id()),
	}
	s.modelInfoService.EXPECT().GetModelInfo(gomock.Any()).Return(modelInfo, nil)
}

func (s *uniterRelationSuite) expectGetRelationUUIDFromKey(key corerelation.Key, relUUID corerelation.UUID) {
	s.relationService.EXPECT().GetRelationUUIDFromKey(gomock.Any(), key).Return(relUUID, nil)
}

func (s *uniterRelationSuite) expectGetRelationDetails(relUUID corerelation.UUID, relID int, relTag names.RelationTag) {
	s.relationService.EXPECT().GetRelationDetails(gomock.Any(), relID).Return(relation.RelationDetails{
		Life: life.Alive,
		UUID: relUUID,
		ID:   relID,
		Key:  relTag.Id(),
		Endpoint: []relation.Endpoint{
			{
				ApplicationName: "wordpress",
				Relation: charm.Relation{
					Name:      "database",
					Role:      charm.RoleRequirer,
					Interface: "mysql",
					Scope:     charm.ScopeGlobal,
				},
			},
			{
				ApplicationName: "mysql",
				Relation: charm.Relation{
					Name:      "mysql",
					Role:      charm.RoleProvider,
					Interface: "mysql",
					Scope:     charm.ScopeGlobal,
				},
			},
		},
	}, nil)
}

func (s *uniterRelationSuite) expectGetRelationDetailsForUnit1(relUUID corerelation.UUID, relID int, relTag names.RelationTag) {
	s.relationService.EXPECT().GetRelationDetailsForUnit(gomock.Any(), relUUID, s.unitTag1.Id()).Return(
		relation.RelationDetails{
			Life: life.Alive,
			UUID: relUUID,
			ID:   relID,
			Key:  relTag.Id(),
			Endpoint: []relation.Endpoint{
				{
					ApplicationName: "wordpress",
					Relation: charm.Relation{
						Name:      "database",
						Role:      charm.RoleRequirer,
						Interface: "mysql",
						Scope:     charm.ScopeGlobal,
					},
				},
				{
					ApplicationName: "mysql",
					Relation: charm.Relation{
						Name:      "mysql",
						Role:      charm.RoleProvider,
						Interface: "mysql",
						Scope:     charm.ScopeGlobal,
					},
				},
			},
		}, nil)

}

func (s *uniterRelationSuite) expectGetRelationDetailsNotFound(relID int) {
	s.relationService.EXPECT().GetRelationDetails(gomock.Any(), relID).Return(relation.RelationDetails{}, errors.NotFound)
}

func (s *uniterRelationSuite) expectGetRelationDetailsUnexpectedAppName(c *gc.C, relID int) {
	s.relationService.EXPECT().GetRelationDetails(gomock.Any(), relID).Return(relation.RelationDetails{
		Life: life.Alive,
		UUID: relationtesting.GenRelationUUID(c),
		ID:   relID,
		Endpoint: []relation.Endpoint{
			{
				ApplicationName: "failure-application",
				Relation: charm.Relation{
					Name:      "database",
					Role:      charm.RoleRequirer,
					Interface: "mysql",
					Scope:     charm.ScopeGlobal,
				},
			},
			{
				ApplicationName: "mysql",
				Relation: charm.Relation{
					Name:      "mysql",
					Role:      charm.RoleProvider,
					Interface: "mysql",
					Scope:     charm.ScopeGlobal,
				},
			},
		},
	}, nil)
}
