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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

const STORE_FILE_NAME = "machineIdInfo.store"

type MachineIdsIdServiceImpl struct {
	IdServiceImpl
	lastTimestamp int64
	machineIdMap map[int64]int64
	storeFilePath string
	storeFile os.File
	sync.Mutex
}

func (mid MachineIdsIdServiceImpl) init(){
	if(!(reflect.TypeOf(mid.machineIdProvider).Name() == "MachineIdsProvider")){
		log.Fatal("The machineIdProvider is not a MachineIdsProvider instance so that Vesta Service refuses to start.")
	}
	initStoreFile()
	initMachineId()
}

func (mid MachineIdsIdServiceImpl) populateId(id Id) {
	mid.supportChangeMachineId(id)
}

func supportChangeMachine(mid MachineIdsIdServiceImpl) error{

	return nil
}

func (mid MachineIdsIdServiceImpl) supportChangeMachineId(id Id) {
	if (id.getMachine() == mid.machineId){
		mid.Lock()
		changeMachineId()
		resetIdPopulator()
		mid.Unlock()
		supportChangeMachineId(id)
	} else {
		id.setMachine(mid.machineId)
		mid.idPopulator.populateId(mid.timer, id, mid.idMeta)
		mid.lastTimestamp = id.getTime()
	}
}

func (mid MachineIdsIdServiceImpl) changeMachineId() {
	mid.machineIdMap[mid.machineId] = mid.lastTimestamp
	storeInFile()
	initMachineId()
}

func (mid MachineIdsIdServiceImpl) resetIdPopulator() {
	if (reflect.TypeOf(mid.idPopulator).Name() == "ResetPopulator") {
		(RestPopulator(mid.idPopulator.(RestPopulator))).reset()
	} else {
		var newIdPopulator= reflect.New(reflect.TypeOf(mid.idPopulator))
		mid.idPopulator = new(newIdPopulator.Type())
	}
}

func (mid MachineIdsIdServiceImpl) initStoreFile() {
	if (len(mid.storeFilePath) == 0) {
		mid.storeFilePath = os.Getenv("user.dir") + "/" + STORE_FILE_NAME
	}
	log.Printf("machineId info store in <[%s]>", mid.storeFilePath)
	if ok , _ := PathExists(mid.storeFilePath); ok {
		bytes, err := ioutil.ReadFile(mid.storeFilePath)
		if err != io.EOF{
			log.Fatal("read machineId info store file occur exception, ", err)
		}
		lines := strings.Split(string(bytes),"/n/r")
		if len(lines) > 0{
			for _, line := range lines{
				kvs := strings.Split(line, ":")
				if len(kvs) == 2{
					if mid.machineIdMap == nil{
						mid.machineIdMap = make(map[int64]int64)
					}
					mid.machineIdMap[int64(kvs[0])] = int64(kvs[1])
				} else {
					log.Fatal(filepath.Abs(mid.storeFilePath) , " has illegal value <[" + line + "]>")
				}
			}
		}
	}
}

func (mid MachineIdsIdServiceImpl) initMachineId() {
	var startId int64 = mid.machineId
	var newMachineId int64 = mid.machineId
	for {
		if (mid.machineIdMap[(newMachineId)] != 0 ) {
			timestamp := mid.timer.genTime()
			if (mid.machineIdMap[newMachineId] < timestamp) {
				mid.machineId = newMachineId
				break
			} else {
				mip := MachineIdsProvider(mid.machineIdProvider)
				newMachineId = (.getNextMachineId()
			}
			if (newMachineId == startId) {
				log.Fatal("No machineId is available")
			}
			mid.validateMachineId(newMachineId)
		} else {
			mid.machineId = newMachineId
			break
		}
	}
	log.Printf("MachineId: %d\n", mid.machineId)
}

func  (mid MachineIdsIdServiceImpl) storeInFile() {
	var datas []string
	for key, value := range mid.machineIdMap{
		datas = append(datas, fmt.Sprintf("%d:%d",key, value))
	}
	data := strings.Join(datas,"\n")
	ioutil.WriteFile(mid.storeFilePath, []byte(data),0644)
}

func  (mid MachineIdsIdServiceImpl) setStoreFilePath(storeFilePath string) {
	mid.storeFilePath = storeFilePath
}