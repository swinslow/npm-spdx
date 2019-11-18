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
	"strings"
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

// IsValidExpression does rudimentary parsing to determine whether
// a string is a valid SPDX license expression, given a particular
// version of the SPDX license list. Note that this algorithm is
// quick-and-dirty and probably needs more attention / expansion.
func IsValidExpression(s string, listIDs map[string]bool) bool {
	// NONE and NOASSERTION are valid
	if s == "NONE" || s == "NOASSERTION" {
		return true
	}

	// if there is only one LPAREN and one RPAREN, and they are
	// at the start and end respectively, strip them off
	if strings.Count(s, "(") == 1 && strings.HasPrefix(s, "(") &&
		strings.Count(s, ")") == 1 && strings.HasSuffix(s, ")") {
		s = strings.TrimPrefix(s, "(")
		s = strings.TrimSuffix(s, ")")
	}

	// ignore all '+' characters, we'll treat them as valid
	strings.ReplaceAll(s, "+", "")

	// if there's ANDs or ORs, split them
	// FIXME currently doing as case-sensitive splits, must be
	// FIXME capitalized; this is not correct
	andStrs := strings.SplitN(s, " AND ", -1)

	// now split each of those by " OR "
	strs := []string{}
	for _, st := range andStrs {
		orStrs := strings.SplitN(st, " OR ", -1)
		for _, orSt := range orStrs {
			strs = append(strs, orSt)
		}
	}

	// finally, strs should now contain the 'processed' set of
	// license IDs, so we can compare each to the license list
	for _, st := range strs {
		if _, ok := listIDs[st]; !ok {
			return false
		}
	}

	return true
}
