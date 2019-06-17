// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package main

import (
	"fmt"
	"log"

	"github.com/swinslow/npm-spdx/pkg/spdxlicenses"
)

func spdx(jsResults, spdxOutput string) {
	llPath := "data/licenses.json"
	elPath := "data/exceptions.json"
	allLics, err := spdxlicenses.ParseJSONLicenses(llPath, elPath)
	if err != nil {
		log.Fatalf("error loading SPDX license IDs from %s, %s: %v", llPath, elPath, err)
	}

	for l := range allLics {
		fmt.Printf("%s\n", l)
	}

	// allResults, err := npm.LoadVersionData(allResults, jsResults)
	// if err != nil {
	// 	log.Fatalf("error loading from %s: %v", jsOutput, err)
	// }

	// js, err := ioutil.ReadFile(pjsFilename)
	// if err != nil {
	// 	log.Fatalf("error reading %s: %v", pjsFilename, err)
	// }

	// manifest, err := npm.ParseManifest(js)
	// if err != nil {
	// 	log.Fatalf("error parsing %s: %v", pjsFilename, err)
	// }

	// ljs, err := ioutil.ReadFile(pljsFilename)
	// if err != nil {
	// 	log.Fatalf("error reading %s: %v", pljsFilename, err)
	// }

	// lockManifest, err := npm.ParseLockManifest(ljs)
	// if err != nil {
	// 	log.Fatalf("error parsing %s: %v", pljsFilename, err)
	// }

	// allResults, err := npm.GetAllVersionData(lockManifest.Dependencies, 500)
	// if err != nil {
	// 	log.Fatalf("error getting version data: %v", err)
	// }

}
