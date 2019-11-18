// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package main

import (
	"log"
	"os"
)

func printMainUsage() {
	log.Fatalf(`
Usage: %s <COMMAND> [arguments]

Commands:
	retrieve    - retrieve dependency info from NPM API and save to disk
	report      - load previously-retrieved dependency info and print summary details
	spdx        - load previously-retrieved dependency info and save as SPDX tag-value file

`, os.Args[0])
}

func checkUsage() {
	if len(os.Args) < 2 {
		printMainUsage()
	}

	command := os.Args[1]
	switch command {

	case "retrieve":
		if len(os.Args) != 5 {
			log.Fatalf(`
Usage: %s retrieve <PACKAGE.JSON> <PACKAGE-LOCK.JSON> <RESULTS.JSON>

PACKAGE.JSON:       path to package.json file for analysis
PACKAGE-LOCK.JSON:  path to package-lock.json file for analysis
RESULTS.JSON:       output path for results of API queries
`, os.Args[0])
		}

	case "report":
		if len(os.Args) != 4 {
			log.Fatalf(`
Usage: %s report <RESULTS.JSON> <SUMMARY.JSON>

RESULTS.JSON:       path to results from API queries (from prior 'retrieve' step)
SUMMARY.JSON:       output path for categorized JSON license results
`, os.Args[0])
		}

	case "spdx":
		if len(os.Args) != 4 {
			log.Fatalf(`
Usage: %s spdx <RESULTS.JSON> <OUTPUT.SPDX>

RESULTS.JSON:       path to results from API queries (from prior 'retrieve' step)
OUTPUT.SPDX:        output path for SPDX tag-value file
`, os.Args[0])
		}

	default:
		printMainUsage()
	}

}
