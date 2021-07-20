// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmhub

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	gomock "github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	path "github.com/juju/juju/charmhub/path"
	"github.com/juju/juju/charmhub/transport"
	"github.com/juju/juju/core/arch"
)

type RefreshSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&RefreshSuite{})

func (s *RefreshSuite) TestRefresh(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	baseURL := MustParseURL(c, "http://api.foo.bar")

	path := path.MakePath(baseURL)
	id := "meshuggah"
	body := transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{{
			InstanceKey: "foo-bar",
			ID:          id,
			Revision:    1,
			Base: transport.Base{
				Name:         "ubuntu",
				Channel:      "20.04",
				Architecture: arch.DefaultArchitecture,
			},
			TrackingChannel: "latest/stable",
		}},
		Actions: []transport.RefreshRequestAction{{
			Action:      "refresh",
			InstanceKey: "foo-bar",
			ID:          &id,
		}},
	}

	config, err := RefreshOne(id, 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	restClient := NewMockRESTClient(ctrl)
	s.expectPost(c, restClient, path, id, body)

	client := NewRefreshClient(path, restClient, &FakeLogger{})
	responses, err := client.Refresh(context.TODO(), config)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(len(responses), gc.Equals, 1)
	c.Assert(responses[0].Name, gc.Equals, id)
}

//	c.Assert(results.Results[0].Error, gc.ErrorMatches, `.* pool "foo" not found`)
func (s *RefreshSuite) TestRefeshConfigValidateArch(c *gc.C) {
	err := s.testRefeshConfigValidate(c, RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: "all",
	})
	c.Assert(err, gc.ErrorMatches, "Architecture.*")
}

func (s *RefreshSuite) TestRefeshConfigValidateSeries(c *gc.C) {
	err := s.testRefeshConfigValidate(c, RefreshBase{
		Name:         "ubuntu",
		Channel:      "all",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, gc.ErrorMatches, "Channel.*")
}

func (s *RefreshSuite) TestRefeshConfigValidateName(c *gc.C) {
	err := s.testRefeshConfigValidate(c, RefreshBase{
		Name:         "all",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, gc.ErrorMatches, "Name.*")
}

func (s *RefreshSuite) TestRefeshConfigValidate(c *gc.C) {
	err := s.testRefeshConfigValidate(c, RefreshBase{
		Name:         "all",
		Channel:      "all",
		Architecture: "all",
	})
	c.Assert(err, gc.ErrorMatches, "Architecture.*, Name.*, Channel.*")
}

func (s *RefreshSuite) testRefeshConfigValidate(c *gc.C, rp RefreshBase) error {
	_, err := DownloadOneFromChannel("meshuggah", "latest/stable", rp)
	c.Assert(err, jc.Satisfies, errors.IsNotValid)
	return err
}

type metadataTransport struct {
	requestHeaders http.Header
	responseBody   string
}

func (t *metadataTransport) Do(req *http.Request) (*http.Response, error) {
	t.requestHeaders = req.Header
	resp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(t.responseBody)),
		ContentLength: int64(len(t.responseBody)),
		Request:       req,
		Header:        http.Header{"Content-Type": []string{"application/json"}},
	}
	return resp, nil
}

func (s *RefreshSuite) TestRefreshMetadata(c *gc.C) {
	baseURL := MustParseURL(c, "http://api.foo.bar")
	path := path.MakePath(baseURL)

	httpTransport := &metadataTransport{
		responseBody: `
{
  "error-list": [],
  "results": [
    {
      "id": "foo",
      "instance-key": "key-foo"
    },
    {
      "id": "bar",
      "instance-key": "key-bar"
    }
  ]
}
`,
	}

	headers := http.Header{"User-Agent": []string{"Test Agent 1.0"}}
	restClient := NewHTTPRESTClient(httpTransport, headers)
	client := NewRefreshClient(path, restClient, &FakeLogger{})

	config1, err := RefreshOne("foo", 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: "amd64",
	})
	c.Assert(err, jc.ErrorIsNil)
	config1 = DefineInstanceKey(c, config1, "key-foo")

	config2, err := RefreshOne("bar", 2, "latest/edge", RefreshBase{
		Name:         "ubuntu",
		Channel:      "14.04",
		Architecture: "amd64",
	})
	c.Assert(err, jc.ErrorIsNil)
	config2 = DefineInstanceKey(c, config2, "key-bar")

	config := RefreshMany(config1, config2)

	response, err := client.Refresh(context.Background(), config)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(httpTransport.requestHeaders[MetadataHeader], jc.SameContents, []string{
		"arch=amd64",
		"name=ubuntu",
		"channel=20.04",
		"channel=14.04",
	})
	c.Assert(httpTransport.requestHeaders["User-Agent"], jc.SameContents, []string{"Test Agent 1.0"})

	c.Assert(response, gc.DeepEquals, []transport.RefreshResponse{
		{ID: "foo", InstanceKey: "key-foo"},
		{ID: "bar", InstanceKey: "key-bar"},
	})
}

