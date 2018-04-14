package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// A GreenHouse is a GreenHouse as seen by the Supervisor
type GreenHouse struct {
	Name          string
	Ip            string
	Port          uint64
	Id            uint64
	LastPicture   []byte
	Configuration string
}

// A GreenHouseInformation represents all the information that should be displayed about a GreenHouse
type GreenHouseInformation struct {
	GreenHouse
	Configuration Configuration
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

// getConfiguration retrieves the configuration of the greenhouse
func (gh *GreenHouse) getConfiguration() (Configuration, error) {
	url := fmt.Sprintf("https://%s:%d/get_configuration", gh.Ip, gh.Port)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return Configuration{}, fmt.Errorf("could not get image from greenhouse: %v", err)
	}
	defer resp.Body.Close()

	// decoder := json.NewDecoder(resp.Body)
	// var cnf Configuration
	// err = decoder.Decode(&cnf)
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Configuration{}, fmt.Errorf("could not read configuration from greenhouse: %v", err)
	}

	var cnf Configuration
	err = json.Unmarshal(b, &cnf)

	if err != nil {
		return Configuration{}, fmt.Errorf("could not decode configuration from greenhouse: %v", err)
	}

	return cnf, nil
}

// A PendingGreenHouse is a GreenHouse that has not yet been accepted
type PendingGreenHouse struct {
	GreenHouse
}
