// Copyright (c) 2014 Miek Gieben. All rights reserved.
// Use of this source code is governed by The GPL License version 2 
// (GPLv2) that can be found in the LICENSE file.

package main

import (
	"net"
	"regexp"
)

type router map[*regexp.Regexp][]net.IP
