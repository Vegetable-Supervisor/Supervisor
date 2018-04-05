package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

// A GreenHouse is a GreenHouse as seen by the Supervisor
type GreenHouse struct {
	Name        string
	Ip          string
	Port        uint64
	Id          uint64
	LastPicture []byte
}

func (gh GreenHouse) String() string {
	return fmt.Sprintf("GreenHouse name:%s, id:%d", gh.Name, gh.Id)
}

func (gh *GreenHouse) getPicture() ([]byte, error) {
	url := fmt.Sprintf("https://%s:%d/camera", gh.Ip, gh.Port)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get image from greenhouse: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read greenhouse picture response: %v", err)
	}

	return data, nil
}

// A PendingGreenHouse is a GreenHouse that has not yet been accepted
type PendingGreenHouse struct {
	GreenHouse
}
