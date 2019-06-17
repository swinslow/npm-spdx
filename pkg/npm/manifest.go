// SPDX-License-Identifier: Apache-2.0
// Copyright The Linux Foundation and npm-spdx contributors.

package npm

import "encoding/json"

// ParseManifest takes a byte slice (from reading a
// package.json file), parses it, and returns a corresponding
// PackageManifest object.
func ParseManifest(js []byte) (*PackageManifest, error) {
	manifest := PackageManifest{}
	err := json.Unmarshal(js, &manifest)
	return &manifest, err
}

// ParseLockManifest takes a byte slice (from reading a
// package-lock.json file), parses it, and returns a corresponding
// PackageLockManifest object.
func ParseLockManifest(js []byte) (*PackageLockManifest, error) {
	lockManifest := PackageLockManifest{}
	err := json.Unmarshal(js, &lockManifest)
	return &lockManifest, err
}
