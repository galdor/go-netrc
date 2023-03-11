# Introduction
The go-netrc library is a small Go package to load netrc files and search for
matching entries.

# Usage
See the example [netrc program](cmd/netrc) to see how to use the API.

# Limitations
We currently do not handle macros (introduced by the `macdef` token), quoted
token values or the default machine.

# Licensing
Go-netrc is open source software distributed under the
[ISC](https://opensource.org/licenses/ISC) license.
