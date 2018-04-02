package main

import (
	"fmt"
	"time"
)

type Supervisor struct {
	usn      string // unique service name
	st       string // service type
	location string // location
	server   string // server
}

func NewSupervisor(location string) Supervisor {
	return Supervisor{
		usn:      fmt.Sprintf("vegetable-supervisor-%v", time.Now().Unix()),
		st:       "vegetable-supervisor",
		location: location,
		server:   "vegetable-supervisor",
	}
}