func (s *RefreshSuite) TestRefreshFailure(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	baseURL := MustParseURL(c, "http://api.foo.bar")

	path := path.MakePath(baseURL)
	name := "meshuggah"

	config, err := RefreshOne(name, 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	restClient := NewMockRESTClient(ctrl)
	s.expectPostFailure(c, restClient)

	client := NewRefreshClient(path, restClient, &FakeLogger{})
	_, err = client.Refresh(context.TODO(), config)
	c.Assert(err, gc.Not(jc.ErrorIsNil))
}

func (s *RefreshSuite) expectPost(c *gc.C, client *MockRESTClient, p path.Path, name string, body interface{}) {
	client.EXPECT().Post(gomock.Any(), p, gomock.Any(), body, gomock.Any()).Do(func(_ context.Context, _ path.Path, _ map[string][]string, _ transport.RefreshRequest, responses *transport.RefreshResponses) {
		responses.Results = []transport.RefreshResponse{{
			InstanceKey: "foo-bar",
			Name:        name,
		}}
	}).Return(RESTResponse{StatusCode: http.StatusOK}, nil)
}

func (s *RefreshSuite) expectPostFailure(c *gc.C, client *MockRESTClient) {
	client.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(RESTResponse{StatusCode: http.StatusInternalServerError}, errors.Errorf("boom"))
}

func DefineInstanceKey(c *gc.C, config RefreshConfig, key string) RefreshConfig {
	switch t := config.(type) {
	case refreshOne:
		t.instanceKey = key
		return t
	case executeOne:
		t.instanceKey = key
		return t
	case installByRevisionOne:
		t.instanceKey = key
		return t
	default:
		c.Fatalf("unexpected config %T", config)
	}
	return nil
}

type RefreshConfigSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&RefreshConfigSuite{})

func (s *RefreshConfigSuite) TestRefreshOneBuild(c *gc.C) {
	id := "foo"
	config, err := RefreshOne(id, 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{{
			InstanceKey: "foo-bar",
			ID:          "foo",
			Revision:    1,
			Base: transport.Base{
				Name:         "ubuntu",
				Channel:      "20.04",
				Architecture: arch.DefaultArchitecture,
			},
			TrackingChannel: "latest/stable",
		}},
		Actions: []transport.RefreshRequestAction{{
			Action:      "refresh",
			InstanceKey: "foo-bar",
			ID:          &id,
		}},
	})
}

func (s *RefreshConfigSuite) TestRefreshOneEnsure(c *gc.C) {
	config, err := RefreshOne("foo", 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	err = config.Ensure([]transport.RefreshResponse{{
		InstanceKey: "foo-bar",
	}})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *RefreshConfigSuite) TestInstallOneBuildRevision(c *gc.C) {
	revision := 1

	name := "foo"
	config, err := InstallOneFromRevision(name, revision)
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{},
		Actions: []transport.RefreshRequestAction{{
			Action:      "install",
			InstanceKey: "foo-bar",
			Name:        &name,
			Revision:    &revision,
		}},
		Fields: []string{"bases", "download", "id", "revision", "version", "resources", "type"},
	})
}

