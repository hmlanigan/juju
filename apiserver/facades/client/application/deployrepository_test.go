// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application_test

import (
	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/facades/client/application"
	"github.com/juju/juju/apiserver/facades/client/application/mocks"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	corecharm "github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	coretesting "github.com/juju/juju/testing"
)

type platformSuite struct {
	backend *mocks.MockBackend
	machine *mocks.MockMachine
	model   *mocks.MockModel
}

var _ = gc.Suite(&platformSuite{})

func (s *platformSuite) TestDeducePlatformSimple(c *gc.C) {
	defer s.setupMocks(c).Finish()
	//model constraint default
	s.backend.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("amd64")}, nil)
	s.model.EXPECT().ModelConfig().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))

	arg := params.DeployFromRepositoryArg{CharmName: "testme"}
	api, err := s.getAPI()
	c.Assert(err, gc.IsNil)
	plat, err := application.DeducePlatformForTest(api, arg)
	c.Assert(err, gc.IsNil)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{Architecture: "amd64"})
}

func (s *platformSuite) TestDeducePlatformArgArchBase(c *gc.C) {
	defer s.setupMocks(c).Finish()

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Cons:      constraints.Value{Arch: strptr("arm64")},
		Base: &params.Base{
			Name:    "ubuntu",
			Channel: "22.10",
		},
	}
	api, err := s.getAPI()
	c.Assert(err, gc.IsNil)
	plat, err := application.DeducePlatformForTest(api, arg)
	c.Assert(err, gc.IsNil)

	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "arm64",
		OS:           "ubuntu",
		Channel:      "22.10",
	})
}

func (s *platformSuite) TestDeducePlatformModelDefaultBase(c *gc.C) {
	defer s.setupMocks(c).Finish()
	//model constraint default
	s.backend.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("amd64")}, nil)
	sConfig := coretesting.FakeConfig()
	sConfig = sConfig.Merge(coretesting.Attrs{
		"default-base": "ubuntu@22.04",
	})
	cfg, err := config.New(config.NoDefaults, sConfig)
	c.Assert(err, jc.ErrorIsNil)
	s.model.EXPECT().ModelConfig().Return(cfg, nil)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
	}
	api, err := s.getAPI()
	c.Assert(err, gc.IsNil)
	plat, err := application.DeducePlatformForTest(api, arg)
	c.Assert(err, gc.IsNil)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "amd64",
		OS:           "ubuntu",
		Channel:      "22.04/stable",
	})
}

func (s *platformSuite) TestDeducePlatformPlacementSimpleFound(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.backend.EXPECT().Machine("0").Return(s.machine, nil)
	s.machine.EXPECT().Base().Return(state.Base{
		OS:      "ubuntu",
		Channel: "18.04",
	})
	hwc := &instance.HardwareCharacteristics{Arch: strptr("arm64")}
	s.machine.EXPECT().HardwareCharacteristics().Return(hwc, nil)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Placement: []*instance.Placement{{
			Directive: "0",
		}},
	}
	api, err := s.getAPI()
	c.Assert(err, gc.IsNil)
	plat, err := application.DeducePlatformForTest(api, arg)
	c.Assert(err, gc.IsNil)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "arm64",
		OS:           "ubuntu",
		Channel:      "18.04",
	})
}

func (s *platformSuite) TestDeducePlatformPlacementSimpleNotFound(c *gc.C) {
	defer s.setupMocks(c).Finish()
	//model constraint default
	s.backend.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("amd64")}, nil)
	s.model.EXPECT().ModelConfig().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))
	s.backend.EXPECT().Machine("0/lxd/0").Return(nil, errors.NotFoundf("machine 0/lxd/0 not found"))

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Placement: []*instance.Placement{{
			Directive: "0/lxd/0",
		}},
	}
	api, err := s.getAPI()
	c.Assert(err, gc.IsNil)
	plat, err := application.DeducePlatformForTest(api, arg)
	c.Assert(err, gc.IsNil)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{Architecture: "amd64"})
}

func (s *platformSuite) TestDeducePlatformPlacementMutipleMatch(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.backend.EXPECT().Machine(gomock.Any()).Return(s.machine, nil).AnyTimes()
	s.machine.EXPECT().Base().Return(state.Base{
		OS:      "ubuntu",
		Channel: "18.04",
	}).Times(3)
	hwc := &instance.HardwareCharacteristics{Arch: strptr("arm64")}
	s.machine.EXPECT().HardwareCharacteristics().Return(hwc, nil).AnyTimes()

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Placement: []*instance.Placement{
			{Directive: "0"},
			{Directive: "1"},
			{Directive: "3"},
		},
	}
	api, err := s.getAPI()
	c.Assert(err, gc.IsNil)
	plat, err := application.DeducePlatformForTest(api, arg)
	c.Assert(err, gc.IsNil)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "arm64",
		OS:           "ubuntu",
		Channel:      "18.04",
	})
}

func (s *platformSuite) TestDeducePlatformPlacementMutipleMatchFail(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.backend.EXPECT().Machine(gomock.Any()).Return(s.machine, nil).AnyTimes()
	s.machine.EXPECT().Base().Return(
		state.Base{
			OS:      "ubuntu",
			Channel: "18.04",
		}).AnyTimes()
	gomock.InOrder(
		s.machine.EXPECT().HardwareCharacteristics().Return(
			&instance.HardwareCharacteristics{Arch: strptr("arm64")},
			nil),
		s.machine.EXPECT().HardwareCharacteristics().Return(
			&instance.HardwareCharacteristics{Arch: strptr("amd64")},
			nil),
	)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Placement: []*instance.Placement{
			{Directive: "0"},
			{Directive: "1"},
		},
	}
	api, err := s.getAPI()
	c.Assert(err, gc.IsNil)
	_, err = application.DeducePlatformForTest(api, arg)
	c.Assert(errors.Is(err, errors.BadRequest), jc.IsTrue, gc.Commentf("%+v", err))
}

func (s *platformSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.model = mocks.NewMockModel(ctrl)
	s.backend = mocks.NewMockBackend(ctrl)
	s.machine = mocks.NewMockMachine(ctrl)
	s.model.EXPECT().Type().Return(state.ModelTypeIAAS)
	return ctrl
}

func (s *platformSuite) getAPI() (*application.APIBase, error) {
	return application.NewAPIBase(
		s.backend,
		nil,
		apiservertesting.FakeAuthorizer{
			Tag: names.NewUserTag("admin"),
		},
		nil,
		nil,
		s.model,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil)
}

func strptr(s string) *string {
	return &s
}
