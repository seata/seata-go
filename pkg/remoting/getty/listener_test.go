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
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package getty

import (
	"context"
	"github.com/seata/seata-go/pkg/protocol/message"
	"github.com/seata/seata-go/pkg/remoting/processor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGettyClientHandler_OnMessage(t *testing.T) {
	var tests = []struct {
		name string
		pkg  interface{}
	}{
		{
			name: "RpcMessage",
			pkg: message.RpcMessage{
				ID:         1,
				Type:       0,
				Codec:      0,
				Compressor: 0,
				HeadMap:    nil,
				Body: message.GlobalBeginRequest{
					Timeout:         3,
					TransactionName: "test",
				},
			},
		},
		{
			name: "Other",
			pkg:  "pkg",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			GetGettyClientHandlerInstance().OnMessage(nil, test.pkg)
		})
	}
}

type mockProcessor struct {
}

func (m mockProcessor) Process(ctx context.Context, rpcMessage message.RpcMessage) error {
	return nil
}

func TestGettyClientHandler_RegisterProcessor(t *testing.T) {
	var tests = []struct {
		name      string
		msgType   message.MessageType
		processor processor.RemotingProcessor
	}{
		{
			name:      "",
			msgType:   message.MessageType_GlobalBegin,
			processor: &mockProcessor{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ins := GetGettyClientHandlerInstance()
			ins.RegisterProcessor(test.msgType, test.processor)
			assert.Equal(t, test.processor, ins.processorMap[test.msgType])
		})
	}
}
