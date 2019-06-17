// Package spdxlicenses does a simple parse of the SPDX
// license-list-data JSON files, and creates a catalog of
// valid license IDs.
// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.
package spdxlicenses

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type licEntry struct {
	LicenseID string `json:"licenseId"`
}

type licList struct {
	Licenses []licEntry `json:"licenses"`
}

type excEntry struct {
	LicenseExceptionID string `json:"licenseExceptionId"`
}

type excList struct {
	Exceptions []excEntry `json:"exceptions"`
}

// ParseJSONLicenses parses the SPDX license-list-data
// licenses.json and exceptions.json files, and returns a
// single set of strings (as a map) containing just the set of
// valid license IDs.
func ParseJSONLicenses(licensesPath, exceptionsPath string) (map[string]bool, error) {
	// load licenses
	ll := licList{}

	js, err := ioutil.ReadFile(licensesPath)
	if err != nil {
		return nil, fmt.Errorf("error reading license list from %s: %v", licensesPath, err)
	}

	err = json.Unmarshal(js, &ll)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling from JSON: %v", err)
	}

	// load exceptions
	el := excList{}

	js, err = ioutil.ReadFile(exceptionsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading exceptions list from %s: %v", exceptionsPath, err)
	}

	err = json.Unmarshal(js, &el)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling from JSON: %v", err)
	}

	// and combine them
	allLics := map[string]bool{}
	for _, l := range ll.Licenses {
		allLics[l.LicenseID] = true
	}
	for _, e := range el.Exceptions {
		allLics[e.LicenseExceptionID] = true
	}

	return allLics, nil
}
