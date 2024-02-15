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

package discovery

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/seata/seata-go/pkg/util/flagext"
	"github.com/seata/seata-go/pkg/util/log"
)

var (
	onceGetNacosClient = &sync.Once{}
)

type NacosRegistryService struct {
	registry *NacosRegistry
	cluster  string
	client   naming_client.INamingClient
}

func newNacosRegistryService(registry *NacosRegistry, vgroupMapping flagext.StringMap, txServiceGroup string) RegistryService {
	if registry == nil {
		log.Fatalf("registry is nil for nacos")
		panic("registry is nil for nacos")
	}

	// get cluster
	cluster, ok := vgroupMapping[txServiceGroup]
	if !ok {
		panic("tx service group lacks vgroup mapping value")
	}

	return &NacosRegistryService{
		registry: registry,
		cluster:  cluster,
	}
}

func (s *NacosRegistryService) getClient() naming_client.INamingClient {
	if s.client == nil {
		onceGetNacosClient.Do(func() {
			addr := strings.Split(s.registry.ServerAddr, ":")
			if len(addr) != 2 {
				panic("nacos server address should be in format ip:port!")
			}
			port, err := strconv.Atoi(addr[1])
			if err != nil {
				panic("nacos server address port is not valid!")
			}
			//create ServerConfig
			sc := []constant.ServerConfig{
				*constant.NewServerConfig(addr[0], uint64(port), constant.WithContextPath(s.registry.ContextPath)),
			}

			//create ClientConfig
			opts := []constant.ClientOption{
				constant.WithNamespaceId(s.registry.Namespace),
				constant.WithTimeoutMs(5000),
				constant.WithNotLoadCacheAtStart(true),
				constant.WithUsername(s.registry.Username),
				constant.WithPassword(s.registry.Password),
			}
			if s.registry.AccessKey != "" {
				opts = append(opts, constant.WithAccessKey(s.registry.AccessKey))
			}
			if s.registry.SecretKey != "" {
				opts = append(opts, constant.WithSecretKey(s.registry.SecretKey))
			}
			cc := constant.NewClientConfig(
				opts...,
			)

			// create naming client
			client, err := clients.NewNamingClient(
				vo.NacosClientParam{
					ServerConfigs: sc,
					ClientConfig:  cc,
				},
			)

			if err != nil {
				panic("error getting nacos naming client: " + err.Error())
			}
			s.client = client
		})

	}
	return s.client
}

func (s *NacosRegistryService) Lookup(key string) ([]*ServiceInstance, error) {
	param := vo.SelectInstancesParam{
		ServiceName: s.registry.Application,
		GroupName:   s.registry.Group,
		Clusters:    []string{s.cluster},
		HealthyOnly: true,
	}

	// @todo add cache for instance, and subscribe to instance change.
	instances, err := s.getClient().SelectInstances(param)
	if err != nil {
		return nil, fmt.Errorf("error selecting nacos instance for key: %s, %w", key, err)
	}
	res := make([]*ServiceInstance, len(instances))
	for i, instance := range instances {
		res[i] = &ServiceInstance{
			Addr: instance.Ip,
			Port: int(instance.Port),
		}
	}

	return res, nil
}

func (s *NacosRegistryService) Close() {
	if s.client != nil {
		s.client.CloseClient()
	}
}
