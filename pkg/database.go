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
)

type Database struct {
	client   *Client
	database *proto.DatabaseDesc
}

func (d *Database) GetCollection(ctx context.Context, name string) (*Collection, error) {
	req := &proto.AdminRequest{
		Request: &proto.AdminRequestUnion{
			Request: &proto.AdminRequestUnion_GetCollection{
				GetCollection: &proto.GetCollectionRequest{
					Name:     name,
					Database: d.database,
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := proto.NewEngulaClient(d.client.conn)
	resp, err := client.Admin(ctx, req)
	if err != nil {
		return nil, convertErr(err)
	}
	collection := resp.GetResponse().GetGetCollection().GetCollection()
	if collection == nil {
		return nil, errors.New("invalid response")
	}
	return &Collection{
		client:     d.client,
		database:   d.database,
		collection: collection,
	}, nil
}

func (d *Database) CreateHashCollection(ctx context.Context, name string, slots uint32) (*Collection, error) {
	req := &proto.AdminRequest{
		Request: &proto.AdminRequestUnion{
			Request: &proto.AdminRequestUnion_CreateCollection{
				CreateCollection: &proto.CreateCollectionRequest{
					Name:     name,
					Database: d.database,
					Partition: &proto.CreateCollectionRequest_Hash{
						Hash: &proto.CreateCollectionRequest_HashPartition{
							Slots: slots,
						},
					},
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := proto.NewEngulaClient(d.client.conn)
	resp, err := client.Admin(ctx, req)
	if err != nil {
		return nil, err
	}
	collection := resp.GetResponse().GetCreateCollection().GetCollection()
	if collection == nil {
		return nil, errors.New("invalid response")
	}
	return &Collection{
		client:     d.client,
		database:   d.database,
		collection: collection,
	}, nil
}

func (d *Database) DeleteCollection(ctx context.Context, name string) error {
	req := &proto.AdminRequest{
		Request: &proto.AdminRequestUnion{
			Request: &proto.AdminRequestUnion_DeleteCollection{
				DeleteCollection: &proto.DeleteCollectionRequest{
					Name:     name,
					Database: d.database,
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client := proto.NewEngulaClient(d.client.conn)
	_, err := client.Admin(ctx, req)
	return err
}
