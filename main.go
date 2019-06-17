// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/swinslow/npm-spdx/pkg/npm"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <PACKAGE.JSON> <PACKAGE-LOCK.JSON>", os.Args[0])
	}

	pjsFilename := os.Args[1]
	pljsFilename := os.Args[2]

	js, err := ioutil.ReadFile(pjsFilename)
	if err != nil {
		log.Fatalf("error reading %s: %v", pjsFilename, err)
	}

	manifest, err := npm.ParseManifest(js)
	if err != nil {
		log.Fatalf("error parsing %s: %v", pjsFilename, err)
	}

	ljs, err := ioutil.ReadFile(pljsFilename)
	if err != nil {
		log.Fatalf("error reading %s: %v", pljsFilename, err)
	}

	lockManifest, err := npm.ParseLockManifest(ljs)
	if err != nil {
		log.Fatalf("error parsing %s: %v", pljsFilename, err)
	}

	fmt.Printf("==========\nmanifest:\n%#v\n==========\n", manifest)
	fmt.Printf("==========\nlockManifest:\n%#v\n==========\n", lockManifest)

	allResults, err := npm.GetAllVersionData(lockManifest.Dependencies, 500)
	if err != nil {
		log.Fatalf("error getting version data: %v", err)
	}

	jsResults, err := json.Marshal(allResults)
	fmt.Println("==========")
	fmt.Printf("%s\n", string(jsResults))
}
