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
package main

import (
	"github.com/apache/incubator-rocketmq-externals/rocketmq-go/api"
	"github.com/apache/incubator-rocketmq-externals/rocketmq-go/api/model"

	"github.com/golang/glog"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	var (
		testTopic = "GoLang"
	)

	var producerConfig = rocketmq_api_model.NewProducerConfig()
	producerConfig.CompressMsgBodyOverHowMuch = 1
	var producer1 = rocketmq_api.NewDefaultMQProducerWithCustomConfig("Test1", producerConfig)
	var consumerConfig = rocketmq_api_model.NewRocketMqConsumerConfig()
	//for test
	consumerConfig.PullInterval = 0
	consumerConfig.ConsumeTimeout = 1
	consumerConfig.ConsumeMessageBatchMaxSize = 16
	consumerConfig.ConsumeFromWhere = rocketmq_api_model.CONSUME_FROM_TIMESTAMP
	consumerConfig.ConsumeTimestamp = time.Now()
	var comsumer1 = rocketmq_api.NewDefaultMQPushConsumerWithCustomConfig(testTopic+"-StyleTang", consumerConfig)
	comsumer1.Subscribe(testTopic, "*")
	comsumer1.RegisterMessageListener(func(msgs []rocketmq_api_model.MessageExt) rocketmq_api_model.ConsumeConcurrentlyResult {
		for _, msg := range msgs {
			glog.Info(msg.BornTimestamp)
		}
		glog.Info("look message len ", len(msgs))
		return rocketmq_api_model.ConsumeConcurrentlyResult{ConsumeConcurrentlyStatus: rocketmq_api_model.CONSUME_SUCCESS, AckIndex: len(msgs)}
	})
	nameServerAddress := "127.0.0.1:9876"
	rocketMqManager := rocketmq_api.InitRocketMQClientInstance(nameServerAddress)
	rocketMqManager.RegisterProducer(producer1)
	rocketMqManager.RegisterConsumer(comsumer1)
	rocketMqManager.Start()
	for i := 0; i < 10000000; i++ {
		var message = &rocketmq_api_model.Message{}
		message.Topic = testTopic
		message.SetKeys([]string{"xxx"})
		message.SetTag("1122")
		message.Body = []byte("hellAXXWord")

		xx, ee := producer1.Send(message)
		if ee != nil {
			glog.Error(ee)
			continue
		}
		glog.V(0).Infof("sendMessageResutl messageId[%s] err[%s]", xx.MsgID(), ee)
	}
	select {}
}