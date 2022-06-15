// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm_test

import (
	"fmt"
	"os"
	"path/filepath"

	corecharm "github.com/juju/charm/v8"
	"github.com/juju/collections/set"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/testcharms"
	"github.com/juju/juju/worker/uniter/charm"
)

// charmReader is a charm.CharmReader that lets us mock out the charmReader we
// deploy to test the Deployers.
type charmReader struct {
	charmArchives map[string]charm.CharmArchive
	stopWaiting   <-chan struct{}
}

// EnableWaitForAbort allows us to test that a Deployer.Stage call passes its abort
// chan down to its CharmReader's Read method. If you call EnableWaitForAbort, the
// next call to Read will block until either the abort chan is closed (in which case
// it will return an error) or the stopWaiting chan is closed (in which case it
// will return the charmArchive).
func (ch *charmReader) EnableWaitForAbort() (stopWaiting chan struct{}) {
	stopWaiting = make(chan struct{})
	ch.stopWaiting = stopWaiting
	return stopWaiting
}

// Read implements the CharmReader interface.
func (ch *charmReader) Read(info charm.CharmInfo, abort <-chan struct{}) (charm.CharmArchive, error) {
	charmArchive, ok := ch.charmArchives[info.URL().String()]
	if !ok {
		return nil, fmt.Errorf("no such charm!")
	}
	if ch.stopWaiting != nil {
		// EnableWaitForAbort is a one-time wait; make sure we clear it.
		defer func() { ch.stopWaiting = nil }()
		select {
		case <-abort:
			return nil, fmt.Errorf("charm read aborted")
		case <-ch.stopWaiting:
			// We can stop waiting for the abort chan and return the charmArchive.
		}
	}
	return charmArchive, nil
}

func (ch *charmReader) AddCustomCharmArchive(c *gc.C, url *corecharm.URL, customize func(path string)) charm.CharmInfo {
	base := c.MkDir()
	dirpath := testcharms.Repo.ClonedDirPath(base, "dummy")
	if customize != nil {
		customize(dirpath)
	}
	dir, err := corecharm.ReadCharmDir(dirpath)
	c.Assert(err, jc.ErrorIsNil)
	err = dir.SetDiskRevision(url.Revision)
	c.Assert(err, jc.ErrorIsNil)
	bunpath := filepath.Join(base, "charmArchive")
	file, err := os.Create(bunpath)
	c.Assert(err, jc.ErrorIsNil)
	defer func() { _ = file.Close() }()
	err = dir.ArchiveTo(file)
	c.Assert(err, jc.ErrorIsNil)
	charmArchive, err := corecharm.ReadCharmArchive(bunpath)
	c.Assert(err, jc.ErrorIsNil)
	return ch.AddCharmArchive(c, url, charmArchive)
}

func (ch *charmReader) AddCharmArchive(c *gc.C, url *corecharm.URL, charmArchive charm.CharmArchive) charm.CharmInfo {
	if ch.charmArchives == nil {
		ch.charmArchives = map[string]charm.CharmArchive{}
	}
	ch.charmArchives[url.String()] = charmArchive
	return &charmInfo{nil, url}
}

type charmInfo struct {
	charm.CharmInfo
	url *corecharm.URL
}

func (info *charmInfo) URL() *corecharm.URL {
	return info.url
}

type mockCharmArchive struct {
	paths  set.Strings
	expand func(dir string) error
}

func (b mockCharmArchive) ArchiveMembers() (set.Strings, error) {
	// TODO(dfc) this looks like set.Strings().Duplicate()
	return set.NewStrings(b.paths.Values()...), nil
}

func (b mockCharmArchive) ExpandTo(dir string) error {
	if b.expand != nil {
		return b.expand(dir)
	}
	return nil
}

func charmURL(revision int) *corecharm.URL {
	baseURL := corecharm.MustParseURL("cs:s/c")
	return baseURL.WithRevision(revision)
}
