package main

// A GreenHouse is a GreenHouse as seen by the Supervisor
type GreenHouse struct {
	name    string
	address string
}

// A PendingGreenHouse is a GreenHouse that has not yet been accepted
type PendingGreenHouse struct {
	GreenHouse
}
