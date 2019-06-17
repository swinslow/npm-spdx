// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package main

import (
	"io/ioutil"
	"log"

	"github.com/swinslow/npm-spdx/pkg/npm"
)

func retrieve(pjsFilename, pljsFilename, jsOutput string) {
	// js, err := ioutil.ReadFile(pjsFilename)
	// if err != nil {
	// 	log.Fatalf("error reading %s: %v", pjsFilename, err)
	// }

	// manifest, err := npm.ParseManifest(js)
	// if err != nil {
	// 	log.Fatalf("error parsing %s: %v", pjsFilename, err)
	// }

	ljs, err := ioutil.ReadFile(pljsFilename)
	if err != nil {
		log.Fatalf("error reading %s: %v", pljsFilename, err)
	}

	lockManifest, err := npm.ParseLockManifest(ljs)
	if err != nil {
		log.Fatalf("error parsing %s: %v", pljsFilename, err)
	}

	allResults, err := npm.GetAllVersionData(lockManifest.Dependencies, 500)
	if err != nil {
		log.Fatalf("error getting version data: %v", err)
	}

	err = npm.SaveVersionData(allResults, jsOutput)
	if err != nil {
		log.Fatalf("error saving to %s: %v", jsOutput, err)
	}
}
