// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package user

import (
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type userSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&userSuite{})

func (s *userSuite) TestIsValidUser(c *gc.C) {
	for i, t := range []struct {
		string string
		expect bool
	}{
		{"", false},
		{"bob", true},
		{"Bob", true},
		{"bOB", true},
		{"b^b", false},
		{"bob1", true},
		{"bob-1", true},
		{"bob+1", true},
		{"bob+", false},
		{"+bob", false},
		{"bob.1", true},
		{"1bob", true},
		{"1-bob", true},
		{"1+bob", true},
		{"1.bob", true},
		{"jim.bob+99-1.", false},
		{"a", false},
		{"0foo", true},
		{"foo bar", false},
		{"bar{}", false},
		{"bar+foo", true},
		{"bar_foo", false},
		{"bar!", false},
		{"bar^", false},
		{"bar*", false},
		{"foo=bar", false},
		{"foo?", false},
		{"[bar]", false},
		{"'foo'", false},
		{"%bar", false},
		{"&bar", false},
		{"#1foo", false},
		{"bar@ram.u", true},
		{"bar@", false},
		{"@local", false},
		{"not/valid", false},
	} {
		c.Logf("test %d: %s", i, t.string)
		c.Assert(IsValidName(t.string), gc.Equals, t.expect, gc.Commentf("%s", t.string))
	}
}

func (s *userSuite) TestNewName(c *gc.C) {
	for i, t := range []struct {
		input   string
		name    string
		isLocal bool
		domain  string
	}{
		{
			input:   "bob",
			name:    "bob",
			isLocal: true,
			domain:  "",
		}, {
			input:   "bob@local",
			name:    "bob",
			isLocal: true,
			domain:  "local",
		}, {
			input:   "bob@foo",
			name:    "bob@foo",
			isLocal: false,
			domain:  "foo",
		}} {
		c.Logf("test %d: %s", i, t.input)
		name, err := NewName(t.input)
		c.Assert(err, jc.ErrorIsNil)
		c.Check(name.Name(), gc.Equals, t.name)
		c.Check(name.IsLocal(), gc.Equals, t.isLocal)
	}
}

func (s *userSuite) TestNewError(c *gc.C) {
	for i, t := range []struct {
		input string
		err   string
	}{
		{
			input: "",
			err:   `user name "" not valid`,
		}, {
			input: "not/valid",
			err:   `user name "not/valid" not valid`,
		}, {
			input: "@foo",
			err:   `user name "@foo" not valid`,
		}} {
		c.Logf("test %d: %s", i, t.input)
		_, err := NewName(t.input)
		c.Assert(err, gc.ErrorMatches, t.err)
	}
}
