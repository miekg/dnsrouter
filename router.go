// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"sync"
)

type router struct {
	route map[string][]string
	sync.RWMutex
}

func NewRouter() *router {
	return &router{route: make(map[string][]string)}
}

func (r *router) Add(dest, re string) error {
	r.Lock()
	defer r.Unlock()
	// For v6 this needs to be [ipv6]:port .
	// Don't care about port here, just if the syntax is OK.
	ip, _, err := net.SplitHostPort(dest)
	if err != nil {
		return err
	}
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("not an IP address %s", dest)
	}
	_, err = regexp.Compile(re)
	if err != nil {
		return err
	}
	if _, ok := r.route[re]; !ok {
		r.route[re] = make([]string, 0)
	}
	// check for doubles
	for _, d := range r.route[re] {
		if d == dest {
			log.Printf("address %s already in list for %s", dest, re)
			return nil
		}
	}
	log.Printf("adding route %s for %s", re, dest)
	r.route[re] = append(r.route[re], dest)
	return nil
}

func (r *router) Remove(dest, re string) error {
	r.Lock()
	defer r.Unlock()
	_, err := regexp.Compile(re)
	if err != nil {
		return err
	}
	if _, ok := r.route[re]; !ok {
		return fmt.Errorf("Regexp %s does not exist", re)
	}
	for i, s := range r.route[re] {
		if s == dest {
			log.Printf("removing %s", s)
			r.route[re] = append(r.route[re][:i], r.route[re][i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *router) RemoveServer(serv string) {
	for rec, servs := range r.route {
		for _, serv1 := range servs {
			if serv1 == serv {
				// TODO(miek): not optimal to convert this back to strings.
				if err := r.Remove(serv, rec); err != nil {
					log.Printf("%s", err)
				}
			}
		}
	}
}

func (r *router) Match(qname string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	for re, dest := range r.route {
		if ok, _ := regexp.Match(re, []byte(qname)); ok {
			return dest, nil
		}
	}
	return nil, fmt.Errorf("No match found for %s", qname)
}

func (r *router) Servers() []string {
	r.RLock()
	defer r.RUnlock()

	s := make([]string, 0, 5)
	// no de-duplication takes place here
	for _, dest := range r.route {
		s = append(s, dest...)
	}
	return s
}
