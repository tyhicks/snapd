// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package snappy

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"

	"github.com/ubuntu-core/snappy/pkg"
)

type SecurityTestSuite struct {
	buildDir              string
	m                     *packageYaml
	scFilterGenCall       []string
	scFilterGenCallReturn []byte
}

var _ = Suite(&SecurityTestSuite{})

func (a *SecurityTestSuite) SetUpTest(c *C) {
	a.buildDir = c.MkDir()
	os.MkdirAll(filepath.Join(a.buildDir, "meta"), 0755)

	a.m = &packageYaml{
		Name:        "foo",
		Version:     "1.0",
		Integration: make(map[string]clickAppHook),
	}
}

func (a *SecurityTestSuite) TestSnappyGetSecurityProfile(c *C) {
	m := packageYaml{
		Name:    "foo",
		Version: "1.0",
	}
	b := Binary{Name: "bin/app"}
	ap, err := getSecurityProfile(&m, b.Name, "/apps/foo.mvo/1.0/")
	c.Assert(err, IsNil)
	c.Check(ap, Equals, "foo.mvo_bin-app_1.0")
}

func (a *SecurityTestSuite) TestSnappyGetSecurityProfileInvalid(c *C) {
	m := packageYaml{
		Name:    "foo",
		Version: "1.0",
	}
	b := Binary{Name: "bin/app"}
	_, err := getSecurityProfile(&m, b.Name, "/apps/foo/1.0/")
	c.Assert(err, Equals, ErrInvalidPart)
}

func (a *SecurityTestSuite) TestSnappyGetSecurityProfileFramework(c *C) {
	m := packageYaml{
		Name:    "foo",
		Version: "1.0",
		Type:    pkg.TypeFramework,
	}
	b := Binary{Name: "bin/app"}
	ap, err := getSecurityProfile(&m, b.Name, "/apps/foo.mvo/1.0/")
	c.Assert(err, IsNil)
	c.Check(ap, Equals, "foo_bin-app_1.0")
}

func (a *SecurityTestSuite) TestSnappyFindUbuntuVersion(c *C) {
	realLsbRelease := lsbRelease
	defer func() { lsbRelease = realLsbRelease }()

	lsbRelease = filepath.Join(c.MkDir(), "mock-lsb-release")
	s := `DISTRIB_RELEASE=18.09`
	err := ioutil.WriteFile(lsbRelease, []byte(s), 0644)
	c.Assert(err, IsNil)

	ver, err := findUbuntuVersion()
	c.Assert(err, IsNil)
	c.Assert(ver, Equals, "18.09")
}

func (a *SecurityTestSuite) TestSnappyFindUbuntuVersionNotFound(c *C) {
	realLsbRelease := lsbRelease
	defer func() { lsbRelease = realLsbRelease }()

	lsbRelease = filepath.Join(c.MkDir(), "mock-lsb-release")
	s := `silly stuff`
	err := ioutil.WriteFile(lsbRelease, []byte(s), 0644)
	c.Assert(err, IsNil)

	_, err = findUbuntuVersion()
	c.Assert(err, Equals, ErrSystemVersionNotFound)
}

func (a *SecurityTestSuite) TestSecurityGenDbusPath(c *C) {
	c.Assert(dbusPath("foo"), Equals, "foo")
	c.Assert(dbusPath("foo bar"), Equals, "foo_20bar")
	c.Assert(dbusPath("foo/bar"), Equals, "foo_2fbar")
}

func (a *SecurityTestSuite) TestSecurityFindWhitespacePrefix(c *C) {
	t := `  ###POLICYGROUPS###`
	c.Assert(findWhitespacePrefix(t, "###POLICYGROUPS###"), Equals, "  ")

	t = `not there`
	c.Assert(findWhitespacePrefix(t, "###POLICYGROUPS###"), Equals, "")
}

func (a *SecurityTestSuite) TestSecurityFindTemplateApparmor(c *C) {
	aaPolicyDir = c.MkDir()
	mockTemplate := filepath.Join(aaPolicyDir, "templates", "mock-templ")
	err := os.MkdirAll(filepath.Dir(mockTemplate), 0755)
	c.Assert(err, IsNil)
	err = ioutil.WriteFile(mockTemplate, []byte(`something`), 0644)
	c.Assert(err, IsNil)

	t, err := findTemplate("mock-templ", "apparmor")
	c.Assert(err, IsNil)
	c.Assert(t, Matches, "something")
}

func (a *SecurityTestSuite) TestSecurityFindTemplateApparmorNotFound(c *C) {
	_, err := findTemplate("not-available-templ", "apparmor")
	c.Assert(err, DeepEquals, &errPolicyNotFound{"template", "not-available-templ"})
}
