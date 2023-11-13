/*
	Copyright 2023 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Input     string
	Output    string
	Remainder string
}

func TestReadElement(t *testing.T) {

	tests := []TestCase{
		{
			Input:     "(hello) world",
			Output:    "(hello)",
			Remainder: "world",
		},
		{
			Input:     "(hello (world)) (ok)",
			Output:    "(hello (world))",
			Remainder: "(ok)",
		},
		{
			Input:     "(hello (world \"something ()\")) (ok)",
			Output:    "(hello (world \"something ()\"))",
			Remainder: "(ok)",
		},
		{
			Input:     "no brackets",
			Output:    "",
			Remainder: "no brackets",
		},
		{
			Input:     "(; Some comment\n(here) ;) (something) there",
			Output:    "(something)",
			Remainder: "there",
		},
	}

	for _, testcase := range tests {

		o, r := ReadElement(testcase.Input)
		assert.Equal(t, testcase.Output, o)
		assert.Equal(t, testcase.Remainder, r)
	}
}

func TestReadString(t *testing.T) {

	tests := []TestCase{
		{
			Input:     "\"hello\" world",
			Output:    "\"hello\"",
			Remainder: "world",
		},
		{
			Input:     "no brackets",
			Output:    "",
			Remainder: "no brackets",
		},
		{
			Input:     "(; Some comment\n(here) ;) \"something\" there",
			Output:    "\"something\"",
			Remainder: "there",
		},
	}

	for _, testcase := range tests {

		o, r := ReadString(testcase.Input)
		assert.Equal(t, testcase.Output, o)
		assert.Equal(t, testcase.Remainder, r)
	}
}

func TestReadToken(t *testing.T) {

	tests := []TestCase{
		{
			Input:     "hello world",
			Output:    "hello",
			Remainder: "world",
		},
		{
			Input:     "123.45 78",
			Output:    "123.45",
			Remainder: "78",
		},
		{
			Input:     "123.45\t78",
			Output:    "123.45",
			Remainder: "78",
		},
		{
			Input:     "123.45\n78",
			Output:    "123.45",
			Remainder: "78",
		},
		{
			Input:     "(; Some comment\n(here) ;) something there",
			Output:    "something",
			Remainder: "there",
		},
	}

	for _, testcase := range tests {

		o, r := ReadToken(testcase.Input)
		assert.Equal(t, testcase.Output, o)
		assert.Equal(t, testcase.Remainder, r)
	}
}
