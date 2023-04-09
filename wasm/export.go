package wasm

import (
	"fmt"
	"strings"
)

type Export struct {
	Source     string
	Identifier string
	Export     string
}

func NewExport(src string) *Export {
	// Skip the "(export", and remove the trailing ")"
	s := strings.TrimLeft(src[7:len(src)-1], " \t\r\n")

	id, s := ReadString(s)
	edata, s := ReadElement(s)

	//  (export "memory" (memory 0))

	return &Export{
		Source:     src,
		Identifier: id,
		Export:     edata,
	}
}

func (d *Export) Write() string {
	return fmt.Sprintf("(export %s %s)", d.Identifier, d.Export)
}

func MergeExports(prefix1 string, exp1 []*Export, prefix2 string, exp2 []*Export) []*Export {
	exports := make([]*Export, 0)

	// Export the combined memory
	exports = append(exports, NewExport("(export \"memory\" (memory 0))"))

	for _, i1 := range exp1 {
		if strings.HasPrefix(i1.Export, "(func") {
			i1.Identifier = string('"') + prefix1 + i1.Identifier[1:]
			i1.Export = "(func $" + prefix1 + i1.Export[7:len(i1.Export)-1] + ")"
			exports = append(exports, i1)
		} else {
			// IGNORE OTHER EXPORTS FOR NOW
		}
	}

	for _, i1 := range exp2 {
		if strings.HasPrefix(i1.Export, "(func") {
			i1.Identifier = string('"') + prefix2 + i1.Identifier[1:]
			i1.Export = "(func $" + prefix2 + i1.Export[7:len(i1.Export)-1] + ")"
			exports = append(exports, i1)
		} else {
			// IGNORE OTHER EXPORTS FOR NOW
		}
	}

	return exports
}
