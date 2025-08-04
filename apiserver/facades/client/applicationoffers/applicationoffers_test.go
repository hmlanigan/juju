// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package applicationoffers

import (
	"context"
	"testing"

	"github.com/juju/names/v6"
	"github.com/juju/tc"
	"go.uber.org/mock/gomock"

	corecrossmodel "github.com/juju/juju/core/crossmodel"
	"github.com/juju/juju/core/model"
	modeltesting "github.com/juju/juju/core/model/testing"
	"github.com/juju/juju/core/permission"
	"github.com/juju/juju/core/user"
	"github.com/juju/juju/domain/access"
	"github.com/juju/juju/internal/uuid"
	"github.com/juju/juju/rpc/params"
)

type offerSuite struct {
	authorizer    *MockAuthorizer
	accessService *MockAccessService
	modelService  *MockModelService
	offerService  *MockOfferService
}

func TestServiceSuite(t *testing.T) {
	tc.Run(t, &offerSuite{})
}

// TestModifyOfferAccess tests a basic call to ModifyOfferAccess by
// a controller admin.
func (s *offerSuite) TestModifyOfferAccess(c *tc.C) {
	s.setupMocks(c).Finish()

	// Arrange
	s.authorizer.EXPECT().HasPermission(gomock.Any(), permission.SuperuserAccess, gomock.Any()).Return(nil)
	s.authorizer.EXPECT().GetAuthTag().Return(names.NewUserTag("admin"))
	modelInfo := model.Model{
		UUID: modeltesting.GenModelUUID(c),
	}
	qualifier := model.QualifierFromUserTag(names.NewUserTag("admin"))
	s.modelService.EXPECT().GetModelByNameAndQualifier(gomock.Any(), "model", qualifier).Return(modelInfo, nil)

	offerURL, _ := corecrossmodel.ParseOfferURL("model.application:db")
	offerUUID := uuid.MustNewUUID()
	s.offerService.EXPECT().GetOfferUUID(gomock.Any(), offerURL).Return(offerUUID, nil)
	userName, _ := user.NewName("simon")
	updateArgs := access.UpdatePermissionArgs{
		AccessSpec: permission.AccessSpec{
			Target: permission.ID{
				ObjectType: permission.Offer,
				Key:        offerUUID.String(),
			},
			Access: permission.ConsumeAccess,
		},
		Change:  permission.Grant,
		Subject: userName,
	}
	s.accessService.EXPECT().UpdatePermission(gomock.Any(), updateArgs).Return(nil)

	args := params.ModifyOfferAccessRequest{
		Changes: []params.ModifyOfferAccess{
			{
				UserTag:  "user-simon",
				Action:   params.GrantOfferAccess,
				Access:   params.OfferConsumeAccess,
				OfferURL: "model.application:db",
			},
		},
	}

	// Act
	results, err := s.offerAPI().ModifyOfferAccess(c.Context(), args)

	// Assert
	c.Assert(err, tc.IsNil)
	c.Assert(results, tc.DeepEquals, params.ErrorResults{Results: []params.ErrorResult{{Error: nil}}})
}

func (s *offerSuite) setupMocks(c *tc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)
	s.accessService = NewMockAccessService(ctrl)
	s.authorizer = NewMockAuthorizer(ctrl)
	s.modelService = NewMockModelService(ctrl)
	s.offerService = NewMockOfferService(ctrl)

	c.Cleanup(func() {
		s.accessService = nil
		s.authorizer = nil
		s.modelService = nil
		s.offerService = nil
	})
	return ctrl
}

func (s *offerSuite) offerAPI() *OffersAPI {
	return &OffersAPI{
		controllerUUID: uuid.MustNewUUID().String(),
		authorizer:     s.authorizer,
		accessService:  s.accessService,
		modelService:   s.modelService,
		offerServiceGetter: func(_ context.Context, _ model.UUID) (OfferService, error) {
			return s.offerService, nil
		},
	}
}
