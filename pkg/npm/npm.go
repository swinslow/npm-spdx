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
	License         interface{}       `json:"license,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}

// RegistryScopedPackage handles the multiple versions that
// are pulled when querying the NPM APIs for a scoped package.
type RegistryScopedPackage struct {
	Versions map[string]*RegistryVersion `json:"versions"`
}

// Dependency contains the processed, version-specific details
// about one dependency used (directly or indirectly) by the
// main package.
type Dependency struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	License         string            `json:"license,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
	IsDirectDep     bool              `json:"isDirectDep,omitempty"`
	IsDirectDevDep  bool              `json:"isDirectDevDep,omitempty"`
}

// DependencyResults maps a dependency's name to its
// Dependency, which itself contains the corresponding
// version-specific details.
type DependencyResults struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	License string                 `json:"license,omitempty"`
	Results map[string]*Dependency `json:"results"`
}

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
