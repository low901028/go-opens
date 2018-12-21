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

import (
	"log"
	"time"
)

type AbstractIdServiceImpl struct {
	AbstractPopulateId
	machineId         int64
	genMethod         int64
	version           int64
	idType            IdType
	idMeta            IdMeta
	idConverter       IdConverter
	machineIdProvider MachineIdProvider
	timer             Timer
}

type AbstractPopulateId interface {
	populateId(id Id)
}

// constructor function
func NewAbstractIdServiceImpl0() *AbstractIdServiceImpl{
	var as = new(AbstractIdServiceImpl)
	as.idType = SECONDS
	return as
}
func NewAbstractIdServiceImpl1(ctype string) *AbstractIdServiceImpl{
	var as = new(AbstractIdServiceImpl)
	t, _ := Parse(ctype)
	as.idType = t
	return as
}
func NewAbstractIdServiceImpl2(ctype int64) *AbstractIdServiceImpl{
	var as = new(AbstractIdServiceImpl)
	t, _ := Parse1(int(ctype))
	as.idType = t
	return as
}
func NewAbstractIdServiceImpl3(ctype IdType) *AbstractIdServiceImpl{
	var as = new(AbstractIdServiceImpl)
	as.idType = ctype
	return as
}

func (as AbstractIdServiceImpl) init() {
	if idMeta, ok := GetIdMeta(as.idType); ok{
		as.setIdMeta(idMeta)
	}

	if as.idConverter == nil {
		as.setIdConverter(IdConverterImpl{})
	}

	if as.timer == nil {
		as.setTimer(SimpleTimer{})
	}
	as.timer.init(as.idMeta, as.idType)

	as.machineId = as.machineIdProvider.getMachineId()
	as.validateMachineId(as.machineId)
}

func (as AbstractIdServiceImpl) genId() int64 {
	var id = Id{}

	id.setMachine(as.machineId)
	id.setGenMethod(as.genMethod)
	id.setType(int64(as.idType.num))
	id.setVersion(as.version)

	as.populateId(id)
	var ret = as.idConverter.convert1(id, as.idMeta)
	log.Printf("Id: %s => %d", id, ret)

	return ret
}

func (as AbstractIdServiceImpl) validateMachineId(machineId int64) {
	if machineId < 0 {
		log.Fatalf("The machine ID is not configured properly %d < 0) so that Vesta Service refuses to start.", machineId)
	} else if machineId >= (1 << as.idMeta.getMachineBits()) {
		log.Fatalf("The machine ID is not configured properly %d >= %d) so that Vesta Service refuses to start.", machineId, (1 << as.idMeta.getMachineBits()))
	}
}

func (as AbstractIdServiceImpl) transTime(time int64) time.Time {
	return as.timer.transTime(time)
}

func (as AbstractIdServiceImpl) expId(id int64) Id {
	return as.idConverter.convert2(id, as.idMeta)
}

func (as AbstractIdServiceImpl) makeId0(time int64, seq int64) int64 {
	return as.makeId1(time, seq, as.machineId)
}
func (as AbstractIdServiceImpl) makeId1(time int64, seq int64, machine int64) int64 {
	return as.makeId2(as.genMethod, time, seq, machine)
}
func (as AbstractIdServiceImpl) makeId2(genMethod int64, time int64, seq int64, machine int64) int64 {
	return as.makeId3(as.idType.value(), genMethod, time, seq, machine);
}
func (as AbstractIdServiceImpl) makeId3(ctype int64, genMethod int64, time int64, seq int64, machine int64) int64 {
	return as.makeId4(as.version, ctype, genMethod, time, seq, machine)
}
func (as AbstractIdServiceImpl) makeId4(version int64, ctype int64, genMethod int64, time int64, seq int64, machine int64) int64 {
	var id = Id{machine, seq, time, genMethod, ctype, version}
	return as.idConverter.convert1(id, as.idMeta)
}

func (as AbstractIdServiceImpl) setMachineId(machineId int64) {
	as.machineId = machineId
}

func (as AbstractIdServiceImpl) setGenMethod(genMethod int64) {
	as.genMethod = genMethod
}

func (as AbstractIdServiceImpl) setVersion(version int64) {
	as.version = version
}

func (as AbstractIdServiceImpl) setIdConverter(idConverter IdConverter) {
	as.idConverter = idConverter
}

func (as AbstractIdServiceImpl) setIdMeta(idMeta IdMeta) {
	as.idMeta = idMeta
}

func (as AbstractIdServiceImpl) setMachineIdProvider(machineIdProvider MachineIdProvider) {
	as.machineIdProvider = machineIdProvider
}

func (as AbstractIdServiceImpl) setTimer(timer Timer) {
	as.timer = timer
}
