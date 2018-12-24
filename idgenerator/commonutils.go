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

import "os"

var (
	SWITCH_ON_EXP  = []string{"ON", "TRUE", "on", "true"}
	SWITCH_OFF_EXP = []string{"OFF", "FALSE", "off", "false"}
)

func IsOn(swtch string) bool {
	for _, sw := range SWITCH_ON_EXP {
		if sw == swtch {
			return true
		}
	}
	return false
}

func isPropKeyOn(key string) bool {
	v := os.Getenv(key)
	for _, sw := range SWITCH_ON_EXP {
		if v == sw {
			return true
		}
	}

	return false
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
