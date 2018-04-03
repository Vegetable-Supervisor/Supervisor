package main

import (
	"fmt"
	"sync"
	"time"
)

type Supervisor struct {
	usn         string              // unique service name
	st          string              // service type
	location    string              // location
	server      string              // server
	greenhouses []GreenHouse        // connected GreenHouses
	pending     []PendingGreenHouse // approval pending GreenHouses
	mutex       *sync.Mutex         // for concurrent access of the greenhouses
}

func NewSupervisor(location string) Supervisor {
	return Supervisor{
		usn:      fmt.Sprintf("vegetable-supervisor-%v", time.Now().Unix()),
		st:       "vegetable-supervisor",
		location: location,
		server:   "vegetable-supervisor",
		mutex:    &sync.Mutex{},
	}
}
