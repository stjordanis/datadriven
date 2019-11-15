// Copyright 2019 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package datadriven

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/errors"
)

func TestNewLineBetweenDirectives(t *testing.T) {
	RunTestFromString(t, `
# Some testing of sensitivity to newlines
foo
----
unknown command

bar
----
unknown command




bar
----
unknown command
`, func(d *TestData) string {
		if d.Input != "sentence" {
			return "unknown command"
		}
		return ""
	})
}

func TestParseLine(t *testing.T) {
	RunTestFromString(t, `
parse
xx +++
----
here: cannot parse directive at column 4: xx +++

parse
xx a=b a=c
----
here: duplicate key in argument list: a=c

parse
xx a=b b=c c=(1,2,3)
----
"xx" [a=b b=c c=(1, 2, 3)]
`, func(d *TestData) string {
		cmd, args, err := ParseLine(d.Input)
		if err != nil {
			return errors.Wrap(err, "here").Error()
		}
		return fmt.Sprintf("%q %+v", cmd, args)
	})
}

func TestArgFormat(t *testing.T) {
	RunTestFromString(t, `
# NB: we allow duplicate args.
# ScanArgs simply picks the first occurrence.
make argTuple=(1, 🍌) argInt=12 argString=greedily,impatient moreIgnore= a,b,c
sentence
----
Did the following: make sentence
1 hungry monkey eats a 🍌
while 12 other monkeys watch greedily,impatient
true I'd say
`, func(d *TestData) string {
		var one int
		var twelve int
		var banana string
		var greedily string
		d.ScanArgs(t, "argTuple", &one, &banana)
		d.ScanArgs(t, "argInt", &twelve)
		d.ScanArgs(t, "argString", &greedily)
		abc := d.HasArg("a,b,c")
		return fmt.Sprintf(`Did the following: %s %s
%d hungry monkey eats a %s
while %d other monkeys watch %s
%v I'd say`,
			d.Cmd, d.Input, one, banana, twelve, greedily, abc,
		)
	})
}
