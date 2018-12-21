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

import "log"

const (
	SYNC_LOCK_IMPL_KEY = "vesta.sync.lock.impl.key"
	ATOMIC_IMPL_KEY    = "vesta.atomic.impl.key"
)

type IdServiceImpl struct {
	AbstractIdServiceImpl
	idPopulator IdPopulator
	iDType      IdType
}

func (ids IdServiceImpl) init() {
	ids.initPopulator()
}

func (ids IdServiceImpl) initPopulator() {
	if (ids.idPopulator != nil){
		log.Println("The " , ids.idPopulator , " is used.")
	} else if (isPropKeyOn(SYNC_LOCK_IMPL_KEY)) {
		log.Println("The SyncIdPopulator is used.")
		ids.idPopulator = SyncIdPopulator{}
	} else if (isPropKeyOn(ATOMIC_IMPL_KEY)) {
		log.Println("The AtomicIdPopulator is used.")
		ids.idPopulator = AtomicIdPopulator{}
	} else {
		log.Println("The default LockIdPopulator is used.")
		ids.idPopulator = LockIdPopulator{}
	}
}

func  (ids IdServiceImpl) populateId(id Id) {
	ids.idPopulator.populateId(ids.timer, id, ids.idMeta)
}

func (ids IdServiceImpl) setIdPopulator(idPopulator IdPopulator) {
	ids.idPopulator = idPopulator
}
