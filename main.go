package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.Uint("port", 8080, "http server port")
	ip := flag.String("ip", "localhost", "http server ip")
	addr := fmt.Sprintf("%s:%d", *ip, *port)

	sv := NewSupervisor(addr)

	err := sv.SSDPLookup()
	fmt.Println(err)

	// // index page
	// http.HandleFunc("/", sv.homeHandler)
	//
	// // pending greenhouses page
	// http.HandleFunc("/pending", sv.pendingHandler)
	//
	// // info
	// http.HandleFunc("/info", sv.infoHandler)
	//
	// // join request handler
	// http.HandleFunc("/join", sv.joinHandler)
	//
	// http.HandleFunc("/camera", sv.getPictureHandler)
	//
	// http.HandleFunc("/push_configuration", sv.pushConfigurationHandler)
	//
	// log.Printf("Starting HTTP Server at %s", addr)
	// http.ListenAndServe(addr, nil)
}
