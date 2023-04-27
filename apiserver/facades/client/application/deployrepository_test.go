// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/juju/charm/v10"
	"github.com/juju/charm/v10/resource"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v3/hash"
	"github.com/kr/pretty"
	gc "gopkg.in/check.v1"

	corecharm "github.com/juju/juju/core/charm"
	coreconfig "github.com/juju/juju/core/config"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	coretesting "github.com/juju/juju/testing"
)

type validatorSuite struct {
	bindings    *MockBindings
	machine     *MockMachine
	model       *MockModel
	repo        *MockRepository
	repoFactory *MockRepositoryFactory
	state       *MockDeployFromRepositoryState
}

var _ = gc.Suite(&deployRepositorySuite{})
var _ = gc.Suite(&validatorSuite{})

func (s *validatorSuite) TestValidateSuccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()
	resultURL, essMeta := s.expectResolveForDeploy(c)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
	}
	dt, errs := s.getValidator().validate(arg)
	c.Assert(errs, gc.HasLen, 0)
	c.Assert(dt, gc.DeepEquals, deployTemplate{
		applicationName: "test-charm",
		charm:           corecharm.NewCharmInfoAdapter(essMeta),
		charmURL:        resultURL,
		numUnits:        1,
		origin:          essMeta.ResolvedOrigin,
	})
}
func (s *validatorSuite) TestValidateIAASAttachStorageFail(c *gc.C) {
	argStorageNames := []string{"one-0"}
	expectedStorageTags := []names.StorageTag{}
	s.testValidateIAASAttachStorage(c, argStorageNames, expectedStorageTags, errors.NotValid)
}

func (s *validatorSuite) TestValidateIAASAttachStorageSuccess(c *gc.C) {
	argStorageNames := []string{"one/0", "two/3"}
	expectedStorageTags := []names.StorageTag{names.NewStorageTag("one/0"), names.NewStorageTag("two/3")}
	s.testValidateIAASAttachStorage(c, argStorageNames, expectedStorageTags, "")
}

func (s *validatorSuite) testValidateIAASAttachStorage(c *gc.C, argStorage []string, expectedStorageTags []names.StorageTag, expectedErr errors.ConstError) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()
	resultURL, essMeta := s.expectResolveForDeploy(c)

	arg := params.DeployFromRepositoryArg{
		CharmName:     "testcharm",
		AttachStorage: argStorage,
	}
	dt, errs := s.iaasDeployFromRepositoryValidator().ValidateArg(arg)
	if expectedErr == "" {
		c.Assert(errs, gc.HasLen, 0)
		c.Assert(dt, gc.DeepEquals, deployTemplate{
			applicationName: "test-charm",
			charm:           corecharm.NewCharmInfoAdapter(essMeta),
			charmURL:        resultURL,
			numUnits:        1,
			origin:          essMeta.ResolvedOrigin,
			attachStorage:   expectedStorageTags,
		})
	} else {
		c.Assert(errs, gc.HasLen, 1)
		c.Assert(errors.Is(errs[0], expectedErr), jc.IsTrue)
	}
}

func (s *validatorSuite) TestValidatePlacementSuccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()
	resultURL, essMeta := s.expectResolveForDeployWithPlacement(c)

	// Placement
	s.state.EXPECT().Machine("0").Return(s.machine, nil).Times(2)
	s.machine.EXPECT().IsLockedForSeriesUpgrade().Return(false, nil)
	s.machine.EXPECT().IsParentLockedForSeriesUpgrade().Return(false, nil)
	s.machine.EXPECT().Base().Return(state.Base{
		OS:      "ubuntu",
		Channel: "22.04",
	})
	hwc := &instance.HardwareCharacteristics{Arch: strptr("amd64")}
	s.machine.EXPECT().HardwareCharacteristics().Return(hwc, nil)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
		Placement: []*instance.Placement{{Directive: "0", Scope: instance.MachineScope}},
	}
	dt, errs := s.getValidator().validate(arg)
	c.Assert(errs, gc.HasLen, 0)
	c.Assert(dt, gc.DeepEquals, deployTemplate{
		applicationName: "test-charm",
		charm:           corecharm.NewCharmInfoAdapter(essMeta),
		charmURL:        resultURL,
		numUnits:        1,
		origin:          essMeta.ResolvedOrigin,
		placement:       arg.Placement,
	})
}

func (s *validatorSuite) TestValidateEndpointBindingSuccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()
	resultURL, essMeta := s.expectResolveForDeploy(c)

	// state bindings
	endpointMap := map[string]string{"to": "from"}
	s.bindings.EXPECT().Map().Return(endpointMap)

	arg := params.DeployFromRepositoryArg{
		CharmName:        "testcharm",
		EndpointBindings: endpointMap,
	}
	dt, errs := s.getValidator().validate(arg)
	c.Assert(errs, gc.HasLen, 0)
	c.Assert(dt, gc.DeepEquals, deployTemplate{
		applicationName: "test-charm",
		charm:           corecharm.NewCharmInfoAdapter(essMeta),
		charmURL:        resultURL,
		endpoints:       endpointMap,
		numUnits:        1,
		origin:          essMeta.ResolvedOrigin,
	})
}

