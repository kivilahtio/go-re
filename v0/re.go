/*

Regexp like a Perl Pumpking in Go!

Rationale:
 #Perl5:
 $str =~ m/This is how (?<we>party!)/g
 $we = $1;
 $we = $+{we}

 #Golang:
 M(str, `m/This is how (?P<we>party!)/g`)
 we := R0.S[1]   // access the first capture group
 we = R0.Z["we"] // access the named capture group

Flags supported:
 - g
 - x
 - m
 - s
 - i

Notes:

 - Named capture groups are overwritten with the last capture group when using global matching
 - No unicode support!

*/
package re

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type RE struct {
	_orig     string         // original regex string
	f         *string        // altered flags string
	n         *string        // altered regex string
	s         *string        // substitution string in substitute-operation
	regex     *regexp.Regexp // compiled regexp after preprocessing the needle parsed
	mode      byte           // s or m or tr
	separator byte           // separate the mode/matcher/substituter/flags components
	captures  bool           // Enable capture group functionality
	nCaptures bool           // Enable named capture groups functionality. The way Go regexp works with mixing named and non-named groups together makes it difficult to distinguish if named groups are actually used or not. This saves a lot of computation.

	g bool // flag g used
	x bool // flag x used

	Matches int               // how many times the regex matched
	S       []string          // $1, $2, ..., $n Captured subpatterns
	Z       map[string]string // %+ Named capture buffers
}

var R0 *RE // The result of the latest regexp operation. Not thread-safe! It could be if Go had thread-local variables or a way to identify the running thread.

var UseRECache bool = true // Enable/Disable the transparent RECache
/*
Inspired by https://github.com/patrickmn/go-cache

Transparently caches given regexps to save on the expensive computation
*/
type RECache struct {
	cache map[string]*RE
	mu    sync.RWMutex
}

func NewRECache() *RECache {
	return &RECache{
		cache: make(map[string]*RE),
	}
}
func (self *RECache) PutStr(regex string, re *RE) {
	self.Put(&regex, re)
}
func (self *RECache) Put(regex *string, re *RE) {
	copy := *re
	self.mu.Lock()
	self.cache[*regex] = &copy
	self.mu.Unlock()
}
func (self *RECache) GetStr(regex string) *RE {
	return self.Get(&regex)
}
func (self *RECache) Get(regex *string) *RE {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.cache[*regex]
}
func (self *RECache) Flush() {
	self.cache = make(map[string]*RE)
}

var regexpCache = NewRECache()

/*
R is useful if you don't know the type of the regex (match/substitute) beforehand and need to dynamically do things.
*/
func R(haystack *string, needle string) bool {
	r := regexParser(&needle)
	if r.mode == 's' {
		s(haystack, r)
		if r.Matches > 0 {
			return true
		}
	} else {
		m(haystack, r)
		if r.Matches > 0 {
			return true
		}
		return false
	}

	return false
}

func M(haystack string, needle string) bool {
	r := regexParser(&needle)
	m(&haystack, r)
	if r.Matches > 0 {
		return true
	} else {
		return false
	}
}

/*
Mr returns *RE, which is thread-safe
*/
func Mr(haystack string, needle string) *RE {
	return m(&haystack, regexParser(&needle))
}

func S(haystack *string, needle string) bool {
	r := regexParser(&needle)
	s(haystack, r)
	if r.Matches > 0 {
		return true
	} else {
		return false
	}
}

/*
Sr returns *RE, which is thread-safe
*/
func Sr(haystack *string, needle string) *RE {
	return s(haystack, regexParser(&needle))
}

/*
Ss returns the substituted string
*/
func Ss(haystack string, needle string) string {
	s(&haystack, regexParser(&needle))
	return haystack
}

