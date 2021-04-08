package re

import (
	"fmt"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var nilCapture []string
var nilNameCap map[string]string // = make(map[string]string, 0)

func Example() {
	if M("kalle ankka", `m/(a.)/g`) {
		fmt.Printf("Plain match got it!\n")
		fmt.Printf("First matching group '%s'!\n", R0.S[1])
	}

	if r := Mr("kalle ankka", `m/(a.)/g`); r.Matches > 0 {
		fmt.Printf("Capture 1: '%s'\n", r.S[1])
		fmt.Printf("Capture 2: '%s'\n", r.S[2])
	}

	str := "This is how we party!\nFooling with our body.\n"
	if S(&str, `s! # Do a substitution operation, start the matching pattern
	(?:party|body) # Non-capturing substitution happening here
	!              # close the matching pattern and begin with the substituting pattern
	code
	!gx            # cinally some flags, x is important to allow documenting the regexp
	`) {
		fmt.Printf("Substitution: '%s'\n", str)
		fmt.Printf("Non-capturing group does not get captured '%d'\n", len(R0.S))
	}

	str = "Geeks shall inherit the properties and methods of object Earth"
	r := Sr(&str, `s!([aeiouy])!!gi`)
	fmt.Printf("Substitution: '%s'\n", str)
	fmt.Printf("Capture 1: '%s'\n", r.S[1])
	fmt.Printf("Capture 2: '%s'\n", r.S[2])
	fmt.Printf("Capture 18: '%s'\n", r.S[18])
	fmt.Printf("Capture 19: '%s'\n", r.S[19])

	str = Ss("Hyppytyynytyydytys", `s!Hyppytyynytyydytys!Kolmivaihekilowattituntimittari!`)
	fmt.Printf("Substitution: '%s'", str)

	// Output: Plain match got it!
	// First matching group 'al'!
	// Capture 1: 'al'
	// Capture 2: 'an'
	// Substitution: 'This is how we code!
	// Fooling with our code.
	// '
	// Non-capturing group does not get captured '0'
	// Substitution: 'Gks shll nhrt th prprts nd mthds f bjct rth'
	// Capture 1: 'e'
	// Capture 2: 'e'
	// Capture 18: 'E'
	// Capture 19: 'a'
	// Substitution: 'Kolmivaihekilowattituntimittari'
}

func ExampleR() {
	str := "Peruspalveluliikelaitoskuntayhtymä"
	if R(&str, `s/laitos/toimi/`) {
		fmt.Printf("%s\n", str) // STDOUT => "Peruspalveluliiketoimikuntayhtymä"
	}
	if R(&str, `m!toimi!`) {
		fmt.Printf("%d\n", R0.Matches) // STDOUT => "1"
	}
	// Output: Peruspalveluliiketoimikuntayhtymä
	// 1
}

func TestUnit(t *testing.T) {
	SetDefaultFailureMode(FailureContinues)
	Convey("matching and capturing", t, func() {
		var r *RE
		_ = r
		r = runMatchTest(" [28] ", `m!(?P<array>\[\d*\])!`, 1, []string{"", "[28]"}, map[string]string{"array": "[28]"})
		r = runMatchTest("[28]", `m!(?P<array>\[\d*\])?!`, 1, []string{"", "[28]"}, map[string]string{"array": "[28]"}) // Why leading whitespace makes this match fail?
		r = runMatchTest("kalle ankka", `m/kalle ankka/`, 1, nilCapture, nilNameCap)
		r = runMatchTest("kalle/ankka", `m/kalle\/ankka/`, 1, nilCapture, nilNameCap)
		r = runMatchTest("kalle ankka", `m/a/`, 1, nilCapture, nilNameCap)
		r = runMatchTest("kalle ankka", `m/a/g`, 3, nilCapture, nilNameCap)
		r = runMatchTest("kalle ankka", `m/(a.)/g`, 2, []string{"", "al", "an"}, nilNameCap)
		r = runMatchTest("bubbelbubbe", `m/(a.)/g`, 0, nilCapture, nilNameCap) // The regexp is cached and now it shouldn't match after a succesful match
		r = runMatchTest("kalle ankka", `m/(?:a.)/g`, 2, nilCapture, nilNameCap)
		r = runMatchTest("kalle ankka", `m!a!g`, 3, nilCapture, nilNameCap)
		r = runMatchTest("kalle#ankka", `m!([ #])!g`, 1, []string{"", "#"}, nilNameCap)
		r = runMatchTest("kaLlE AnKka", `m!(?P<aleph>[e])!gi`, 1, []string{"", "E"}, map[string]string{"aleph": "E"})
	})
	Convey("substitution", t, func() {
		var r *RE
		_ = r

		r = runSubstTest("kalle ankka", `s/kalle ankka/minni hiiri/`, `minni hiiri`, 1, nilCapture, nilNameCap)
		r = runSubstTest("kalle/ankka", `s/\///`, `kalleankka`, 1, nilCapture, nilNameCap)
		r = runSubstTest("kalle ankka", `s!a!u!`, `kulle ankka`, 1, nilCapture, nilNameCap)
		r = runSubstTest("kalle ankka", `s!a!u!g`, `kulle unkku`, 3, nilCapture, nilNameCap)
		r = runSubstTest("kalle ankka", `s/(a.)/--/g`, `k--le --kka`, 2, []string{"", "al", "an"}, nilNameCap)
		r = runSubstTest("kaLlE AnKka", `s!(?P<aleph>[e])!!gi`, "kaLl AnKka", 1, []string{"", "E"}, map[string]string{"aleph": "E"})
	})
	Convey("Ignore whitespace regex, ISO8601 parser", t, func() {
		r := Mr(`
kalle: "ankka" - 2021-12-31 23:59:59.0123+0200Z
paavo: "pesus" - 2020-10-31T21:39:39.4321+0230`,
			`m/
			# this is a comment which encompasses the whole line
				^
				(?P<username>\w+) # This is a comment too
				:\s+
				"(?P<surname>\w+)"
				\s+ - \s+
				(?P<year>\d{4})
				-
				(?P<month>\d{2})
				-
				(?P<day>\d{2})
				[T #]
				(?P<hour>\d{2})
				:
				(?P<minute>\d{2})
				:
				(?P<second>\d{2})
				(?:\.
				(?P<decimal>\d{1,4}))?
				(?:(?P<timezone>[+-]\d{2,4})Z?)?
				$
			/xgms`)
		So(r.S[1], ShouldEqual, "kalle")
		So(r.S[2], ShouldEqual, "ankka")
		So(r.S[3], ShouldEqual, "2021")
		So(r.S[4], ShouldEqual, "12")
		So(r.S[5], ShouldEqual, "31")
		So(r.S[6], ShouldEqual, "23")
		So(r.S[7], ShouldEqual, "59")
		So(r.S[8], ShouldEqual, "59")
		So(r.S[9], ShouldEqual, "0123")
		So(r.S[10], ShouldEqual, "+0200")
		So(r.S[11], ShouldEqual, "paavo")
		//...//
		So(r.S[20], ShouldEqual, "+0230")
		So(r.Z["username"], ShouldEqual, "paavo")
		So(r.Z["surname"], ShouldEqual, "pesus")
		So(r.Z["year"], ShouldEqual, "2020")
		So(r.Z["month"], ShouldEqual, "10")
		So(r.Z["day"], ShouldEqual, "31")
		So(r.Z["hour"], ShouldEqual, "21")
		So(r.Z["minute"], ShouldEqual, "39")
		So(r.Z["second"], ShouldEqual, "39")
		So(r.Z["decimal"], ShouldEqual, "4321")
		So(r.Z["timezone"], ShouldEqual, "+0230")
	})
	Convey("regexp cache", t, func() {
		So(len(regexpCache.cache), ShouldBeGreaterThan, 10)
		regexpCache.Flush()
		So(len(regexpCache.cache), ShouldEqual, 0)
	})
}

func runMatchTest(haystack string, needle string, matches int, captures []string, nCaptures map[string]string) *RE {
	var r *RE
	Convey(fmt.Sprintf(`"%s" match "%s" '%d' times! Captures='%d', Named captures='%d'`, needle, haystack, matches, len(captures)-1, len(nCaptures)), func() {
		r = Mr(haystack, needle)

		So(r.Matches, ShouldEqual, matches)
		So(R0.Matches, ShouldEqual, matches)
		So(r.S, ShouldResemble, captures)
		So(R0.S, ShouldResemble, captures)
		So(r.Z, ShouldResemble, nCaptures)
		So(R0.Z, ShouldResemble, nCaptures)
	})

	return r
}

func runSubstTest(haystackA string, needle string, expected string, matches int, captures []string, nCaptures map[string]string) *RE {
	haystack := &haystackA
	var r *RE
	Convey(fmt.Sprintf(`"%s" substituted "%s" '%d' times! Captures='%d', Named captures='%d'`, needle, *haystack, matches, len(captures)-1, len(nCaptures)), func() {
		r = Sr(haystack, needle)

		So(*haystack, ShouldEqual, expected)
		So(r.Matches, ShouldEqual, matches)
		So(R0.Matches, ShouldEqual, matches)
		So(r.S, ShouldResemble, captures)
		So(R0.S, ShouldResemble, captures)
		So(r.Z, ShouldResemble, nCaptures)
		So(R0.Z, ShouldResemble, nCaptures)
	})

	return r
}

func BenchmarkRECache_on_SimpleRE(b *testing.B) {
	UseRECache = true
	benchmarkRECacheComplexRE(b)
}
func BenchmarkRECache_on_ComplexRE(b *testing.B) {
	UseRECache = true
	benchmarkRECacheComplexRE(b)
}
func BenchmarkRECache_off_SimpleRE(b *testing.B) {
	UseRECache = false
	benchmarkRECacheComplexRE(b)
}
func BenchmarkRECache_off_ComplexRE(b *testing.B) {
	UseRECache = false
	benchmarkRECacheComplexRE(b)
}
func benchmarkRECacheSimpleRE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		M("kalle ankka", `m/kalle ankka/`)
		if R0.Matches == 0 {
			b.Errorf("BenchmarkRECacheSimpleRE regexp doesnt match?")
		}
	}
}
func benchmarkRECacheComplexRE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		M(`
kalle: "ankka" - 2021-12-31 23:59:59.0123+0200Z
paavo: "pesus" - 2020-10-31T21:39:39.4321+0230`,
			`m/
			# this is a comment which encompasses the whole line
				^
				(?P<username>\w+) # This is a comment too
				:\s+
				"(?P<surname>\w+)"
				\s+ - \s+
				(?P<year>\d{4})
				-
				(?P<month>\d{2})
				-
				(?P<day>\d{2})
				[T ]
				(?P<hour>\d{2})
				:
				(?P<minute>\d{2})
				:
				(?P<second>\d{2})
				(?:\.
				(?P<decimal>\d{1,4}))?
				(?:(?P<timezone>[+-]\d{2,4})Z?)?
				$
			/xgms`)
		if R0.Matches == 0 {
			b.Errorf("BenchmarkRECacheComplexRE regexp doesnt match?")
		}
	}
}
func BenchmarkNativeGolangRegexNamedCaptureGroup(b *testing.B) {
	parseISO8601Regexp := regexp.MustCompile(
		`^` +
			`(?P<username>\w+)` +
			`:\s+` +
			`"(?P<surname>\w+)"` +
			`\s+-\s+` +
			`(?P<year>\d{4})` +
			`-` +
			`(?P<month>\d{2})` +
			`-` +
			`(?P<day>\d{2})` +
			`[T ]` +
			`(?P<hour>\d{2})` +
			`:` +
			`(?P<minute>\d{2})` +
			`:` +
			`(?P<second>\d{2})` +
			`(?:\.` +
			`(?P<decimal>\d{1,4}))?` +
			`(?:(?P<timezone>[+-]\d{2,4})Z?)?` +
			`$`)
	for i := 0; i < b.N; i++ {
		benchmarkNativeGolangRegexNamedCaptureGroup(parseISO8601Regexp, b)
	}
}
func benchmarkNativeGolangRegexNamedCaptureGroup(parseISO8601Regexp *regexp.Regexp, b *testing.B) {
	capture, err := parseNamedCaptureGroupsRegex(`kalle: "ankka" - 2021-12-31 23:59:59.0123+0200Z`, parseISO8601Regexp)
	if err != nil {
		b.Errorf("%+v\n", err)
	}
	if len(capture) == 0 {
		b.Errorf("No captures?\n")
	}
}
func parseNamedCaptureGroupsRegex(textRow string, regexpStr *regexp.Regexp) (map[string]string, error) {
	subexpNames := regexpStr.SubexpNames()
	namedCaptureGroups := make(map[string]string, len(subexpNames))
	matches := regexpStr.FindAllStringSubmatch(textRow, -1)
	if matches == nil {
		return namedCaptureGroups, fmt.Errorf("Unable to parse text '%s' using regexp '%s'\n", textRow, regexpStr)
	}
	for _, match := range matches {
		for groupIdx, group := range match {
			if group != "" {
				namedCaptureGroups[subexpNames[groupIdx]] = group
			}
		}
	}
	return namedCaptureGroups, nil
}