func (s *validatorSuite) TestResolveCharm(c *gc.C) {
	defer s.setupMocks(c).Finish()
	curl := charm.MustParseURL("testcharm")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	resArgs := corecharm.ResolveForDeployArg{
		URL:    curl,
		Origin: origin,
	}

	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	resolvedOrigin := corecharm.Origin{
		Source:   "charm-hub",
		Type:     "charm",
		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
		Revision: intptr(4),
	}
	deployData := corecharm.ResolvedDataForDeploy{
		URL:               resultURL,
		EssentialMetadata: corecharm.EssentialMetadata{ResolvedOrigin: resolvedOrigin},
	}
	s.repo.EXPECT().ResolveForDeploy(resolveForDeployMatcher{c: c, expectedArgs: resArgs}).Return(deployData, nil)

	obtained, err := s.getValidator().resolveCharm(curl, origin, false, false, constraints.Value{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(obtained.URL, gc.DeepEquals, resultURL)
	c.Assert(obtained.EssentialMetadata.ResolvedOrigin, gc.DeepEquals, resolvedOrigin)
}

func (s *validatorSuite) TestBaseSelectionArchAll(c *gc.C) {
	defer s.setupMocks(c).Finish()
	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	resolvedOrigin := corecharm.Origin{
		Source:   "charm-hub",
		Type:     "charm",
		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "all", OS: "ubuntu", Channel: "22.04"},
		Revision: intptr(4),
	}
	validator := s.getValidator()
	baseSelectionFunc := validator.createBaseSelectionFunc(resultURL, false, false, constraints.Value{})

	supportedSeries := []string{"jammy", "focal"}
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("arm64")}, nil)

	selectionOrigin, err := baseSelectionFunc(origin, resolvedOrigin, supportedSeries)
	c.Assert(err, jc.ErrorIsNil)
	expectedOrigin := resolvedOrigin
	expectedOrigin.Platform.Architecture = "arm64"
	c.Assert(selectionOrigin, gc.DeepEquals, expectedOrigin)
}

func (s *validatorSuite) TestBaseSelectionUnsupportedSeriesErrorForce(c *gc.C) {
	defer s.setupMocks(c).Finish()
	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	origin := corecharm.Origin{
		Source:  "charm-hub",
		Channel: &charm.Channel{Risk: "stable"},
		//Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	resolvedOrigin := corecharm.Origin{
		Source:   "charm-hub",
		Type:     "charm",
		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
		Revision: intptr(4),
	}
	supportedSeries := []string{"focal"}
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)

	validator := s.getValidator()
	baseSelectionFunc := validator.createBaseSelectionFunc(resultURL, true, false, constraints.Value{})

	selectionOrigin, err := baseSelectionFunc(origin, resolvedOrigin, supportedSeries)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(selectionOrigin, gc.DeepEquals, resolvedOrigin)
}

func (s *validatorSuite) TestBaseSelectionUnsupportedSeriesError(c *gc.C) {
	defer s.setupMocks(c).Finish()
	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
	}
	resolvedOrigin := corecharm.Origin{
		Source:   "charm-hub",
		Type:     "charm",
		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
		Revision: intptr(4),
	}
	supportedSeries := []string{"focal"}
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)

	validator := s.getValidator()
	baseSelectionFunc := validator.createBaseSelectionFunc(resultURL, false, false, constraints.Value{})

	_, err := baseSelectionFunc(origin, resolvedOrigin, supportedSeries)
	//c.Assert(err, jc.ErrorIsNil)
	//c.Assert(selectionOrigin, gc.DeepEquals, resolvedOrigin)
	//	defer s.setupMocks(c).Finish()
	//	curl := charm.MustParseURL("testcharm")
	//	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	//	origin := corecharm.Origin{
	//		Source:   "charm-hub",
	//		Channel:  &charm.Channel{Risk: "stable"},
	//		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
	//	}
	//	resolvedOrigin := corecharm.Origin{
	//		Source:   "charm-hub",
	//		Type:     "charm",
	//		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
	//		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
	//		Revision: intptr(4),
	//	}
	//	supportedSeries := []string{"focal"}
	//	newErr := charm.NewUnsupportedSeriesError("jammy", supportedSeries)
	//	s.repo.EXPECT().ResolveWithPreferredChannel(curl, origin).Return(resultURL, resolvedOrigin, supportedSeries, newErr)
	//
	//	_, err := s.getValidator().resolveCharm(curl, origin, false, false, constraints.Value{})
	c.Assert(err, gc.ErrorMatches, `series "jammy" not supported by charm, supported series are: focal. Use --force to deploy the charm anyway.`)
}

// Is this test in the correct place?
//func (s *validatorSuite) TestResolveCharmExplicitBaseErrorWhenUserImageID(c *gc.C) {
//	defer s.setupMocks(c).Finish()
//	curl := charm.MustParseURL("testcharm")
//	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
//	origin := corecharm.Origin{
//		Source:   "charm-hub",
//		Channel:  &charm.Channel{Risk: "stable"},
//		Platform: corecharm.Platform{Architecture: "amd64"},
//	}
//	resolvedOrigin := corecharm.Origin{
//		Source:   "charm-hub",
//		Type:     "charm",
//		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
//		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04/stable"},
//		Revision: intptr(4),
//	}
//	supportedSeries := []string{"jammy", "focal"}
//	s.repo.EXPECT().ResolveWithPreferredChannel(curl, origin).Return(resultURL, resolvedOrigin, supportedSeries, nil)
//	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))
//	s.state.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("arm64")}, nil)
//
//	_, err := s.getValidator().resolveCharm(curl, origin, false, false, constraints.Value{ImageID: strptr("ubuntu-bf2")})
//	c.Assert(err, gc.ErrorMatches, `base must be explicitly provided when image-id constraint is used`)
//}

// Is this test in the correct place?
//func (s *validatorSuite) TestResolveCharmExplicitBaseErrorWhenModelImageID(c *gc.C) {
//	defer s.setupMocks(c).Finish()
//	curl := charm.MustParseURL("testcharm")
//	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
//	origin := corecharm.Origin{
//		Source:   "charm-hub",
//		Channel:  &charm.Channel{Risk: "stable"},
//		Platform: corecharm.Platform{Architecture: "amd64"},
//	}
//	resolvedOrigin := corecharm.Origin{
//		Source:   "charm-hub",
//		Type:     "charm",
//		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
//		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04/stable"},
//		Revision: intptr(4),
//	}
//	supportedSeries := []string{"jammy", "focal"}
//	s.repo.EXPECT().ResolveWithPreferredChannel(curl, origin).Return(resultURL, resolvedOrigin, supportedSeries, nil)
//	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))
//	s.state.EXPECT().ModelConstraints().Return(constraints.Value{
//		Arch:    strptr("arm64"),
//		ImageID: strptr("ubuntu-bf2"),
//	}, nil)
//
//	_, err := s.getValidator().resolveCharm(curl, origin, false, false, constraints.Value{})
//	c.Assert(err, gc.ErrorMatches, `base must be explicitly provided when image-id constraint is used`)
//}

