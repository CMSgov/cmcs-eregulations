package main

import (
	"encoding/json"
	"fmt"
)

type Link struct {
	Href string   `json:"href"`
	Name string   `json:"name"`
	Regs []string `json:"regs"`
}

type Regulation struct {
	Header string `json:"header"`
	Links  []Link `json:"links"`
}

func buildRegulation(header string, guidances []Guidance) Regulation {
	links := make([]Link, 0)
	for _, guidance := range guidances {
		linkField := Link{
			Href: guidance.link,
			Name: guidance.name,
			Regs: guidance.regs,
		}
		links = append(links, linkField)
	}

	regulation := Regulation{
		Header: header,
		Links:  links,
	}

	return regulation
}

func toJSON(file []byte, header string, guidances []Guidance) ([]byte, error) {
	var regulations []Regulation

	if len(file) > 0 {
		err := json.Unmarshal(file, &regulations)
		if err != nil {
			fmt.Println(err)
			return []byte{}, err
		}
	}

	newReg := buildRegulation(header, guidances)
	regulations = append(regulations, newReg)

	regsJSON, err := json.MarshalIndent(regulations, "", " ")

	if err != nil {
		fmt.Println("An error has occured :: ", err)
		return []byte{}, err
	}

	return regsJSON, nil
}
