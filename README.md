# npm-spdx

npm-spdx is a Go program that queries the NPM API to gather declared license
information for dependencies from a `package.json` manifest. It generates an
[SPDX](https://spdx.org) document containing that information and the
corresponding subdependency relationships.

## Example

See the [examples](examples/) directory for a quick usage example.

## Usage

Compile with `go build`, then:

### Step 1: Obtain license data from NPM

You will need the `package.json` file for your NPM-based project, as well as the
corresponding `package-lock.json` file (to determine which specific versions of
which subdependencies were installed).

Then, retrieve the declared dependency license info by calling `npm-spdx retrieve`:

`./npm-spdx retrieve <PACKAGE.JSON> <PACKAGE-LOCK.JSON> <RESULTS.JSON>`

This will pull the results and save them to the file specified in
`<RESULTS.JSON>`, which will be used in the next steps.

### Step 2: Create SPDX document from results.json

Now, generate the SPDX document by calling `npm-spdx spdx`:

`./npm-spdx spdx <RESULTS.JSON> <OUTPUT.SPDX>`

This will read in the `results.json` file you obtained from Step 1, and process
it into an SPDX version 2.1 document that will be saved to the file specified in
`<OUTPUT.SPDX>`.

### (optional) Step 3: Create summary json file

You can also optionally process the results into a JSON file with dependencies
categorized by license expression. The resulting JSON file might be easier to
use for certain policy or automation processes. You can generate this by calling
`npm-spdx report`:

`./npm-spdx report <RESULTS.JSON> <SUMMARY.JSON>`

This will read in the `results.json` file you obtained from Step 1, and process
it into a JSON file that will be saved to the file specified in
`<SUMMARY.JSON>`.

### Using Docker

Alternatively, you can build and run npm-spdx using Docker.
Clone this repository and run `docker build -t npm-spdx .`
Then you can use npm-spdx by running:
`docker run --rm -v <PROJECT_PATH>:<CONTAINER_PATH> npm-spdx ...`

## License

npm-spdx is available under the [Apache License, version 2.0](LICENSE).

Copyright The Linux Foundation and npm-spdx contributors.

SPDX-License-Identifier: Apache-2.0