func (s *validatorSuite) TestCreateOrigin(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
		Revision:  intptr(7),
	}
	curl, origin, defaultBase, err := s.getValidator().createOrigin(arg)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(curl, gc.DeepEquals, charm.MustParseURL("ch:testcharm-7"))
	c.Assert(origin, gc.DeepEquals, corecharm.Origin{
		Source:   "charm-hub",
		Revision: intptr(7),
		Channel:  &corecharm.DefaultChannel,
		Platform: corecharm.Platform{Architecture: "amd64"},
	})
	c.Assert(defaultBase, jc.IsFalse)
}

func (s *validatorSuite) TestCreateOriginChannel(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
		Revision:  intptr(7),
		Channel:   strptr("yoga/candidate"),
	}
	curl, origin, defaultBase, err := s.getValidator().createOrigin(arg)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(curl, gc.DeepEquals, charm.MustParseURL("ch:testcharm-7"))
	expectedChannel := corecharm.MustParseChannel("yoga/candidate")
	c.Assert(origin, gc.DeepEquals, corecharm.Origin{
		Source:   "charm-hub",
		Revision: intptr(7),
		Channel:  &expectedChannel,
		Platform: corecharm.Platform{Architecture: "amd64"},
	})
	c.Assert(defaultBase, jc.IsFalse)
}

func (s *validatorSuite) TestGetCharm(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()
	resultURL, essMeta := s.expectResolveForDeploy(c)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
	}
	obtainedData, obtainedCharm, err := s.getValidator().getCharm(arg)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(obtainedData.EssentialMetadata.ResolvedOrigin, gc.DeepEquals, essMeta.ResolvedOrigin)
	c.Assert(obtainedCharm, gc.DeepEquals, corecharm.NewCharmInfoAdapter(essMeta))
	c.Assert(obtainedData.URL.String(), gc.Equals, resultURL.String())
}

func (s *validatorSuite) TestGetCharmAlreadyDeployed(c *gc.C) {
	ctrl := s.setupMocks(c)
	defer ctrl.Finish()
	s.expectSimpleValidate()
	curl := charm.MustParseURL("testcharm")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	resArgs := corecharm.ResolveForDeployArg{
		URL:    curl,
		Origin: origin,
	}

	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	resolvedOrigin := corecharm.Origin{
		Source:   "charm-hub",
		Type:     "charm",
		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
		Revision: intptr(4),
	}
	expMeta := &charm.Meta{
		Name: "test-charm",
	}
	expManifest := new(charm.Manifest)
	expConfig := new(charm.Config)
	essMeta := corecharm.EssentialMetadata{
		Meta:           expMeta,
		Manifest:       expManifest,
		Config:         expConfig,
		ResolvedOrigin: resolvedOrigin,
	}
	deployData := corecharm.ResolvedDataForDeploy{
		URL:               resultURL,
		EssentialMetadata: essMeta,
	}

	s.repo.EXPECT().ResolveForDeploy(resolveForDeployMatcher{c: c, expectedArgs: resArgs}).Return(deployData, nil)
	ch := NewMockCharm(ctrl)
	s.state.EXPECT().Charm(gomock.Any()).Return(ch, nil)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
	}
	obtainedData, obtainedCharm, err := s.getValidator().getCharm(arg)

	c.Assert(err, jc.ErrorIsNil)
	c.Assert(obtainedData.EssentialMetadata.ResolvedOrigin, gc.DeepEquals, essMeta.ResolvedOrigin)
	c.Assert(obtainedCharm, gc.NotNil)
	c.Assert(obtainedData.URL.String(), gc.Equals, resultURL.String())
}

func (s *validatorSuite) TestGetCharmFindsBundle(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()
	// resolveCharm
	curl := charm.MustParseURL("testcharm")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	resArgs := corecharm.ResolveForDeployArg{
		URL:    curl,
		Origin: origin,
	}

	resolvedOrigin := corecharm.Origin{
		Type: "bundle",
	}
	deployData := corecharm.ResolvedDataForDeploy{
		EssentialMetadata: corecharm.EssentialMetadata{ResolvedOrigin: resolvedOrigin},
	}

	s.repo.EXPECT().ResolveForDeploy(resolveForDeployMatcher{c: c, expectedArgs: resArgs}).Return(deployData, nil)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
	}
	_, _, err := s.getValidator().getCharm(arg)
	c.Assert(errors.Is(err, errors.BadRequest), jc.IsTrue)
}

func (s *validatorSuite) TestGetCharmNoJujuControllerCharm(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()

	curl := charm.MustParseURL("testcharm")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	resArgs := corecharm.ResolveForDeployArg{
		URL:    curl,
		Origin: origin,
	}

	resolvedOrigin := corecharm.Origin{
		Type: "charm",
	}
	expMeta := &charm.Meta{
		Name: "juju-controller",
	}
	essMeta := corecharm.EssentialMetadata{
		Meta:           expMeta,
		ResolvedOrigin: resolvedOrigin,
	}
	deployData := corecharm.ResolvedDataForDeploy{
		URL:               charm.MustParseURL("ch:amd64/jammy/juju-qa-test-4"),
		EssentialMetadata: essMeta,
	}
	s.repo.EXPECT().ResolveForDeploy(resolveForDeployMatcher{c: c, expectedArgs: resArgs}).Return(deployData, nil)
	s.state.EXPECT().Charm(gomock.Any()).Return(nil, errors.NotFoundf("charm"))

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
	}
	_, _, err := s.getValidator().getCharm(arg)
	c.Assert(errors.Is(err, errors.NotSupported), jc.IsTrue, gc.Commentf("%+v", err))
}

func (s *validatorSuite) TestDeducePlatformSimple(c *gc.C) {
	defer s.setupMocks(c).Finish()
	//model constraint default
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("amd64")}, nil)
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))

	arg := params.DeployFromRepositoryArg{CharmName: "testme"}
	plat, usedModelDefaultBase, err := s.getValidator().deducePlatform(arg)
	c.Assert(err, gc.IsNil)
	c.Assert(usedModelDefaultBase, jc.IsFalse)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{Architecture: "amd64"})
}

