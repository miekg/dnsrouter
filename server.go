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
	stop         chan bool
}

// Newserver returns a new server.
func NewServer(client *etcd.Client, addr string) *server {
	if addr == "" {
		addr = "127.0.0.1:53"
	}
	return &server{client: client, addr: addr, group: new(sync.WaitGroup), router: NewRouter(), stop: make(chan bool)}
}

// Run is a blocking operation that starts the server listening on the DNS ports.
func (s *server) Run() error {
	mux := dns.NewServeMux()
	mux.Handle(".", s)

	s.group.Add(2)
	go s.run(mux, "tcp")
	go s.run(mux, "udp")

	// Healthchecking.
	log.Printf("enabling health checking")
	go func() {
		for {
			time.Sleep(5 * 1e9)
			s.HealthCheck()
		}
	}()

	// Set a Watch and check for changes.
	log.Printf("setting watch")
	ch := make(chan *etcd.Response)
	go func() {
		go s.client.Watch("/dnsrouter", 0, true, ch, s.stop)
		for {
			select {
			case n := <-ch:
				s.Update(n)
			}
		}
	}()
	log.Printf("getting initial list")
	n, err := s.client.Get("/dnsrouter/", false, true)
	if err == nil {
		s.Update(n)
	}
	log.Printf("ready for queries")
	s.group.Wait()
	return nil
}

// Stop stops a server.
func (s *server) Stop() {
	s.stop <- true
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

	if q.Qtype == dns.TypeIXFR || q.Qtype == dns.TypeAXFR {
		m := new(dns.Msg)
		m.SetRcode(req, dns.RcodeNotImplemented)
		w.WriteMsg(m)
		return
	}

	allServers, err := s.router.Match(name)
	if err != nil || len(allServers) == 0 {
		m := new(dns.Msg)
		m.SetRcode(req, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}
	serv := allServers[int(dns.Id())%len(allServers)]
	log.Printf("routing %s to %s", name, serv)

	c := new(dns.Client)
	ret, _, err := c.Exchange(req, serv)	 // serv has the port
	if err != nil {
		m := new(dns.Msg)
		m.SetRcode(req, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}
	w.WriteMsg(ret)
}
