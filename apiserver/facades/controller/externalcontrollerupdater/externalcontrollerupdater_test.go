// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package externalcontrollerupdater_test

import (
	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common"
	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facades/controller/externalcontrollerupdater"
	"github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/rpc/params"
	coretesting "github.com/juju/juju/testing"
)

var _ = gc.Suite(&CrossControllerSuite{})

type CrossControllerSuite struct {
	coretesting.BaseSuite

	watcher             *mockStringsWatcher
	externalControllers *mockExternalControllers
	resources           *common.Resources
	auth                testing.FakeAuthorizer
}

func (s *CrossControllerSuite) SetUpTest(c *gc.C) {

	s.BaseSuite.SetUpTest(c)
	s.auth = testing.FakeAuthorizer{Controller: true}
	s.resources = common.NewResources()
	s.AddCleanup(func(*gc.C) { s.resources.StopAll() })
	s.watcher = newMockStringsWatcher()
	s.AddCleanup(func(*gc.C) { s.watcher.Stop() })
	s.externalControllers = &mockExternalControllers{
		watcher: s.watcher,
	}
}

func (s *CrossControllerSuite) TestNewAPINonController(c *gc.C) {
	s.auth.Controller = false
	_, err := externalcontrollerupdater.NewAPI(s.auth, s.resources, s.externalControllers, nil)
	c.Assert(err, gc.Equals, apiservererrors.ErrPerm)
}

func (s *CrossControllerSuite) TestExternalControllerInfo(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	ecService := NewMockEcService(ctrl)

	ctrlTag, err := names.ParseControllerTag(coretesting.ControllerTag.String())
	c.Assert(err, jc.ErrorIsNil)
	ecService.EXPECT().Controller(gomock.Any(), ctrlTag.Id()).Return(&crossmodel.ControllerInfo{
		ControllerTag: coretesting.ControllerTag,
		Alias:         "foo",
		Addrs:         []string{"bar"},
		CACert:        "baz",
	}, nil)

	modelTag, err := names.ParseControllerTag("controller-" + coretesting.ModelTag.Id())
	c.Assert(err, jc.ErrorIsNil)
	ecService.EXPECT().Controller(gomock.Any(), modelTag.Id()).Return(nil, errors.NotFoundf("external controller with UUID deadbeef-0bad-400d-8000-4b1d0d06f00d"))

	api, err := externalcontrollerupdater.NewAPI(s.auth, s.resources, s.externalControllers, ecService)
	c.Assert(err, jc.ErrorIsNil)
	results, err := api.ExternalControllerInfo(params.Entities{
		Entities: []params.Entity{
			{coretesting.ControllerTag.String()},
			{"controller-" + coretesting.ModelTag.Id()},
			{"machine-42"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.ExternalControllerInfoResults{
		[]params.ExternalControllerInfoResult{{
			Result: &params.ExternalControllerInfo{
				ControllerTag: coretesting.ControllerTag.String(),
				Alias:         "foo",
				Addrs:         []string{"bar"},
				CACert:        "baz",
			},
		}, {
			Error: &params.Error{
				Code:    "not found",
				Message: `external controller with UUID deadbeef-0bad-400d-8000-4b1d0d06f00d not found`,
			},
		}, {
			Error: &params.Error{Message: `"machine-42" is not a valid controller tag`},
		}},
	})
}

func (s *CrossControllerSuite) TestSetExternalControllerInfo(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	ecService := NewMockEcService(ctrl)

	firstControllerTag := coretesting.ControllerTag.String()
	firstControllerTagParsed, err := names.ParseControllerTag(firstControllerTag)
	c.Assert(err, jc.ErrorIsNil)
	secondControllerTag := "controller-" + coretesting.ModelTag.Id()
	secondControllerTagParsed, err := names.ParseControllerTag(secondControllerTag)
	c.Assert(err, jc.ErrorIsNil)

	ecService.EXPECT().UpdateExternalController(gomock.Any(), crossmodel.ControllerInfo{
		ControllerTag: firstControllerTagParsed,
		Alias:         "foo",
		Addrs:         []string{"bar"},
		CACert:        "baz",
	})
	ecService.EXPECT().UpdateExternalController(gomock.Any(), crossmodel.ControllerInfo{
		ControllerTag: secondControllerTagParsed,
		Alias:         "qux",
		Addrs:         []string{"quux"},
		CACert:        "quuz",
	})

	api, err := externalcontrollerupdater.NewAPI(s.auth, s.resources, s.externalControllers, ecService)
	c.Assert(err, jc.ErrorIsNil)

	results, err := api.SetExternalControllerInfo(params.SetExternalControllersInfoParams{
		[]params.SetExternalControllerInfoParams{{
			params.ExternalControllerInfo{
				ControllerTag: firstControllerTag,
				Alias:         "foo",
				Addrs:         []string{"bar"},
				CACert:        "baz",
			},
		}, {
			params.ExternalControllerInfo{
				ControllerTag: secondControllerTag,
				Alias:         "qux",
				Addrs:         []string{"quux"},
				CACert:        "quuz",
			},
		}, {
			params.ExternalControllerInfo{
				ControllerTag: "machine-42",
			},
		}},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.ErrorResults{
		[]params.ErrorResult{
			{nil},
			{nil},
			{Error: &params.Error{Message: `"machine-42" is not a valid controller tag`}},
		},
	})
}

func (s *CrossControllerSuite) TestWatchExternalControllers(c *gc.C) {
	api, err := externalcontrollerupdater.NewAPI(s.auth, s.resources, s.externalControllers, nil)
	c.Assert(err, jc.ErrorIsNil)

	s.watcher.changes <- []string{"a", "b"} // initial value
	results, err := api.WatchExternalControllers()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.StringsWatchResults{
		[]params.StringsWatchResult{{
			StringsWatcherId: "1",
			Changes:          []string{"a", "b"},
		}},
	})
	c.Assert(s.resources.Get("1"), gc.Equals, s.watcher)
}

func (s *CrossControllerSuite) TestWatchControllerInfoError(c *gc.C) {
	s.watcher.tomb.Kill(errors.New("nope"))
	close(s.watcher.changes)

	api, err := externalcontrollerupdater.NewAPI(s.auth, s.resources, s.externalControllers, nil)
	c.Assert(err, jc.ErrorIsNil)

	results, err := api.WatchExternalControllers()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.StringsWatchResults{
		[]params.StringsWatchResult{{
			Error: &params.Error{Message: "nope"},
		}},
	})
	c.Assert(s.resources.Get("1"), gc.IsNil)
}
