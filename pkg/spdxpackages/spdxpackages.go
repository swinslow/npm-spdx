// Package spdxpackages contains functions to work with
// SPDX's tools-golang to generate an SPDX document.
// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.
package spdxpackages

import (
	"fmt"
	"log"
	"time"

	"github.com/spdx/tools-golang/v0/spdx"
	"github.com/swinslow/npm-spdx/pkg/npm"
	"github.com/swinslow/npm-spdx/pkg/spdxlicenses"
)

// BuildSPDXDocument takes the processed license and dependency data
// from a previously-generated results.json file, and returns an SPDX
// document based on them, together with the relevant relationship details.
func BuildSPDXDocument(dr *npm.DependencyResults) (*spdx.Document2_1, error) {
	// load valid license IDs
	llPath := "data/licenses.json"
	elPath := "data/exceptions.json"
	allLics, err := spdxlicenses.ParseJSONLicenses(llPath, elPath)
	if err != nil {
		log.Fatalf("error loading SPDX license IDs from %s, %s: %v", llPath, elPath, err)
	}

	// build creation info section
	// FIXME namespace should be unique, see SPDX 2.1 spec section 2.5
	namespace := fmt.Sprintf("https://spdx.org/spdxdocs/%s-%s", dr.Name, dr.Version)
	ci := buildCreationInfoSection(dr.Name, namespace)

	// build collection of package sections, looking to results for
	// what we actually installed; also build relationship sections
	// at the same time
	pkgs := []*spdx.Package2_1{}
	rlns := []*spdx.Relationship2_1{}
	ols := []*spdx.OtherLicense2_1{}

	// also track which converted "other licenses" we have created
	convertedLics := map[string]bool{}

	// build entry for main package
	lic := dr.License
	if lic == "" {
		lic = "NOASSERTION"
	}
	// make sure the main package license is valid too!
	if !spdxlicenses.IsValidExpression(lic, allLics) {
		lic = spdxlicenses.ConvertToLicenseRef(lic)
		convertedLics[lic] = true
		ol := buildOtherLicense(lic, dr.License)
		ols = append(ols, ol)
	}

	mainPkg := buildPackageSection(dr.Name, dr.Version, "NOASSERTION", lic)
	pkgs = append(pkgs, mainPkg)

	// also add DESCRIBES relationship for main package
	mainRln := &spdx.Relationship2_1{
		RefA:         "SPDXRef-DOCUMENT",
		RefB:         mainPkg.PackageSPDXIdentifier,
		Relationship: "DESCRIBES",
	}
	rlns = append(rlns, mainRln)

	for _, rp := range dr.Results {
		// convert license if needed
		pkgLic := rp.License
		if !spdxlicenses.IsValidExpression(pkgLic, allLics) {
			// first see if we've already got an OtherLicense entry
			// for this one
			foundMatch := false
			for _, ol := range ols {
				if ol.ExtractedText == pkgLic {
					pkgLic = ol.LicenseIdentifier
					foundMatch = true
				}
			}

			if !foundMatch {
				// convert it to LicenseRef format
				pkgLic = spdxlicenses.ConvertToLicenseRef(pkgLic)

				// and now make sure that we find a free ID
				if _, ok := convertedLics[pkgLic]; ok {
					var newConverted string
					for i := 1; ; i++ {
						newConverted = fmt.Sprintf("%s-%d", pkgLic, i)
						if _, ok := convertedLics[newConverted]; !ok {
							break
						}
					}
					pkgLic = newConverted
				}

				// and now add it to the list and build other license section
				convertedLics[pkgLic] = true
				ol := buildOtherLicense(pkgLic, rp.License)
				ols = append(ols, ol)
			}

		}

		// FIXME for now, don't fill in PackageDownloadLocation
		pkg := buildPackageSection(rp.Name, rp.Version, "NOASSERTION", pkgLic)
		pkgs = append(pkgs, pkg)

		// build relationships
		for depName := range rp.Dependencies {
			// get the dependency version that was actually installed
			if depRp, ok := dr.Results[depName]; ok {
				rln := buildDependencyRelationship(rp.Name, rp.Version, depName, depRp.Version)
				rlns = append(rlns, rln)
			}
		}

		// also add relationship if it's a direct dependency (main or dev)
		// of the main package
		if rp.IsDirectDep {
			rln := buildDependencyRelationship(dr.Name, dr.Version, rp.Name, rp.Version)
			rlns = append(rlns, rln)
		}
		if rp.IsDirectDevDep {
			rln := buildDevDependencyRelationship(dr.Name, dr.Version, rp.Name, rp.Version)
			rlns = append(rlns, rln)
		}
	}

	doc := &spdx.Document2_1{
		CreationInfo:  ci,
		Packages:      pkgs,
		Relationships: rlns,
		OtherLicenses: ols,
	}

	return doc, nil
}

