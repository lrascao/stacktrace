// Copyright 2016 Palantir Technologies
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

package stacktrace_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"runtime"
	"strings"
	"path/filepath"

	"github.com/stretchr/testify/assert"

	"github.com/lrascao/stacktrace"
	"github.com/lrascao/stacktrace/cleanpath"
)

func TestFormat(t *testing.T) {
	plainErr := errors.New("plain")

	_, filename, _, _ := runtime.Caller(0)
	cleanpath.GoPath = filepath.Dir(filename)
	stacktraceErr := stacktrace.Propagate(plainErr, "decorated")

	for _, test := range []struct {
		format             stacktrace.Format
		specifier          string
		expectedPlain      string
		expectedStacktrace string
	}{
		{
			format:             stacktrace.FormatFull,
			specifier:          "%v",
			expectedPlain:      "plain",
			expectedStacktrace: "decorated\n --- at %s/format_test.go:## (TestFormat) ---\nCaused by: plain",
		},
		{
			format:             stacktrace.FormatFull,
			specifier:          "%q",
			expectedPlain:      "\"plain\"",
			expectedStacktrace: "\"decorated\\n --- at %s/format_test.go:## (TestFormat) ---\\nCaused by: plain\"",
		},
		{
			format:             stacktrace.FormatFull,
			specifier:          "%105s",
			expectedPlain:      "                                                                                                    plain",
			expectedStacktrace: "decorated\n --- at %s/format_test.go:## (TestFormat) ---\nCaused by: plain",
		},
		{
			format:             stacktrace.FormatFull,
			specifier:          "%#s",
			expectedPlain:      "plain",
			expectedStacktrace: "decorated: plain",
		},
		{
			format:             stacktrace.FormatBrief,
			specifier:          "%v",
			expectedPlain:      "plain",
			expectedStacktrace: "decorated: plain",
		},
		{
			format:             stacktrace.FormatBrief,
			specifier:          "%q",
			expectedPlain:      "\"plain\"",
			expectedStacktrace: "\"decorated: plain\"",
		},
		{
			format:             stacktrace.FormatBrief,
			specifier:          "%20s",
			expectedPlain:      "               plain",
			expectedStacktrace: "    decorated: plain",
		},
		{
			format:             stacktrace.FormatBrief,
			specifier:          "%+s",
			expectedPlain:      "plain",
			expectedStacktrace: "decorated\n --- at %s/format_test.go:## (TestFormat) ---\nCaused by: plain",
		},
	} {
		stacktrace.DefaultFormat = test.format

		actualPlain := fmt.Sprintf(test.specifier, plainErr)
		assert.Equal(t, test.expectedPlain, actualPlain)

		actualStacktrace := fmt.Sprintf(test.specifier, stacktraceErr)

		expectedStacktrace := test.expectedStacktrace
		if strings.Count(test.expectedStacktrace, "%s") > 0 {
			expectedStacktrace = fmt.Sprintf(test.expectedStacktrace, cleanpath.GoPath)
		}

		digits := regexp.MustCompile(`\d`)
		actualStacktrace = digits.ReplaceAllString(actualStacktrace, "#")
		assert.Equal(t, expectedStacktrace, actualStacktrace)
	}
}
