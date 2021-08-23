// Copyright 2021 The Rode Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"log"
	"sync"

	"github.com/rode/rode/common"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc"
)

type rodeClient struct {
	sync.Once
	config *common.ClientConfig
	v1alpha1.RodeClient
	userAgent string
}

var clientInitErr error
func (r *rodeClient) init() error {
	r.Once.Do(func() {
		log.Println("[DEBUG] Rode client init")
		rode, err := common.NewRodeClient(
			r.config,
			grpc.WithUserAgent(r.userAgent),
		)

		if err != nil {
			log.Printf("[ERROR] An error occurred initializing Rode client: %s\n", err)
		}

		r.RodeClient = rode
		clientInitErr = err
		log.Println("[DEBUG] Rode client init successful")
	})

	return clientInitErr
}
