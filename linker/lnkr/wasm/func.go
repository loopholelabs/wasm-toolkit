package wasm

import (
	"fmt"
	"strconv"
	"strings"
)

type Func struct {
	Source             string
	OriginalIdentifier string
	Identifier         string
	Params             []string
	Locals             []string
	Type               string
	Result             string
	Instructions       []string
}

func NewFunc(src string) *Func {
	// Skip the "(func", and remove the trailing ")"
	s := strings.TrimLeft(src[5:len(src)-1], " \t\r\n")

	id, s := ReadToken(s)

	fun := &Func{
		Source:             src,
		Identifier:         id,
		OriginalIdentifier: id,
		Locals:             make([]string, 0),
		Params:             make([]string, 0),
	}

	for {
		if s == "" {
			break
		}
		if s[0] == '(' {
			var el string
			var eType string

			el, s = ReadElement(s)
			eType, _ = ReadToken(el[1:])
			// This can be param/local/type/result
			if eType == "type" {
				fun.Type = el
			} else if eType == "result" {
				fun.Result = el
			} else if eType == "local" {
				fun.Locals = append(fun.Locals, el)
			} else if eType == "param" {
				fun.Params = append(fun.Params, el)
			} else {
				panic("Unknown token in func")
			}
		} else {
			break
		}
	}

	s = s + "\n" // Add a newline at the end to make process simpler

	for {
		p := strings.Index(s, "\n")
		if p == -1 {
			break
		}

		fun.Instructions = append(fun.Instructions, s[0:p])

		s = s[p+1:]
	}

	return fun
}

func (d *Func) AdjustType(adj int) {
	if d.Type == "" {
		return
	}
	if strings.HasPrefix(d.Type, "(type ") {
		n, err := strconv.Atoi(d.Type[6 : len(d.Type)-1])
		if err != nil {
			panic("Error parsing type")
		}
		d.Type = fmt.Sprintf("(type %d)", n+adj)
	}
}

func (d *Func) PrefixGlobals(prefix string) {
	for index, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		if strings.HasPrefix(i, "global.get ") || strings.HasPrefix(i, "global.set ") {
			// Rewrite with a prefix
			d.Instructions[index] = i[:11] + "$" + prefix + i[12:]
		}
	}
}

func (d *Func) FixCallIndirectType(adj int) {
	for index, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		if strings.HasPrefix(i, "call_indirect (type ") {
			v, err := strconv.Atoi(i[20 : len(i)-1])
			if err != nil {
				panic("Error fixing up call_indirect")
			}
			d.Instructions[index] = fmt.Sprintf("call_indirect (type %d)", v+adj)
		}
	}
}

func (d *Func) PrefixCalls(prefix string, funcs []*Func) {
	for index, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		if strings.HasPrefix(i, "call ") {
			known := false
			target := i[5:]

			for _, f := range funcs {
				if target == f.OriginalIdentifier {
					known = true
				}
			}

			if known {
				// Rewrite with a prefix IF it's a known function
				d.Instructions[index] = i[:5] + "$" + prefix + i[6:]
			}
		}
	}
}

func (d *Func) LinkCalls(funcs []*Func) {
	for index, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		if strings.HasPrefix(i, "call ") {
			target := i[5:]

			for _, f := range funcs {
				if target == f.OriginalIdentifier {
					d.Instructions[index] = i[:5] + f.Identifier
					break
				}
			}
		}
	}
}

func (d *Func) FixMemoryInstr(inst string, fn string) {
	for index, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		if i == inst {
			d.Instructions[index] = "call " + fn
		}
	}
}

func (d *Func) FixMemoryInstrOffsetAlign(inst string, fn string) {
	ins := make([]string, 0)
	for _, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		// Either exactly the instruction, or it has space afterwards...
		if strings.HasPrefix(i, inst+" ") || (i == inst) {
			offset := 0
			align := 1
			var err error
			// Read the offset if there is one...
			rest := strings.TrimLeft(i[len(inst):], " \t\r\n")
			if strings.HasPrefix(rest, "offset=") {
				offsetString := rest[7:]
				p := strings.Index(offsetString, " ")
				if p != -1 {
					offsetString = offsetString[:p]
					rest = strings.TrimLeft(rest[7+p:], " \t\r\n")
				} else {
					rest = "" // No more!
				}
				offset, err = strconv.Atoi(offsetString)
				if err != nil {
					panic(fmt.Sprintf("Error parsing offset arg %s", i))
				}
			}

			if strings.HasPrefix(rest, "align=") {
				alignString := rest[6:]
				align, err = strconv.Atoi(alignString)
				if err != nil {
					panic(fmt.Sprintf("Error parsing offset arg %s", i))
				}
			}

			ins = append(ins, fmt.Sprintf("i32.const %d", offset))
			ins = append(ins, fmt.Sprintf("i32.const %d", align))
			ins = append(ins, "call "+fn)
		} else {
			ins = append(ins, i)
		}
	}

	d.Instructions = ins
}

func (d *Func) AdjustLoad(glob string) {
	ins := make([]string, 0)

	for _, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		if strings.HasPrefix(i, "i32.load") ||
			strings.HasPrefix(i, "i64.load") ||
			strings.HasPrefix(i, "f32.load") ||
			strings.HasPrefix(i, "f64.load") {
			ins = append(ins, "global.get "+glob)
			ins = append(ins, "i32.add")
			ins = append(ins, i)
		} else {
			ins = append(ins, i)
		}
	}

	d.Instructions = ins
}

// Rewrite any host calls so that memory is correct
func (d *Func) AdjustOutcalls(imps []*Import) {
	ins := make([]string, 0)

	for _, i := range d.Instructions {
		i = strings.Trim(i, " \t\r\n")
		if strings.HasPrefix(i, "call ") {
			known := false
			target := i[5:]

			// Check if it's one of the imports
			for _, imp := range imps {
				importName := imp.GetFuncName()
				if target == importName {
					known = true
				}
			}

			if known {
				// Rewrite with a swapin/swapout if it's a host call
				ins = append(ins, "call $m2_swapin")
				ins = append(ins, i)
				ins = append(ins, "call $m2_swapout")
			} else {
				ins = append(ins, i)
			}
		} else {
			ins = append(ins, i)
		}
	}
	d.Instructions = ins
}

func (d *Func) Write() string {
	f := fmt.Sprintf("(func %s", d.Identifier)

	if d.Type != "" {
		f = f + " " + d.Type
	}

	for _, p := range d.Params {
		f = f + " " + p
	}

	if d.Result != "" {
		f = f + " " + d.Result
	}

	f = f + "\n"

	for _, l := range d.Locals {
		f = f + l + "\n"
	}

	for _, i := range d.Instructions {
		f = f + i + "\n"
	}

	return f + ")"
}
