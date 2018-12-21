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
	"errors"
	"fmt"
)

// IdType: id type
type IdType struct {
	name string
	num  int
}

var IllegalArgumentErr error = errors.New("Illegal IdType name, available names are seconds and milliseconds")
var SECONDS = IdType{"seconds", 0}
var MILLISECONDS = IdType{"milliseconds", 1}
var SHORTID = IdType{"short_id", 2}

var idTypes = []IdType{SECONDS, MILLISECONDS, SHORTID}

func (ctype IdType) value() int64{
	switch ctype {
	case SECONDS:
		return 0
	case MILLISECONDS:
		return 1
	case SHORTID:
		return 2
	default:
		return 0
	}
}

func Parse(name string) (t IdType, err error) {
	for _, it := range idTypes {
		if it.name == name {
			t, err = it, nil
			return
		}
	}
	err = fmt.Errorf("Illegal IdType name <[%s]>, available names are seconds and milliseconds", name)
	return
}

//
func Parse1(num int) (t IdType, err error) {
	for _, it := range idTypes {
		if it.num == num {
			t, err = it, nil
			return
		}
	}
	err = fmt.Errorf("Illegal IdType value <[%d]>, available values are 0 (for seconds) and 1 (for milliseconds)", num)
	return
}
