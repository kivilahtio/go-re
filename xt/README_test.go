package xt

import (
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*
Compliance testing for https://go.dev/about
*/
func TestREADME(t *testing.T) {
	var readme []byte
	var err error

	Convey("README.md is readable", t, func() {
		readme, err = ioutil.ReadFile("../README.md")
		So(err, ShouldBeNil)
	})

	Convey("README.md has a pkg.go.dev badge", t, func() {
		So(strings.Contains(string(readme), "[![Go Reference]"), ShouldBeTrue)
	})
}
