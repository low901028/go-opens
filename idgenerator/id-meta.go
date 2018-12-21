// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package idgenerator

// IdMeta: id meta data
type IdMeta struct {
	machineBits   byte
	seqBits       byte
	timeBits      byte
	genMethodBits byte
	typeBits      byte
	versionBits   byte
}

// machine bits
func (id *IdMeta) getMachineBits() byte {
	return id.machineBits
}
func (id *IdMeta) getMachineBitsMask() int64 {
	return int64(-1) ^ int64(-1)<<id.machineBits
}

// SeqBits
func (id *IdMeta) getSeqBitsStartPos() byte {
	return id.machineBits
}
func (id *IdMeta) getSeqBitsMask() int64 {
	return int64(-1) ^ int64(-1)<<id.seqBits
}

// TimeBits
func (id *IdMeta) getTimeBits() byte {
	return id.timeBits
}
func (id *IdMeta) getTimeBitsStartPos() int {
	return int(id.machineBits + id.seqBits)
}
func (id *IdMeta) getTimeBitsMask() int64 {
	return int64(-1) ^ int64(-1)<<id.timeBits
}

// GenMethodBits
func (id *IdMeta) getGenMethodBitsStartPos() int64 {
	return int64(id.machineBits + id.seqBits + id.timeBits)
}
func (id *IdMeta) getGenMethodBitsMask() int64 {
	return int64(-1) ^ int64(-1)<<id.genMethodBits
}

// TypeBits
func (id *IdMeta) getTypeBitsStartPos() int64 {
	return int64(id.machineBits + id.seqBits + id.timeBits + id.genMethodBits)
}
func (id *IdMeta) getTypeBitsMask() int64 {
	return int64(-1) ^ int64(-1)<<id.typeBits
}

// VersionBits
func (id *IdMeta) getVersionBitsStartPos() int64 {
	return int64(id.machineBits + id.seqBits + id.timeBits + id.genMethodBits + id.typeBits)
}
func (id *IdMeta) getVersionBitsMask() int64 {
	return int64(-1) ^ int64(-1)<<id.versionBits
}
