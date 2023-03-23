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

package main

import (
	"github.com/seata/seata-go/pkg/client"
	"github.com/seata/seata-go/pkg/registry"
	"net"
	"time"
)

func main() {
	cfg := client.LoadPath("/Users/ali/Desktop/GO/vader/seata-go/testdata/conf/seatago.yml")
	register, _ := registry.GetRegistry(&cfg.RegistryConfig)
	address := net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 9001,
	}
	register.RegisterServiceInstance(address)
	register.RegisterServiceInstance(address)
	for i := 0; i < 10; i++ {
		time.Sleep(10000 * time.Second)
	}
}
