[![Go Reference](https://pkg.go.dev/badge/github.com/kivilahtio/go-re.svg)](https://pkg.go.dev/github.com/kivilahtio/go-re)

# Regexp like a Perl Pumpking in Go! 

## Installation

```
go get github.com/kivilahtio/go-re/v0
```

## Usage

See godoc examples

### Example ISO8601 parser
```
require . "github.com/kivilahtio/go-re/v0"

func init() {
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

    _ = r.Z["timezone"] == R0.Z["timezone"]
}
```
