// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/version/v2"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/api/client/client"
	jujunames "github.com/juju/juju/juju/names"
	jujutesting "github.com/juju/juju/juju/testing"
	"github.com/juju/juju/state"
	coretesting "github.com/juju/juju/testing"
)

type clientSuite struct {
	jujutesting.JujuConnSuite
}

var _ = gc.Suite(&clientSuite{})

// TODO(jam) 2013-08-27 http://pad.lv/1217282
// Right now most of the direct tests for client.Client behavior are in
// apiserver/client/*_test.go
func (s *clientSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
}

func (s *clientSuite) TestCloseMultipleOk(c *gc.C) {
	client := client.NewClient(s.APIState)
	c.Assert(client.Close(), gc.IsNil)
	c.Assert(client.Close(), gc.IsNil)
	c.Assert(client.Close(), gc.IsNil)
}

func (s *clientSuite) TestUploadToolsOtherModel(c *gc.C) {
	otherSt, otherAPISt := s.otherModel(c)
	defer otherSt.Close()
	defer otherAPISt.Close()
	client := client.NewClient(otherAPISt)
	newVersion := version.MustParseBinary("5.4.3-ubuntu-amd64")
	var called bool

	// build fake tools
	expectedTools, _ := coretesting.TarGz(
		coretesting.NewTarFile(jujunames.Jujud, 0777, "jujud contents "+newVersion.String()))

	// UploadTools does not use the facades, so instead of patching the
	// facade call, we set up a fake endpoint to test.
	defer fakeAPIEndpoint(c, client, modelEndpoint(c, otherAPISt, "tools"), "POST",
		func(w http.ResponseWriter, r *http.Request) {
			called = true

			c.Assert(r.URL.Query(), gc.DeepEquals, url.Values{
				"binaryVersion": []string{"5.4.3-ubuntu-amd64"},
				"series":        []string{""},
			})
			defer r.Body.Close()
			obtainedTools, err := ioutil.ReadAll(r.Body)
			c.Assert(err, jc.ErrorIsNil)
			c.Assert(obtainedTools, gc.DeepEquals, expectedTools)
		},
	).Close()

	// We don't test the error or tools results as we only wish to assert that
	// the API client POSTs the tools archive to the correct endpoint.
	client.UploadTools(bytes.NewReader(expectedTools), newVersion)
	c.Assert(called, jc.IsTrue)
}

func (s *clientSuite) otherModel(c *gc.C) (*state.State, api.Connection) {
	otherSt := s.Factory.MakeModel(c, nil)
	info := s.APIInfo(c)
	model, err := otherSt.Model()
	c.Assert(err, jc.ErrorIsNil)
	info.ModelTag = model.ModelTag()
	apiState, err := api.Open(info, api.DefaultDialOpts())
	c.Assert(err, jc.ErrorIsNil)
	return otherSt, apiState
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
func modelEndpoint(c *gc.C, apiState api.Connection, destination string) string {
	modelTag, ok := apiState.ModelTag()
	c.Assert(ok, jc.IsTrue)
	return path.Join("/model", modelTag.Id(), destination)
}

func (s *clientSuite) TestConnectStreamRequiresSlashPathPrefix(c *gc.C) {
	reader, err := s.APIState.ConnectStream("foo", nil)
	c.Assert(err, gc.ErrorMatches, `cannot make API path from non-slash-prefixed path "foo"`)
	c.Assert(reader, gc.Equals, nil)
}

func (s *clientSuite) TestConnectStreamErrorBadConnection(c *gc.C) {
	s.PatchValue(&api.WebsocketDial, func(_ api.WebsocketDialer, _ string, _ http.Header) (base.Stream, error) {
		return nil, fmt.Errorf("bad connection")
	})
	reader, err := s.APIState.ConnectStream("/", nil)
	c.Assert(err, gc.ErrorMatches, "bad connection")
	c.Assert(reader, gc.IsNil)
}

func (s *clientSuite) TestConnectStreamErrorNoData(c *gc.C) {
	s.PatchValue(&api.WebsocketDial, func(_ api.WebsocketDialer, _ string, _ http.Header) (base.Stream, error) {
		return api.NewFakeStreamReader(&bytes.Buffer{}), nil
	})
	reader, err := s.APIState.ConnectStream("/", nil)
	c.Assert(err, gc.ErrorMatches, "unable to read initial response: EOF")
	c.Assert(reader, gc.IsNil)
}

func (s *clientSuite) TestConnectStreamErrorBadData(c *gc.C) {
	s.PatchValue(&api.WebsocketDial, func(_ api.WebsocketDialer, _ string, _ http.Header) (base.Stream, error) {
		return api.NewFakeStreamReader(strings.NewReader("junk\n")), nil
	})
	reader, err := s.APIState.ConnectStream("/", nil)
	c.Assert(err, gc.ErrorMatches, "unable to unmarshal initial response: .*")
	c.Assert(reader, gc.IsNil)
}

func (s *clientSuite) TestConnectStreamErrorReadError(c *gc.C) {
	s.PatchValue(&api.WebsocketDial, func(_ api.WebsocketDialer, _ string, _ http.Header) (base.Stream, error) {
		err := fmt.Errorf("bad read")
		return api.NewFakeStreamReader(&badReader{err}), nil
	})
	reader, err := s.APIState.ConnectStream("/", nil)
	c.Assert(err, gc.ErrorMatches, "unable to read initial response: bad read")
	c.Assert(reader, gc.IsNil)
}

func (s *clientSuite) TestConnectControllerStreamRejectsRelativePaths(c *gc.C) {
	reader, err := s.APIState.ConnectControllerStream("foo", nil, nil)
	c.Assert(err, gc.ErrorMatches, `path "foo" is not absolute`)
	c.Assert(reader, gc.IsNil)
}

func (s *clientSuite) TestConnectControllerStreamRejectsModelPaths(c *gc.C) {
	reader, err := s.APIState.ConnectControllerStream("/model/foo", nil, nil)
	c.Assert(err, gc.ErrorMatches, `path "/model/foo" is model-specific`)
	c.Assert(reader, gc.IsNil)
}

func (s *clientSuite) TestConnectControllerStreamAppliesHeaders(c *gc.C) {
	catcher := api.UrlCatcher{}
	headers := http.Header{}
	headers.Add("thomas", "cromwell")
	headers.Add("anne", "boleyn")
	s.PatchValue(&api.WebsocketDial, catcher.RecordLocation)

	_, err := s.APIState.ConnectControllerStream("/something", nil, headers)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(catcher.Headers().Get("thomas"), gc.Equals, "cromwell")
	c.Assert(catcher.Headers().Get("anne"), gc.Equals, "boleyn")
}

func (s *clientSuite) TestAbortCurrentUpgrade(c *gc.C) {
	cl := client.NewClient(s.APIState)
	someErr := errors.New("random")
	cleanup := client.PatchClientFacadeCall(cl,
		func(request string, args interface{}, response interface{}) error {
			c.Assert(request, gc.Equals, "AbortCurrentUpgrade")
			c.Assert(args, gc.IsNil)
			c.Assert(response, gc.IsNil)
			return someErr
		},
	)
	defer cleanup()

	err := cl.AbortCurrentUpgrade()
	c.Assert(err, gc.Equals, someErr) // Confirms that the correct facade was called
}

// badReader raises err when Read is called.
type badReader struct {
	err error
}

func (r *badReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}
