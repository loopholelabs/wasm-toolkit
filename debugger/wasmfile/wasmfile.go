package wasmfile

import (
	"debug/dwarf"
	"io/ioutil"
)

const WasmHeader uint32 = 0x6d736100
const WasmVersion uint32 = 0x00000001

type ValType byte

const (
	ValI32  ValType = 0x7f
	ValI64  ValType = 0x7e
	ValF32  ValType = 0x7d
	ValF64  ValType = 0x7c
	ValNone ValType = 0x40
)

var valTypeToByte map[string]ValType
var byteToValType map[ValType]string

func init() {
	valTypeToByte = make(map[string]ValType)
	valTypeToByte["i32"] = ValI32
	valTypeToByte["i64"] = ValI64
	valTypeToByte["f32"] = ValF32
	valTypeToByte["f64"] = ValF64

	byteToValType = make(map[ValType]string)
	byteToValType[ValI32] = "i32"
	byteToValType[ValI64] = "i64"
	byteToValType[ValF32] = "f32"
	byteToValType[ValF64] = "f64"
}

const (
	LimitTypeMin    byte = 0x00
	LimitTypeMinMax byte = 0x01
)

type ExportType byte

const (
	ExportFunc   ExportType = 0
	ExportTable  ExportType = 1
	ExportMem    ExportType = 2
	ExportGlobal ExportType = 3
)

const FuncTypePrefix byte = 0x60

const TableTypeFuncref byte = 0x70

type SectionId byte

const (
	SectionCustom    SectionId = 0
	SectionType      SectionId = 1
	SectionImport    SectionId = 2
	SectionFunction  SectionId = 3
	SectionTable     SectionId = 4
	SectionMemory    SectionId = 5
	SectionGlobal    SectionId = 6
	SectionExport    SectionId = 7
	SectionStart     SectionId = 8
	SectionElem      SectionId = 9
	SectionCode      SectionId = 10
	SectionData      SectionId = 11
	SectionDataCount SectionId = 12
)

type FunctionEntry struct {
	TypeIndex int
}

type TypeEntry struct {
	Param  []ValType
	Result []ValType
}

type CustomEntry struct {
	Name string
	Data []byte
}

type ExportEntry struct {
	Name  string
	Type  ExportType
	Index int
}

type ImportEntry struct {
	Module string
	Name   string
	Type   ExportType
	Index  int
}

type TableEntry struct {
	TableType byte
	LimitMin  int
	LimitMax  int
}

type GlobalEntry struct {
	Type       ValType
	Mut        byte
	Expression []*Expression
}

type MemoryEntry struct {
	LimitMin int
	LimitMax int
}

type CodeEntry struct {
	Locals         []ValType
	CodeSectionPtr uint64
	ExprData       []byte
	Expression     []*Expression
}

type DataEntry struct {
	MemIndex int
	Offset   []*Expression
	Data     []byte
}

type ElemEntry struct {
	TableIndex int
	Offset     []*Expression
	Indexes    []uint64
}

type LineInfo struct {
	Filename   string
	Linenumber int
}

type WasmFile struct {
	Function []*FunctionEntry
	Type     []*TypeEntry
	Custom   []*CustomEntry
	Export   []*ExportEntry
	Import   []*ImportEntry
	Table    []*TableEntry
	Global   []*GlobalEntry
	Memory   []*MemoryEntry
	Code     []*CodeEntry
	Data     []*DataEntry
	Elem     []*ElemEntry

	dwarfData   *dwarf.Data
	lineNumbers map[uint64]LineInfo

	functionNames map[int]string
}

// Create a new WasmFile from a file
func New(filename string) (*WasmFile, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	wf := &WasmFile{}
	err = wf.DecodeBinary(data)
	return wf, err
}

func (wf *WasmFile) GetCustomSectionData(name string) []byte {
	for _, c := range wf.Custom {
		if c.Name == name {
			return c.Data
		}
	}
	return nil
}
