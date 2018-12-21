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

type IdMetaFactory struct {}

// MaxPeak
var maxPeakIdMeta = IdMeta{
		byte(10),
		byte(20),
		byte(30),
		byte(1),
		byte(1),
		byte(1),
}

// minGranularity
var minGranularityIdMeta =  IdMeta{
		byte(10),
		byte(10),
		byte(40),
		byte(1),
		byte(1),
		byte(1),
}

// shortId
var shortIdMeta =  IdMeta{
		byte(10),
		byte(10),
		byte(30),
		byte(1),
		byte(1),
		byte(1),
}

func GetIdMeta(ctype IdType) (IdMeta, bool){
	var idMeta IdMeta
	var isOk = false
	if(SECONDS == (ctype)){
		idMeta = maxPeakIdMeta
		isOk = true
	} else if(MILLISECONDS == (ctype)){
		idMeta = minGranularityIdMeta
		isOk = true
	} else if(SHORTID == ctype){
		idMeta = shortIdMeta
		isOk = true
	}
	return idMeta, isOk
}
