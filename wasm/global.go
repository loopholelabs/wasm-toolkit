package wasm

import (
	"fmt"
	"strings"
)

type Global struct {
	Source     string
	Identifier string
	Type       string
	Value      string
}

func NewGlobal(src string) *Global {
	// Skip the "(global", and remove the trailing ")"
	s := strings.TrimLeft(src[7:len(src)-1], " \t\r\n")
	s = SkipComment(s)

	var id string = ""
	if s[0] == '$' { // Optional ID
		id, s = ReadToken(s)
	}

	s = SkipComment(s)

	// Type can either be eg "i32" but can also be eg "(i32 mut)"
	var tdata string
	if s[0] == '(' {
		tdata, s = ReadElement(s)
	} else {
		tdata, s = ReadToken(s)
	}

	vdata, s := ReadElement(s)
	//  (global $__stack_pointer (mut i32) (i32.const 65536))

	return &Global{
		Source:     src,
		Identifier: id,
		Type:       tdata,
		Value:      vdata,
	}
}

func (d *Global) Write() string {
	if d.Identifier == "" {
		return fmt.Sprintf("(global %s %s)", d.Type, d.Value)
	}
	return fmt.Sprintf("(global %s %s %s)", d.Identifier, d.Type, d.Value)
}

func MergeGlobals(prefix1 string, glob1 []*Global, prefix2 string, glob2 []*Global) []*Global {
	globals := make([]*Global, 0)

	for _, i1 := range glob1 {
		if i1.Identifier != "" {
			i1.Identifier = string('$') + prefix1 + i1.Identifier[1:]
		}
		globals = append(globals, i1)
	}

	for _, i1 := range glob2 {
		if i1.Identifier != "" {
			i1.Identifier = string('$') + prefix2 + i1.Identifier[1:]
		}
		globals = append(globals, i1)
	}

	return globals
}
