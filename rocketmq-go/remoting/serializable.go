/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
package remoting

import (
	"github.com/apache/incubator-rocketmq-externals/rocketmq-go/model/constant"
	"github.com/golang/glog"
	"bytes"
	"encoding/binary"
)

type SerializerHandler struct {
	serializer Serializer //which serializer this client use, depend on  constant.USE_HEADER_SERIALIZETYPE
}

type Serializer interface {
	EncodeHeaderData(request *RemotingCommand) []byte
	DecodeRemoteCommand(header, body []byte) *RemotingCommand
}

func NewSerializerHandler() SerializerHandler {
	serializerHandler := SerializerHandler{}
	switch constant.USE_HEADER_SERIALIZETYPE {
	case constant.JSON_SERIALIZE:
		serializerHandler.serializer = &JsonSerializer{}
		break

	case constant.ROCKETMQ_SERIALIZE:
		serializerHandler.serializer = &RocketMqSerializer{}
		break
	default:
		panic("illeage serializer type");
	}
	return serializerHandler
}
func (self *SerializerHandler) EncodeHeader(request *RemotingCommand) []byte {
	length := 4
	headerData := self.serializer.EncodeHeaderData(request)
	length += len(headerData)
	if request.Body != nil {
		length += len(request.Body)
	}
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(length)) // len
	binary.Write(buf, binary.BigEndian, int32(len(headerData) | (int(constant.USE_HEADER_SERIALIZETYPE) << 24))) // header len
	buf.Write(headerData)
	var look = buf.Bytes()
	return look
	return self.serializer.EncodeHeaderData(request)
}

func (self *SerializerHandler) DecodeRemoteCommand(headerSerializableType byte, header, body []byte) *RemotingCommand {
	//  todo singleton
	var serializer Serializer
	switch headerSerializableType {
	case constant.JSON_SERIALIZE:
		serializer = &JsonSerializer{}
		break
	case constant.ROCKETMQ_SERIALIZE:
		serializer = &RocketMqSerializer{}
		break
	default:
		glog.Error("Unknow headerSerializableType", headerSerializableType)
	}
	return serializer.DecodeRemoteCommand(header, body);
}