func m(haystack *string, r *RE) *RE {
	R0 = r
	if strings.Contains(*r.f, "g") {
		captureGroups := r.regex.FindAllStringSubmatch(*haystack, -1)
		if captureGroups == nil {
			return r
		}
		r.Matches = len(captureGroups)
		if r.captures {
			namedCaptureGroups := r.regex.SubexpNames()
			if r.nCaptures && len(namedCaptureGroups) > 1 {
				r.Z = make(map[string]string, len(namedCaptureGroups))
			}
			r.S = make([]string, len(captureGroups)*(len(captureGroups[0])-1)+1)
			for i, captures := range captureGroups {
				captureGroup(r, captures, i, namedCaptureGroups)
			}
		}
	} else {
		captures := r.regex.FindStringSubmatch(*haystack)
		if captures == nil {
			return r
		}
		r.Matches = 1
		if r.captures {
			namedCaptureGroups := r.regex.SubexpNames()
			if r.nCaptures && len(namedCaptureGroups) > 1 {
				r.Z = make(map[string]string, len(namedCaptureGroups))
			}
			r.Z = make(map[string]string, len(namedCaptureGroups))
			r.S = make([]string, len(captures))
			captureGroup(r, captures, 0, namedCaptureGroups)
		}
	}
	return r
}

func s(haystack *string, r *RE) *RE {
	R0 = r
	result := []byte{}
	if strings.Contains(*r.f, "g") {
		if r.captures {
			captureGroups := r.regex.FindAllStringSubmatch(*haystack, -1)
			if captureGroups == nil {
				return r
			}
			namedCaptureGroups := r.regex.SubexpNames()
			if r.nCaptures && len(namedCaptureGroups) > 1 {
				r.Z = make(map[string]string, len(namedCaptureGroups))
			}
			r.S = make([]string, len(captureGroups)*(len(captureGroups[0])-1)+1)
			for i, captures := range captureGroups {
				captureGroup(r, captures, i, namedCaptureGroups)
			}
		}
		// For each match of the regex in the content.
		r.Matches = len(r.regex.FindAllStringSubmatchIndex(*haystack, -1))
		if r.Matches > 0 {
			*haystack = r.regex.ReplaceAllString(*haystack, *r.s)
		}
	} else {
		if r.captures {
			captures := r.regex.FindStringSubmatch(*haystack)
			if captures == nil {
				return r
			}
			namedCaptureGroups := r.regex.SubexpNames()
			if r.nCaptures && len(namedCaptureGroups) > 1 {
				r.Z = make(map[string]string, len(namedCaptureGroups))
			}
			r.Z = make(map[string]string, len(namedCaptureGroups))
			r.S = make([]string, len(captures))
			captureGroup(r, captures, 0, namedCaptureGroups)
		}

		captureIdxs := r.regex.FindStringSubmatchIndex(*haystack)
		if len(captureIdxs) > 0 {
			r.Matches = 1
			result = r.regex.ExpandString(nil, *r.s, *haystack, captureIdxs)
			*haystack = (*haystack)[:captureIdxs[0]] + string(result) + (*haystack)[captureIdxs[1]:]
		}
	}
	return r
}

func captureGroup(r *RE, captures []string, captureGroupsIteration int, namedCaptureGroups []string) {
	for j := 1; j < len(captures); j++ {
		r.S[(captureGroupsIteration*(len(captures)-1))+j] = captures[j]
		if j < len(namedCaptureGroups) {
			if namedCaptureGroup := namedCaptureGroups[j]; namedCaptureGroup != "" && captures[j] != "" {
				r.Z[namedCaptureGroup] = captures[j]
			}
		}
	}
}

