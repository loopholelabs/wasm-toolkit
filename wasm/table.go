package wasm

import (
	"fmt"
	"strconv"
	"strings"
)

type Table struct {
	Source string
	Size   int
	Limit  int
	Type   string
}

func NewTable(src string) *Table {
	// Skip the "(table", and remove the trailing ")"
	s := strings.TrimLeft(src[6:len(src)-1], " \t\r\n")

	xsize, s := ReadToken(s)
	size, err := strconv.Atoi(xsize)
	if err != nil {
		panic("Invalid size")
	}

	xlimit, s := ReadToken(s)
	limit, err := strconv.Atoi(xlimit)
	if err != nil {
		panic(fmt.Sprintf("Invalid limit \"%s\"", xlimit))
	}

	ty, s := ReadToken(s)

	//  (table (;0;) 3 3 funcref)

	return &Table{
		Source: src,
		Size:   size,
		Limit:  limit,
		Type:   ty,
	}
}

func (d *Table) Write() string {
	return fmt.Sprintf("(table %d %d %s)", d.Size, d.Limit, d.Type)
}
