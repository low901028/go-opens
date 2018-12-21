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
	"fmt"
	"strings"
)

// Id: generator unique id
// 最大峰值型: 秒级有序
// version  ctype  genMethod time   seq   machine
// ----------------------------------------------
// |  63  |   62  |  60-61 | 30-59 | 10-29 | 0-9|
// ---------------------------------------------
// 最小粒度型: 毫秒级有序
// version  ctype  genMethod time   seq   machine
// -----------------------------------------------
// |  63  |   62  |  60-61 | 20-59 | 10-19 | 0-9|
// -----------------------------------------------
type Id struct {
	machine   int64
	seq       int64
	time      int64
	genMethod int64
	ctype     int64
	version   int64
}

// seq
func (id *Id) getSeq() int64 {
	return id.seq
}
func (id *Id) setSeq(seq int64) {
	id.seq = seq
}

// time
func (id *Id) getTime() int64 {
	return id.time
}
func (id *Id) setTime(time int64) {
	id.time = time
}

// machine
func (id *Id) getMachine() int64 {
	return id.machine
}
func (id *Id) setMachine(machine int64) {
	id.machine = machine
}

// genMethod
func (id *Id) getGenMethod() int64 {
	return id.genMethod
}
func (id *Id) setGenMethod(genMethod int64) {
	id.genMethod = genMethod
}

// ctype
func (id *Id) getType() int64 {
	return id.ctype
}
func (id *Id) setType(ctype int64) {
	id.ctype = ctype
}

// version
func (id *Id) getVersion() int64 {
	return id.version
}
func (id *Id) setVersion(version int64) {
	id.version = version
}

func (id *Id) String() string {
	var result string
	result = strings.Join([]string{
		fmt.Sprintf("machine=%d", id.machine),
		fmt.Sprintf("seq=%d", id.seq),
		fmt.Sprintf("time=%d", id.time),
		fmt.Sprintf("genMethod=%d", id.genMethod),
		fmt.Sprintf("type=%d", id.ctype),
		fmt.Sprintf("version=%d", id.version),
	}, ",")

	return fmt.Sprintf("[%s]", result)
}
