<div align="center">

# GURL

### A simple `curl` alternative that works almost everywhere.

<img width="100%" alt="GURL" src="https://github.com/matveynator/gurl/blob/master/gurl.png?raw=true" />
<br>

![HTTP](https://img.shields.io/badge/HTTP-yes-green)
![HTTPS](https://img.shields.io/badge/HTTPS-yes-green)
![TLS](https://img.shields.io/badge/TLS-built--in-green)
![Dependencies](https://img.shields.io/badge/external_dependencies-none-green)
![Cross Platform](https://img.shields.io/badge/cross--platform-yes-green)

[![Go Report Card](https://goreportcard.com/badge/github.com/matveynator/gurl)](https://goreportcard.com/report/github.com/matveynator/gurl)

</div>

---

## What is GURL?

**GURL** is a small, standalone command-line HTTP/HTTPS client written in Go.

It is designed for situations where you just need to download a file, call an HTTP endpoint, send POST data, or make a simple request — without installing `curl`, OpenSSL, or a collection of shared libraries.

TLS support is built into the binary through Go's standard library.

This makes GURL especially useful for:

* minimal Linux installations
* recovery environments
* old or unusual systems
* containers
* embedded systems
* automation scripts
* servers where `curl` or OpenSSL is unavailable
* copying a single binary between machines

Just download **one executable** and run it.

```bash
gurl https://example.com
```

If no protocol is specified, GURL uses `http://`:

```bash
gurl example.com
```

---

# Downloads / Скачать


<details>
<summary>
  <img width="42" alt="linux" src="https://github.com/user-attachments/assets/bf3141b6-4c93-4fd6-b2d1-421b79876dcb" />
  <b><big>Linux</big></b>
  <sub>amd64 / arm64 / 386 / ARM / MIPS / RISC-V / PPC / s390x / LoongArch</sub>
</summary>

<br>

### amd64 / x86_64

[download](http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl)

```bash
sudo curl -L -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; sudo chmod +x /usr/local/bin/gurl
```

Or, if you already have an older GURL:

```bash
sudo gurl -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; sudo chmod +x /usr/local/bin/gurl
```

### arm64 / aarch64

[download](http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl)

```bash
sudo curl -L -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl; sudo chmod +x /usr/local/bin/gurl
```

### 386 / x86

[download](http://files.zabiyaka.net/gurl/latest/linux/386/gurl)

### ARM

[download](http://files.zabiyaka.net/gurl/latest/linux/arm/gurl)

### LoongArch 64

[download](http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl)

### MIPS

[download](http://files.zabiyaka.net/gurl/latest/linux/mips/gurl)

### MIPS little-endian

[download](http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl)

### MIPS64

[download](http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl)

### MIPS64 little-endian

[download](http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl)

### PowerPC 64

[download](http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl)

### PowerPC 64 little-endian

[download](http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl)

### RISC-V 64

[download](http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl)

### IBM s390x

[download](http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl)

</details>

---

<details>
<summary>
  <img width="36" alt="mac" src="https://github.com/user-attachments/assets/946102b8-f043-494d-809a-a589e536ee9a" />
  <b><big>macOS</big></b>
  <sub>Intel / Apple Silicon</sub>
</summary>

<br>

### Intel / amd64

[download](http://files.zabiyaka.net/gurl/latest/mac/amd64/gurl)

```bash
sudo mkdir -p /usr/local/bin; sudo curl -L -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/mac/amd64/gurl; sudo chmod +x /usr/local/bin/gurl
```

### Apple Silicon / arm64

[download](http://files.zabiyaka.net/gurl/latest/mac/arm64/gurl)

```bash
sudo mkdir -p /usr/local/bin; sudo curl -L -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/mac/arm64/gurl; sudo chmod +x /usr/local/bin/gurl
```

</details>

---

<details>
<summary>
  <img width="42" alt="windows" src="https://github.com/user-attachments/assets/f6044001-95b0-4500-a4f6-1c3b08eb65fb" />
  <b><big>Windows</big></b>
  <sub>amd64 / arm64 / 386 / ARM</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/windows/amd64/gurl.exe)

```powershell
$p="$env:ProgramFiles\gurl\gurl.exe"; New-Item -ItemType Directory -Force -Path (Split-Path $p) | Out-Null; Invoke-WebRequest -Uri "http://files.zabiyaka.net/gurl/latest/windows/amd64/gurl.exe" -OutFile $p; & $p -V
```

### arm64

[download](http://files.zabiyaka.net/gurl/latest/windows/arm64/gurl.exe)

```powershell
$p="$env:ProgramFiles\gurl\gurl.exe"; New-Item -ItemType Directory -Force -Path (Split-Path $p) | Out-Null; Invoke-WebRequest -Uri "http://files.zabiyaka.net/gurl/latest/windows/arm64/gurl.exe" -OutFile $p; & $p -V
```

### 386

[download](http://files.zabiyaka.net/gurl/latest/windows/386/gurl.exe)

### ARM

[download](http://files.zabiyaka.net/gurl/latest/windows/arm/gurl.exe)

</details>

---

<details>
<summary>
  <img width="42" alt="freebsd" src="https://github.com/user-attachments/assets/d35baaac-d296-41b1-a281-55dc761328e9" />
  <b><big>FreeBSD</big></b>
  <sub>amd64 / arm64 / 386 / ARM / RISC-V</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/freebsd/amd64/gurl)

```bash
sudo fetch -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/freebsd/amd64/gurl; sudo chmod +x /usr/local/bin/gurl
```

### arm64

[download](http://files.zabiyaka.net/gurl/latest/freebsd/arm64/gurl)

```bash
sudo fetch -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/freebsd/arm64/gurl; sudo chmod +x /usr/local/bin/gurl
```

### 386

[download](http://files.zabiyaka.net/gurl/latest/freebsd/386/gurl)

### ARM

[download](http://files.zabiyaka.net/gurl/latest/freebsd/arm/gurl)

### RISC-V 64

[download](http://files.zabiyaka.net/gurl/latest/freebsd/riscv64/gurl)

</details>

---

<details>
<summary>
  <img width="42" alt="openbsd" src="https://github.com/user-attachments/assets/11633d7e-5744-46da-ad2f-6e49c69e51de" />
  <b><big>OpenBSD</big></b>
  <sub>amd64 / arm64 / 386 / ARM / PPC64 / RISC-V</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/amd64/gurl)

```bash
sudo ftp -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/openbsd/amd64/gurl; sudo chmod +x /usr/local/bin/gurl
```

### arm64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/arm64/gurl)

```bash
sudo ftp -o /usr/local/bin/gurl http://files.zabiyaka.net/gurl/latest/openbsd/arm64/gurl; sudo chmod +x /usr/local/bin/gurl
```

### 386

[download](http://files.zabiyaka.net/gurl/latest/openbsd/386/gurl)

### ARM

[download](http://files.zabiyaka.net/gurl/latest/openbsd/arm/gurl)

### PowerPC 64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/ppc64/gurl)

### RISC-V 64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/riscv64/gurl)

</details>

---

<details>
<summary>
  <b><big>Other platforms</big></b>
  <sub>NetBSD / Android / Solaris / Plan 9 / Illumos / DragonFlyBSD / AIX / WebAssembly</sub>
</summary>

<br>

### Android

**arm64** · [download](http://files.zabiyaka.net/gurl/latest/android/arm64/gurl)

### NetBSD

**amd64** · [download](http://files.zabiyaka.net/gurl/latest/netbsd/amd64/gurl)  
**386** · [download](http://files.zabiyaka.net/gurl/latest/netbsd/386/gurl)  
**arm** · [download](http://files.zabiyaka.net/gurl/latest/netbsd/arm/gurl)  
**arm64** · [download](http://files.zabiyaka.net/gurl/latest/netbsd/arm64/gurl)

### Solaris

**amd64** · [download](http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl)

### Plan 9

**amd64** · [download](http://files.zabiyaka.net/gurl/latest/plan9/amd64/gurl)  
**386** · [download](http://files.zabiyaka.net/gurl/latest/plan9/386/gurl)  
**arm** · [download](http://files.zabiyaka.net/gurl/latest/plan9/arm/gurl)

### Illumos

**amd64** · [download](http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl)

### DragonFlyBSD

**amd64** · [download](http://files.zabiyaka.net/gurl/latest/dragonfly/amd64/gurl)

### AIX

**ppc64** · [download](http://files.zabiyaka.net/gurl/latest/aix/ppc64/gurl)

### WebAssembly

**js/wasm** · [download](http://files.zabiyaka.net/gurl/latest/js/wasm/gurl)

### WASI

**wasip1/wasm** · [download](http://files.zabiyaka.net/gurl/latest/wasip1/wasm/gurl)

</details>

---

**All GitHub releases:**

https://github.com/matveynator/gurl/releases

---

## Why are the binary download links HTTP?

This is intentional.

One of GURL's use cases is bootstrapping old or minimal systems where HTTPS tools, CA certificates, OpenSSL, or even `curl` may not yet be available.

The initial GURL binary can therefore be downloaded over plain HTTP.

Once installed, GURL itself supports both HTTP and HTTPS:

```bash
gurl https://example.com
```

---

# Usage / Как пользоваться

Basic syntax:

```text
gurl [options] <URL>
```

Simple GET request:

```bash
gurl https://example.com
```

A URL without a scheme automatically uses HTTP:

```bash
gurl example.com
```

Equivalent to:

```bash
gurl http://example.com
```

---

## Examples / Примеры

### Download a file

```bash
gurl -o file.zip https://example.com/file.zip
```

or:

```bash
gurl --output file.zip https://example.com/file.zip
```

---

### Use GURL in shell scripts

```bash
gurl https://example.com/script.sh | bash
```

For example:

```bash
gurl https://raw.githubusercontent.com/matveynator/sysadminscripts/main/label | bash
```

---

### POST data

```bash
gurl -d "key=value&key2=value2" https://example.com
```

or:

```bash
gurl --data "key=value&key2=value2" https://example.com
```

When `-d` / `--data` is used, GURL automatically performs a POST request.

---

### Send JSON

```bash
gurl -H "Content-Type: application/json" -d '{"key":"value"}' https://example.com/api
```

---

### Multipart form

Text fields and files can be supplied in one `-F` argument separated by `&`:

```bash
gurl -F "name=test&file=@/tmp/file.txt" https://example.com/upload
```

---

### Custom HTTP method

```bash
gurl -X DELETE https://example.com/api/item
```

```bash
gurl -X PUT https://example.com/api/item
```

---

### Custom header

```bash
gurl -H "Authorization: Bearer TOKEN" https://example.com/api
```

or:

```bash
gurl --header "Authorization: Bearer TOKEN" https://example.com/api
```

---

### Cookies

```bash
gurl -b "session_id=abc123" https://example.com
```

or:

```bash
gurl --cookie "session_id=abc123" https://example.com
```

---

### HEAD request

```bash
gurl -I https://example.com
```

or:

```bash
gurl --head https://example.com
```

---

### Custom User-Agent

```bash
gurl -A "MyClient/1.0" https://example.com
```

or:

```bash
gurl --useragent "MyClient/1.0" https://example.com
```

---

### Timeout

```bash
gurl -m 10s https://example.com
```

or:

```bash
gurl --timeout 10s https://example.com
```

The default timeout is **30 seconds**.

---

### Ignore TLS certificate verification

Useful for self-signed certificates:

```bash
gurl -k https://192.168.1.1
```

or:

```bash
gurl --unsafe https://192.168.1.1
```

> **Warning:** disabling certificate verification reduces connection security.

---

### Fail on HTTP errors

```bash
gurl --fail https://example.com/not-found
```

For HTTP status `400` and above, GURL exits with error code `22`.

This is useful in shell scripts:

```bash
gurl --fail https://example.com/file || echo "Download failed"
```

---

### Redirects

GURL follows HTTP redirects by default.

```bash
gurl https://example.com
```

The compatible `-L` / `--location` option is also available.

---

### Show version

```bash
gurl -V
```

or:

```bash
gurl --version
```

---

## Flags / Флаги

```text
-V, --version       show GURL version

-m, --timeout       request timeout
                    default: 30s

-A, --useragent     custom User-Agent
                    default: GURL

-k, --unsafe        disable TLS certificate verification

-d, --data          send POST data

-F                   multipart form:
                     key=value
                     key=@file
                     multiple fields separated with "&"

-b, --cookie        send Cookie header

-I, --head          perform HEAD request

-H, --header        send a custom HTTP header

-o, --output        save response body to a file

-L, --location      follow redirects
                    enabled by default

--fail               return an error for HTTP status >= 400

-X                    custom HTTP method
                     default: GET
```

---

# Common Problems Solved by GURL

These are typical situations where GURL is useful:

* I need to download one file but `curl` is not installed.
* I need HTTPS but OpenSSL is unavailable.
* I need a standalone HTTP client with no external runtime dependencies.
* I need to bootstrap a minimal server.
* I need an HTTP client for an old or unusual operating system.
* I need the same command-line utility on Linux, BSD, macOS, and Windows.
* I need a small binary I can copy to another machine.
* I need GET and POST requests without installing a large package.
* I need to download scripts during server provisioning.
* I need to send JSON from a shell script.
* I need to upload a file using multipart/form-data.
* I need to call an HTTP API from a rescue environment.
* I need HTTP downloads before HTTPS tooling is configured.
* I need a curl-like utility compiled with TLS support already included.
* I need a simple HTTP client for automation.

In short:

> “I just need `curl`, but it isn't there.”

That's what **GURL** is for.

One binary.  
No external libraries.  
HTTP and HTTPS.  
Runs on a lot of platforms.

---

## Build from source

Clone the repository:

```bash
git clone https://github.com/matveynator/gurl.git
cd gurl
```

Build:

```bash
go build -o gurl gurl.go
```

Run:

```bash
./gurl https://example.com
```

---

## Source

https://github.com/matveynator/gurl

Contributions and fixes are welcome.
