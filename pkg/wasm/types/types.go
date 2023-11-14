/*
	Copyright 2023 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package types

// Value types
type ValType byte

const (
	ValNone ValType = 0x40
	ValF64  ValType = 0x7c
	ValF32  ValType = 0x7d
	ValI64  ValType = 0x7e
	ValI32  ValType = 0x7f
)

var ValTypeToByte map[string]ValType = map[string]ValType{
	"i32":  ValI32,
	"i64":  ValI64,
	"f32":  ValF32,
	"f64":  ValF64,
	"none": ValNone,
}

var ByteToValType map[ValType]string = map[ValType]string{
	ValI32:  "i32",
	ValI64:  "i64",
	ValF32:  "f32",
	ValF64:  "f64",
	ValNone: "none",
}

// Section IDs
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

// Limit types
type LimitType byte

const (
	LimitTypeMin    LimitType = 0x00
	LimitTypeMinMax LimitType = 0x01
)

// Export type
type ExportType byte

const (
	ExportFunc   ExportType = 0
	ExportTable  ExportType = 1
	ExportMem    ExportType = 2
	ExportGlobal ExportType = 3
)

// Other
const FuncTypePrefix byte = 0x60

const TableTypeFuncref byte = 0x70
