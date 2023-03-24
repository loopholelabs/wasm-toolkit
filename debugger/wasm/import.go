package wasm

import (
	"fmt"
	"strconv"
	"strings"
)

type Import struct {
	Source      string
	Identifier1 string
	Identifier2 string
	Import      string
}

func NewImport(src string) *Import {
	// Skip the "(import", and remove the trailing ")"
	s := strings.TrimLeft(src[7:len(src)-1], " \t\r\n")

	id1, s := ReadString(s)
	id2, s := ReadString(s)

	idata, s := ReadElement(s)

	//  (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 0)))

	return &Import{
		Source:      src,
		Identifier1: id1,
		Identifier2: id2,
		Import:      idata,
	}
}

func (d *Import) GetFuncName() string {
	name := ""
	if strings.HasPrefix(d.Import, "(func ") {
		name = d.Import[6 : len(d.Import)-1] // Remove the '(func' and the ')'
		p := strings.Index(name, " ")
		if p != -1 {
			name = name[:p]
		}
	}
	return name
}

func (d *Import) Write() string {
	return fmt.Sprintf("(import %s %s %s)", d.Identifier1, d.Identifier2, d.Import)
}

func MergeImports(imp1 []*Import, imp2 []*Import, m2_type_offset int) []*Import {
	imports := make([]*Import, 0)

	// Note that imp2 types will have been moved, so need adjusting

	// Add all from mod1
	for _, i1 := range imp1 {
		imports = append(imports, i1)
	}

	for _, i1 := range imp2 {
		// Check if we already have it
		dupe := false
		for _, m := range imports {
			if i1.Identifier1 == m.Identifier1 && i1.Identifier2 == m.Identifier2 {
				dupe = true
				break
			}
		}
		if !dupe {
			// Adjust the type ID
			// eg (func $runtime.fd_write (type 0))
			if strings.HasPrefix(i1.Import, "(func ") {
				newI := i1.Import[6 : len(i1.Import)-1]
				id, newI := ReadToken(newI)

				// Now we just have "(type 0)"
				if strings.HasPrefix(newI, "(type ") {
					newI = newI[6 : len(newI)-1]
					v, err := strconv.Atoi(newI)
					if err != nil {
						panic("Error converting type")
					}
					i1.Import = fmt.Sprintf("(func %s (type %d))", id, v+m2_type_offset)
				}
			}
			imports = append(imports, i1)
		}
	}

	return imports
}
