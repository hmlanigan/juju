// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmigration

import (
	"testing"

	"github.com/juju/description/v9"
	"github.com/juju/tc"
	"go.uber.org/mock/gomock"

	"github.com/juju/juju/core/network"
	"github.com/juju/juju/domain/network/internal"
	loggertesting "github.com/juju/juju/internal/logger/testing"
)

func TestImportLinkLayerSuite(t *testing.T) {
	tc.Run(t, &importLinkLayerSuite{})
}

type importLinkLayerSuite struct {
	migrationService *MockMigrationService
}

func (s *importLinkLayerSuite) TestImportLinkLayerDevices(c *tc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	model := description.NewModel(description.ModelArgs{})
	dArgs := description.LinkLayerDeviceArgs{
		Name:        "test-device",
		MTU:         1500,
		ProviderID:  "net-lxdbr0",
		MachineID:   "77",
		Type:        "ethernet",
		MACAddress:  "00:16:3e:ad:4e:01",
		IsAutoStart: true,
		IsUp:        true,
	}
	model.AddLinkLayerDevice(dArgs)

	args := []internal.ImportLinkLayerDevice{
		{
			IsAutoStart: dArgs.IsAutoStart,
			IsEnabled:   dArgs.IsUp,
			MTU:         ptr(int64(dArgs.MTU)),
			MachineID:   dArgs.MachineID,
			MACAddress:  ptr(dArgs.MACAddress),
			Name:        dArgs.Name,
			ProviderID:  ptr(dArgs.ProviderID),
			Type:        network.EthernetDevice,
		},
	}
	s.migrationService.EXPECT().ImportLinkLayerDevices(gomock.Any(), lldArgMatcher{c: c, expected: args}).Return(nil)

	// Act
	op := s.newImportOperation(c)
	err := op.Execute(c.Context(), model)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
}

func (s *importLinkLayerSuite) TestImportLinkLayerDevicesOptionalValues(c *tc.C) {
	// Arrange: ensure input not containing an MTU, ProviderID, nor
	// MACAddress values to see nil values in the data passed to the
	// service.
	defer s.setupMocks(c).Finish()
	model := description.NewModel(description.ModelArgs{})
	dArgs := description.LinkLayerDeviceArgs{
		Name:        "test-device",
		MachineID:   "77",
		Type:        "ethernet",
		IsAutoStart: true,
		IsUp:        true,
	}
	model.AddLinkLayerDevice(dArgs)

	args := []internal.ImportLinkLayerDevice{
		{
			IsAutoStart: dArgs.IsAutoStart,
			IsEnabled:   dArgs.IsUp,
			MachineID:   dArgs.MachineID,
			Name:        dArgs.Name,
			Type:        network.EthernetDevice,
		},
	}
	s.migrationService.EXPECT().ImportLinkLayerDevices(gomock.Any(), lldArgMatcher{c: c, expected: args}).Return(nil)

	// Act
	op := s.newImportOperation(c)
	err := op.Execute(c.Context(), model)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
}

func (s *importLinkLayerSuite) TestRollbackLinkLayerDevices(c *tc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	model := description.NewModel(description.ModelArgs{})
	dArgs := description.LinkLayerDeviceArgs{}
	model.AddLinkLayerDevice(dArgs)
	s.migrationService.EXPECT().DeleteImportedLinkLayerDevices(gomock.Any()).Return(nil)

	// Act
	op := s.newImportOperation(c)
	err := op.Rollback(c.Context(), model)

	// Assert
	c.Assert(err, tc.ErrorIsNil)
}

func (s *importLinkLayerSuite) TestRollbackLinkLayerDevicesNoData(c *tc.C) {
	// Arrange
	defer s.setupMocks(c).Finish()
	model := description.NewModel(description.ModelArgs{})

	// Act
	op := s.newImportOperation(c)
	err := op.Rollback(c.Context(), model)

	// Assert: with no link layer device data, there is no failure.
	c.Assert(err, tc.ErrorIsNil)
}

func (s *importLinkLayerSuite) setupMocks(c *tc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.migrationService = NewMockMigrationService(ctrl)

	c.Cleanup(func() {
		s.migrationService = nil
	})

	return ctrl
}

func (s *importLinkLayerSuite) newImportOperation(c *tc.C) *importLinkLayerOperation {
	return &importLinkLayerOperation{
		migrationService: s.migrationService,
		logger:           loggertesting.WrapCheckLog(c),
	}
}

// lldArgMatcher verifies the args for ImportLinkLayerDevice.
type lldArgMatcher struct {
	c        *tc.C
	expected []internal.ImportLinkLayerDevice
}

func (m lldArgMatcher) Matches(x interface{}) bool {
	input, ok := x.([]internal.ImportLinkLayerDevice)
	if !ok {
		return false
	}
	// UUIDs are assigned in the code under test. Ensure they exist, then
	// remove it to enable SameContents checks over the other fields.
	for i, in := range input {
		m.c.Check(in.UUID, tc.Not(tc.Equals), "")
		out := in
		out.UUID = ""
		input[i] = out
	}
	return m.c.Check(input, tc.SameContents, m.expected)
}

func (lldArgMatcher) String() string {
	return "matches args for ImportLinkLayerDevice"
}
