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
	"database/sql"
	"fmt"
	"log"
)

type Type string
type TypeEnum struct {
	PROPERTY, IP_CONFIGURABLE, DB Type
}

// "user:admin1234@tcp(127.0.0.1:3306)/vesta-id-generator",
var DBService = IdServiceFactoryBean{
	providerType: "DB",
	dbUrl:        "127.0.0.1:3306",
	dbName:       "vesta_id_generator",
	dbUser:       "root",
	dbPassword:   "admin1234",
}

type IdServiceFactoryBean struct {
	providerType Type
	machineId    int64
	ips          string
	dbUrl        string
	dbName       string
	dbUser       string
	dbPassword   string
	genMethod    int64
	ctype        int64
	version      int64
	idService    IdService
}

func (idFac IdServiceFactoryBean) Init() {
	fmt.Println("=====here=====")
	//idFac.providerType = "DB"
	switch idFac.providerType {
	case "PROPERTY":
		idFac.idService = idFac.constructPropertyIdService(idFac.machineId).(IdService)
		break
	case "IP_CONFIGURABLE":
		idFac.idService = idFac.constructIpConfigurableIdService(idFac.ips).(IdService)
		break
	case "DB":
		idFac.idService = idFac.constructDbIdService(idFac.dbUrl, idFac.dbName, idFac.dbUser, idFac.dbPassword).(IdService)
		break

	}
}

func (idFac IdServiceFactoryBean) constructPropertyIdService(machineId int64) interface{} {
	log.Println("Construct Property IdService machineId {}", machineId)

	propertyMachineIdProvider := PropertyMachineIdProvider{
		machineId: machineId,
	}

	var idServiceImpl IdServiceImpl
	if idFac.ctype != -1 {
		ctype, _ := Parse1(int(idFac.ctype))
		idServiceImpl = IdServiceImpl{
			iDType: ctype,
		}
	} else {
		idServiceImpl = IdServiceImpl{}
	}

	idServiceImpl.setMachineIdProvider(propertyMachineIdProvider)
	if idFac.genMethod != -1 {
		idServiceImpl.setGenMethod(idFac.genMethod)
	}

	if idFac.version != -1 {
		idServiceImpl.setVersion(idFac.version)
	}
	idServiceImpl.init()

	return idServiceImpl
}

func (idFac IdServiceFactoryBean) constructIpConfigurableIdService(ips string) interface{} {
	log.Printf("Construct Ip Configurable IdService ips {}", idFac.ips)
	ipConfigurableMachineIdProvider := IpConfigurableMachineIdProvider{}
	ipConfigurableMachineIdProvider.setIps(ips)

	var idServiceImpl IdServiceImpl
	if idFac.ctype != -1 {
		ctype, _ := Parse1(int(idFac.ctype))
		idServiceImpl = IdServiceImpl{
			iDType: ctype,
		}
	} else {
		idServiceImpl = IdServiceImpl{}
	}

	idServiceImpl.setMachineIdProvider(ipConfigurableMachineIdProvider)
	if idFac.genMethod != -1 {
		idServiceImpl.setGenMethod(idFac.genMethod)
	}
	if idFac.version != -1 {
		idServiceImpl.setVersion(idFac.version)
	}
	idServiceImpl.init()

	return idServiceImpl
}

func (idFac IdServiceFactoryBean) constructDbIdService(dbUrl string, dbName string, dbUser string, dbPassword string) interface{} {
	dbMachineIdProvider := DbMachineIdProvider{}

	dbMachineIdProvider.db = DbWorker{
		Dsn: fmt.Sprintf("%v:%v@tcp(%v)/%v", dbUser, dbPassword, dbUrl, dbName), // "user:admin1234@tcp(127.0.0.1:3306)/vesta_id_generator",
	}
	db, _ := sql.Open("mysql", dbMachineIdProvider.db.Dsn)
	dbMachineIdProvider.db.Db = db
	dbMachineIdProvider.init()

	var idServiceImpl IdServiceImpl
	if idFac.ctype != -1 {
		ctype, _ := Parse1(int(idFac.ctype))
		idServiceImpl = IdServiceImpl{
			iDType: ctype,
		}
	} else {
		idServiceImpl = IdServiceImpl{}
	}

	idServiceImpl.setMachineIdProvider(dbMachineIdProvider)
	if idFac.genMethod != -1 {
		idServiceImpl.setGenMethod(idFac.genMethod)
	}
	if idFac.version != -1 {
		idServiceImpl.setVersion(idFac.version)
	}
	idServiceImpl.init()

	return idServiceImpl
}
