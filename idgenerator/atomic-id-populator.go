// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
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

import (
	"sync/atomic"
	"unsafe"
)

type Variant struct {
	sequence      int64
	lastTimestamp int64
}

type AtomicIdPopulator struct {
	Variant
}

//
type AtomicReference struct {
	variant unsafe.Pointer
}

var aRef = AtomicReference{}

var variant = Variant{
	sequence:      0,
	lastTimestamp: -1,
}

func (idPopulator AtomicIdPopulator) populateId(timer Timer, id Id, idMeta IdMeta) {
	var varOld, varNew Variant
	var timestamp, sequence int64

	for {
		// Save the old variant
		va := (*Variant)(atomic.LoadPointer(&aRef.variant))
		varOld = *va

		// populate the current variant
		timestamp = timer.genTime()
		timer.validateTimestamp1(varOld.lastTimestamp, timestamp)

		sequence = varOld.sequence

		if timestamp == varOld.lastTimestamp {
			sequence++
			sequence &= idMeta.getSeqBitsMask()
			if sequence == 0 {
				timestamp = timer.tillNextTimeUnit(varOld.lastTimestamp)
			}
		} else {
			sequence = 0
		}

		// Assign the current variant by the atomic tools
		varNew = Variant{}
		varNew.sequence = sequence
		varNew.lastTimestamp = timestamp

		if atomic.CompareAndSwapPointer(&aRef.variant, &varOld, &varNew) {
			id.setSeq(sequence)
			id.setTime(timestamp)
			break
		}
	}
}

//
func (idPopulator AtomicIdPopulator) reset() {
	aRef = AtomicReference{}
}
