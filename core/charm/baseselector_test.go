// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/series"
	"github.com/juju/juju/version"
)

type baseSelectorSuite struct{}

var _ = gc.Suite(&baseSelectorSuite{})

var (
	bionic      = series.MustParseBaseFromString("ubuntu@18.04")
	cosmic      = series.MustParseBaseFromString("ubuntu@18.10")
	disco       = series.MustParseBaseFromString("ubuntu@19.04")
	jammy       = series.MustParseBaseFromString("ubuntu@22.04")
	precise     = series.MustParseBaseFromString("ubuntu@14.04")
	utopic      = series.MustParseBaseFromString("ubuntu@16.10")
	vivid       = series.MustParseBaseFromString("ubuntu@17.04")
	latest      = series.LatestLTSBase()
	jujuDefault = version.DefaultSupportedLTSBase()
)

func (s *baseSelectorSuite) TestCharmSeries(c *gc.C) {

	deploySeriesTests := []struct {
		title string
		BaseSelector
		expectedBase series.Base
		err          string
	}{
		{
			// Simple selectors first, no supported bases, check we're validating
			title: "juju deploy simple   # no default base, no supported base",
			BaseSelector: BaseSelector{
				Conf: defaultBase{},
			},
			err: "charm does not define any bases, not valid",
		}, {
			title: "juju deploy simple   # default base set",
			BaseSelector: BaseSelector{
				Conf:               defaultBase{"ubuntu@18.04", true},
				SupportedBases:     []series.Base{bionic, cosmic},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			expectedBase: bionic,
		},
		{
			title: "juju deploy simple with old base  # default base set, no supported base",
			BaseSelector: BaseSelector{
				Conf:               defaultBase{"ubuntu@15.10", true},
				SupportedBases:     []series.Base{bionic, cosmic},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			err: `base "ubuntu@15.10" not supported by charm, the charm supported series are: ubuntu@18.04, ubuntu@18.10`,
		},
		{
			title: "juju deploy simple --base=ubuntu@14.04   # no supported base",
			BaseSelector: BaseSelector{
				BaseFlag:           precise,
				Conf:               defaultBase{},
				SupportedBases:     []series.Base{bionic, cosmic},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			err: "base: \"ubuntu@14.04\" not supported",
		}, {
			title: "juju deploy simple --base=ubuntu@18.04   # default base set, no supported base, no supported juju base",
			BaseSelector: BaseSelector{
				BaseFlag: bionic,
				Conf:     defaultBase{"ubuntu@15.10", true},
			},
			err: "charm does not define any bases, not valid",
		},
		{
			title: "juju deploy simple --base=ubuntu@18.04   # user provided base takes presedence over default base ",
			BaseSelector: BaseSelector{
				BaseFlag:           bionic,
				Conf:               defaultBase{"ubuntu@15.10", true},
				SupportedBases:     []series.Base{bionic, cosmic},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			expectedBase: bionic,
		},

		// Now charms with supported base.

		{
			title: "juju deploy multiseries   # use charm default, nothing specified, no default base",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{bionic, cosmic},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			expectedBase: bionic,
		},
		{
			title: "juju deploy multiseries with invalid base  # use charm default, nothing specified, no default base",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{precise, bionic, cosmic},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			expectedBase: bionic,
		},
		{
			title: "juju deploy multiseries with invalid series  # use charm default, nothing specified, no default base",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{precise},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			err: `the charm defined base(s) "ubuntu@14.04" not supported`,
		},
		{
			title: "juju deploy multiseries   # use charm defaults used if default base doesn't match, nothing specified",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{bionic, cosmic},
				Conf:               defaultBase{"ubuntu@15.10", true},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			err: `base "ubuntu@15.10" not supported by charm, the charm supported series are: ubuntu@18.04, ubuntu@18.10`,
		},
		{
			title: "juju deploy multiseries   # use model base defaults if supported by charm",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{bionic, cosmic, disco},
				Conf:               defaultBase{"ubuntu@19.04", true},
				SupportedJujuBases: []series.Base{bionic, cosmic, disco},
			},
			expectedBase: disco,
		},
		{
			title: "juju deploy multiseries   # use model base defaults if supported by charm",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{bionic, cosmic, disco},
				Conf:               defaultBase{"ubuntu@19.04", true},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			err: `base: "ubuntu@19.04" not supported`,
		},
		{
			title: "juju deploy multiseries --base=ubuntu@18.04   # use supported requested",
			BaseSelector: BaseSelector{
				BaseFlag:           bionic,
				SupportedBases:     []series.Base{utopic, vivid, bionic},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			expectedBase: bionic,
		},
		{
			title: "juju deploy multiseries --base=ubuntu@18.04   # use supported requested",
			BaseSelector: BaseSelector{
				BaseFlag:           bionic,
				SupportedBases:     []series.Base{cosmic, bionic},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			expectedBase: bionic,
		},
		{
			title: "juju deploy multiseries --base=ubuntu@18.04   # unsupported requested",
			BaseSelector: BaseSelector{
				BaseFlag:           bionic,
				SupportedBases:     []series.Base{utopic, vivid},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic},
			},
			err: `base "ubuntu@18.04" not supported by charm, the charm supported series are: ubuntu@16.10, ubuntu@17.04`,
		},
		{
			title: "juju deploy multiseries    # fallback to series.LatestLTSBase()",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{utopic, vivid, latest},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic, latest},
			},
			expectedBase: latest,
		},
		{
			title: "juju deploy multiseries    # fallback to version.DefaultSupportedLTSBase()",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{utopic, vivid, jujuDefault},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic, jujuDefault},
			},
			expectedBase: jujuDefault,
		},
		{
			title: "juju deploy multiseries    # prefer series.LatestLTSBase() to  version.DefaultSupportedLTSBase()",
			BaseSelector: BaseSelector{
				SupportedBases:     []series.Base{utopic, vivid, jujuDefault, latest},
				Conf:               defaultBase{},
				SupportedJujuBases: []series.Base{bionic, cosmic, jujuDefault, latest},
			},
			expectedBase: jujuDefault,
		},
	}

	// Use bionic for LTS for all calls.
	previous := series.SetLatestLtsForTesting("bionic")
	defer series.SetLatestLtsForTesting(previous)

	for i, test := range deploySeriesTests {
		c.Logf("test %d [%s]", i, test.title)
		test.BaseSelector.Logger = &noOpLogger{}
		base, err := test.BaseSelector.CharmBase()
		if test.err != "" {
			c.Check(err, gc.ErrorMatches, test.err)
		} else {
			c.Check(err, jc.ErrorIsNil)
			c.Check(base, gc.Equals, test.expectedBase)
		}
	}
}

