// Package npm contains data structure for use in loading
// and parsing NPM manifests (e.g., package.json and
// package-lock.json files), and for retrieving data from
// the NPM API.
// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.
package npm

// RegistryVersion contains relevant data for the NPM API's
// response to a GET call for a particular version of an
// NPM package.
type RegistryVersion struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	License         string            `json:"license,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}

// RegistryPackage maps a package's version strings to the
// corresponding version-specific details.
type RegistryPackage map[string]*RegistryVersion

// RegistryResults maps a dependency's name to its
// RegistryPackage, which itself contains the corresponding
// version-specific details.
type RegistryResults map[string]RegistryPackage

// PackageManifest represents the data from a package.json
// file.
type PackageManifest struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	License         string            `json:"license,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}

// PackageLockDependency represents an entry within the
// "dependencies" object in a package-lock.json file.
type PackageLockDependency struct {
	Version   string            `json:"version"`
	Resolved  string            `json:"resolved"`
	Integrity string            `json:"integrity"`
	Dev       bool              `json:"dev,omitempty"`
	Requires  map[string]string `json:"requires,omitempty"`
}

// PackageLockManifest represents the data from a
// package-lock.json file.
type PackageLockManifest struct {
	Name         string                            `json:"name"`
	Version      string                            `json:"version"`
	Dependencies map[string]*PackageLockDependency `json:"dependencies,omitempty"`
}
