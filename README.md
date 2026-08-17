<div align="center">
# GURL
### When CURL says your SSL library is too old — use GURL.
One file. Zero dependencies.
<img width="100%" alt="GURL" src="https://github.com/matveynator/gurl/releases/download/v64/gurl.png" />
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

If no protocol is specified, GURL uses http://:

gurl example.com

⸻

Downloads / Скачать

Choose your platform and architecture below.

Every target has a direct download link and a ready-to-copy bootstrap command.

The command downloads GURL into the current directory first, so it does not require root, sudo, or a writable /usr/local/bin.

The binary URLs intentionally use HTTP. This is the bootstrap path for old or minimal systems where HTTPS tools, CA certificates, OpenSSL, or the installed TLS stack may be unavailable or obsolete. After downloading GURL, use GURL itself for HTTPS.

Note: a CPU architecture does not guarantee that a particular downloader is installed. Linux and Android therefore try several common downloaders automatically. AIX has no universal base-system HTTP downloader, so its block uses wget or curl when available. WebAssembly and WASI binaries are downloaded from the host system.

<details>
<summary>
  <img width="42" alt="Linux" src="https://github.com/user-attachments/assets/bf3141b6-4c93-4fd6-b2d1-421b79876dcb" />
  <b><big>Linux</big></b>
  <sub>amd64 / arm64 / 386 / ARM / LoongArch / MIPS / PPC64 / RISC-V / s390x</sub>
</summary>
<br>

amd64 / x86_64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

arm64 / aarch64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

386 / x86

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/386/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/386/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/386/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

ARM

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/arm/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/arm/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/arm/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

LoongArch 64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

MIPS

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mips/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mips/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/mips/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

MIPS little-endian

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

MIPS64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

MIPS64 little-endian

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

PowerPC 64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

PowerPC 64 little-endian

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

RISC-V 64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V

IBM s390x

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <img width="36" alt="macOS" src="https://github.com/user-attachments/assets/946102b8-f043-494d-809a-a589e536ee9a" />
  <b><big>macOS</big></b>
  <sub>Intel / Apple Silicon</sub>
</summary>
<br>

Intel / amd64

download⁠￼

/usr/bin/curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/mac/amd64/gurl && chmod +x gurl && ./gurl -V

Apple Silicon / arm64

download⁠￼

/usr/bin/curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/mac/arm64/gurl && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <img width="42" alt="Windows" src="https://github.com/user-attachments/assets/f6044001-95b0-4500-a4f6-1c3b08eb65fb" />
  <b><big>Windows</big></b>
  <sub>amd64 / arm64 / 386 / ARM</sub>
</summary>
<br>

The commands use System.Net.WebClient instead of Invoke-WebRequest, so they also work with older Windows PowerShell versions.

amd64 / x86_64

download⁠￼