func (s *validatorSuite) TestDeducePlatformArgArchBase(c *gc.C) {
	defer s.setupMocks(c).Finish()

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Cons:      constraints.Value{Arch: strptr("arm64")},
		Base: &params.Base{
			Name:    "ubuntu",
			Channel: "22.10",
		},
	}
	plat, usedModelDefaultBase, err := s.getValidator().deducePlatform(arg)
	c.Assert(err, gc.IsNil)
	c.Assert(usedModelDefaultBase, jc.IsFalse)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "arm64",
		OS:           "ubuntu",
		Channel:      "22.10",
	})
}

func (s *validatorSuite) TestDeducePlatformModelDefaultBase(c *gc.C) {
	defer s.setupMocks(c).Finish()
	//model constraint default
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)
	sConfig := coretesting.FakeConfig()
	sConfig = sConfig.Merge(coretesting.Attrs{
		"default-base": "ubuntu@22.04",
	})
	cfg, err := config.New(config.NoDefaults, sConfig)
	c.Assert(err, jc.ErrorIsNil)
	s.model.EXPECT().Config().Return(cfg, nil)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
	}
	plat, usedModelDefaultBase, err := s.getValidator().deducePlatform(arg)
	c.Assert(err, gc.IsNil)
	c.Assert(usedModelDefaultBase, jc.IsTrue)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "amd64",
		OS:           "ubuntu",
		Channel:      "22.04/stable",
	})
}

func (s *validatorSuite) TestDeducePlatformPlacementSimpleFound(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)
	s.state.EXPECT().Machine("0").Return(s.machine, nil)
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
	plat, usedModelDefaultBase, err := s.getValidator().deducePlatform(arg)
	c.Assert(err, gc.IsNil)
	c.Assert(usedModelDefaultBase, jc.IsFalse)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "arm64",
		OS:           "ubuntu",
		Channel:      "18.04",
	})
}

func (s *validatorSuite) TestDeducePlatformPlacementSimpleNotFound(c *gc.C) {
	defer s.setupMocks(c).Finish()
	//model constraint default
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("amd64")}, nil)
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig()))
	s.state.EXPECT().Machine("0/lxd/0").Return(nil, errors.NotFoundf("machine 0/lxd/0 not found"))

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Placement: []*instance.Placement{{
			Directive: "0/lxd/0",
		}},
	}
	plat, usedModelDefaultBase, err := s.getValidator().deducePlatform(arg)
	c.Assert(err, gc.IsNil)
	c.Assert(usedModelDefaultBase, jc.IsFalse)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{Architecture: "amd64"})
}

func (s *validatorSuite) TestDeducePlatformPlacementMutipleMatch(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)
	s.state.EXPECT().Machine(gomock.Any()).Return(s.machine, nil).Times(3)
	s.machine.EXPECT().Base().Return(state.Base{
		OS:      "ubuntu",
		Channel: "18.04",
	}).Times(3)
	hwc := &instance.HardwareCharacteristics{Arch: strptr("arm64")}
	s.machine.EXPECT().HardwareCharacteristics().Return(hwc, nil).Times(3)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
		Placement: []*instance.Placement{
			{Directive: "0"},
			{Directive: "1"},
			{Directive: "3"},
		},
	}
	plat, usedModelDefaultBase, err := s.getValidator().deducePlatform(arg)
	c.Assert(err, gc.IsNil)
	c.Assert(usedModelDefaultBase, jc.IsFalse)
	c.Assert(plat, gc.DeepEquals, corecharm.Platform{
		Architecture: "arm64",
		OS:           "ubuntu",
		Channel:      "18.04",
	})
}

func (s *validatorSuite) TestDeducePlatformPlacementMutipleMatchFail(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)
	s.state.EXPECT().Machine(gomock.Any()).Return(s.machine, nil).AnyTimes()
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
	_, _, err := s.getValidator().deducePlatform(arg)
	c.Assert(errors.Is(err, errors.BadRequest), jc.IsTrue, gc.Commentf("%+v", err))
}

var configYaml = `
testme:
  optionOne: one
  optionTwo: 8
`[1:]

func (s *validatorSuite) TestAppCharmSettings(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.model.EXPECT().Type().Return(state.ModelTypeIAAS)

	cfg := charm.NewConfig()
	cfg.Options = map[string]charm.Option{
		"optionOne": {
			Type:        "string",
			Description: "option one",
		},
		"optionTwo": {
			Type:        "int",
			Description: "option two",
		},
	}

	appCfgSchema, _, err := applicationConfigSchema(state.ModelTypeIAAS)
	c.Assert(err, jc.ErrorIsNil)

	expectedAppConfig, err := coreconfig.NewConfig(map[string]interface{}{"trust": true}, appCfgSchema, nil)
	c.Assert(err, jc.ErrorIsNil)

	appConfig, charmConfig, err := s.getValidator().appCharmSettings("testme", true, cfg, configYaml)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(appConfig, gc.DeepEquals, expectedAppConfig)
	c.Assert(charmConfig["optionOne"], gc.DeepEquals, "one")
	c.Assert(charmConfig["optionTwo"], gc.DeepEquals, int64(8))
}

// The purpose of the resolveResourcesArgsMatcher is
// to compare the slices of resource.Resource, b/c the
// order is non-deterministic.
type resolveResourcesArgsMatcher struct {
	c        *gc.C
	expected *[]resource.Resource
}

func (m resolveResourcesArgsMatcher) String() string {
	return "match ResolveResources arg map"
}

func (m resolveResourcesArgsMatcher) Matches(x interface{}) bool {
	obtainedSlice, ok := x.([]resource.Resource)
	if !ok {
		return false
	}

	m.c.Assert(obtainedSlice, gc.HasLen, len(*m.expected))
	// Unfortunately the jc.SameContents don't work here
	// because resource.Resource is unhashable
	for _, r := range obtainedSlice {
		found := false
		for _, exR := range *m.expected {
			if reflect.DeepEqual(r, exR) {
				found = true
				break
			}
		}
		m.c.Assert(found, gc.Equals, true)
	}
	return true
}

