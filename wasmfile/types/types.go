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

type ValType byte

const (
	ValI32  ValType = 0x7f
	ValI64  ValType = 0x7e
	ValF32  ValType = 0x7d
	ValF64  ValType = 0x7c
	ValNone ValType = 0x40
)

var ValTypeToByte map[string]ValType
var ByteToValType map[ValType]string

func init() {
	ValTypeToByte = make(map[string]ValType)
	ValTypeToByte["i32"] = ValI32
	ValTypeToByte["i64"] = ValI64
	ValTypeToByte["f32"] = ValF32
	ValTypeToByte["f64"] = ValF64
	ValTypeToByte["none"] = ValNone

	ByteToValType = make(map[ValType]string)
	ByteToValType[ValI32] = "i32"
	ByteToValType[ValI64] = "i64"
	ByteToValType[ValF32] = "f32"
	ByteToValType[ValF64] = "f64"
	ByteToValType[ValNone] = "none"
}

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
