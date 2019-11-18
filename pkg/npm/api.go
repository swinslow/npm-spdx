// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetVersionData queries the NPM API for a specified package
// and version, and returns a RegistryVersion object with
// that data.
func GetVersionData(pkg string, ver string) (*RegistryVersion, error) {
	// apparently the NPM API doesn't allow retrieving just one version
	// info for scoped packages. If it's a scoped package, go to a
	// different function to parse it from the full response with
	// all versions
	if strings.HasPrefix(pkg, "@") {
		return GetVersionDataForScopedPackage(pkg, ver)
	}

	fmt.Printf("Getting data for %s/%s\n", pkg, ver)
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

// GetVersionDataForScopedPackage queries the NPM API for a
// specified _scoped_ package and version, and returns a
// RegistryVersion object with that data.
func GetVersionDataForScopedPackage(pkg string, ver string) (*RegistryVersion, error) {
	fmt.Printf("Getting data for scoped package %s, version %s\n", pkg, ver)
	url := "https://registry.npmjs.org/" + url.PathEscape(pkg)

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
	rsp := RegistryScopedPackage{}
	err = json.Unmarshal(b, &rsp)
	if err != nil {
		return nil, fmt.Errorf("while getting %s/%s: error parsing NPM registry data: %v", pkg, ver, err)
	}

	rver, ok := rsp.Versions[ver]
	if !ok {
		return nil, fmt.Errorf("while getting %s/%s: version %s not found", pkg, ver, ver)
	}

	return rver, nil
}

// GetAllDependencies takes a map of names to PackageLockDependencies
// and, for each one, retrieves the corresponding RegistryVersion object.
// It compiles that information into a DependencyResults object, which
// it then returns.
// It also takes a PackageManifest (e.g., a parsed package.json file) so
// that it can note which dependencies are direct or direct dev.
// The ms parameter configures the sleep time period in milliseconds
// for a pause between each API call, to avoid overloading the API.
func GetAllDependencies(deps map[string]*PackageLockDependency, manifest *PackageManifest, ms int) (map[string]*Dependency, error) {
	allDeps := map[string]*Dependency{}
	for depName, depData := range deps {
		// check whether we already have a Dependency for this
		// (which shouldn't happen...)
		existingDep, ok := allDeps[depName]
		if ok {
			return allDeps, fmt.Errorf("while getting %s/%s: already got dependency %s (version %s)", depName, depData.Version, existingDep.Name, existingDep.Version)
		}
		d := &Dependency{}

		// query the API and get data
		rver, err := GetVersionData(depName, depData.Version)
		if err != nil {
			return nil, err
		}

		// translate RegistryVersion into a Dependency object
		d.Name = rver.Name
		d.Version = rver.Version
		d.Dependencies = rver.Dependencies
		d.DevDependencies = rver.DevDependencies
		// default to NOASSERTION in case we can't fill it in below
		d.License = "NOASSERTION"

		// also translate the license field, which could be a string or
		// a JSON object (thanks npm)
		if t, ok := rver.License.(string); ok {
			// it's just a string, hooray
			d.License = t
		} else if t, ok := rver.License.(map[string]interface{}); ok {
			// it's an object, we can look for an appropriate field
			if val, ok := t["type"]; ok {
				// let's also make sure val is a string...
				if s, ok := val.(string); ok {
					d.License = s
				}
			}
		}

		// also note whether it's a direct dependency and/or direct dev dep
		if _, ok := manifest.Dependencies[depName]; ok {
			d.IsDirectDep = true
		}
		if _, ok := manifest.DevDependencies[depName]; ok {
			d.IsDirectDevDep = true
		}

		// and finally, add it to the local results
		allDeps[depName] = d

		// and sleep for a bit
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	return allDeps, nil
}

// SaveResults saves the retrieved version data to disk as a
// JSON file representation of a DependencyResults object.
func SaveResults(dr *DependencyResults, filename string) error {
	jsResults, err := json.Marshal(dr)
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}

	err = ioutil.WriteFile(filename, jsResults, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON to %s: %v", filename, err)
	}

	return nil
}

// LoadResults reads the saved version data from disk as a
// JSON file representation of a DependencyResults object.
func LoadResults(filename string) (*DependencyResults, error) {
	js, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filename, err)
	}

	dr := &DependencyResults{}
	err = json.Unmarshal(js, &dr)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling from JSON: %v", err)
	}

	return dr, nil
}
