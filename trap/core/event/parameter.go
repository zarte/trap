/*
 * Trap
 * An anti-pryer server for better privacy
 *
 * This file is a part of Trap project
 *
 * Copyright 2016 Rain Lee <raincious@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package event

import (
	"github.com/raincious/trap/trap/core/types"
)

const (
	DATA_TYPE_INT16 = iota
	DATA_TYPE_INT32
	DATA_TYPE_INT64

	DATA_TYPE_UINT16
	DATA_TYPE_UINT32
	DATA_TYPE_UINT64

	DATA_TYPE_STRING
	DATA_TYPE_BYTES
)

type Parameter struct {
	dataType int

	int16Data types.Int16
	int32Data types.Int32
	int64Data types.Int64

	uint16Data types.UInt16
	uint32Data types.UInt32
	uint64Data types.UInt64

	strData   types.String
	bytesData []byte
}

// Int 16s
func (p *Parameter) setInt16(i types.Int16) {
	p.int16Data = i

	p.dataType = DATA_TYPE_INT16
}

func (p *Parameter) GetInt16() types.Int16 {
	if p.dataType != DATA_TYPE_INT16 {
		return 0
	}

	return p.int16Data
}

// Int 32s
func (p *Parameter) setInt32(i types.Int32) {
	p.int32Data = i

	p.dataType = DATA_TYPE_INT32
}

func (p *Parameter) GetInt32() types.Int32 {
	if p.dataType != DATA_TYPE_INT32 {
		return 0
	}

	return p.int32Data
}

// Int 64s
func (p *Parameter) setInt64(i types.Int64) {
	p.int64Data = i

	p.dataType = DATA_TYPE_INT64
}

func (p *Parameter) GetInt64() types.Int64 {
	if p.dataType != DATA_TYPE_INT64 {
		return 0
	}

	return p.int64Data
}

// UInt 16s
func (p *Parameter) setUInt16(i types.UInt16) {
	p.uint16Data = i

	p.dataType = DATA_TYPE_UINT16
}

func (p *Parameter) GetUInt16() types.UInt16 {
	if p.dataType != DATA_TYPE_UINT16 {
		return 0
	}

	return p.uint16Data
}

// UInt 32s
func (p *Parameter) setUInt32(i types.UInt32) {
	p.uint32Data = i

	p.dataType = DATA_TYPE_UINT32
}

func (p *Parameter) GetUInt32() types.UInt32 {
	if p.dataType != DATA_TYPE_UINT32 {
		return 0
	}

	return p.uint32Data
}

// UInt 64s
func (p *Parameter) setUInt64(i types.UInt64) {
	p.uint64Data = i

	p.dataType = DATA_TYPE_UINT64
}

func (p *Parameter) GetUInt64() types.UInt64 {
	if p.dataType != DATA_TYPE_UINT64 {
		return 0
	}

	return p.uint64Data
}

// String
func (p *Parameter) setStr(s types.String) {
	p.strData = s

	p.dataType = DATA_TYPE_STRING
}

func (p *Parameter) GetStr() types.String {
	if p.dataType != DATA_TYPE_STRING {
		return ""
	}

	return p.strData
}

// Bytes
func (p *Parameter) setBytes(b []byte) {
	p.bytesData = b

	p.dataType = DATA_TYPE_BYTES
}

func (p *Parameter) GetBytes() []byte {
	if p.dataType != DATA_TYPE_BYTES {
		return []byte("")
	}

	return p.bytesData
}

// Exporter
func (p *Parameter) String() types.String {
	switch p.dataType {
	case DATA_TYPE_INT16:
		return p.int16Data.String()

	case DATA_TYPE_INT32:
		return p.int32Data.String()

	case DATA_TYPE_INT64:
		return p.int64Data.String()

	case DATA_TYPE_UINT16:
		return p.uint16Data.String()

	case DATA_TYPE_UINT32:
		return p.uint32Data.String()

	case DATA_TYPE_UINT64:
		return p.uint64Data.String()

	case DATA_TYPE_STRING:
		return p.strData

	case DATA_TYPE_BYTES:
		return types.String(string(p.bytesData[:]))
	}

	return types.String("")
}

type Parameters map[types.String]Parameter

func (p Parameters) AddInt16(key types.String, val types.Int16) Parameters {
	item := Parameter{}

	item.setInt16(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) AddInt32(key types.String, val types.Int32) Parameters {
	item := Parameter{}

	item.setInt32(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) AddInt64(key types.String, val types.Int64) Parameters {
	item := Parameter{}

	item.setInt64(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) AddUInt16(key types.String, val types.UInt16) Parameters {
	item := Parameter{}

	item.setUInt16(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) AddUInt32(key types.String, val types.UInt32) Parameters {
	item := Parameter{}

	item.setUInt32(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) AddUInt64(key types.String, val types.UInt64) Parameters {
	item := Parameter{}

	item.setUInt64(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) AddString(key types.String, val types.String) Parameters {
	item := Parameter{}

	item.setStr(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) AddBytes(key types.String, val []byte) Parameters {
	item := Parameter{}

	item.setBytes(val)

	p["$(("+key+"))"] = item

	return p
}

func (p Parameters) Parse(format types.String,
	labels []types.String) types.String {
	for _, label := range labels {
		if _, ok := p[label]; !ok {
			format = format.Replace(label, "")

			continue
		}

		val := p[label]

		format = format.Replace(label, val.String())
	}

	return format
}
