// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package checkers_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	gc "launchpad.net/gocheck"

	jc "launchpad.net/juju-core/testing/checkers"
)

type FileSuite struct{}

var _ = gc.Suite(&FileSuite{})

func (s *FileSuite) TestIsNonEmptyFile(c *gc.C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, gc.IsNil)
	fmt.Fprintf(file, "something")
	file.Close()

	c.Assert(file.Name(), jc.IsNonEmptyFile)
}

func (s *FileSuite) TestIsNonEmptyFileWithEmptyFile(c *gc.C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, gc.IsNil)
	file.Close()

	result, message := jc.IsNonEmptyFile.Check([]interface{}{file.Name()}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, file.Name()+" is empty")
}

func (s *FileSuite) TestIsNonEmptyFileWithMissingFile(c *gc.C) {
	name := filepath.Join(c.MkDir(), "missing")

	result, message := jc.IsNonEmptyFile.Check([]interface{}{name}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, name+" does not exist")
}

func (s *FileSuite) TestIsNonEmptyFileWithNumber(c *gc.C) {
	result, message := jc.IsNonEmptyFile.Check([]interface{}{42}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, "obtained value is not a string and has no .String(), int:42")
}

func (s *FileSuite) TestIsDirectory(c *gc.C) {
	dir := c.MkDir()
	c.Assert(dir, jc.IsDirectory)
}

func (s *FileSuite) TestIsDirectoryMissing(c *gc.C) {
	absentDir := filepath.Join(c.MkDir(), "foo")

	result, message := jc.IsDirectory.Check([]interface{}{absentDir}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, absentDir+" does not exist")
}

func (s *FileSuite) TestIsDirectoryWithFile(c *gc.C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, gc.IsNil)
	file.Close()

	result, message := jc.IsDirectory.Check([]interface{}{file.Name()}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, file.Name()+" is not a directory")
}

func (s *FileSuite) TestIsDirectoryWithNumber(c *gc.C) {
	result, message := jc.IsDirectory.Check([]interface{}{42}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, "obtained value is not a string and has no .String(), int:42")
}

func (s *FileSuite) TestDoesNotExist(c *gc.C) {
	absentDir := filepath.Join(c.MkDir(), "foo")
	c.Assert(absentDir, jc.DoesNotExist)
}

func (s *FileSuite) TestDoesNotExistWithPath(c *gc.C) {
	dir := c.MkDir()
	result, message := jc.DoesNotExist.Check([]interface{}{dir}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, dir+" exists")
}

func (s *FileSuite) TestDoesNotExistWithNumber(c *gc.C) {
	result, message := jc.DoesNotExist.Check([]interface{}{42}, nil)
	c.Assert(result, jc.IsFalse)
	c.Assert(message, gc.Equals, "obtained value is not a string and has no .String(), int:42")
}
