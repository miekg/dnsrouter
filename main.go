// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

var (
	machines = strings.Split(os.Getenv("ETCD_MACHINES"), ",") // List of URLs to etcd
	dnsaddr  = os.Getenv("DNS_ADDR")                          // Listen on this address
)

func NewClient() (client *etcd.Client) {
	if len(machines) == 1 && machines[0] == "" {
		machines[0] = "http://127.0.0.1:4001"
	}
	client = etcd.NewClient(machines)
	client.SyncCluster()
	return client
}

func main() {
	s := NewServer(NewClient(), dnsaddr)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
