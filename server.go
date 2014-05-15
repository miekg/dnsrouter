// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"log"
	"strings"
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
	router       *router
}

// Newserver returns a new server.
func NewServer(client *etcd.Client, addr string) *server {
	if addr == "" {
		addr = "127.0.0.1:53"
	}
	return &server{client: client, addr: addr, group: new(sync.WaitGroup), router: NewRouter()}
}

// Run is a blocking operation that starts the server listening on the DNS ports.
func (s *server) Run() error {
	mux := dns.NewServeMux()
	mux.Handle(".", s)

	s.group.Add(2)
	go s.run(mux, "tcp")
	go s.run(mux, "udp")
	//	log.Printf("connected to etcd cluster at %s", machines)

	// Setup healthchecking
	// Get first list of servers

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

func (s *server) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	name := strings.ToLower(q.Name)

	allServers, err := s.router.Match(name)
	if err != nil {
		m := new(dns.Msg)
		m.SetRcode(req, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}
	serv := allServers[int(dns.Id())%len(allServers)]

	c := new(dns.Client)
	ret, _, err := c.Exchange(req, serv+":53")
	if err != nil {
		m := new(dns.Msg)
		m.SetRcode(req, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}
	w.WriteMsg(ret)
}
