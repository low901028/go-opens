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
	"log"
	"strings"
)

type IpConfigurableMachineIdProvider struct {
	machineId int64
	ipMap     map[string]int64
}

func (ip IpConfigurableMachineIdProvider) init() {
	var ips string = GetHostIp()

	if "" == ips || len(ips) == 0 {
		var msg string = "Fail to get host IP address. Stop to initialize the IpConfigurableMachineIdProvider provider."

		log.Fatal(msg)
	}

	if _, ok := ip.ipMap[ips]; !ok {
		msg := fmt.Sprintf("Fail to configure ID for host IP address %s. Stop to initialize the IpConfigurableMachineIdProvider provider.", ips)

		log.Fatal(msg)
	}

	ip.machineId = ip.ipMap[ips]
	log.Printf("IpConfigurableMachineIdProvider.init ip {} id {}", ips, ip.machineId)
}

func (ip IpConfigurableMachineIdProvider) setIps(ips string) {
	if "" == ips || len(ips) == 0 {
		ipArray := strings.Split(ips, ",")

		for i := 0; i < len(ipArray); i++ {
			ip.ipMap[ipArray[i]] = int64(i)
		}
	}
}

func (ip IpConfigurableMachineIdProvider) getMachineId() int64 {
	return ip.machineId
}

func (ip IpConfigurableMachineIdProvider) setMachineId(machineId int64) {
	ip.machineId = machineId
}
