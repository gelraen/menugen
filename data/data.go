package data

import (
	"encoding/json"
)

type Provider interface {
	Dishes() ([]Dish, error)
	Restore([]byte) error
}

type Dish struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url,omitempty"`
}

type dataDump struct {
	Dishes []Dish `json:"dishes"`
}

func Dump(d Provider) ([]byte, error) {
	dishes, err := d.Dishes()
	if err != nil {
		return nil, err
	}
	dump := &dataDump{Dishes: dishes}
	r, err := json.MarshalIndent(dump, "", "  ")
	if err != nil {
		return nil, err
	}
	return r, nil
}