func (s *validatorSuite) TestResolveResourcesSuccess(c *gc.C) {
	defer s.setupMocks(c).Finish()
	curl := charm.MustParseURL("testcharm")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	// Resource 1 : file upload from client
	meta1 := resource.Meta{
		Name:        "foo-resource",
		Type:        resource.TypeFile,
		Path:        "foo.txt",
		Description: "bar",
	}
	res := resource.Resource{
		Meta:     meta1,
		Origin:   resource.OriginUpload,
		Revision: -1,
	}
	// Resource 2 : store resource with --resource <revision> flag
	meta2 := resource.Meta{
		Name:        "foo-resource2",
		Type:        resource.TypeFile,
		Path:        "foo.txt",
		Description: "bar",
	}
	res2 := resource.Resource{
		Meta:     meta2,
		Origin:   resource.OriginStore,
		Revision: 3,
	}
	// Resource 3 : store resource without the --resource flag
	// (revision is reported by the store)
	meta3 := resource.Meta{
		Name:        "foo-resource3",
		Type:        resource.TypeFile,
		Path:        "foo.txt",
		Description: "bar",
	}
	res3 := resource.Resource{
		Meta:     meta3,
		Origin:   resource.OriginStore,
		Revision: -1,
	}

	resMeta := map[string]resource.Meta{"foo-file": meta1, "foo-file2": meta2, "store-file-res": meta3}
	resArgs := []resource.Resource{res, res2, res3}
	// Note that for the Resource 3, in the args res3 has revision -1, and the result below has revision 4
	r4 := resource.Resource{
		Meta:     meta3,
		Origin:   resource.OriginStore,
		Revision: 4,
	}
	resResult := []resource.Resource{res, res2, r4}

	dataResources := map[string]resource.Resource{
		"foo-file": {
			Meta:     resource.Meta{Name: "foo-file", Type: 1, Path: "foo.txt", Description: "foo resource."},
			Origin:   2,
			Revision: 2,
			Fingerprint: resource.Fingerprint{
				Fingerprint: hash.Fingerprint{
					//sum: {0xa9, 0x3b, 0xa5, 0xaa, 0x9c, 0x39, 0x87, 0xa6, 0x72, 0xfe, 0xc9, 0xc9, 0x77, 0xea, 0xef, 0x5c, 0x54, 0xba, 0x89, 0x51, 0x9b, 0xda, 0xf1, 0xb1, 0xf3, 0x3b, 0xe5, 0x1a, 0x31, 0x45, 0xbc, 0x13, 0x69, 0xb9, 0x1e, 0x15, 0x13, 0xfe, 0x94, 0x4f, 0x33, 0xd7, 0xc2, 0x5f, 0x5e, 0xd0, 0xad, 0x9c},
				},
			},
			Size: 13,
		},
	}
	// First one of below is the file upload for Resource 1, the second is the revision for Resource 2e
	deployResArg := map[string]string{"foo-file": "bar", "foo-file2": "3"}
	deployData := corecharm.ResolvedDataForDeploy{
		URL:               curl,
		EssentialMetadata: corecharm.EssentialMetadata{ResolvedOrigin: origin},
		Resources:         dataResources,
	}

	s.repo.EXPECT().ResolveResourcesForDeploy(resolveResourcesArgsMatcher{c: c, expected: &resArgs}, dataResources, corecharm.CharmID{URL: curl, Origin: origin}).Return(resResult, nil)
	//s.repo.EXPECT().ResolveResources(resolveResourcesArgsMatcher{c: c, expected: &resArgs}, corecharm.CharmID{URL: curl, Origin: origin}).Return(resResult, nil)
	resources, pendingResourceUploads, resolveResErr := s.getValidator().resolveResources(deployData, deployResArg, resMeta)
	pendUp := &params.PendingResourceUpload{
		Name:     "foo-resource",
		Type:     "file",
		Filename: "bar",
	}
	c.Assert(resolveResErr, jc.ErrorIsNil)
	c.Assert(resources, gc.DeepEquals, resResult)
	c.Assert(pendingResourceUploads, gc.DeepEquals, []*params.PendingResourceUpload{pendUp})
}

func (s *validatorSuite) TestCaasDeployFromRepositoryValidator(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectSimpleValidate()
	// resolveCharm
	curl := charm.MustParseURL("testcharm")
	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64"},
	}
	resolvedOrigin := corecharm.Origin{
		Source:   "charm-hub",
		Type:     "charm",
		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
		Revision: intptr(4),
	}
	resolveArg := corecharm.ResolveForDeployArg{
		URL:    curl,
		Origin: origin,
	}
	// getCharm
	expMeta := &charm.Meta{
		Name: "test-charm",
	}
	expManifest := new(charm.Manifest)
	expConfig := new(charm.Config)
	essMeta := corecharm.EssentialMetadata{
		Meta:           expMeta,
		Manifest:       expManifest,
		Config:         expConfig,
		ResolvedOrigin: resolvedOrigin,
	}
	data := corecharm.ResolvedDataForDeploy{
		URL:               resultURL,
		EssentialMetadata: essMeta,
		Resources:         nil,
	}
	s.repo.EXPECT().ResolveForDeploy(resolveForDeployMatcher{c: c, expectedArgs: resolveArg}).Return(data, nil)
	//supportedSeries := []string{"jammy", "focal"}
	//s.repo.EXPECT().ResolveWithPreferredChannel(curl, origin).Return(resultURL, resolvedOrigin, supportedSeries, nil)
	//s.repo.EXPECT().ResolveResources(nil, corecharm.CharmID{URL: resultURL, Origin: resolvedOrigin}).Return(nil, nil)
	s.state.EXPECT().Charm(gomock.Any()).Return(nil, errors.NotFoundf("charm"))

	//s.repo.EXPECT().GetEssentialMetadata(corecharm.MetadataRequest{
	//	CharmURL: resultURL,
	//	Origin:   resolvedOrigin,
	//}).Return([]corecharm.EssentialMetadata{essMeta}, nil)
	//s.state.EXPECT().ModelConstraints().Return(constraints.Value{Arch: strptr("arm64")}, nil)

	arg := params.DeployFromRepositoryArg{
		CharmName: "testcharm",
	}

	obtainedDT, errs := s.caasDeployFromRepositoryValidator(c).ValidateArg(arg)
	c.Assert(errs, gc.HasLen, 0)
	c.Logf("%s", pretty.Sprint(obtainedDT))
	c.Assert(obtainedDT, gc.DeepEquals, deployTemplate{
		applicationName: "test-charm",
		charm:           corecharm.NewCharmInfoAdapter(essMeta),
		charmURL:        resultURL,
		numUnits:        1,
		origin:          resolvedOrigin,
	})
}

