// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"

	"github.com/juju/description/v8"
	"github.com/juju/errors"
	"github.com/juju/mgo/v3/txn"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/internal/uuid"
)

type MigrationImportTasksSuite struct{}

var _ = gc.Suite(&MigrationImportTasksSuite{})

func (s *MigrationImportTasksSuite) TestImportApplicationOffers(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	offerUUID, err := uuid.NewUUID()
	c.Assert(err, jc.ErrorIsNil)
	offerUUID2, err := uuid.NewUUID()
	c.Assert(err, jc.ErrorIsNil)

	runner := ImportApplicationOfferRunner{
		OfferUUID: offerUUID.String(),
		model:     NewMockApplicationOfferInput(ctrl),
		runner:    NewMockTransactionRunner(ctrl),
	}

	entity := s.applicationOffer(ctrl)
	offerDoc := applicationOfferDoc{
		DocID:                  fmt.Sprintf("%s:%s", offerUUID.String(), "offer-name-foo"),
		OfferUUID:              offerUUID.String(),
		OfferName:              "offer-name-foo",
		ApplicationName:        "foo",
		ApplicationDescription: "foo app description",
		Endpoints: map[string]string{
			"db": "db",
		},
	}
	secondOfferDoc := offerDoc
	secondOfferDoc.DocID = fmt.Sprintf("%s:%s", offerUUID2.String(), "second-offer")
	secondOfferDoc.OfferUUID = offerUUID2.String()
	secondOfferDoc.OfferName = "second-offer"

	entity.EXPECT().ApplicationName().Return(offerDoc.ApplicationName).Times(2)
	runner.Add(runner.applicationOffers(entity, entity))
	runner.Add(runner.docID(offerDoc.OfferName, offerDoc.DocID))
	runner.Add(runner.docID(secondOfferDoc.OfferName, secondOfferDoc.DocID))
	runner.Add(runner.applicationOfferDoc(offerDoc, entity))
	runner.Add(runner.applicationOfferDoc(secondOfferDoc, entity))

	refOp := txn.Op{
		Assert: txn.DocMissing,
	}
	runner.Add(runner.applicationOffersRefOp(refOp, 2))

	entity.EXPECT().ACL().Return(map[string]string{"fred": "consume"}).Times(2)
	runner.Add(runner.transaction([]applicationOfferDoc{offerDoc, secondOfferDoc}, refOp))

	err = runner.Run(ctrl)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *MigrationImportTasksSuite) TestImportApplicationOffersTransactionFailure(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	offerUUID, err := uuid.NewUUID()
	c.Assert(err, jc.ErrorIsNil)

	runner := ImportApplicationOfferRunner{
		OfferUUID: offerUUID.String(),
		model:     NewMockApplicationOfferInput(ctrl),
		runner:    NewMockTransactionRunner(ctrl),
	}

	entity := s.applicationOffer(ctrl)
	offerDoc := applicationOfferDoc{
		DocID:                  fmt.Sprintf("%s:%s", offerUUID.String(), "offer-name-foo"),
		OfferUUID:              offerUUID.String(),
		OfferName:              "offer-name-foo",
		ApplicationName:        "foo",
		ApplicationDescription: "foo app description",
		Endpoints: map[string]string{
			"db": "db",
		},
	}

	entity.EXPECT().ACL().Return(map[string]string{})
	entity.EXPECT().ApplicationName().Return(offerDoc.ApplicationName)
	runner.Add(runner.applicationOffers(entity))
	runner.Add(runner.applicationOfferDoc(offerDoc, entity))
	runner.Add(runner.docID(offerDoc.OfferName, offerDoc.DocID))

	refOp := txn.Op{
		Assert: txn.DocMissing,
	}
	runner.Add(runner.applicationOffersRefOp(refOp, 1))
	runner.Add(runner.transactionWithError(errors.New("fail")))

	err = runner.Run(ctrl)
	c.Assert(err, gc.ErrorMatches, "fail")
}

type ImportApplicationOfferRunner struct {
	OfferUUID string
	model     *MockApplicationOfferInput
	runner    *MockTransactionRunner
	ops       []func(*gomock.Controller)
}

func (s *ImportApplicationOfferRunner) Add(fn func(*gomock.Controller)) {
	s.ops = append(s.ops, fn)
}

func (s *ImportApplicationOfferRunner) Run(ctrl *gomock.Controller) error {
	for _, v := range s.ops {
		v(ctrl)
	}

	m := ImportApplicationOffer{}
	return m.Execute(s.model, s.runner)
}

