// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"github.com/miekg/dns"
	"log"
)

const HealthQuery = "id.server." // ClassCHAOS, TXT

// healthcheck just removes a server, they need to re-register to get queries again.

func (s *server) HealthCheck() {
	c := new(dns.Client)
	c.Net = "tcp"

	m := new(dns.Msg)
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{HealthQuery, dns.TypeTXT, dns.ClassCHAOS}

	// doing this in the loop is not the best idea
	for _, serv := range s.router.Servers() {
		log.Printf("healthchecking %s", serv)
		if !check(c, m, serv+":53") {
			// do it again
			if !check(c, m, serv+":53") {
				log.Printf("healthcheck failed for %s", serv)
				s.router.RemoveServer(serv)
			}
		}
	}
}

func check(c *dns.Client, m *dns.Msg, addr string) bool {
	m.Id = dns.Id()
	in, _, err := c.Exchange(m, addr)
	if err != nil {
		return false
	}
	if in.Rcode != dns.RcodeSuccess {
		return false
	}
	return true
}
