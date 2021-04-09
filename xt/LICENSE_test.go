/*
This file is part of go-re

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Artistic License, Version 2.0 (the "License");
    https://www.perlfoundation.org/artistic-license-20
*/

package xt

import (
	"io/ioutil"
	"testing"

	"github.com/google/licensecheck"
	. "github.com/smartystreets/goconvey/convey"
)

/*
Compliance testing for https://go.dev/about
*/
func TestRedistributableLicense(t *testing.T) {
	var license []byte
	var err error

	Convey("LICENSE is readable", t, func() {
		license, err = ioutil.ReadFile("../LICENSE")
		So(err, ShouldBeNil)
	})

	Convey("LICENSE is Redistributable by pkg.go.dev", t, func() {
		cov := licensecheck.Scan(license)
		So(cov.Percent, ShouldBeGreaterThan, 90)
		So(cov.Match[0].ID, ShouldEqual, "Artistic-2.0")
	})
}
