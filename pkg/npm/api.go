// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GetVersionData queries the NPM API for a specified package
// and version, and returns a RegistryVersion object with
// that data.
func GetVersionData(pkg string, ver string) (*RegistryVersion, error) {
	url := "https://registry.npmjs.org/" + pkg + "/" + ver

	// get data from NPM API
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("while getting %s/%s: error retrieving NPM registry data: %v", pkg, ver, err)
	}

	// read response body
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("while getting %s/%s: error reading NPM registry data: %v", pkg, ver, err)
	}

	// parse JSON response
	rver := RegistryVersion{}
	err = json.Unmarshal(b, &rver)
	if err != nil {
		return nil, fmt.Errorf("while getting %s/%s: error parsing NPM registry data: %v", pkg, ver, err)
	}

	return &rver, nil
}

// GetAllVersionData takes a map of names to PackageLockDependencies
// and, for each one, retrieves the corresponding RegistryVersion object.
// It compiles that information into a RegistryResults object, which
// it then returns.
// The sleep parameter configures the time period in milliseconds
// for a pause between each API call, to avoid overloading the API.
func GetAllVersionData(deps map[string]*PackageLockDependency, ms int) (RegistryResults, error) {
	allDeps := RegistryResults{}
	for depName, depData := range deps {
		// check whether we already have a RegistryPackage for this
		var pkg RegistryPackage
		pkg, ok := allDeps[depName]
		if !ok {
			pkg = RegistryPackage{}
			allDeps[depName] = pkg
		}

		// query the API and get data
		rver, err := GetVersionData(depName, depData.Version)
		if err != nil {
			return nil, err
		}

		// add it to the local results
		pkg[depData.Version] = rver

		// and sleep for a bit
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	return allDeps, nil
}

// SaveVersionData saves the retrieved version data to disk as a
// JSON file representation of a RegistryResults object.
func SaveVersionData(allResults RegistryResults, filename string) error {
	jsResults, err := json.Marshal(allResults)
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}

	err = ioutil.WriteFile(filename, jsResults, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to %s: %v", filename, err)
	}

	return nil
}

// LoadVersionData reads the saved version data from disk as a
// JSON file representation of a RegistryResults object.
func LoadVersionData(filename string) (RegistryResults, error) {
	js, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filename, err)
	}

	allResults := RegistryResults{}
	err = json.Unmarshal(js, &allResults)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling from JSON: %v", err)
	}

	return allResults, nil
}