func (s *ImportApplicationOfferRunner) applicationOffers(entities ...description.ApplicationOffer) func(ctrl *gomock.Controller) {
	return func(ctrl *gomock.Controller) {
		s.model.EXPECT().Offers().Return(entities)
	}
}

func (s *ImportApplicationOfferRunner) applicationOfferDoc(offerDoc applicationOfferDoc, entity description.ApplicationOffer) func(ctrl *gomock.Controller) {
	return func(ctrl *gomock.Controller) {
		s.model.EXPECT().MakeApplicationOfferDoc(entity).Return(offerDoc, nil)
	}
}

func (s *ImportApplicationOfferRunner) applicationOffersRefOp(op txn.Op, cnt int) func(ctrl *gomock.Controller) {
	return func(ctrl *gomock.Controller) {
		s.model.EXPECT().MakeApplicationOffersRefOp("foo", cnt).Return(op, nil)
	}
}

func (s *MigrationImportTasksSuite) applicationOffer(ctrl *gomock.Controller) *MockApplicationOffer {
	return NewMockApplicationOffer(ctrl)
}

func (s *ImportApplicationOfferRunner) docID(offerName, docID string) func(ctrl *gomock.Controller) {
	return func(ctrl *gomock.Controller) {
		s.model.EXPECT().DocID(offerName).Return(docID)
	}
}

func (s *ImportApplicationOfferRunner) transaction(offerDocs []applicationOfferDoc, ops ...txn.Op) func(ctrl *gomock.Controller) {
	return func(ctrl *gomock.Controller) {
		useOps := make([]txn.Op, 0)

		for _, doc := range offerDocs {
			useOps = append(useOps, []txn.Op{
				{
					C:      applicationOffersC,
					Id:     doc.DocID,
					Assert: txn.DocMissing,
					Insert: doc,
				},
			}...)
		}
		useOps = append(useOps, ops...)

		s.runner.EXPECT().RunTransaction(useOps).Return(nil)
	}
}

func (s *ImportApplicationOfferRunner) transactionWithError(err error) func(ctrl *gomock.Controller) {
	return func(ctrl *gomock.Controller) {
		s.runner.EXPECT().RunTransaction(gomock.Any()).Return(err)
	}
}

func (s *MigrationImportTasksSuite) TestImportRemoteApplications(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	status := s.status(ctrl)

	entity0 := s.remoteApplication(ctrl, status)
	entities := []description.RemoteApplication{
		entity0,
	}

	appDoc := &remoteApplicationDoc{
		Name:      "remote-application",
		URL:       "me/model.rainbow",
		OfferUUID: "offer-uuid",
		Endpoints: []remoteEndpointDoc{
			{Name: "db", Interface: "mysql"},
			{Name: "db-admin", Interface: "mysql-root"},
		},
	}
	statusDoc := statusDoc{}
	statusOp := txn.Op{}

	model := NewMockRemoteApplicationsInput(ctrl)
	model.EXPECT().RemoteApplications().Return(entities)
	model.EXPECT().MakeRemoteApplicationDoc(entity0).Return(appDoc)
	model.EXPECT().NewRemoteApplication(appDoc).Return(&RemoteApplication{
		doc: *appDoc,
	})
	model.EXPECT().MakeStatusDoc(status).Return(statusDoc)
	model.EXPECT().MakeStatusOp("c#remote-application", statusDoc).Return(statusOp)
	model.EXPECT().DocID("remote-application").Return("c#remote-application")

	runner := NewMockTransactionRunner(ctrl)
	runner.EXPECT().RunTransaction([]txn.Op{
		{
			C:      applicationsC,
			Id:     "remote-application",
			Assert: txn.DocMissing,
		},
		{
			C:      remoteApplicationsC,
			Id:     "c#remote-application",
			Assert: txn.DocMissing,
			Insert: appDoc,
		},
		statusOp,
	}).Return(nil)

	m := ImportRemoteApplications{}
	err := m.Execute(model, runner)
	c.Assert(err, jc.ErrorIsNil)
}

