// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

func (s *server) Update(e *etcd.Response) {
	// process the first and then loop over nodes
	parts := strings.SplitN(e.Node.Value, ",", 2)
	if len(parts) != 2 {
		log.Printf("Unable to parse node %s with value %s", e.Node.Key, e.Node.Value)
	} else {
		if err := s.router.Add(parts[0], parts[1]); err != nil {
			log.Printf("Unable to add %s", err)
		}
	}
	for _, n := range e.Node.Nodes {
		parts := strings.SplitN(n.Value, ",", 2)
		if len(parts) != 2 {
			log.Printf("Unable to parse node %s with value %s", n.Key, n.Value)
		} else {
			if err := s.router.Add(parts[0], parts[1]); err != nil {
				log.Printf("Unable to add %s", err)
			}
		}
	}
}
