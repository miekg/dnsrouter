// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"log"
	"sync"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/miekg/dns"
)

type server struct {
	client       *etcd.Client
	addr         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	group        *sync.WaitGroup
}

// Newserver returns a new server.
func NewServer(client *etcd.Client, addr string) *server {
	if addr == "" {
		addr = "127.0.0.1:53"
	}
	return &server{client: client, addr: addr, group: new(sync.WaitGroup)}
}

// Run is a blocking operation that starts the server listening on the DNS ports.
func (s *server) Run() error {
	mux := dns.NewServeMux()
	mux.Handle(".", s)

	s.group.Add(2)
	go s.run(mux, "tcp")
	go s.run(mux, "udp")
	log.Printf("connected to etcd cluster at %s", machines)

	s.group.Wait()
	return nil
}

// Stop stops a server.
func (s *server) Stop() {
	s.group.Done()
	s.group.Done()
}

func (s *server) run(mux *dns.ServeMux, net string) {
	defer s.group.Done()

	server := &dns.Server{
		Addr:         s.addr,
		Net:          net,
		Handler:      mux,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// ServeDNS is the handler for DNS requests, responsible for parsing DNS request, possibly forwarding
// it to a real dns server and returning a response.
func (s *server) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	//q := req.Question[0]
	//name := strings.ToLower(q.Name)
}
