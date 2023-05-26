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

package encoding

import (
	"encoding/binary"
	"io"
)

func DecodeSleb128(b []byte) (s int64, n int) {
	result := int64(0)
	shift := 0
	ptr := 0
	for {
		by := b[ptr]
		ptr++
		result = result | (int64(by&0x7f) << shift)
		shift += 7
		if (by & 0x80) == 0 {
			if shift < 64 && (by&0x40) != 0 {
				return result | (^0 << shift), ptr
			}
			return result, ptr
		}
	}
}

func AppendSleb128(buf []byte, val int64) []byte {
	for {
		b := val & 0x7f
		val = val >> 7
		if (val == 0 && b&0x40 == 0) ||
			(val == -1 && b&0x40 != 0) {
			buf = append(buf, byte(b))
			return buf
		}
		buf = append(buf, byte(b|0x80))
	}
}

func WriteString(w io.Writer, s string) error {
	data := []byte(s)
	err := WriteUvarint(w, uint64(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func WriteUvarint(w io.Writer, v uint64) error {
	b := binary.AppendUvarint(make([]byte, 0), v)
	_, err := w.Write(b)
	return err
}

func WriteVarint(w io.Writer, v int64) error {
	b := AppendSleb128(make([]byte, 0), v)
	_, err := w.Write(b)
	return err
}
