// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package main

import (
	"log"
	"os"
)

func main() {
	checkUsage()

	command := os.Args[1]
	switch command {

	case "retrieve":
		pjsFilename := os.Args[2]
		pljsFilename := os.Args[3]
		jsOutput := os.Args[4]
		retrieve(pjsFilename, pljsFilename, jsOutput)

	case "report":
		jsResults := os.Args[2]
		jsReportOutput := os.Args[3]
		report(jsResults, jsReportOutput)

	case "spdx":
		jsResults := os.Args[2]
		spdxOutput := os.Args[3]
		spdx(jsResults, spdxOutput)

	default:
		log.Fatalf("No command specified")
	}
}
