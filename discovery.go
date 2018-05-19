package main

import (
	"fmt"
	"time"

	ssdp "github.com/bcurren/go-ssdp"
)

func (sv *Supervisor) SSDPLookup() error {
	offers, err := ssdp.Search("greenhouse", 3*time.Second)

	if err != nil {
		return fmt.Errorf("error during SSDP lookup: %v", err)
	}

	for _, offer := range offers {
		// fmt.Println(offer.Location)
		// fmt.Println(offer.ResponseAddr)
		ipaddr := offer.ResponseAddr.IP.String()
		port := offer.Location.Port()
		url := fmt.Sprintf("https://%s:%s", ipaddr, port)
		fmt.Println(url)
	}

	return nil
}