// A Remote Application with a missing status field is a valid remote
// application and should be correctly imported.
func (s *MigrationImportTasksSuite) TestImportRemoteApplicationsWithMissingStatusField(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	entity0 := s.remoteApplication(ctrl, nil)

	entities := []description.RemoteApplication{
		entity0,
	}

	appDoc := &remoteApplicationDoc{
		Name:      "remote-application",
		URL:       "me/model.rainbow",
		OfferUUID: "offer-uuid",
		Endpoints: []remoteEndpointDoc{
			{Name: "db", Interface: "mysql"},
			{Name: "db-admin", Interface: "mysql-root"},
		},
	}

	model := NewMockRemoteApplicationsInput(ctrl)
	model.EXPECT().RemoteApplications().Return(entities)
	model.EXPECT().MakeRemoteApplicationDoc(entity0).Return(appDoc)
	model.EXPECT().NewRemoteApplication(appDoc).Return(&RemoteApplication{
		doc: *appDoc,
	})
	model.EXPECT().DocID("remote-application").Return("c#remote-application")

	runner := NewMockTransactionRunner(ctrl)
	runner.EXPECT().RunTransaction([]txn.Op{
		{
			C:      applicationsC,
			Id:     "remote-application",
			Assert: txn.DocMissing,
		},
		{
			C:      remoteApplicationsC,
			Id:     "c#remote-application",
			Assert: txn.DocMissing,
			Insert: appDoc,
		},
	}).Return(nil)

	m := ImportRemoteApplications{}
	err := m.Execute(model, runner)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *MigrationImportTasksSuite) remoteApplication(ctrl *gomock.Controller, status description.Status) description.RemoteApplication {
	entity := NewMockRemoteApplication(ctrl)
	entity.EXPECT().Status().Return(status)
	return entity
}

func (s *MigrationImportTasksSuite) status(ctrl *gomock.Controller) description.Status {
	entity := NewMockStatus(ctrl)
	return entity
}

func (s *MigrationImportTasksSuite) TestImportRemoteEntities(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	entity0 := s.remoteEntity(ctrl, "application-app2", "xxx-yyy-ccc")
	entity1 := s.remoteEntity(ctrl, "applicationoffer-offer3", "aaa-bbb-zzz")

	entities := []description.RemoteEntity{
		entity0,
		entity1,
	}

	model := NewMockRemoteEntitiesInput(ctrl)
	model.EXPECT().RemoteEntities().Return(entities)
	model.EXPECT().OfferUUIDForApp("app2").Return("uuid2", nil)
	model.EXPECT().OfferUUID("offer3").Return("uuid3", true)
	model.EXPECT().DocID("applicationoffer-uuid2").Return("doc-uuid2")
	model.EXPECT().DocID("applicationoffer-uuid3").Return("doc-uuid3")

	runner := NewMockTransactionRunner(ctrl)
	runner.EXPECT().RunTransaction([]txn.Op{
		{
			C:      remoteEntitiesC,
			Id:     "doc-uuid2",
			Assert: txn.DocMissing,
			Insert: &remoteEntityDoc{
				DocID: "doc-uuid2",
				Token: "xxx-yyy-ccc",
			},
		},
		{
			C:      remoteEntitiesC,
			Id:     "doc-uuid3",
			Assert: txn.DocMissing,
			Insert: &remoteEntityDoc{
				DocID: "doc-uuid3",
				Token: "aaa-bbb-zzz",
			},
		},
	}).Return(nil)

	m := ImportRemoteEntities{}
	err := m.Execute(model, runner)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *MigrationImportTasksSuite) TestImportRemoteEntitiesWithNoEntities(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	entities := []description.RemoteEntity{}

	model := NewMockRemoteEntitiesInput(ctrl)
	model.EXPECT().RemoteEntities().Return(entities)

	runner := NewMockTransactionRunner(ctrl)
	// No call to RunTransaction if there are no operations.

	m := ImportRemoteEntities{}
	err := m.Execute(model, runner)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *MigrationImportTasksSuite) TestImportRemoteEntitiesWithTransactionRunnerReturnsError(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	entity0 := s.remoteEntity(ctrl, "application-uuid2", "xxx-yyy-ccc")

	entities := []description.RemoteEntity{
		entity0,
	}

	model := NewMockRemoteEntitiesInput(ctrl)
	model.EXPECT().RemoteEntities().Return(entities)
	model.EXPECT().OfferUUIDForApp("uuid2").Return("offeruuid2", nil)
	model.EXPECT().DocID("applicationoffer-offeruuid2").Return("doc-uuid2")

	runner := NewMockTransactionRunner(ctrl)
	runner.EXPECT().RunTransaction([]txn.Op{
		{
			C:      remoteEntitiesC,
			Id:     "doc-uuid2",
			Assert: txn.DocMissing,
			Insert: &remoteEntityDoc{
				DocID: "doc-uuid2",
				Token: "xxx-yyy-ccc",
			},
		},
	}).Return(errors.New("fail"))

	m := ImportRemoteEntities{}
	err := m.Execute(model, runner)
	c.Assert(err, gc.ErrorMatches, "fail")
}

