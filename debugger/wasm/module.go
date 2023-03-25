package wasm

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type WasmModule struct {
	filename string

	Types   []*Type
	Imports []*Import
	Funcs   []*Func
	Tables  []*Table
	Memorys []*Memory
	Globals []*Global
	Exports []*Export
	Elems   []*Elem
	Datas   []*Data
}

func NewModule(file string) *WasmModule {
	return &WasmModule{
		filename: file,

		Types:   make([]*Type, 0),
		Imports: make([]*Import, 0),
		Funcs:   make([]*Func, 0),
		Tables:  make([]*Table, 0),
		Memorys: make([]*Memory, 0),
		Globals: make([]*Global, 0),
		Exports: make([]*Export, 0),
		Elems:   make([]*Elem, 0),
		Datas:   make([]*Data, 0),
	}
}

func (wm *WasmModule) Parse() {
	data, err := ioutil.ReadFile(wm.filename)
	if err != nil {
		log.Fatal(err)
	}

	text := string(data)

	// Read the module
	moduleText, _ := ReadElement(text)

	moduleType, _ := ReadToken(moduleText[1:])

	//fmt.Printf("Module %s is %d bytes\n", wm.filename, len(moduleText))

	// Now read all the individual elements from within the module...

	text = text[len(moduleType)+1:]

	for {
		text = strings.TrimLeft(text, " \t\r\n") // Skip to next bit
		// End of the module?
		if text[0] == ')' {
			break
		}

		e, _ := ReadElement(text)
		eType, _ := ReadToken(e[1:])

		if eType == "data" {
			d := NewData(e)
			wm.Datas = append(wm.Datas, d)
		} else if eType == "elem" {
			el := NewElem(e)
			wm.Elems = append(wm.Elems, el)
		} else if eType == "export" {
			ex := NewExport(e)
			wm.Exports = append(wm.Exports, ex)
		} else if eType == "func" {
			f := NewFunc(e)
			wm.Funcs = append(wm.Funcs, f)
		} else if eType == "global" {
			g := NewGlobal(e)
			wm.Globals = append(wm.Globals, g)
		} else if eType == "import" {
			i := NewImport(e)
			wm.Imports = append(wm.Imports, i)
		} else if eType == "memory" {
			mem := NewMemory(e)
			wm.Memorys = append(wm.Memorys, mem)
		} else if eType == "table" {
			t := NewTable(e)
			wm.Tables = append(wm.Tables, t)
		} else if eType == "type" {
			t := NewType(e)
			wm.Types = append(wm.Types, t)
		} else {
			panic(fmt.Sprintf("Unknown element \"%s\"", eType))
		}

		// Skip over this element
		text = text[len(e):]
	}
	/*
	   fmt.Printf("Parsed wat file. %d Data, %d Elem, %d Export, %d Func, %d Global, %d Import, %d Memory, %d Table, %d type\n",

	   	len(wm.Datas),
	   	len(wm.Elems),
	   	len(wm.Exports),
	   	len(wm.Funcs),
	   	len(wm.Globals),
	   	len(wm.Imports),
	   	len(wm.Memorys),
	   	len(wm.Tables),
	   	len(wm.Types),

	   )
	*/
}

func (m *WasmModule) Write() string {
	d := "(module\n"

	// Write out the various things...

	for _, t := range m.Types {
		d = d + t.Write() + "\n"
	}

	for _, i := range m.Imports {
		d = d + i.Write() + "\n"
	}

	for _, f := range m.Funcs {
		d = d + f.Write() + "\n"
	}

	for _, t := range m.Tables {
		d = d + t.Write() + "\n"
	}

	for _, m := range m.Memorys {
		d = d + m.Write() + "\n"
	}

	for _, g := range m.Globals {
		d = d + g.Write() + "\n"
	}

	for _, e := range m.Exports {
		d = d + e.Write() + "\n"
	}

	for _, e := range m.Elems {
		d = d + e.Write() + "\n"
	}

	for _, da := range m.Datas {
		d = d + da.Write() + "\n"
	}

	return d + ")"
}
