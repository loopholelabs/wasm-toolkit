package wasm

import (
	"fmt"
	"strconv"
	"strings"
)

type Elem struct {
	Source string
	Offset string
	Type   string
	Elems  []string
}

func NewElem(src string) *Elem {
	// Skip the "(elem", and remove the trailing ")"
	s := strings.TrimLeft(src[5:len(src)-1], " \t\r\n")

	offset, s := ReadElement(s)
	etype, s := ReadToken(s)

	// Now read the functions
	e := make([]string, 0)

	var t string
	for {
		t, s = ReadToken(s)
		if len(t) == 0 {
			break
		}
		e = append(e, t)
	}

	// (elem (;0;) (i32.const 1) func $runtime.memequal $runtime.hash32)

	return &Elem{
		Source: src,
		Offset: offset,
		Type:   etype,
		Elems:  e,
	}
}

func (d *Elem) AdjustOffset(adj int) {
	if strings.HasPrefix(d.Offset, "(i32.const") {
		loc, err := strconv.Atoi(d.Offset[11 : len(d.Offset)-1])
		if err != nil {
			panic("Wasm.elem: Error parsing location")
		}
		d.Offset = fmt.Sprintf("(i32.const %d)", loc+adj)
	} else {
		panic("Wasm.elem: Can only have i32.const location")
	}
}

func (d *Elem) Write() string {
	v := fmt.Sprintf("(elem %s %s", d.Offset, d.Type)

	for _, e := range d.Elems {
		v = v + " " + e
	}

	v = v + ")"
	return v
}