func (s *MigrationImportTasksSuite) remoteEntity(ctrl *gomock.Controller, id, token string) *MockRemoteEntity {
	entity := NewMockRemoteEntity(ctrl)
	entity.EXPECT().ID().Return(id).AnyTimes()
	entity.EXPECT().Token().Return(token)
	return entity
}

func (s *MigrationImportTasksSuite) TestImportRelationNetworks(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	entity0 := s.relationNetwork(ctrl, "ctrl-uuid-2", "xxx-yyy-ccc", []string{"10.0.1.0/16"})
	entity1 := s.relationNetwork(ctrl, "ctrl-uuid-3", "aaa-bbb-zzz", []string{"10.0.0.1/24"})

	entities := []description.RelationNetwork{
		entity0,
		entity1,
	}

	model := NewMockRelationNetworksInput(ctrl)
	model.EXPECT().RelationNetworks().Return(entities)
	model.EXPECT().DocID("ctrl-uuid-2").Return("ctrl-uuid-2")
	model.EXPECT().DocID("ctrl-uuid-3").Return("ctrl-uuid-3")

	runner := NewMockTransactionRunner(ctrl)
	runner.EXPECT().RunTransaction([]txn.Op{
		{
			C:      relationNetworksC,
			Id:     "ctrl-uuid-2",
			Assert: txn.DocMissing,
			Insert: relationNetworksDoc{
				Id:          "ctrl-uuid-2",
				RelationKey: "xxx-yyy-ccc",
				CIDRS:       []string{"10.0.1.0/16"},
			},
		},
		{
			C:      relationNetworksC,
			Id:     "ctrl-uuid-3",
			Assert: txn.DocMissing,
			Insert: relationNetworksDoc{
				Id:          "ctrl-uuid-3",
				RelationKey: "aaa-bbb-zzz",
				CIDRS:       []string{"10.0.0.1/24"},
			},
		},
	}).Return(nil)

	m := ImportRelationNetworks{}
	err := m.Execute(model, runner)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *MigrationImportTasksSuite) TestImportRelationNetworksWithNoEntities(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	entities := []description.RelationNetwork{}

	model := NewMockRelationNetworksInput(ctrl)
	model.EXPECT().RelationNetworks().Return(entities)

	runner := NewMockTransactionRunner(ctrl)
	// No call to RunTransaction if there are no operations.

	m := ImportRelationNetworks{}
	err := m.Execute(model, runner)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *MigrationImportTasksSuite) TestImportRelationNetworksWithTransactionRunnerReturnsError(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	entity0 := s.relationNetwork(ctrl, "ctrl-uuid-2", "xxx-yyy-ccc", []string{"10.0.1.0/16"})

	entities := []description.RelationNetwork{
		entity0,
	}

	model := NewMockRelationNetworksInput(ctrl)
	model.EXPECT().RelationNetworks().Return(entities)
	model.EXPECT().DocID("ctrl-uuid-2").Return("ctrl-uuid-2")

	runner := NewMockTransactionRunner(ctrl)
	runner.EXPECT().RunTransaction([]txn.Op{
		{
			C:      relationNetworksC,
			Id:     "ctrl-uuid-2",
			Assert: txn.DocMissing,
			Insert: relationNetworksDoc{
				Id:          "ctrl-uuid-2",
				RelationKey: "xxx-yyy-ccc",
				CIDRS:       []string{"10.0.1.0/16"},
			},
		},
	}).Return(errors.New("fail"))

	m := ImportRelationNetworks{}
	err := m.Execute(model, runner)
	c.Assert(err, gc.ErrorMatches, "fail")
}

func (s *MigrationImportTasksSuite) relationNetwork(ctrl *gomock.Controller, id, key string, cidrs []string) *MockRelationNetwork {
	entity := NewMockRelationNetwork(ctrl)
	entity.EXPECT().ID().Return(id)
	entity.EXPECT().RelationKey().Return(key)
	entity.EXPECT().CIDRS().Return(cidrs)
	return entity
}
