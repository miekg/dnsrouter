// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"regexp"
	"sync"
)

type router struct {
	route map[*regexp.Regexp][]string
	sync.RWMutex
}

func NewRouter() *router {
	return &router{route: make(map[*regexp.Regexp][]string)}
}

func (r *router) Add(re, dest string) error {
	r.Lock()
	defer r.Unlock()
	if net.ParseIP(dest) == nil {
		return fmt.Errorf("not an IP address %s", dest)
	}
	rec, err := regexp.Compile(re)
	if err != nil {
		return err
	}
	if _, ok := r.route[rec]; !ok {
		r.route[rec] = make([]string, 0)
	}
	// check for doubles
	for _, d := range r.route[rec] {
		if d == dest {
			return fmt.Errorf("IP address %s already in list for %s", dest, re)
		}
	}
	r.route[rec] = append(r.route[rec], dest)
	return nil
}

func (r *router) Remove(re, dest string) error {
	r.Lock()
	defer r.Unlock()
	if net.ParseIP(dest) == nil {
		return fmt.Errorf("not an IP address %s", dest)
	}
	rec, err := regexp.Compile(re)
	if err != nil {
		return err
	}
	if _, ok := r.route[rec]; !ok {
		return fmt.Errorf("Regexp %s does not exist", re)
	}
	// remove from list
	return nil
}

func (r *router) Match(qname string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	for re, dest := range r.route {
		if re.MatchString(qname) {
			return dest, nil
		}
	}
	return nil, fmt.Errorf("No match found for %s", qname)
}