type resolveForDeployMatcher struct {
	c            *gc.C
	expectedArgs corecharm.ResolveForDeployArg
}

func (m resolveForDeployMatcher) String() string {
	return "match ResolveForDeployArg"
}

func (m resolveForDeployMatcher) Matches(x interface{}) bool {
	oA, ok := x.(corecharm.ResolveForDeployArg)
	if !ok {
		return false
	}

	eA := m.expectedArgs
	m.c.Assert(oA.URL, gc.DeepEquals, eA.URL)
	m.c.Assert(oA.Origin, gc.DeepEquals, eA.Origin)
	m.c.Assert(oA.BaseSelectionFunc, gc.NotNil)
	return true
}

func (s *validatorSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.bindings = NewMockBindings(ctrl)
	s.machine = NewMockMachine(ctrl)
	s.model = NewMockModel(ctrl)
	s.repo = NewMockRepository(ctrl)
	s.repoFactory = NewMockRepositoryFactory(ctrl)
	s.state = NewMockDeployFromRepositoryState(ctrl)
	return ctrl
}

func (s *validatorSuite) getValidator() *deployFromRepositoryValidator {
	s.repoFactory.EXPECT().GetCharmRepository(gomock.Any()).Return(s.repo, nil).AnyTimes()
	return &deployFromRepositoryValidator{
		model:       s.model,
		state:       s.state,
		repoFactory: s.repoFactory,
		newStateBindings: func(st state.EndpointBinding, givenMap map[string]string) (Bindings, error) {
			return s.bindings, nil
		},
	}
}

func (s *validatorSuite) caasDeployFromRepositoryValidator(c *gc.C) caasDeployFromRepositoryValidator {
	return caasDeployFromRepositoryValidator{
		validator: s.getValidator(),
		caasPrecheckFunc: func(dt deployTemplate) error {
			// Do a quick check to ensure the expected deployTemplate
			// has been passed.
			c.Assert(dt.applicationName, gc.Equals, "test-charm")
			return nil
		},
	}
}

func (s *validatorSuite) iaasDeployFromRepositoryValidator() iaasDeployFromRepositoryValidator {
	return iaasDeployFromRepositoryValidator{
		validator: s.getValidator(),
	}
}

func (s *validatorSuite) expectSimpleValidate() {
	// createOrigin
	s.state.EXPECT().ModelConstraints().Return(constraints.Value{}, nil)
	s.model.EXPECT().Config().Return(config.New(config.UseDefaults, coretesting.FakeConfig())).AnyTimes()
}

func (s *validatorSuite) expectResolveForDeploy(c *gc.C) (*charm.URL, corecharm.EssentialMetadata) {
	return s.expectResolveForDeployWithPlatform(c, corecharm.Platform{Architecture: "amd64"})
}