func regexParser(needle *string) *RE {
	var r *RE
	if UseRECache {
		r = regexpCache.Get(needle)
		if r != nil {
			return r
		}
	}

	r = &RE{
		_orig: *needle,
		mode:  'm',
	}

	sbM := strings.Builder{}
	sbM.Grow(len(r._orig))
	sbS := strings.Builder{}
	sbS.Grow(len(r._orig))
	sbF := strings.Builder{}
	sbF.Grow(3)
	var sb *strings.Builder

	i := 0
	switch (r._orig)[i] {
	case 'm':
		r.mode = 'm'
		i++
	case 't':
		r.mode = 't'
		if (r._orig)[i+1] == 'r' {
			i++
		}
		i++
	case 's':
		r.mode = 's'
		i++
	}

	r.separator = (r._orig)[i]
	i++

	var mode byte = 'm'
	sb = &sbM
	for ; i < len(r._orig); i++ {
		switch (r._orig)[i] {
		case r.separator:
			switch mode {
			case 'm':
				if r.mode == 'm' {
					mode = 'f'
					sb = &sbF
				} else {
					mode = 's'
					sb = &sbS
				}
			case 's':
				mode = 'f'
				sb = &sbF
			case 'f':
				fmt.Printf("Separator '%d' after flags-field in regex '%s'", r.separator, r._orig)
			}
			continue
		case '\\':
			sb.WriteByte((r._orig)[i])
			i++ // Skip the escaping backslash and the character being escaped
			sb.WriteByte((r._orig)[i])
			continue
		case '(':
			if (r._orig)[i+1] == '?' && (r._orig)[i+2] == ':' { // (?:) is a non-capturing group

			} else if (r._orig)[i+1] == '?' && (r._orig)[i+2] == 'P' && (r._orig)[i+3] == '<' { // (?P<) named capture groups
				r.nCaptures = true
				r.captures = true
			} else {
				r.captures = true
			}
		}
		sb.WriteByte((r._orig)[i])
	}
	if mode != 'f' && r._orig[len(r._orig)-1] != r.separator {
		fmt.Printf("Ending terminator '%c' not found for regexp string '%s'!\n", r.separator, r._orig)
	}
	rn, rs, rf := sbM.String(), sbS.String(), sbF.String()
	r.n, r.s, r.f = &rn, &rs, &rf

	flagHandler_x(r)
	flagHandlerGoNative(r)

	r.regex = regexp.MustCompile(*r.n)
	if UseRECache {
		regexpCache.Put(needle, r)
	}
	return r
}

func flagHandlerGoNative(r *RE) {
	sb := strings.Builder{}
	if strings.Contains(*r.f, "i") {
		sb.WriteByte('i')
	}
	if strings.Contains(*r.f, "m") {
		sb.WriteByte('m')
	}
	if strings.Contains(*r.f, "s") {
		sb.WriteByte('s')
	}
	if sb.Len() > 0 {
		*r.n = `(?` + sb.String() + `)` + *r.n
	}
}
func flagHandler_x(r *RE) {
	if strings.Contains(*r.f, "x") {
		flagHandler_x_(r, r.n)
		flagHandler_x_(r, r.s)
	}
}
func flagHandler_x_(r *RE, regstr *string) {
	sb := strings.Builder{}
	sb.Grow(len(*regstr))

	var inBracketedCharacterClass bool
	var inComment bool
	for i := 0; i < len(*regstr); i++ {
		if inComment {
			if (*regstr)[i] == '\n' {
				inComment = false
				continue
			}
			continue
		}

		switch (*regstr)[i] {
		case '[':
			if inBracketedCharacterClass {
				// This is an illegal regex, but keep removing whitespace to detect the real problem
			} else {
				inBracketedCharacterClass = true
			}
		case ']':
			if inBracketedCharacterClass {
				inBracketedCharacterClass = false
			} else {
				// This is an illegal regex, but keep removing whitespace to detect the real problem
			}
		case '\\':
			sb.WriteByte((*regstr)[i])
			i++ // Skip the escaping backslash and the character being escaped
			sb.WriteByte((*regstr)[i])
			continue
		case '#':
			if inBracketedCharacterClass {

			} else {
				inComment = true
				continue
			}
		case '\t', '\n', '\v', '\f', '\r', ' ':
			if !inBracketedCharacterClass {
				continue
			}
		}

		sb.WriteByte((*regstr)[i])
	}
	str := sb.String()
	*regstr = str
}
