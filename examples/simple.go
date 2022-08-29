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

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	engula_go "github.com/w41ter/engula-go/pkg"
)

var (
	Addrs   = flag.String("addrs", "", "The addrs of cluster proxy servers")
	DB      = flag.String("db", "test", "The database of operations")
	Collect = flag.String("collection", "test1", "The collection of operations")

	ErrAddrsIsEmpty = errors.New("ADDRS IS EMPTY")
)

func main() {
	flag.Parse()
	fmt.Printf("addrs %s\n", *Addrs)
	addrs := filterEmpty(strings.Split(*Addrs, ","))
	fmt.Printf("addrs %v\n", addrs)
	if len(addrs) == 0 {
		panic(ErrAddrsIsEmpty)
	}

	builder := engula_go.NewStaticResolverBuilder(addrs)
	cfg := engula_go.DefaultConfig()
	cfg.Endpoints = builder.Endpoints()
	cfg.ResolverBuilder = builder
	client, err := engula_go.New(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Create client success\n")

	db, err := client.GetDatabase(context.TODO(), *DB)
	if err == engula_go.ErrNotFound {
		db, err = client.CreateDatabase(context.TODO(), *DB)
	}

	if err != nil {
		panic(err)
	}

	collection, err := db.GetCollection(context.TODO(), *Collect)
	if err == engula_go.ErrNotFound {
		collection, err = db.CreateHashCollection(context.TODO(), *Collect, 3)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("Open collection success\n")

	key := "key1"
	value := "value"
	got, err := collection.Get(context.TODO(), []uint8(key))
	if err != nil {
		panic(err)
	}
	fmt.Printf("get key %s found %v\n", key, got)

	err = collection.Put(context.TODO(), []uint8(key), []uint8(value))
	if err != nil {
		panic(err)
	}

	got, err = collection.Get(context.TODO(), []uint8(key))
	if err != nil {
		panic(err)
	}
	fmt.Printf("after put value %s, get key %s found %v\n", value, key, got)

	err = collection.Delete(context.TODO(), []uint8(key))
	if err != nil {
		panic(err)
	}
}

func filterEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
