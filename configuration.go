package main

import (
	"fmt"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)

// A Configuration represents the configuration of a GreenHouse.
type Configuration struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewConfigurationFromForm constructs a Configuration from the received HTML Form.
// It does input sanitization and validation
func NewConfigurationFromForm(r *http.Request) (Configuration, error) {
	p := bluemonday.StrictPolicy()

	name := p.Sanitize(r.PostFormValue("name"))
	if len(name) <= 0 {
		return Configuration{}, fmt.Errorf("name should be non empty")
	}

	description := p.Sanitize(r.PostFormValue("description"))

	c := Configuration{
		Name:        name,
		Description: description,
	}
	return c, nil
}