func getSPDXID(pkg string, ver string) string {
	return fmt.Sprintf("SPDXRef-%s-%s", pkg, ver)
}

func getNpmURL(pkg string, ver string) string {
	return fmt.Sprintf("https://www.npmjs.com/package/%s/v/%s", pkg, ver)
}

func buildCreationInfoSection(mainPackageName string, namespace string) *spdx.CreationInfo2_1 {
	// get current time in UTC
	location, _ := time.LoadLocation("UTC")
	locationTime := time.Now().In(location)
	created := locationTime.Format("2006-01-02T15:04:05Z")

	ci := &spdx.CreationInfo2_1{
		SPDXVersion:        "SPDX-2.1",
		DataLicense:        "CC0-1.0",
		SPDXIdentifier:     "SPDXRef-DOCUMENT",
		DocumentName:       mainPackageName,
		DocumentNamespace:  namespace,
		LicenseListVersion: "3.5",
		CreatorTools:       []string{"github.com/swinslow/npm-spdx"},
		Created:            created,
	}

	return ci
}

func buildPackageSection(pkgName string, pkgVer string, url string, licDeclared string) *spdx.Package2_1 {
	if licDeclared == "" {
		licDeclared = "NOASSERTION"
	}
	pkg := &spdx.Package2_1{
		PackageName:             pkgName,
		PackageSPDXIdentifier:   getSPDXID(pkgName, pkgVer),
		PackageVersion:          pkgVer,
		PackageDownloadLocation: url,
		FilesAnalyzed:           false,
		PackageHomePage:         getNpmURL(pkgName, pkgVer),
		PackageLicenseConcluded: "NOASSERTION",
		PackageLicenseDeclared:  licDeclared,
		PackageCopyrightText:    "NOASSERTION",
		PackageExternalReferences: []*spdx.PackageExternalReference2_1{
			&spdx.PackageExternalReference2_1{
				Category: "PACKAGE-MANAGER",
				RefType:  "npm",
				Locator:  fmt.Sprintf("%s@%s", pkgName, pkgVer),
			},
		},
	}

	return pkg
}

func buildDependencyRelationship(pkgName, pkgVer, depName, depVer string) *spdx.Relationship2_1 {
	pkgID := getSPDXID(pkgName, pkgVer)
	depID := getSPDXID(depName, depVer)
	rln := &spdx.Relationship2_1{
		RefA:         depID,
		RefB:         pkgID,
		Relationship: "PREREQUISITE_FOR",
	}

	return rln
}

func buildDevDependencyRelationship(pkgName, pkgVer, depName, depVer string) *spdx.Relationship2_1 {
	pkgID := getSPDXID(pkgName, pkgVer)
	depID := getSPDXID(depName, depVer)
	rln := &spdx.Relationship2_1{
		RefA:         depID,
		RefB:         pkgID,
		Relationship: "BUILD_TOOL_OF",
	}

	return rln
}

func buildOtherLicense(converted, orig string) *spdx.OtherLicense2_1 {
	cmt := fmt.Sprintf("Represents the license expression '%s' which is not on the SPDX License List", orig)

	ol := &spdx.OtherLicense2_1{
		LicenseIdentifier: converted,
		ExtractedText:     orig,
		LicenseName:       orig,
		LicenseComment:    cmt,
	}

	return ol
}
