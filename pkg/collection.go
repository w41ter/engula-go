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

type Collection struct {
	client     *Client
	database   *proto.DatabaseDesc
	collection *proto.CollectionDesc
}

func (c *Collection) Get(ctx context.Context, key []uint8) ([]uint8, error) {
	client := proto.NewEngulaClient(c.client.conn)
	req := &proto.DatabaseRequest{
		Database: c.database,
		Request: &proto.CollectionRequest{
			Collection: c.collection,
			Request: &proto.CollectionRequestUnion{
				Request: &proto.CollectionRequestUnion_Get{
					Get: &proto.GetRequest{
						Key: key,
					},
				},
			},
		},
	}
	resp, err := client.Database(ctx, req)
	if err != nil {
		return nil, err
	}
	get_resp := resp.GetResponse().GetResponse().GetGet()
	if get_resp == nil {
		return nil, errors.New("invalid response")
	}
	return get_resp.GetValue(), nil
}

func (c *Collection) Put(ctx context.Context, key, value []uint8) error {
	client := proto.NewEngulaClient(c.client.conn)
	req := &proto.DatabaseRequest{
		Database: c.database,
		Request: &proto.CollectionRequest{
			Collection: c.collection,
			Request: &proto.CollectionRequestUnion{
				Request: &proto.CollectionRequestUnion_Put{
					Put: &proto.PutRequest{
						Key:   key,
						Value: value,
					},
				},
			},
		},
	}
	_, err := client.Database(ctx, req)
	return err
}

func (c *Collection) Delete(ctx context.Context, key []uint8) error {
	client := proto.NewEngulaClient(c.client.conn)
	req := &proto.DatabaseRequest{
		Database: c.database,
		Request: &proto.CollectionRequest{
			Collection: c.collection,
			Request: &proto.CollectionRequestUnion{
				Request: &proto.CollectionRequestUnion_Delete{
					Delete: &proto.DeleteRequest{
						Key: key,
					},
				},
			},
		},
	}
	_, err := client.Database(ctx, req)
	return err
}