func (s *RefreshConfigSuite) TestInstallOneBuildRevisionResources(c *gc.C) {
	// Tests InstallOne by revision with specific resources.
	revision := 1

	name := "foo"
	config, err := InstallOneFromRevision(name, revision)
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")
	config, ok := AddResource(config, "testme", 3)
	c.Assert(ok, jc.IsTrue)

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{},
		Actions: []transport.RefreshRequestAction{{
			Action:      "install",
			InstanceKey: "foo-bar",
			Name:        &name,
			Revision:    &revision,
			ResourceRevisions: []transport.RefreshResourceRevision{
				{Name: "testme", Revision: 3},
			},
		}},
		Fields: []string{"bases", "download", "id", "revision", "version", "resources", "type"},
	})
}

func (s *RefreshConfigSuite) TestAddResourceFail(c *gc.C) {
	config, err := RefreshOne("testingID", 7, "latest/edge", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)
	_, ok := AddResource(config, "testme", 3)
	c.Assert(ok, jc.IsFalse)
}

func (s *RefreshConfigSuite) TestInstallOneBuildChannel(c *gc.C) {
	channel := "latest/stable"

	name := "foo"
	config, err := InstallOneFromChannel(name, channel, RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{},
		Actions: []transport.RefreshRequestAction{{
			Action:      "install",
			InstanceKey: "foo-bar",
			Name:        &name,
			Channel:     &channel,
			Base: &transport.Base{
				Name:         "ubuntu",
				Channel:      "20.04",
				Architecture: arch.DefaultArchitecture,
			},
		}},
	})
}

func (s *RefreshConfigSuite) TestInstallOneWithPartialPlatform(c *gc.C) {
	channel := "latest/stable"

	name := "foo"
	config, err := InstallOneFromChannel(name, channel, RefreshBase{
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{},
		Actions: []transport.RefreshRequestAction{{
			Action:      "install",
			InstanceKey: "foo-bar",
			Name:        &name,
			Channel:     &channel,
			Base: &transport.Base{
				Name:         NotAvailable,
				Channel:      NotAvailable,
				Architecture: arch.DefaultArchitecture,
			},
		}},
	})
}

func (s *RefreshConfigSuite) TestInstallOneWithMissingArch(c *gc.C) {
	channel := "latest/stable"

	name := "foo"
	config, err := InstallOneFromChannel(name, channel, RefreshBase{})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	_, _, err = config.Build()
	c.Assert(errors.IsNotValid(err), jc.IsTrue)
}

func (s *RefreshConfigSuite) TestInstallOneEnsure(c *gc.C) {
	config, err := InstallOneFromChannel("foo", "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	err = config.Ensure([]transport.RefreshResponse{{
		InstanceKey: "foo-bar",
	}})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *RefreshConfigSuite) TestInstallOneFromChannelEnsure(c *gc.C) {
	config, err := InstallOneFromChannel("foo", "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	err = config.Ensure([]transport.RefreshResponse{{
		InstanceKey: "foo-bar",
	}})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *RefreshConfigSuite) TestDownloadOneEnsure(c *gc.C) {
	config, err := DownloadOne("foo", 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	err = config.Ensure([]transport.RefreshResponse{{
		InstanceKey: "foo-bar",
	}})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *RefreshConfigSuite) TestDownloadOneFromChannelBuild(c *gc.C) {
	channel := "latest/stable"
	id := "foo"
	config, err := DownloadOneFromChannel(id, channel, RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{},
		Actions: []transport.RefreshRequestAction{{
			Action:      "download",
			InstanceKey: "foo-bar",
			ID:          &id,
			Channel:     &channel,
			Base: &transport.Base{
				Name:         "ubuntu",
				Channel:      "20.04",
				Architecture: arch.DefaultArchitecture,
			},
		}},
	})
}

func (s *RefreshConfigSuite) TestDownloadOneFromChannelBuildK8s(c *gc.C) {
	channel := "latest/stable"
	id := "foo"
	config, err := DownloadOneFromChannel(id, channel, RefreshBase{
		Name:         "kubernetes",
		Channel:      "kubernetes",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{},
		Actions: []transport.RefreshRequestAction{{
			Action:      "download",
			InstanceKey: "foo-bar",
			ID:          &id,
			Channel:     &channel,
			Base: &transport.Base{
				Name:         "ubuntu",
				Channel:      "20.04",
				Architecture: arch.DefaultArchitecture,
			},
		}},
	})
}

