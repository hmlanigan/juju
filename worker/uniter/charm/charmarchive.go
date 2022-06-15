// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm

import (
	"net/url"
	"os"
	"path"

	"github.com/juju/charm/v8"
	"github.com/juju/errors"

	"github.com/juju/juju/downloader"
)

// Download exposes the downloader.Download methods needed here.
type Downloader interface {
	// Download starts a new charm archive download, waits for it to
	// complete, and returns the local name of the file.
	Download(req downloader.Request) (string, error)
}

// CharmArchiveDir is responsible for storing and retrieving charm archive
// identified by state charms.
type CharmArchiveDir struct {
	path       string
	downloader Downloader
	logger     Logger
}

// NewCharmArchiveDir returns a new CharmArchiveDir which uses path for storage.
func NewCharmArchiveDir(path string, dlr Downloader, logger Logger) *CharmArchiveDir {
	if dlr == nil {
		dlr = downloader.New(downloader.NewArgs{
			HostnameVerification: false,
		})
	}
	return &CharmArchiveDir{
		path:       path,
		downloader: dlr,
		logger:     logger,
	}
}

// Read returns a charm charmArchive from the directory. If no charmArchive exists yet,
// one will be downloaded and validated and copied into the directory before
// being returned. Downloads will be aborted if a value is received on abort.
func (d *CharmArchiveDir) Read(info CharmInfo, abort <-chan struct{}) (CharmArchive, error) {
	path := d.archivePath(info)
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := d.download(info, path, abort); err != nil {
			return nil, err
		}
	}
	return charm.ReadCharmArchive(path)
}

// download fetches the supplied charm and checks that it has the correct sha256
// hash, then copies it into the directory. If a value is received on abort, the
// download will be stopped.
func (d *CharmArchiveDir) download(info CharmInfo, target string, abort <-chan struct{}) (err error) {
	// First download...
	curl, err := url.Parse(info.URL().String())
	if err != nil {
		return errors.Annotate(err, "could not parse charm URL")
	}
	expectedSha256, err := info.ArchiveSha256()
	req := downloader.Request{
		URL:       curl,
		TargetDir: downloadsPath(d.path),
		Verify:    downloader.NewSha256Verifier(expectedSha256),
		Abort:     abort,
	}
	d.logger.Infof("downloading %s from API server", info.URL())
	filename, err := d.downloader.Download(req)
	if err != nil {
		return errors.Annotatef(err, "failed to download charm %q from API server", info.URL())
	}
	defer errors.DeferredAnnotatef(&err, "downloaded but failed to copy charm to %q from %q", target, filename)

	// ...then move the right location.
	if err := os.MkdirAll(d.path, 0755); err != nil {
		return errors.Trace(err)
	}
	if err := os.Rename(filename, target); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// archivePath returns the path to the location where the verified charm
// archive identified by info will be, or has been, saved.
func (d *CharmArchiveDir) archivePath(info CharmInfo) string {
	return d.charmURLPath(info.URL())
}

// charmURLPath returns the path to the location where the verified charm
// charmArchive identified by url will be, or has been, saved.
func (d *CharmArchiveDir) charmURLPath(url *charm.URL) string {
	return path.Join(d.path, charm.Quote(url.String()))
}

// ClearDownloads removes any entries in the temporary charmArchive download
// directory. It is intended to be called on uniter startup.
func ClearDownloads(charmArchiveDir string) error {
	downloadDir := downloadsPath(charmArchiveDir)
	err := os.RemoveAll(downloadDir)
	return errors.Annotate(err, "unable to clear charmArchive downloads")
}

// downloadsPath returns the path to the directory into which charms are
// downloaded.
func downloadsPath(bunsDir string) string {
	return path.Join(bunsDir, "downloads")
}
