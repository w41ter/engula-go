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
	"errors"

	"github.com/w41ter/engula-go/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	cfg    *Config
	conn   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func New(cfg *Config) (*Client, error) {
	if len(cfg.Endpoints) == 0 {
		return nil, ErrNoAvailableEndpoints
	}

	if cfg.Context == nil {
		cfg.Context = context.Background()
	}
	ctx, cancel := context.WithCancel(cfg.Context)
	conn, err := grpc.DialContext(ctx, cfg.Endpoints[0],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithResolvers(cfg.ResolverBuilder))
	if err != nil {
		cancel()
		return nil, err
	}

	client := &Client{
		cfg:    cfg,
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
	}

	return client, nil
}

func (c *Client) Close() {
	c.cancel()
	c.conn.Close()
}

func (c *Client) CreateDatabase(ctx context.Context, name string) (*Database, error) {
	req := &proto.AdminRequest{
		Request: &proto.AdminRequestUnion{
			Request: &proto.AdminRequestUnion_CreateDatabase{
				CreateDatabase: &proto.CreateDatabaseRequest{
					Name: name,
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := proto.NewEngulaClient(c.conn)
	resp, err := client.Admin(ctx, req)
	if err != nil {
		return nil, err
	}
	database := resp.GetResponse().GetCreateDatabase().GetDatabase()
	if database == nil {
		return nil, errors.New("invalid response")
	}

	return &Database{
		database: database,
		client:   c,
	}, nil
}

func (c *Client) GetDatabase(ctx context.Context, name string) (*Database, error) {
	req := &proto.AdminRequest{
		Request: &proto.AdminRequestUnion{
			Request: &proto.AdminRequestUnion_GetDatabase{
				GetDatabase: &proto.GetDatabaseRequest{
					Name: name,
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := proto.NewEngulaClient(c.conn)
	resp, err := client.Admin(ctx, req)
	if err != nil {
		return nil, convertErr(err)
	}
	database := resp.GetResponse().GetGetDatabase().GetDatabase()
	if database == nil {
		return nil, errors.New("invalid response")
	}

	return &Database{
		database: database,
		client:   c,
	}, nil
}

func (c *Client) DeleteDatabase(ctx context.Context, name string) error {
	req := &proto.AdminRequest{
		Request: &proto.AdminRequestUnion{
			Request: &proto.AdminRequestUnion_DeleteDatabase{
				DeleteDatabase: &proto.DeleteDatabaseRequest{
					Name: name,
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := proto.NewEngulaClient(c.conn)
	_, err := client.Admin(ctx, req)
	return err
}