(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/amd64/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V

arm64

download⁠￼

(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/arm64/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V

386 / x86

download⁠￼

(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/386/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V

ARM

download⁠￼

(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/arm/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V
</details>

⸻

<details>
<summary>
  <img width="42" alt="FreeBSD" src="https://github.com/user-attachments/assets/d35baaac-d296-41b1-a281-55dc761328e9" />
  <b><big>FreeBSD</big></b>
  <sub>amd64 / arm64 / 386 / ARM / RISC-V</sub>
</summary>
<br>

amd64

download⁠￼

fetch -o gurl http://files.zabiyaka.net/gurl/latest/freebsd/amd64/gurl && chmod +x gurl && ./gurl -V

arm64

download⁠￼

fetch -o gurl http://files.zabiyaka.net/gurl/latest/freebsd/arm64/gurl && chmod +x gurl && ./gurl -V

386

download⁠￼

fetch -o gurl http://files.zabiyaka.net/gurl/latest/freebsd/386/gurl && chmod +x gurl && ./gurl -V

ARM

download⁠￼

fetch -o gurl http://files.zabiyaka.net/gurl/latest/freebsd/arm/gurl && chmod +x gurl && ./gurl -V

RISC-V 64

download⁠￼

fetch -o gurl http://files.zabiyaka.net/gurl/latest/freebsd/riscv64/gurl && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <img width="42" alt="OpenBSD" src="https://github.com/user-attachments/assets/11633d7e-5744-46da-ad2f-6e49c69e51de" />
  <b><big>OpenBSD</big></b>
  <sub>amd64 / arm64 / 386 / ARM / PPC64 / RISC-V</sub>
</summary>
<br>

amd64

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/openbsd/amd64/gurl && chmod +x gurl && ./gurl -V

arm64

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/openbsd/arm64/gurl && chmod +x gurl && ./gurl -V

386

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/openbsd/386/gurl && chmod +x gurl && ./gurl -V

ARM

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/openbsd/arm/gurl && chmod +x gurl && ./gurl -V

PowerPC 64

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/openbsd/ppc64/gurl && chmod +x gurl && ./gurl -V

RISC-V 64

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/openbsd/riscv64/gurl && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>Android</big></b>
  <sub>arm64</sub>
</summary>
<br>

arm64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; elif command -v toybox >/dev/null 2>&1; then toybox wget -O gurl http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; else echo 'No HTTP downloader found (wget/Toybox/BusyBox/curl).' >&2; exit 1; fi && chmod 755 gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>NetBSD</big></b>
  <sub>amd64 / 386 / ARM / arm64</sub>
</summary>
<br>

amd64

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/netbsd/amd64/gurl && chmod +x gurl && ./gurl -V

386

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/netbsd/386/gurl && chmod +x gurl && ./gurl -V

ARM

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/netbsd/arm/gurl && chmod +x gurl && ./gurl -V

arm64

download⁠￼

ftp -o gurl http://files.zabiyaka.net/gurl/latest/netbsd/arm64/gurl && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>Solaris</big></b>
  <sub>amd64</sub>
</summary>
<br>

amd64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl; elif [ -x /usr/sfw/bin/wget ]; then /usr/sfw/bin/wget -O gurl http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl; else echo 'No HTTP downloader found. Install/copy wget or curl, or use the direct download link above.' >&2; exit 1; fi && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>Plan 9</big></b>
  <sub>amd64 / 386 / ARM</sub>
</summary>
<br>

Uses Plan 9 hget.

amd64

download⁠￼

hget http://files.zabiyaka.net/gurl/latest/plan9/amd64/gurl >gurl
./gurl -V

386

download⁠￼

hget http://files.zabiyaka.net/gurl/latest/plan9/386/gurl >gurl
./gurl -V

ARM

download⁠￼

hget http://files.zabiyaka.net/gurl/latest/plan9/arm/gurl >gurl
./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>Illumos</big></b>
  <sub>amd64</sub>
</summary>
<br>

amd64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl; elif [ -x /usr/sfw/bin/wget ]; then /usr/sfw/bin/wget -O gurl http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl; else echo 'No HTTP downloader found. Install/copy wget or curl, or use the direct download link above.' >&2; exit 1; fi && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>DragonFly BSD</big></b>
  <sub>amd64</sub>
</summary>
<br>

amd64

download⁠￼

fetch -o gurl http://files.zabiyaka.net/gurl/latest/dragonfly/amd64/gurl && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>AIX</big></b>
  <sub>ppc64</sub>
</summary>
<br>

AIX base ftp is an FTP client and does not provide a universal HTTP bootstrap path. The command below uses wget or curl if one is already installed.

PowerPC 64

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl http://files.zabiyaka.net/gurl/latest/aix/ppc64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl http://files.zabiyaka.net/gurl/latest/aix/ppc64/gurl; else echo 'Base AIX ftp does not fetch HTTP URLs. Copy/install wget or curl, or use the direct download link above.' >&2; exit 1; fi && chmod +x gurl && ./gurl -V
</details>

⸻

<details>
<summary>
  <b><big>WebAssembly</big></b>
  <sub>js / wasm</sub>
</summary>
<br>

Run this command on the host system where the .wasm file will be stored.

js / wasm

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl.wasm http://files.zabiyaka.net/gurl/latest/js/wasm/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl.wasm http://files.zabiyaka.net/gurl/latest/js/wasm/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl.wasm http://files.zabiyaka.net/gurl/latest/js/wasm/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi
</details>

⸻

<details>
<summary>
  <b><big>WASI</big></b>
  <sub>wasip1 / wasm</sub>
</summary>
<br>

Run this command on the host system where the .wasm file will be stored.

wasip1 / wasm

download⁠￼

if command -v wget >/dev/null 2>&1; then wget -O gurl.wasm http://files.zabiyaka.net/gurl/latest/wasip1/wasm/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O gurl.wasm http://files.zabiyaka.net/gurl/latest/wasip1/wasm/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o gurl.wasm http://files.zabiyaka.net/gurl/latest/wasip1/wasm/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi
</details>

⸻

All GitHub releases:

https://github.com/matveynator/gurl/releases

⸻

Why are the binary download links HTTP?

This is intentional.

One of GURL’s use cases is bootstrapping old or minimal systems where HTTPS tools, CA certificates, OpenSSL, or even curl may not yet be available.

The initial GURL binary can therefore be downloaded over plain HTTP.

Once installed, GURL itself supports both HTTP and HTTPS:

gurl https://example.com

⸻

Usage / Как пользоваться

Basic syntax:

gurl [options] <URL>

Simple GET request:

gurl https://example.com

A URL without a scheme automatically uses HTTP:

gurl example.com

Equivalent to:

gurl http://example.com

⸻

Examples / Примеры

Download a file

gurl -o file.zip https://example.com/file.zip

or:

gurl --output file.zip https://example.com/file.zip

⸻

Use GURL in shell scripts

gurl https://example.com/script.sh | bash

For example:

gurl https://raw.githubusercontent.com/matveynator/sysadminscripts/main/label | bash

⸻

POST data

gurl -d "key=value&key2=value2" https://example.com

or:

gurl --data "key=value&key2=value2" https://example.com

When -d / --data is used, GURL automatically performs a POST request.

⸻

Send JSON

gurl -H "Content-Type: application/json" -d '{"key":"value"}' https://example.com/api

⸻

Multipart form

Text fields and files can be supplied in one -F argument separated by &:

gurl -F "name=test&file=@/tmp/file.txt" https://example.com/upload

⸻

Custom HTTP method

gurl -X DELETE https://example.com/api/item
gurl -X PUT https://example.com/api/item

⸻

Custom header

gurl -H "Authorization: Bearer TOKEN" https://example.com/api

or:

gurl --header "Authorization: Bearer TOKEN" https://example.com/api

⸻

Cookies

gurl -b "session_id=abc123" https://example.com

or:

gurl --cookie "session_id=abc123" https://example.com

⸻

HEAD request

gurl -I https://example.com

or:

gurl --head https://example.com

⸻

Custom User-Agent

gurl -A "MyClient/1.0" https://example.com

or:

gurl --useragent "MyClient/1.0" https://example.com

⸻

Timeout

gurl -m 10s https://example.com

or:

gurl --timeout 10s https://example.com

The default timeout is 30 seconds.

⸻

Ignore TLS certificate verification

Useful for self-signed certificates:

gurl -k https://192.168.1.1

or:

gurl --unsafe https://192.168.1.1

Warning: disabling certificate verification reduces connection security.

⸻

Fail on HTTP errors

gurl --fail https://example.com/not-found

For HTTP status 400 and above, GURL exits with error code 22.

This is useful in shell scripts:

gurl --fail https://example.com/file || echo "Download failed"

⸻

Redirects

GURL follows HTTP redirects by default.

gurl https://example.com

The compatible -L / --location option is also available.

⸻

Show version

gurl -V

or:

gurl --version

⸻

Flags / Флаги

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

⸻

Common Problems Solved by GURL

These are typical situations where GURL is useful:

* I need to download one file but curl is not installed.
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

“I just need curl, but it isn’t there.”

That’s what GURL is for.

One binary.
No external libraries.
HTTP and HTTPS.
Runs on a lot of platforms.

⸻

Build from source

Clone the repository:

git clone https://github.com/matveynator/gurl.git
cd gurl

Build:

go build -o gurl gurl.go

Run:

./gurl https://example.com

⸻

Source

https://github.com/matveynator/gurl

Contributions and fixes are welcome.