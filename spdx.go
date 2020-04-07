// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package main

import (
	"log"
	"os"

	"github.com/spdx/tools-golang/tvsaver"
	"github.com/swinslow/npm-spdx/pkg/npm"
	"github.com/swinslow/npm-spdx/pkg/spdxpackages"
)

func spdx(jsResults, spdxOutput string) {
	// load results from JSON file
	dr, err := npm.LoadResults(jsResults)
	if err != nil {
		log.Fatalf("error loading from %s: %v", jsResults, err)
	}

	// create SPDX document from results
	doc, err := spdxpackages.BuildSPDXDocument(dr)
	if err != nil {
		log.Fatalf("error building SPDX document from %s: %v", jsResults, err)
	}

	// save SPDX document out to disk
	w, err := os.Create(spdxOutput)
	if err != nil {
		log.Fatalf("error creating file at %s for SPDX document: %v", spdxOutput, err)
	}
	defer w.Close()

	err = tvsaver.Save2_1(doc, w)
	if err != nil {
		log.Fatalf("error saving SPDX document to %s: %v", spdxOutput, err)
		return
	}
}
