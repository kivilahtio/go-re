/*
This file is part of go-re

Copyright © 2021 Technology Innovation Institute, United Arab Emirates

Licensed under the Artistic License, Version 2.0 (the "License");
    https://www.perlfoundation.org/artistic-license-20
*/

package t

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	re0 "github.com/kivilahtio/go-re/v0"
	re1 "github.com/kivilahtio/go-re/v1"
	re2 "github.com/kivilahtio/go-re/v2"
)

/*
Compliance testing for https://go.dev/about
*/
func TestSemanticImportVersions(t *testing.T) {
	Convey("v0", t, func() {
		So(re0.M("Regexp like a Perl Pumpking in Go!", `/Perl/`), ShouldBeTrue)
	})
	Convey("v1", t, func() {
		So(re1.M("Regexp like a Perl Pumpking in Go!", `/Perl/`), ShouldBeTrue)
	})
	Convey("v2", t, func() {
		So(re2.M("Regexp like a Perl Pumpking in Go!", `/Perl/`), ShouldBeTrue)
	})
}
