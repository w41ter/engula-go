// Copyright 2022 The Engula Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engula_go

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type Config struct {
	Context         context.Context
	DialTimeout     time.Duration
	DialOptions     []grpc.DialOption
	Endpoints       []string
	ResolverBuilder resolver.Builder
}

func DefaultConfig() *Config {
	return &Config{
		Context:         context.Background(),
		DialTimeout:     time.Second * 3,
		DialOptions:     nil,
		Endpoints:       nil,
		ResolverBuilder: nil,
	}
}