func (s *validatorSuite) expectResolveForDeployWithPlacement(c *gc.C) (*charm.URL, corecharm.EssentialMetadata) {
	return s.expectResolveForDeployWithPlatform(c, corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"})
}

func (s *validatorSuite) expectResolveForDeployWithPlatform(c *gc.C, platform corecharm.Platform) (*charm.URL, corecharm.EssentialMetadata) {
	curl := charm.MustParseURL("testcharm")
	origin := corecharm.Origin{
		Source:   "charm-hub",
		Channel:  &charm.Channel{Risk: "stable"},
		Platform: platform,
	}
	resArgs := corecharm.ResolveForDeployArg{
		URL:    curl,
		Origin: origin,
	}

	resultURL := charm.MustParseURL("ch:amd64/jammy/testcharm-4")
	resolvedOrigin := corecharm.Origin{
		Source:   "charm-hub",
		Type:     "charm",
		Channel:  &charm.Channel{Track: "default", Risk: "stable"},
		Platform: corecharm.Platform{Architecture: "amd64", OS: "ubuntu", Channel: "22.04"},
		Revision: intptr(4),
	}
	expMeta := &charm.Meta{
		Name: "test-charm",
	}
	expManifest := new(charm.Manifest)
	expConfig := new(charm.Config)
	essMeta := corecharm.EssentialMetadata{
		Meta:           expMeta,
		Manifest:       expManifest,
		Config:         expConfig,
		ResolvedOrigin: resolvedOrigin,
	}
	deployData := corecharm.ResolvedDataForDeploy{
		URL:               resultURL,
		EssentialMetadata: essMeta,
	}

	s.repo.EXPECT().ResolveForDeploy(resolveForDeployMatcher{c: c, expectedArgs: resArgs}).Return(deployData, nil)
	s.state.EXPECT().Charm(gomock.Any()).Return(nil, errors.NotFoundf("charm"))
	return resultURL, essMeta
}

func strptr(s string) *string {
	return &s
}

func intptr(i int) *int {
	return &i
}

type deployRepositorySuite struct {
	application *MockApplication
	charm       *MockCharm
	state       *MockDeployFromRepositoryState
	validator   *MockDeployFromRepositoryValidator
}

func (s *deployRepositorySuite) TestDeployFromRepositoryAPI(c *gc.C) {
	defer s.setupMocks(c).Finish()
	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
	}
	template := deployTemplate{
		applicationName: "metadata-name",
		charm:           corecharm.NewCharmInfoAdapter(corecharm.EssentialMetadata{}),
		charmURL:        charm.MustParseURL("ch:amd64/jammy/testme-5"),
		endpoints:       map[string]string{"to": "from"},
		numUnits:        1,
		origin: corecharm.Origin{
			Source:   "charm-hub",
			Revision: intptr(5),
			Channel:  &charm.Channel{Risk: "stable"},
			Platform: corecharm.MustParsePlatform("amd64/ubuntu/22.04"),
		},
		placement: []*instance.Placement{{Directive: "0", Scope: instance.MachineScope}},
	}
	s.validator.EXPECT().ValidateArg(arg).Return(template, nil)
	info := state.CharmInfo{
		Charm: template.charm,
		ID:    charm.MustParseURL("ch:amd64/jammy/testme-5"),
	}

	s.state.EXPECT().AddCharmMetadata(info).Return(s.charm, nil)

	addAppArgs := state.AddApplicationArgs{
		Name: "metadata-name",
		// the app.Charm is casted into a state.Charm in the code
		// we mock it separately here (s.charm above), the test works
		// thanks to the addApplicationArgsMatcher used below
		Charm: &state.Charm{},
		CharmOrigin: &state.CharmOrigin{
			Source:   "charm-hub",
			Revision: intptr(5),
			Channel: &state.Channel{
				Risk: "stable",
			},
			Platform: &state.Platform{
				Architecture: "amd64",
				OS:           "ubuntu",
				Channel:      "22.04",
			},
		},
		Devices:          map[string]state.DeviceConstraints{},
		EndpointBindings: map[string]string{"to": "from"},
		NumUnits:         1,
		Placement:        []*instance.Placement{{Directive: "0", Scope: instance.MachineScope}},
		Resources:        map[string]string{},
		Storage:          map[string]state.StorageConstraints{},
	}
	s.state.EXPECT().AddApplication(addApplicationArgsMatcher{c: c, expectedArgs: addAppArgs}).Return(s.application, nil)

	deployFromRepositoryAPI := s.getDeployFromRepositoryAPI()

	obtainedInfo, resources, errs := deployFromRepositoryAPI.DeployFromRepository(arg)
	c.Assert(errs, gc.HasLen, 0)
	c.Assert(resources, gc.HasLen, 0)
	c.Assert(obtainedInfo, gc.DeepEquals, params.DeployFromRepositoryInfo{
		Architecture:     "amd64",
		Base:             params.Base{Name: "ubuntu", Channel: "22.04"},
		Channel:          "stable",
		EffectiveChannel: nil,
		Name:             "metadata-name",
		Revision:         5,
	})
}

// The reason for this matcher is that the AddApplicationArgs.Charm is
// obtained by casting application.Charm into a state.Charm, but we
// can't do that cast with a MockCharm
type addApplicationArgsMatcher struct {
	c            *gc.C
	expectedArgs state.AddApplicationArgs
}

func (m addApplicationArgsMatcher) String() string {
	return "match AddApplicationArgs"
}

func (m addApplicationArgsMatcher) Matches(x interface{}) bool {

	oA, ok := x.(state.AddApplicationArgs)
	if !ok {
		return false
	}

	eA := m.expectedArgs
	// Check everything but the Charm
	m.c.Assert(oA.Name, gc.DeepEquals, eA.Name)
	m.c.Assert(oA.ApplicationConfig, gc.DeepEquals, eA.ApplicationConfig)
	m.c.Assert(oA.NumUnits, gc.DeepEquals, eA.NumUnits)
	m.c.Assert(oA.Constraints, gc.DeepEquals, eA.Constraints)
	m.c.Assert(oA.Storage, gc.DeepEquals, eA.Storage)
	m.c.Assert(oA.Devices, gc.DeepEquals, eA.Devices)
	m.c.Assert(eA.AttachStorage, gc.DeepEquals, eA.AttachStorage)
	m.c.Assert(oA.EndpointBindings, gc.DeepEquals, eA.EndpointBindings)
	m.c.Assert(oA.CharmConfig, gc.DeepEquals, eA.CharmConfig)
	m.c.Assert(oA.Placement, gc.DeepEquals, eA.Placement)
	m.c.Assert(oA.Resources, gc.DeepEquals, eA.Resources)
	return true
}

