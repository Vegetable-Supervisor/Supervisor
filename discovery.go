package main

import (
	"log"

	ssdp "github.com/koron/go-ssdp"
)

func (sv *Supervisor) SSDPAdvertiser() {
	_, err := ssdp.Advertise(
		sv.st,       // send as "ST"
		sv.usn,      // send as "USN"
		sv.location, // send as "LOCATION"
		sv.server,   // send as "SERVER"
		0)

	if err != nil {
		log.Fatalf("could not setup discover handler: %v", err)
	}

	// run Advertiser infinitely.
	quit := make(chan bool)
	<-quit
}
