package wasm

import (
	"fmt"
	"strconv"
	"strings"
)

type Memory struct {
	Source string
	Size   int
}

func NewMemory(src string) *Memory {
	// Skip the "(memory", and remove the trailing ")"
	s := strings.TrimLeft(src[7:len(src)-1], " \t\r\n")

	len, s := ReadToken(s)

	size, err := strconv.Atoi(len)
	if err != nil {
		panic("Invalid memory size")
	}

	// (memory (;0;) 2)
	return &Memory{
		Source: src,
		Size:   size,
	}
}

func (d *Memory) Write() string {
	return fmt.Sprintf("(memory %d)", d.Size)
}