func (s *deployRepositorySuite) TestAddPendingResourcesForDeployFromRepositoryAPI(c *gc.C) {
	defer s.setupMocks(c).Finish()
	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
	}
	pendUp := &params.PendingResourceUpload{
		Name:     "foo-resource",
		Type:     "file",
		Filename: "bar",
	}
	meta := resource.Meta{
		Name:        "foo-resource",
		Type:        resource.TypeFile,
		Path:        "foo.txt",
		Description: "bar",
	}
	r := resource.Resource{
		Meta:   meta,
		Origin: resource.OriginUpload,
	}

	template := deployTemplate{
		applicationName: "metadata-name",
		charm:           corecharm.NewCharmInfoAdapter(corecharm.EssentialMetadata{}),
		charmURL:        charm.MustParseURL("ch:amd64/jammy/testme-5"),
		endpoints:       map[string]string{"to": "from"},
		numUnits:        1,
		origin: corecharm.Origin{
			Source:   "charm-hub",
			Revision: intptr(5),
			Channel:  &charm.Channel{Risk: "stable"},
			Platform: corecharm.MustParsePlatform("amd64/ubuntu/22.04"),
		},
		placement:              []*instance.Placement{{Directive: "0", Scope: instance.MachineScope}},
		resources:              map[string]string{"foo-file": "bar"},
		pendingResourceUploads: []*params.PendingResourceUpload{pendUp},
		resolvedResources:      []resource.Resource{r},
	}
	s.validator.EXPECT().ValidateArg(arg).Return(template, nil)
	info := state.CharmInfo{
		Charm: template.charm,
		ID:    charm.MustParseURL("ch:amd64/jammy/testme-5"),
	}

	s.state.EXPECT().AddCharmMetadata(info).Return(s.charm, nil)

	s.state.EXPECT().AddPendingResource("metadata-name", r).Return("3", nil)

	addAppArgs := state.AddApplicationArgs{
		Name: "metadata-name",
		// the app.Charm is casted into a state.Charm in the code
		// we mock it separately here (s.charm above), the test works
		// thanks to the addApplicationArgsMatcher used below
		Charm: &state.Charm{},
		CharmOrigin: &state.CharmOrigin{
			Source:   "charm-hub",
			Revision: intptr(5),
			Channel: &state.Channel{
				Risk: "stable",
			},
			Platform: &state.Platform{
				Architecture: "amd64",
				OS:           "ubuntu",
				Channel:      "22.04",
			},
		},
		Devices:          map[string]state.DeviceConstraints{},
		EndpointBindings: map[string]string{"to": "from"},
		NumUnits:         1,
		Placement:        []*instance.Placement{{Directive: "0", Scope: instance.MachineScope}},
		Resources:        map[string]string{"foo-resource": "3"},
		Storage:          map[string]state.StorageConstraints{},
	}
	s.state.EXPECT().AddApplication(addApplicationArgsMatcher{c: c, expectedArgs: addAppArgs}).Return(s.application, nil)

	deployFromRepositoryAPI := s.getDeployFromRepositoryAPI()

	obtainedInfo, resources, errs := deployFromRepositoryAPI.DeployFromRepository(arg)
	c.Assert(errs, gc.HasLen, 0)
	c.Assert(resources, gc.HasLen, 1)
	c.Assert(obtainedInfo, gc.DeepEquals, params.DeployFromRepositoryInfo{
		Architecture:     "amd64",
		Base:             params.Base{Name: "ubuntu", Channel: "22.04"},
		Channel:          "stable",
		EffectiveChannel: nil,
		Name:             "metadata-name",
		Revision:         5,
	})

	c.Assert(resources, gc.DeepEquals, []*params.PendingResourceUpload{pendUp})
}

func (s *deployRepositorySuite) TestRemovePendingResourcesWhenDeployErrors(c *gc.C) {
	defer s.setupMocks(c).Finish()
	arg := params.DeployFromRepositoryArg{
		CharmName: "testme",
	}
	pendUp := &params.PendingResourceUpload{
		Name:     "foo-resource",
		Type:     "file",
		Filename: "bar",
	}
	meta := resource.Meta{
		Name:        "foo-resource",
		Type:        resource.TypeFile,
		Path:        "foo.txt",
		Description: "bar",
	}
	r := resource.Resource{
		Meta:   meta,
		Origin: resource.OriginUpload,
	}
	template := deployTemplate{
		applicationName: "metadata-name",
		charm:           corecharm.NewCharmInfoAdapter(corecharm.EssentialMetadata{}),
		charmURL:        charm.MustParseURL("ch:amd64/jammy/testme-5"),
		endpoints:       map[string]string{"to": "from"},
		numUnits:        1,
		origin: corecharm.Origin{
			Source:   "charm-hub",
			Revision: intptr(5),
			Channel:  &charm.Channel{Risk: "stable"},
			Platform: corecharm.MustParsePlatform("amd64/ubuntu/22.04"),
		},
		placement:              []*instance.Placement{{Directive: "0", Scope: instance.MachineScope}},
		resources:              map[string]string{"foo-file": "bar"},
		pendingResourceUploads: []*params.PendingResourceUpload{pendUp},
		resolvedResources:      []resource.Resource{r},
	}
	s.validator.EXPECT().ValidateArg(arg).Return(template, nil)
	info := state.CharmInfo{
		Charm: template.charm,
		ID:    charm.MustParseURL("ch:amd64/jammy/testme-5"),
	}

	s.state.EXPECT().AddCharmMetadata(info).Return(s.charm, nil)

	s.state.EXPECT().AddPendingResource("metadata-name", r).Return("3", nil)

	addAppArgs := state.AddApplicationArgs{
		Name: "metadata-name",
		// the app.Charm is casted into a state.Charm in the code
		// we mock it separately here (s.charm above), the test works
		// thanks to the addApplicationArgsMatcher used below
		Charm: &state.Charm{},
		CharmOrigin: &state.CharmOrigin{
			Source:   "charm-hub",
			Revision: intptr(5),
			Channel: &state.Channel{
				Risk: "stable",
			},
			Platform: &state.Platform{
				Architecture: "amd64",
				OS:           "ubuntu",
				Channel:      "22.04",
			},
		},
		Devices:          map[string]state.DeviceConstraints{},
		EndpointBindings: map[string]string{"to": "from"},
		NumUnits:         1,
		Placement:        []*instance.Placement{{Directive: "0", Scope: instance.MachineScope}},
		Resources:        map[string]string{"foo-resource": "3"},
		Storage:          map[string]state.StorageConstraints{},
	}

	s.state.EXPECT().RemovePendingResources("metadata-name", map[string]string{"foo-resource": "3"})

	s.state.EXPECT().AddApplication(addApplicationArgsMatcher{c: c, expectedArgs: addAppArgs}).Return(s.application,
		errors.New("fail"))

	deployFromRepositoryAPI := s.getDeployFromRepositoryAPI()

	obtainedInfo, resources, errs := deployFromRepositoryAPI.DeployFromRepository(arg)
	c.Assert(errs, gc.HasLen, 1)
	c.Assert(resources, gc.HasLen, 0)
	c.Assert(obtainedInfo, gc.DeepEquals, params.DeployFromRepositoryInfo{})
}

func (s *deployRepositorySuite) getDeployFromRepositoryAPI() *DeployFromRepositoryAPI {
	return &DeployFromRepositoryAPI{
		state:      s.state,
		validator:  s.validator,
		stateCharm: func(Charm) *state.Charm { return nil },
	}
}

func (s *deployRepositorySuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.charm = NewMockCharm(ctrl)
	s.state = NewMockDeployFromRepositoryState(ctrl)
	s.validator = NewMockDeployFromRepositoryValidator(ctrl)
	return ctrl
}
