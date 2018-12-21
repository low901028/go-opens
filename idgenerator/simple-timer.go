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
	"log"
	"time"
)

const EPOCH int64 = 1514736000000

type SimpleTimer struct {
	idMeta  IdMeta
	idType  IdType
	maxTime int64
	epoch   int64
}

func (st SimpleTimer) init(idMeta IdMeta, idType IdType) {
	st.idMeta = idMeta
	st.maxTime = (1 << idMeta.getTimeBits()) - 1
	st.idType = idType
	st.epoch = EPOCH
	st.genTime()
	st.timerUsedLog()
}

func (st SimpleTimer) timerUsedLog() {
	var expirationDate = st.transTime(st.maxTime)
	days := ((expirationDate.UnixNano() - time.Now().UnixNano()) / (1000 * 60 * 60 * 24))
	log.Printf("The current time bit length is %v, the expiration date is %s, this can be used for %d days.",
		st.idMeta.getTimeBits(), expirationDate.Format(time.UnixDate), days)
}

func (st SimpleTimer) transTime(t int64) time.Time {
	if st.idType == MILLISECONDS {
		return time.Unix(0, t+st.epoch)
	} else {
		return time.Unix(0, (t*1000 + st.epoch))
	}
}

func (st SimpleTimer) validateTimestamp1(lastTimestamp int64, timestamp int64) {
	if timestamp < lastTimestamp {
		log.Printf("Clock moved backwards.  Refusing to generate id for %d second/milisecond.",
			lastTimestamp-timestamp)
	}
	return
}

func (st SimpleTimer) validateTimestamp2(timestamp int64) {
	if timestamp < st.maxTime {
		log.Printf("The current timestamp (%s >= %s) has overflowed, Vesta Service will be terminate.", timestamp, st.maxTime)
	}
	return
}

func (st SimpleTimer) tillNextTimeUnit(lastTimestamp int64) int64 {
	log.Printf("Ids are used out during %d. Waiting till next second/milisencond.",
		lastTimestamp)

	var timestamp int64 = st.genTime()
	for timestamp <= lastTimestamp {
		timestamp = st.genTime()
	}

	log.Printf("Next second/milisencond %d is up.", timestamp)

	return timestamp
}

func (st SimpleTimer) genTime() int64 {
	var t int64
	if st.idType == MILLISECONDS {
		t = (time.Now().UnixNano() - st.epoch)
	} else {
		t = (time.Now().UnixNano() - st.epoch) / 1000
	}
	st.validateTimestamp2(t)
	return t
}
