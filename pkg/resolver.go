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
	"fmt"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

type StaticResolverBuilder struct {
	scheme   string
	max_conn uint16 // Max num of conns for each server
	addrs    []string
}

func NewStaticResolverBuilder(max_conn uint16, addrs []string) *StaticResolverBuilder {
	return &StaticResolverBuilder{
		scheme:   "engula",
		max_conn: max_conn,
		addrs:    addrs,
	}
}

func (b *StaticResolverBuilder) Endpoints() []string {
	// 'default' is target to fix: 'malformed headers: malformed authority (b""): empty string'
	return []string{"engula://default/endpoints"}
}

func (b *StaticResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &StaticResolver{
		max_conn: b.max_conn,
		target:   target,
		cc:       cc,
		addrs:    b.addrs,
	}
	r.start()
	return r, nil
}

func (b *StaticResolverBuilder) Scheme() string {
	return b.scheme
}

type StaticResolver struct {
	max_conn uint16
	target   resolver.Target
	cc       resolver.ClientConn
	addrs    []string
}

func (r *StaticResolver) start() {
	addrs := make([]resolver.Address, 0, len(r.addrs)*int(r.max_conn))
	for _, addr := range r.addrs {
		for i := 0; i < int(r.max_conn); i++ {
			addrs = append(addrs, resolver.Address{Addr: addr, Attributes: attributes.New("tag", fmt.Sprintf("%d", i))})
		}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *StaticResolver) ResolveNow(opt resolver.ResolveNowOptions) {}

func (r *StaticResolver) Close() {}
