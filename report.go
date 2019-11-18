// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/swinslow/npm-spdx/pkg/npm"
	"github.com/swinslow/npm-spdx/pkg/spdxlicenses"
)

type packageVersion struct {
	Pkg            string `json:"package"`
	Ver            string `json:"version"`
	IsDirectDep    bool   `json:"isDirectDep,omitempty"`
	IsDirectDevDep bool   `json:"isDirectDevDep,omitempty"`
}

type licEntry struct {
	LicID  string           `json:"id"`
	IsSPDX bool             `json:"valid"`
	Deps   []packageVersion `json:"dependencies"`
}

func report(jsResults string, jsReportOutput string) {
	// load valid license IDs
	llPath := "data/licenses.json"
	elPath := "data/exceptions.json"
	allLics, err := spdxlicenses.ParseJSONLicenses(llPath, elPath)
	if err != nil {
		log.Fatalf("error loading SPDX license IDs from %s, %s: %v", llPath, elPath, err)
	}

	// load results
	dr, err := npm.LoadResults(jsResults)
	if err != nil {
		log.Fatalf("error loading from %s: %v", jsResults, err)
	}

	// analyze
	lics := map[string]*licEntry{}
	for p, pData := range dr.Results {
		l := pData.License
		if l == "" {
			l = "NOASSERTION"
		}
		le, ok := lics[l]
		if !ok {
			// create new empty licEntry
			le = &licEntry{}
			le.LicID = l

			le.IsSPDX = spdxlicenses.IsValidExpression(l, allLics)

			le.Deps = []packageVersion{}
			lics[l] = le
		}

		// add this version
		pv := packageVersion{
			Pkg:            p,
			Ver:            pData.Version,
			IsDirectDep:    pData.IsDirectDep,
			IsDirectDevDep: pData.IsDirectDevDep,
		}
		le.Deps = append(le.Deps, pv)
	}

	// create JSON output
	js, err := json.Marshal(&lics)
	if err != nil {
		log.Fatalf("error marshalling summarized licenses to JSON: %v", err)
	}

	// and write to disk
	err = ioutil.WriteFile(jsReportOutput, js, 0644)
	if err != nil {
		log.Fatalf("error writing JSON to %s: %v", jsReportOutput, err)
	}
}
