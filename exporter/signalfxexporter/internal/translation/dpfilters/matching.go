// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dpfilters // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters"

import (
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

// Contains all of the logic for glob and regex based filtering.

func isGlobbed(s string) bool {
	return strings.ContainsAny(s, "*?[]{}!")
}

func isRegex(s string) bool {
	return len(s) > 2 && s[0] == '/' && s[len(s)-1] == '/'
}

// remove the bracketing slashes for a regex.
func stripSlashes(s string) string {
	return s[1 : len(s)-1]
}

// stripNegation checks if a string is prefixed with "!"
// and will returned the stripped string and true if so
// else, return original value and false.
func stripNegation(value string) (string, bool) {
	if strings.HasPrefix(value, "!") {
		return value[1:], true
	}
	return value, false
}

type matcher interface {
	// Returns whether the string matched and whether it was a negated match.
	Matches(s string) (bool, bool)
}

type regexMatcher struct {
	re      *regexp.Regexp
	negated bool
}

var _ matcher = (*regexMatcher)(nil)

func (m *regexMatcher) Matches(s string) (bool, bool) {
	return m.re.MatchString(s), m.negated
}

type globMatcher struct {
	glob    glob.Glob
	negated bool
}

var _ matcher = &globMatcher{}

func (m *globMatcher) Matches(s string) (bool, bool) {
	return m.glob.Match(s), m.negated
}
