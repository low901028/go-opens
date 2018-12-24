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
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

type DbWorker struct {
	sync.Mutex
	Dsn string
	Db  *sql.DB
}

type DbMachineIdProvider struct {
	machineId int64
	db        DbWorker
}

func (did DbMachineIdProvider) init() {
	var ip string = GetHostIp()
	if ip == "" || len(ip) <= 0 {
		msg := "Fail to get host IP address. Stop to initialize the DbMachineIdProvider provider."
		log.Fatal(msg)
	}
	var id int64
	id = did.db.queryData(`select ID from DB_MACHINE_ID_PROVIDER where IP = ?`, ip)

	if id != 0 {
		did.machineId = id
		return
	}

	result := did.db.update(`update DB_MACHINE_ID_PROVIDER set IP = ? where IP is null limit 1`, ip)
	if result <= 0 || result > 1 {
		msg := fmt.Sprintf("Fail to allocte ID for host IP address {}. The {} records are updated. Stop to initialize the DbMachineIdProvider provider.",
			ip, result)

		log.Fatal(msg)
	}

	id = did.db.queryData(`select ID from DB_MACHINE_ID_PROVIDER where IP = ?`, ip)
	if id == 0 {
		msg := fmt.Sprintf("Fail to get ID from DB for host IP address %v after allocation. Stop to initialize the DbMachineIdProvider provider.", ip)

		log.Fatal(msg)
	}
	//
	did.machineId = id
}

func (did DbMachineIdProvider) getMachineId() int64 {
	return did.machineId
}

func (did DbMachineIdProvider) setMachineId(machineId int64) {
	did.machineId = machineId
}

// =======================================mysql operation=================================
func (dbw DbWorker) queryData(sql string, parament string) int64 {
	// 1、statement
	stmt, _ := dbw.Db.Prepare(sql)
	defer stmt.Close()

	rows, err := stmt.Query(parament)
	defer rows.Close()
	if err != nil {
		log.Fatal("query data error: %v\n", err)
		return 0
	}
	var ip int64
	for rows.Next() {
		rows.Scan(&ip)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err.Error())
	}
	return ip
}

func (dbw DbWorker) update(sql string, parament string) int64 {
	dbw.Db.Begin()
	stmt, _ := dbw.Db.Prepare(sql)
	defer stmt.Close()

	res, err := stmt.Exec(parament)
	if err != nil {
		log.Printf("update data error: %v\n", err)
		return 0
	}
	count, _ := res.RowsAffected()
	return count
}

func (dbw DbWorker) insertData(sql string, parament string) {
	stmt, _ := dbw.Db.Prepare(sql)
	defer stmt.Close()

	_, err := stmt.Exec(parament)
	if err != nil {
		log.Printf("insert data error: %v\n", err)
		return
	}
}
