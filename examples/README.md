SPDX-License-Identifier: Apache-2.0

Copyright The Linux Foundation and npm-spdx contributors.

# Usage example

The user starts with `package.json` and `package-lock.json` files that they
provide. Then, running from the npm-spdx root directory:

`./npm-spdx retrieve examples/package.json examples/package-lock.json examples/results.json`
* runs the **retrieve** command
* queries the npm API to retrieve license info
* generates the `results.json` file which is used in subsequent steps

`./npm-spdx spdx examples/results.json examples/sbom.spdx`
* runs the **spdx** command
* loads the `results.json` file
* generates an SPDX document (tag-value, version 2.1) at `sbom.spdx`

`./npm-spdx report examples/results.json examples/summary.json`
* runs the **report** command
* loads the `results.json` file
* generates a reformatted JSON file at `summary.json` with dependencies
  categorized by license
