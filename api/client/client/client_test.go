// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/version/v2"
	gc "gopkg.in/check.v1"
	"gopkg.in/httprequest.v1"

	"github.com/juju/juju/api"
	basemocks "github.com/juju/juju/api/base/mocks"
	"github.com/juju/juju/api/client/client"
	"github.com/juju/juju/api/client/client/mocks"
	corestatus "github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"
	coretools "github.com/juju/juju/tools"
)

type clientSuite struct {
}

var _ = gc.Suite(&clientSuite{})

func (s *clientSuite) TestStatus(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)

	pattern := []string{"ubuntu/0"}
	args := params.StatusParams{Patterns: pattern}
	result := new(params.FullStatus)
	hostname := "juju-testme"
	results := params.FullStatus{
		Model:               params.ModelStatusInfo{},
		Machines:            map[string]params.MachineStatus{"0": {Hostname: hostname}},
		Applications:        nil,
		RemoteApplications:  nil,
		Offers:              nil,
		Relations:           nil,
		ControllerTimestamp: nil,
		Branches:            nil,
	}
	mockFacadeCaller.EXPECT().FacadeCall("FullStatus", args, result).SetArg(2, results).Return(nil)

	testClient := client.NewClientFromCaller(mockFacadeCaller)
	obtainedStatus, err := testClient.Status(pattern)
	c.Assert(err, jc.ErrorIsNil) // Confirms that the correct facade was called
	c.Assert(obtainedStatus.Model.Type, gc.Equals, "iaas")
	m, ok := obtainedStatus.Machines["0"]
	c.Assert(ok, jc.IsTrue)
	c.Assert(m.Hostname, gc.Equals, hostname)
}

func (s *clientSuite) TestStatusHistory(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)

	tag := names.NewMachineTag("0")
	result := new(params.StatusHistoryResults)
	results := params.StatusHistoryResults{
		Results: []params.StatusHistoryResult{{
			History: params.History{
				Statuses: []params.DetailedStatus{
					{
						Kind: corestatus.KindMachine.String(),
					},
				},
			},
			Error: nil,
		},
		}}
	mockFacadeCaller.EXPECT().FacadeCall("StatusHistory", gomock.AssignableToTypeOf(params.StatusHistoryRequests{}), result).SetArg(2, results).DoAndReturn(
		func(request string, args, response interface{}) error {
			p, ok := args.(params.StatusHistoryRequests)
			c.Assert(ok, jc.IsTrue)
			c.Assert(p.Requests, gc.HasLen, 1)
			c.Assert(p.Requests[0].Kind, gc.Equals, corestatus.KindMachine.String())
			c.Assert(p.Requests[0].Filter.Size, gc.Equals, 10)

			return nil
		})

	testClient := client.NewClientFromCaller(mockFacadeCaller)

	obtainedHistory, err := testClient.StatusHistory(corestatus.KindMachine, tag, corestatus.StatusHistoryFilter{Size: 10})
	c.Assert(err, jc.ErrorIsNil) // Confirms that the correct facade was called
	c.Assert(obtainedHistory, gc.HasLen, 1)
}

func (s *clientSuite) TestWatchAll(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)

	result := new(params.AllWatcherId)
	results := params.AllWatcherId{AllWatcherId: "watcher"}
	mockFacadeCaller.EXPECT().FacadeCall("WatchAll", nil, result).SetArg(2, results).Return(nil)

	testClient := client.NewClientFromCallerAndConnection(mockFacadeCaller, mockConnection)
	w, err := testClient.WatchAll()
	c.Assert(err, jc.ErrorIsNil)
	//c.Assert(w, jc.Satisfies, *api.AllWatcher)
	//	var logWriter raftleasestore.Logger
	//	c.Assert(args[1], gc.Implements, &logWriter)
	var allWatcher *api.AllWatch
	c.Assert(w, gc.Implements, allWatcher)
}

func (s *clientSuite) TestSetModelAgentVersion(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)

	ver := version.MustParse("2.9.100")
	args := params.SetModelAgentVersion{
		Version:             ver,
		AgentStream:         "test-stream",
		IgnoreAgentVersions: true,
	}

	mockFacadeCaller.EXPECT().FacadeCall("SetModelAgentVersion", args, nil).Return(nil)

	testClient := client.NewClientFromCaller(mockFacadeCaller)
	err := testClient.SetModelAgentVersion(ver, "test-stream", true)
	c.Assert(err, jc.ErrorIsNil) // Confirms that the correct facade was called
}

func (s *clientSuite) TestAbortCurrentUpgrade(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)
	someErr := errors.New("random")
	mockFacadeCaller.EXPECT().FacadeCall("AbortCurrentUpgrade", nil, nil).Return(someErr)

	testClient := client.NewClientFromCaller(mockFacadeCaller)
	err := testClient.AbortCurrentUpgrade()
	c.Assert(err, gc.Equals, someErr) // Confirms that the correct facade was called
}

func (s *clientSuite) TestFindTools(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)

	ver := version.MustParseBinary("2.9.100-ubuntu-amd64")
	args := params.FindToolsParams{
		MajorVersion: ver.Major,
		MinorVersion: ver.Minor,
		Arch:         ver.Arch,
		OSType:       "ubuntu",
		AgentStream:  "test-stream",
	}
	list := coretools.List{
		&coretools.Tools{Version: ver},
	}
	results := params.FindToolsResult{
		List: list,
	}
	result := new(params.FindToolsResult)
	mockFacadeCaller.EXPECT().FacadeCall("FindTools", args, result).SetArg(2, results).Return(nil)

	testClient := client.NewClientFromCaller(mockFacadeCaller)
	obtained, err := testClient.FindTools(2, 9, "ubuntu", "amd64", "test-stream")
	c.Assert(err, jc.ErrorIsNil) // Confirms that the correct facade was called
	c.Assert(obtained.List, gc.DeepEquals, list)
}

func (s *clientSuite) TestUploadTools(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockDoer := mocks.NewMockDoer(ctrl)
	mockAPICallCloser := basemocks.NewMockAPICallCloser(ctrl)
	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)
	mockConnection := mocks.NewMockConnection(ctrl)

	mockConnection.EXPECT().HTTPClient().Return(&httprequest.Client{Doer: mockDoer}, nil)
	mockFacadeCaller.EXPECT().RawAPICaller().Return(mockAPICallCloser)
	mockAPICallCloser.EXPECT().Context().Return(context.TODO())

	ver := version.MustParseBinary("2.9.100-ubuntu-amd64")

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("/tools?binaryVersion=%s&series=", ver),
		nil,
	)
	c.Assert(err, jc.ErrorIsNil)
	req.Header.Set("Content-Type", "application/x-tar-gz")
	req = req.WithContext(context.TODO())

	resp := &http.Response{
		Request:    req,
		StatusCode: http.StatusCreated,
		Header:     http.Header{},
		Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"tools": [{"version": "%s"}]}`, ver))),
	}
	resp.Header.Set("Content-Type", "application/json")
	mockDoer.EXPECT().Do(req).Return(resp, nil)

	testClient := client.NewClientFromCallerAndConnection(mockFacadeCaller, mockConnection)
	result, err := testClient.UploadTools(
		nil, ver,
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, gc.DeepEquals, coretools.List{
		{Version: ver},
	})
}
