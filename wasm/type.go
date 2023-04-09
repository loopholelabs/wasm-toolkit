package wasm

import (
	"fmt"
	"strings"
)

type Type struct {
	Source string
	Type   string
}

func NewType(src string) *Type {
	// Skip the "(type", and remove the trailing ")"
	s := strings.TrimLeft(src[5:len(src)-1], " \t\r\n")

	tdata, s := ReadElement(s)

	//   (type (;0;) (func (param i32 i32 i32 i32) (result i32)))

	return &Type{
		Source: src,
		Type:   tdata,
	}
}

func (d *Type) Write() string {
	return fmt.Sprintf("(type %s)", d.Type)
}

func MergeTypes(prefix1 string, type1 []*Type, prefix2 string, type2 []*Type) []*Type {
	types := make([]*Type, 0)

	for _, i1 := range type1 {
		types = append(types, i1)
	}

	for _, i1 := range type2 {
		types = append(types, i1)
	}

	return types
}
