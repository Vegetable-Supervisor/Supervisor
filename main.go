package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	port := flag.Uint("port", 8080, "http server port")
	ip := flag.String("ip", "localhost", "http server ip")
	addr := fmt.Sprintf("%s:%d", *ip, *port)
	sv := NewSupervisor(addr)
	// register route handlers
	http.HandleFunc("/", rootHandler)

	log.Printf("Starting SSDP advertisement as %s", sv.usn)
	go sv.SSDPAdvertiser()

	log.Printf("Starting HTTP Server at %s", addr)
	http.ListenAndServe(addr, nil)

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Vegetable Supervisor")
}

//
// func getPicture() []byte {
//
// 	resp, err := http.Get("http://127.0.0.1:5000/picture")
// 	if err != nil {
// 		fmt.Printf("could not get picture: %v", err)
// 		return nil
// 	}
// 	defer resp.Body.Close()
//
// 	data, err := ioutil.ReadAll(resp.Body)
//
// 	if err != nil {
// 		log.Fatalf("could not read response: %v", err)
// 	}
// 	return data
// }
