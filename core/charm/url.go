// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm

import (
	"fmt"
	"strings"

	"github.com/juju/charm/v8"
	"github.com/juju/errors"

	"github.com/juju/juju/core/series"
)

// CharmURLSeriesToBase converts a charm url that contains a series, to one
// that contains a base.
func CharmURLSeriesToBase(url *charm.URL) (*charm.URL, error) {
	if url.Series == "" {
		return url, nil
	}

	// Kubernetes is not a valid series for a base.  Instead use the latest
	// LTS version of ubuntu.
	urlSeries := url.Series
	if urlSeries == "kubernetes" {
		urlSeries = series.LatestLts()
	}

	baseNameType, err := series.GetOSFromSeries(urlSeries)
	if err != nil {
		return nil, errors.Annotatef(err, "os name invalid")
	}

	baseVersion, err := series.SeriesVersion(urlSeries)
	if err != nil {
		return nil, errors.Annotatef(err, "version invalid")
	}

	baseName := strings.ToLower(baseNameType.String())
	return url.WithSeries("").WithBase(fmt.Sprintf("%s:%s", baseName, baseVersion)), nil
}
