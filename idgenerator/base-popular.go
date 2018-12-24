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

type BasePopular struct {
	sequence      int64
	lastTimestamp int64
}

func (bp BasePopular) populateId(timer Timer, id Id, idMeta IdMeta) {
	var timestamp int64 = timer.genTime()
	timer.validateTimestamp1(bp.lastTimestamp, timestamp)

	if timestamp == bp.lastTimestamp {
		bp.sequence++
		bp.sequence &= idMeta.getSeqBitsMask()
		if bp.sequence == 0 {
			timestamp = timer.tillNextTimeUnit(bp.lastTimestamp)
		}
	} else {
		bp.lastTimestamp = timestamp
		bp.sequence = 0
	}

	id.setSeq(bp.sequence)
	id.setTime(timestamp)
}

func (bp BasePopular) reset() {
	bp.sequence = 0
	bp.lastTimestamp = -1
}
