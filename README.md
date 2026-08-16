# GURL

[![Go Reference](https://pkg.go.dev/badge/github.com/matveynator/gurl.svg)](https://pkg.go.dev/github.com/matveynator/gurl)
[![Go Report Card](https://goreportcard.com/badge/github.com/matveynator/gurl)](https://goreportcard.com/report/github.com/matveynator/gurl)
[![Coverage: 90.6%](https://img.shields.io/badge/coverage-90.6%25-brightgreen)](COVERAGE.md)
[![Security](https://github.com/matveynator/gurl/actions/workflows/security.yml/badge.svg)](https://github.com/matveynator/gurl/actions/workflows/security.yml)

GURL is a curl-like command-line HTTP client written in Go. Release binaries
are built with `CGO_ENABLED=0`, so they do not need a runtime OpenSSL, libcurl,
or other shared-library installation. This makes GURL useful on minimal and
older operating-system installations whose bundled curl can no longer
negotiate current TLS versions.

GURL implements a focused subset of curl's interface. It is not a drop-in
replacement for every curl protocol or option.

## Installation

Download the binary for your operating system and architecture from the
[latest GitHub release](https://github.com/matveynator/gurl/releases/latest).
For example, on Linux amd64:

```sh
curl -LO https://github.com/matveynator/gurl/releases/latest/download/gurl_linux_amd64
chmod +x gurl_linux_amd64
sudo install gurl_linux_amd64 /usr/local/bin/gurl
```

The release includes SHA-256 checksums. To build from source with a supported Go
toolchain:

```sh
go install github.com/matveynator/gurl@latest
```

## Usage

Make a request and write the response body to standard output:

```sh
gurl https://example.com
```

Save a response:

```sh
gurl --output response.json https://example.com/api
```

Send JSON and repeat headers when needed:

```sh
gurl --request POST \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer TOKEN' \
  --data '{"enabled":true}' \
  https://example.com/api
```

Upload multipart fields and files:

```sh
gurl --form 'name=report' --form 'document=@report.pdf' https://example.com/upload
```

Inspect headers, follow redirects, or fail on HTTP error responses:

```sh
gurl --head https://example.com
gurl --location https://example.com/redirect
gurl --fail https://example.com/missing
```

Run `gurl -help` for the complete option list. Long options and their common curl
short forms are supported, including `-A`, `-b`, `-d`, `-F`, `-H`, `-I`, `-k`,
`-L`, `-m`, `-o`, `-s`, `-v`, `-V`, and `-X`.

## TLS and old systems

The static release binaries use Go's TLS implementation instead of dynamically
loading OpenSSL. They still require a kernel supported by the Go version used
for the release and a usable system certificate store. Systems without current
CA certificates can provide an updated trust store through their normal OS
configuration. The `--unsafe` (`-k`) option disables certificate verification
and should only be used for controlled diagnostics.

## Supported release targets

Release binaries are produced for `amd64` and `arm64` on Linux, macOS, Windows,
FreeBSD, and OpenBSD. The standard Go toolchain may support additional targets
when building from source.

## Development

Run the local quality checks with:

```sh
go test -race -cover ./...
go vet ./...
```

Build every supported release target with:

```sh
go run ./scripts/crosscompile -output dist -version dev
```

Security reports are handled according to [SECURITY.md](SECURITY.md). GURL is
distributed under the [BSD 3-Clause License](LICENSE).