func (s *RefreshConfigSuite) TestDownloadOneFromChannelEnsure(c *gc.C) {
	config, err := DownloadOneFromChannel("foo", "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config = DefineInstanceKey(c, config, "foo-bar")

	err = config.Ensure([]transport.RefreshResponse{{
		InstanceKey: "foo-bar",
	}})
	c.Assert(err, jc.ErrorIsNil)
}

func (s *RefreshConfigSuite) TestRefreshManyBuildContextNotNil(c *gc.C) {
	id1 := "foo"
	config1, err := DownloadOneFromRevision(id1, 1, RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)
	config1 = DefineInstanceKey(c, config1, "foo-bar")

	id2 := "bar"
	config2, err := DownloadOneFromChannel(id2, "latest/edge", RefreshBase{
		Name:         "ubuntu",
		Channel:      "14.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)
	config2 = DefineInstanceKey(c, config2, "foo-baz")
	config := RefreshMany(config1, config2)

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req.Context, gc.NotNil)
}

func (s *RefreshConfigSuite) TestRefreshManyBuild(c *gc.C) {
	id1 := "foo"
	config1, err := RefreshOne(id1, 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)
	config1 = DefineInstanceKey(c, config1, "foo-bar")

	id2 := "bar"
	config2, err := RefreshOne(id2, 2, "latest/edge", RefreshBase{
		Name:         "ubuntu",
		Channel:      "14.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)
	config2 = DefineInstanceKey(c, config2, "foo-baz")

	channel := "1/stable"

	name3 := "baz"
	config3, err := InstallOneFromChannel(name3, "1/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "disco",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)

	config3 = DefineInstanceKey(c, config3, "foo-taz")

	name4 := "forty-two"
	rev4 := 42
	config4, err := InstallOneFromRevision(name4, rev4)
	c.Assert(err, jc.ErrorIsNil)

	config4 = DefineInstanceKey(c, config4, "foo-two")

	config := RefreshMany(config1, config2, config3, config4)

	req, _, err := config.Build()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(req, gc.DeepEquals, transport.RefreshRequest{
		Context: []transport.RefreshRequestContext{{
			InstanceKey: "foo-bar",
			ID:          "foo",
			Revision:    1,
			Base: transport.Base{
				Name:         "ubuntu",
				Channel:      "20.04",
				Architecture: arch.DefaultArchitecture,
			},
			TrackingChannel: "latest/stable",
		}, {
			InstanceKey: "foo-baz",
			ID:          "bar",
			Revision:    2,
			Base: transport.Base{
				Name:         "ubuntu",
				Channel:      "14.04",
				Architecture: arch.DefaultArchitecture,
			},
			TrackingChannel: "latest/edge",
		}},
		Actions: []transport.RefreshRequestAction{{
			Action:      "refresh",
			InstanceKey: "foo-bar",
			ID:          &id1,
		}, {
			Action:      "refresh",
			InstanceKey: "foo-baz",
			ID:          &id2,
		}, {
			Action:      "install",
			InstanceKey: "foo-taz",
			Name:        &name3,
			Base: &transport.Base{
				Name:         "ubuntu",
				Channel:      "19.04",
				Architecture: arch.DefaultArchitecture,
			},
			Channel: &channel,
		}, {
			Action:      "install",
			InstanceKey: "foo-two",
			Name:        &name4,
			Revision:    &rev4,
		}},
		Fields: []string{"bases", "download", "id", "revision", "version", "resources", "type"},
	})
}

func (s *RefreshConfigSuite) TestRefreshManyEnsure(c *gc.C) {
	config1, err := RefreshOne("foo", 1, "latest/stable", RefreshBase{
		Name:         "ubuntu",
		Channel:      "20.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)
	config1 = DefineInstanceKey(c, config1, "foo-bar")

	config2, err := RefreshOne("bar", 2, "latest/edge", RefreshBase{
		Name:         "ubuntu",
		Channel:      "14.04",
		Architecture: arch.DefaultArchitecture,
	})
	c.Assert(err, jc.ErrorIsNil)
	config2 = DefineInstanceKey(c, config2, "foo-baz")

	config := RefreshMany(config1, config2)

	err = config.Ensure([]transport.RefreshResponse{{
		InstanceKey: "foo-bar",
	}, {
		InstanceKey: "foo-baz",
	}})
	c.Assert(err, jc.ErrorIsNil)
}
