package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Uint("port", 8080, "http server port")
	ip := flag.String("ip", "localhost", "http server ip")
	addr := fmt.Sprintf("%s:%d", *ip, *port)

	sv := NewSupervisor(addr)

	// register route handlers

	// default handler
	// http.Handle("/", http.FileServer(http.Dir("public/")))

	// index
	http.HandleFunc("/", sv.homeHandler)

	// pending greenhouses page
	http.HandleFunc("/pending", sv.pendingHandler)

	// info
	http.HandleFunc("/info", sv.infoHandler)

	// Special Purpose Handlers

	// join request handler
	http.HandleFunc("/join", sv.joinHandler)

	http.HandleFunc("/camera", sv.getPictureHandler)

	log.Printf("Starting SSDP advertisement as %s", sv.usn)
	go sv.SSDPAdvertiser()

	log.Printf("Starting HTTP Server at %s", addr)
	http.ListenAndServe(addr, nil)
}
