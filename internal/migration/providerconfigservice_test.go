// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migration

import (
	"testing"

	"github.com/juju/tc"
	"go.uber.org/mock/gomock"

	"github.com/juju/juju/cloud"
	coremodel "github.com/juju/juju/core/model"
	"github.com/juju/juju/core/providertracker"
	"github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/internal/testhelpers"
	coretesting "github.com/juju/juju/internal/testing"
	"github.com/juju/juju/internal/uuid"
)

type getEphemeralProviderConfigSuite struct {
	testhelpers.IsolationSuite
	modelConfigService     *MockModelConfigService
	controllerCloudService *MockControllerCloudService
	providerGetter         *MockProviderConfigServicesGetter
	providerConfigService  *MockProviderConfigServices
}

func TestGetEphemeralProviderConfigSuite(t *testing.T) {
	tc.Run(t, &getEphemeralProviderConfigSuite{})
}

func (s *getEphemeralProviderConfigSuite) setupMocks(c *tc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.controllerCloudService = NewMockControllerCloudService(ctrl)
	s.modelConfigService = NewMockModelConfigService(ctrl)
	s.providerGetter = NewMockProviderConfigServicesGetter(ctrl)
	s.providerConfigService = NewMockProviderConfigServices(ctrl)

	s.providerGetter.EXPECT().ServicesForModel(gomock.Any(), gomock.Any()).Return(s.providerConfigService, nil)
	s.providerConfigService.EXPECT().Config().Return(s.modelConfigService)
	s.providerConfigService.EXPECT().Cloud().Return(s.controllerCloudService)

	c.Cleanup(func() {
		s.controllerCloudService = nil
		s.modelConfigService = nil
		s.providerGetter = nil
		s.providerConfigService = nil
	})

	return ctrl
}

func (s *getEphemeralProviderConfigSuite) TestGetEphemeralProviderConfig(c *tc.C) {
	defer s.setupMocks(c).Finish()

	// Arrange
	cloudName := "test-cloud"
	cloudRegion := "test-region"
	cloud := &cloud.Cloud{
		Name: cloudName,
		Type: "kubernetes",
		Regions: []cloud.Region{
			{
				Name: cloudRegion,
			},
		},
	}
	s.controllerCloudService.EXPECT().Cloud(gomock.Any(), cloudName).Return(cloud, nil)
	modelConfig := coretesting.CustomModelConfig(c,
		coretesting.Attrs{
			"apt-mirror": "http://mirror",
		},
	)
	s.modelConfigService.EXPECT().ModelConfig(gomock.Any()).Return(modelConfig, nil)

	// Arrange: struct under test
	controllerUUID := tc.Must(c, uuid.NewUUID)
	epcp := ephemeralProviderConfigProvider{
		cloudCredentials: nil,
		cloudName:        cloudName,
		cloudRegion:      cloudRegion,
		controllerUUID:   controllerUUID.String(),
		modelType:        "caas",
		modelUUID:        tc.Must(c, coremodel.NewUUID),
		servicesGetter:   s.providerGetter,
	}

	// Act
	cfg, err := epcp.GetEphemeralProviderConfig(c.Context())

	// Assert
	c.Assert(err, tc.ErrorIsNil)
	c.Check(cfg, tc.DeepEquals, providertracker.EphemeralProviderConfig{
		ModelType:   coremodel.CAAS,
		ModelConfig: modelConfig,
		CloudSpec: cloudspec.CloudSpec{
			Type:   "kubernetes",
			Name:   cloudName,
			Region: cloudRegion,
		},
		ControllerUUID: controllerUUID,
	})
}
