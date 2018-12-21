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

type IdConverterImpl struct{}

func (idConverter IdConverterImpl) convert1(id Id, idMeta IdMeta) int64 {
	return doConvert1(id, idMeta)
}

func doConvert1(id Id, idMeta IdMeta) int64 {
	var ret int64 = 0
	ret |= id.getMachine()

	ret |= id.getSeq() << idMeta.getSeqBitsStartPos()

	ret |= id.getTime() << byte(idMeta.getTimeBitsStartPos())

	ret |= id.getGenMethod() << byte(idMeta.getGenMethodBitsStartPos())

	ret |= id.getType() << byte(idMeta.getTypeBitsStartPos())

	ret |= id.getVersion() << byte(idMeta.getVersionBitsStartPos())

	return ret
}

func (idConverter IdConverterImpl) convert2(id int64, idMeta IdMeta) Id {
	return doConvert2(id, idMeta)
}

func doConvert2(id int64, idMeta IdMeta) Id {
	ret := Id{}

	ret.setMachine(id & idMeta.getMachineBitsMask())

	ret.setSeq((id >> byte(idMeta.getSeqBitsStartPos())) & idMeta.getSeqBitsMask())

	ret.setTime((id >> byte(idMeta.getTimeBitsStartPos())) & idMeta.getTimeBitsMask())

	ret.setGenMethod((id >> byte(idMeta.getGenMethodBitsStartPos())) & idMeta.getGenMethodBitsMask())

	ret.setType((id >> byte(idMeta.getTypeBitsStartPos())) & idMeta.getTypeBitsMask())

	ret.setVersion((id >> byte(idMeta.getVersionBitsStartPos())) & idMeta.getVersionBitsMask())

	return ret
}
