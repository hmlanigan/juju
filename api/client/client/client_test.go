// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client_test

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"

	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/version/v2"
	gc "gopkg.in/check.v1"

	basemocks "github.com/juju/juju/api/base/mocks"
	"github.com/juju/juju/api/client/client"
	jujunames "github.com/juju/juju/juju/names"
	coretesting "github.com/juju/juju/testing"
)

type clientSuite struct {
}

var _ = gc.Suite(&clientSuite{})

func (s *clientSuite) TestUploadToolsOtherModel(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)
	testClient := client.NewClientFromCaller(mockFacadeCaller)

	newVersion := version.MustParseBinary("5.4.3-ubuntu-amd64")
	var called bool

	// build fake tools
	expectedTools, _ := coretesting.TarGz(
		coretesting.NewTarFile(jujunames.Jujud, 0777, "jujud contents "+newVersion.String()))

	// UploadTools does not use the facades, so instead of patching the
	// facade call, we set up a fake endpoint to test.
	defer func() {
		_ = fakeAPIEndpoint(c, testClient, modelEndpoint(c, "tools"), "POST",
			func(w http.ResponseWriter, r *http.Request) {
				called = true

				c.Assert(r.URL.Query(), gc.DeepEquals, url.Values{
					"binaryVersion": []string{"5.4.3-ubuntu-amd64"},
					"series":        []string{""},
				})
				defer func() { _ = r.Body.Close() }()
				obtainedTools, err := ioutil.ReadAll(r.Body)
				c.Assert(err, jc.ErrorIsNil)
				c.Assert(obtainedTools, gc.DeepEquals, expectedTools)
			},
		).Close()
	}()

	// We don't test the error or tools results as we only wish to assert that
	// the API client POSTs the tools archive to the correct endpoint.
	_, _ = testClient.UploadTools(bytes.NewReader(expectedTools), newVersion)
	c.Assert(called, jc.IsTrue)
}

func fakeAPIEndpoint(c *gc.C, cl *client.Client, address, method string, handle func(http.ResponseWriter, *http.Request)) net.Listener {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	c.Assert(err, jc.ErrorIsNil)

	mux := http.NewServeMux()
	mux.HandleFunc(address, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			handle(w, r)
		}
	})
	go func() {
		http.Serve(lis, mux)
	}()
	client.SetServerAddress(cl, "http", lis.Addr().String())
	return lis
}

// modelEndpoint returns "/model/<model-uuid>/<destination>"
func modelEndpoint(c *gc.C, destination string) string {
	return path.Join("/model", coretesting.ModelTag.String(), destination)
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
