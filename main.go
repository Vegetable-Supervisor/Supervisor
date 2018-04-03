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
	http.Handle("/", http.FileServer(http.Dir("public/")))

	// join request handler
	http.HandleFunc("/join", sv.joinHandler)

	log.Printf("Starting SSDP advertisement as %s", sv.usn)
	go sv.SSDPAdvertiser()

	log.Printf("Starting HTTP Server at %s", addr)
	http.ListenAndServe(addr, nil)

}

func (sv *Supervisor) joinHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/join" || r.Method != "POST" {
		http.Error(w, "404 not found.", http.StatusNotFound)
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "could not parse form: %v", err)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		log.Printf("join request from address: %s, empty name")
		return
	}

	remote := r.RemoteAddr
	log.Printf("join request from address: %s, name: %s", remote, name)

	sv.mutex.Lock()
	sv.pending = append(sv.pending, PendingGreenHouse{
		GreenHouse: GreenHouse{
			name:    name,
			address: remote,
		},
	})
	sv.mutex.Unlock()

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "")
}

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
