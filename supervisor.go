package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Supervisor struct {
	usn                string                       // unique service name
	st                 string                       // service type
	location           string                       // location
	server             string                       // server
	greenhouses        map[uint64]GreenHouse        // connected GreenHouses
	pending            map[uint64]PendingGreenHouse // approval pending GreenHouses
	mutex              *sync.Mutex                  // for concurrent access of the greenhouses
	next_greenhouse_id uint64                       // next id for greenhouses
}

func NewSupervisor(location string) Supervisor {
	return Supervisor{
		usn:                fmt.Sprintf("vegetable-supervisor-%v", time.Now().Unix()),
		st:                 "vegetable-supervisor",
		location:           location,
		server:             "vegetable-supervisor",
		greenhouses:        make(map[uint64]GreenHouse),
		pending:            make(map[uint64]PendingGreenHouse),
		mutex:              &sync.Mutex{},
		next_greenhouse_id: 0,
	}
}

func (sv *Supervisor) homeHandler(w http.ResponseWriter, r *http.Request) {
	pageTemplate := template.Must(template.ParseFiles("public/templates/index.html"))
	sv.mutex.Lock()
	err := pageTemplate.Execute(w, sv.greenhouses)
	sv.mutex.Unlock()
	if err != nil {
		log.Printf("could not execute template: %v", err)
	}
}

func (sv *Supervisor) pendingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/pending" {
		http.Error(w, "404 not found.", http.StatusNotFound)
	}

	switch r.Method {
	case "GET":
		pageTemplate := template.Must(template.ParseFiles("public/templates/pending.html"))
		sv.mutex.Lock()
		err := pageTemplate.Execute(w, sv.pending)
		sv.mutex.Unlock()
		if err != nil {
			log.Fatalf("could not execute template: %v", err)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "error parsing form in /pending: %v", err)
			return
		}
		// ghName := r.FormValue("greenhouse_name")

		ghId, err := strconv.ParseUint(r.FormValue("greenhouse_id"), 0, 64)
		if err != nil {
			fmt.Fprintf(w, "greenhouse id must be unsigned integer")
		}
		action := r.FormValue("action")

		if action == "accept" {
			// accept pending greenhouse
			sv.mutex.Lock()

			gh, ok := sv.pending[ghId]
			if !ok {
				// not pending
				fmt.Fprintf(w, "greenhouse already accepted or not existing")
				sv.mutex.Unlock()
				return
			}

			delete(sv.pending, gh.Id)

			sv.greenhouses[ghId] = gh.GreenHouse
			log.Printf("starting %v", gh)

			sv.mutex.Unlock()
			// fmt.Fprintf(w, "successfully added greenhouse")
		} else if action == "deny" {
			sv.mutex.Lock()
			gh, ok := sv.pending[ghId]
			if !ok {
				// not pending
				fmt.Fprintf(w, "greenhouse already accepted or not existing")
				sv.mutex.Unlock()
				return
			}

			delete(sv.pending, gh.Id)
			sv.mutex.Unlock()

			// fmt.Fprintf(w, "successfully denied greenhouse")
		}

		pageTemplate := template.Must(template.ParseFiles("public/templates/pending.html"))
		sv.mutex.Lock()
		err = pageTemplate.Execute(w, sv.pending)
		sv.mutex.Unlock()
		if err != nil {
			log.Fatalf("could not execute template: %v", err)
		}

		// log.Printf("Post from website! r.PostFrom = %v\n", r.PostForm)
	default:
		http.Error(w, "404 not found.", http.StatusNotFound)
	}

}

func (sv *Supervisor) infoHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	ghId, err := strconv.ParseUint(queryValues.Get("id"), 0, 64)
	if err != nil {
		http.Error(w, "bad greenhouse id.", http.StatusNotFound)
		return
	}

	sv.mutex.Lock()
	_, ok := sv.greenhouses[ghId]
	sv.mutex.Unlock()

	if !ok {
		// not in accepted greenhouses, might be pending
		sv.pendingInfoHandler(w, r)
		return
	}

	fmt.Fprintf(w, "OK :)")

}

func (sv *Supervisor) pendingInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK pending:)")
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

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("http request remote address is not of the form ip:port: %v")
		return
	}

	log.Printf("join request from address: %s, name: %s", ip, name)

	ghId := sv.next_greenhouse_id

	sv.mutex.Lock()

	if gh, ok := sv.pending[ghId]; ok {
		sv.mutex.Unlock()
		log.Fatalf("tried to add an existing greehouse :%v", gh)
	}

	sv.pending[ghId] = PendingGreenHouse{
		GreenHouse: GreenHouse{
			Name: name,
			Ip:   ip,
			Port: 5000,
			Id:   ghId,
		},
	}

	sv.mutex.Unlock()

	sv.next_greenhouse_id += 1

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "")
}

func (sv *Supervisor) getPictureHandler(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	ghId, err := strconv.ParseUint(queryValues.Get("id"), 0, 64)
	if err != nil {
		http.Error(w, "bad greenhouse id.", http.StatusNotFound)
		return
	}

	sv.mutex.Lock()
	gh, ok := sv.greenhouses[ghId]
	sv.mutex.Unlock()

	if !ok {
		http.Error(w, "this greenhouse does not exist or has not been accepted.", http.StatusNotFound)
		return
	}

	data, err := gh.getPicture()
	if err != nil {
		http.Error(w, "Could not get requested camera image.", http.StatusNotFound)
		log.Printf("could not get picture of %v: %v", gh, err)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Could not send requested camera image.", http.StatusNotFound)
		return
	}

	// success
}