func (s *baseSelectorSuite) TestValidate(c *gc.C) {
	deploySeriesTests := []struct {
		title    string
		selector BaseSelector
		err      string
	}{
		{
			title: "should fail when image-id constraint is used and no base is explicitly set",
			selector: BaseSelector{
				Conf: defaultBase{
					explicit: false,
				},
				UsingImageID: true,
			},
			err: "base must be explicitly provided when image-id constraint is used",
		},
		{
			title: "should return no errors when using image-id and base flag",
			selector: BaseSelector{
				Conf: defaultBase{
					explicit: false,
				},
				BaseFlag:     jammy,
				UsingImageID: true,
			},
		},
		{
			title: "should return no errors when using image-id and charms url base is set",
			selector: BaseSelector{
				Conf: defaultBase{
					explicit: false,
				},
				UsingImageID: true,
			},
		},
		{
			title: "should return no errors when using image-id and explicit base from conf",
			selector: BaseSelector{
				Conf: defaultBase{
					explicit: true,
				},
				UsingImageID: true,
			},
		},
	}

	for i, test := range deploySeriesTests {
		c.Logf("test %d [%s]", i, test.title)
		test.selector.Logger = &noOpLogger{}
		err := test.selector.validate()
		if test.err != "" {
			c.Check(err, gc.ErrorMatches, test.err)
		} else {
			c.Check(err, jc.ErrorIsNil)
		}
	}
}
