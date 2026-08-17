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

It is designed for situations where you just need to download a new software on an old system, call an HTTPS endpoint, send POST data, or make a simple request — without installing `curl`, OpenSSL, or a collection of shared libraries.

TLS support is built into the binary through Go's standard library.

This makes GURL especially useful for:

* installing new software on old operating systems 
* servers where `curl` or OpenSSL is unavailable or outdated

Just download **one executable** and run it.

```bash
gurl https://example.com
```

If no protocol is specified, GURL uses `http://`:

```bash
gurl example.com
```

---

<details>
<summary>
  <b><big>Universal installer / Универсальная установка</big></b>
  <br><sub>auto-detect OS/ARCH · root → /bin · user → ~/.local/bin</sub>
</summary>

<br>

For Unix-like systems you can use one bootstrap installer.

It detects the operating system and CPU architecture, downloads the correct GURL binary using an available native downloader, and selects the installation path automatically:

* if the installer runs as **root** → `/bin/gurl`
* if it runs as a normal user → `~/.local/bin/gurl`

After installation it immediately prints ready-to-copy examples showing how to run GURL.

### Install

Normal user:

```sh
sh install-gurl.sh
```

System-wide:

```sh
sudo sh install-gurl.sh
```

### `install-gurl.sh`

```sh
#!/bin/sh

set -u

BASE_URL="${GURL_BASE_URL:-http://files.zabiyaka.net/gurl/latest}"

have() {
    command -v "$1" >/dev/null 2>&1
}

die() {
    echo "gurl-install: $*" >&2
    exit 1
}

say() {
    echo "gurl-install: $*"
}

# Detect OS

UNAME_S="$(uname -s 2>/dev/null || echo unknown)"

case "$UNAME_S" in
    Linux)
        if [ -n "${ANDROID_ROOT:-}" ] ||
           [ -n "${ANDROID_DATA:-}" ] ||
           [ -n "${TERMUX_VERSION:-}" ]; then
            OS="android"
        else
            OS="linux"
        fi
        ;;
    Darwin) OS="mac" ;;
    FreeBSD) OS="freebsd" ;;
    OpenBSD) OS="openbsd" ;;
    NetBSD) OS="netbsd" ;;
    DragonFly) OS="dragonfly" ;;
    SunOS)
        SUN_VERSION="$(uname -v 2>/dev/null || true)"
        if echo "$SUN_VERSION" | grep -i illumos >/dev/null 2>&1; then
            OS="illumos"
        elif [ -r /etc/release ] &&
             grep -Ei 'illumos|OpenIndiana|OmniOS|SmartOS' /etc/release >/dev/null 2>&1; then
            OS="illumos"
        else
            OS="solaris"
        fi
        ;;
    AIX) OS="aix" ;;
    *) die "unsupported operating system: $UNAME_S" ;;
esac

# Detect architecture

if [ "$OS" = "aix" ]; then
    ARCH="ppc64"
else
    UNAME_M="$(uname -m 2>/dev/null || echo unknown)"

    case "$UNAME_M" in
        x86_64|amd64) ARCH="amd64" ;;
        i386|i486|i586|i686|x86) ARCH="386" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv4*|armv5*|armv6*|armv7*|armv8l|arm) ARCH="arm" ;;
        loongarch64|loong64) ARCH="loong64" ;;
        mips64el|mips64le) ARCH="mips64le" ;;
        mips64) ARCH="mips64" ;;
        mipsel|mipsle) ARCH="mipsle" ;;
        mips) ARCH="mips" ;;
        ppc64le|powerpc64le) ARCH="ppc64le" ;;
        ppc64|powerpc64) ARCH="ppc64" ;;
        riscv64) ARCH="riscv64" ;;
        s390x) ARCH="s390x" ;;
        *) die "unsupported CPU architecture: $UNAME_M" ;;
    esac
fi

# Verify that this build exists

case "$OS/$ARCH" in
    linux/amd64|linux/arm64|linux/386|linux/arm|linux/loong64|\
    linux/mips|linux/mipsle|linux/mips64|linux/mips64le|\
    linux/ppc64|linux/ppc64le|linux/riscv64|linux/s390x|\
    mac/amd64|mac/arm64|\
    freebsd/amd64|freebsd/arm64|freebsd/386|freebsd/arm|freebsd/riscv64|\
    openbsd/amd64|openbsd/arm64|openbsd/386|openbsd/arm|openbsd/ppc64|openbsd/riscv64|\
    android/arm64|\
    netbsd/amd64|netbsd/386|netbsd/arm|netbsd/arm64|\
    solaris/amd64|illumos/amd64|dragonfly/amd64|aix/ppc64)
        ;;
    *) die "GURL binary is not available for $OS/$ARCH" ;;
esac

URL="$BASE_URL/$OS/$ARCH/gurl"

say "detected platform: $OS/$ARCH"
say "download: $URL"

# Temporary directory

TMPBASE="${TMPDIR:-/tmp}"
TMP="$TMPBASE/gurl-install.$$"
DOWNLOAD="$TMP/gurl"

umask 077
mkdir "$TMP" 2>/dev/null || die "cannot create temporary directory: $TMP"

cleanup() {
    rm -rf "$TMP" >/dev/null 2>&1 || true
}

trap cleanup EXIT HUP INT TERM

# Download using the most suitable available tool

download() {
    if [ "$OS" = "freebsd" ] || [ "$OS" = "dragonfly" ]; then
        if have fetch; then
            say "using fetch"
            fetch -o "$DOWNLOAD" "$URL" && return 0
        fi
    fi

    if [ "$OS" = "openbsd" ] || [ "$OS" = "netbsd" ]; then
        if have ftp; then
            say "using ftp"
            ftp -o "$DOWNLOAD" "$URL" && return 0
        fi
    fi

    if [ -x /usr/sfw/bin/wget ]; then
        say "using /usr/sfw/bin/wget"
        /usr/sfw/bin/wget -O "$DOWNLOAD" "$URL" && return 0
    fi

    if have wget; then
        say "using wget"
        wget -O "$DOWNLOAD" "$URL" && return 0
    fi

    if have busybox && busybox wget --help >/dev/null 2>&1; then
        say "using BusyBox wget"
        busybox wget -O "$DOWNLOAD" "$URL" && return 0
    fi

    if have toybox && toybox wget --help >/dev/null 2>&1; then
        say "using Toybox wget"
        toybox wget -O "$DOWNLOAD" "$URL" && return 0
    fi

    if have curl; then
        say "using curl"
        curl -fL -o "$DOWNLOAD" "$URL" && return 0
    fi

    die "no usable HTTP downloader found.

Tried:
  fetch
  ftp
  wget
  /usr/sfw/bin/wget
  BusyBox wget
  Toybox wget
  curl

Direct URL:
  $URL"
}

download

[ -s "$DOWNLOAD" ] || die "downloaded file is empty"

chmod 755 "$DOWNLOAD" || die "cannot make downloaded GURL executable"

# Root -> /bin, normal user -> ~/.local/bin

IS_ROOT=0

if have id && [ "$(id -u 2>/dev/null)" = "0" ]; then
    IS_ROOT=1
fi

if [ "$IS_ROOT" = "1" ]; then
    DESTDIR="/bin"
else
    [ -n "${HOME:-}" ] || die "HOME is not defined"
    DESTDIR="$HOME/.local/bin"
fi

DEST="$DESTDIR/gurl"

mkdir -p "$DESTDIR" || die "cannot create $DESTDIR"
cp "$DOWNLOAD" "$DEST" || die "cannot install $DEST"
chmod 755 "$DEST" || die "cannot chmod $DEST"

say "installed: $DEST"

"$DEST" -V || die "GURL was installed, but version test failed"

# Print immediate help

echo
echo "============================================================"
echo "GURL is ready"
echo "============================================================"
echo
echo "Installed:"
echo "  $DEST"
echo
echo "Try it now:"
echo
echo "  $DEST https://example.com"
echo
echo "Download a file:"
echo
echo "  $DEST -o file.zip https://example.com/file.zip"
echo

if [ "$IS_ROOT" = "1" ]; then
    echo "System-wide command:"
    echo
    echo "  gurl https://example.com"
else
    case ":${PATH:-}:" in
        *":$DESTDIR:"*)
            echo "GURL is already in your PATH."
            echo
            echo "Run:"
            echo
            echo "  gurl https://example.com"
            ;;
        *)
            echo "$DESTDIR is not in your PATH yet."
            echo
            echo "For this shell, run:"
            echo
            echo "  export PATH=\"$DESTDIR:\$PATH\""
            echo
            echo "Then:"
            echo
            echo "  gurl https://example.com"
            echo
            echo "To make it permanent, add the export line to your shell profile."
            ;;
    esac
fi

echo
echo "Help:"
echo "  gurl --help"
echo
echo "Version:"
echo "  gurl -V"
echo
```

---

</details>

---

# Downloads / Скачать

Choose your platform and architecture below.

> The binary URLs intentionally use **HTTP**. This is the bootstrap path for old or minimal systems where HTTPS tools, CA certificates, OpenSSL, or the installed TLS stack may be unavailable or obsolete. After downloading GURL, use GURL itself for HTTPS.


<details>
<summary>
  <img width="42" alt="Linux" src="https://github.com/user-attachments/assets/bf3141b6-4c93-4fd6-b2d1-421b79876dcb" />
  <b><big>Linux</big></b>
  <br><sub>amd64 / arm64 / 386 / ARM / LoongArch / MIPS / PPC64 / RISC-V / s390x</sub>
</summary>

<br>

### amd64 / x86_64

[download](http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/amd64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### arm64 / aarch64

[download](http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/arm64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### 386 / x86

[download](http://files.zabiyaka.net/gurl/latest/linux/386/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/386/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/386/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/386/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### ARM

[download](http://files.zabiyaka.net/gurl/latest/linux/arm/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/arm/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/arm/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/arm/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### LoongArch 64

[download](http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/loong64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### MIPS

[download](http://files.zabiyaka.net/gurl/latest/linux/mips/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### MIPS little-endian

[download](http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mipsle/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### MIPS64

[download](http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### MIPS64 little-endian

[download](http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/mips64le/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### PowerPC 64

[download](http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/ppc64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### PowerPC 64 little-endian

[download](http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/ppc64le/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### RISC-V 64

[download](http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/riscv64/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### IBM s390x

[download](http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/linux/s390x/gurl; else echo 'No HTTP downloader found (wget, BusyBox wget, or curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <img width="36" alt="macOS" src="https://github.com/user-attachments/assets/946102b8-f043-494d-809a-a589e536ee9a" />
  <b><big>macOS</big></b>
  <br><sub>Intel / Apple Silicon</sub>
</summary>

<br>

### Intel / amd64

[download](http://files.zabiyaka.net/gurl/latest/mac/amd64/gurl)

```bash
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; /usr/bin/curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/mac/amd64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### Apple Silicon / arm64

[download](http://files.zabiyaka.net/gurl/latest/mac/arm64/gurl)

```bash
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; /usr/bin/curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/mac/arm64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <img width="42" alt="Windows" src="https://github.com/user-attachments/assets/f6044001-95b0-4500-a4f6-1c3b08eb65fb" />
  <b><big>Windows</big></b>
  <br><sub>amd64 / arm64 / 386 / ARM</sub>
</summary>

<br>

The command uses `System.Net.WebClient`, which works with older Windows PowerShell versions and does not require `Invoke-WebRequest`.

### amd64 / x86_64

[download](http://files.zabiyaka.net/gurl/latest/windows/amd64/gurl.exe)

```powershell
(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/amd64/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V
```

### arm64

[download](http://files.zabiyaka.net/gurl/latest/windows/arm64/gurl.exe)

```powershell
(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/arm64/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V
```

### 386 / x86

[download](http://files.zabiyaka.net/gurl/latest/windows/386/gurl.exe)

```powershell
(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/386/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V
```

### ARM

[download](http://files.zabiyaka.net/gurl/latest/windows/arm/gurl.exe)

```powershell
(New-Object System.Net.WebClient).DownloadFile("http://files.zabiyaka.net/gurl/latest/windows/arm/gurl.exe", "$PWD\gurl.exe"); .\gurl.exe -V
```

</details>

---

<details>
<summary>
  <img width="42" alt="FreeBSD" src="https://github.com/user-attachments/assets/d35baaac-d296-41b1-a281-55dc761328e9" />
  <b><big>FreeBSD</big></b>
  <br><sub>amd64 / arm64 / 386 / ARM / RISC-V</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/freebsd/amd64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; fetch -o "$DEST" http://files.zabiyaka.net/gurl/latest/freebsd/amd64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### arm64

[download](http://files.zabiyaka.net/gurl/latest/freebsd/arm64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; fetch -o "$DEST" http://files.zabiyaka.net/gurl/latest/freebsd/arm64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### 386

[download](http://files.zabiyaka.net/gurl/latest/freebsd/386/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; fetch -o "$DEST" http://files.zabiyaka.net/gurl/latest/freebsd/386/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### ARM

[download](http://files.zabiyaka.net/gurl/latest/freebsd/arm/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; fetch -o "$DEST" http://files.zabiyaka.net/gurl/latest/freebsd/arm/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### RISC-V 64

[download](http://files.zabiyaka.net/gurl/latest/freebsd/riscv64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; fetch -o "$DEST" http://files.zabiyaka.net/gurl/latest/freebsd/riscv64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <img width="42" alt="OpenBSD" src="https://github.com/user-attachments/assets/11633d7e-5744-46da-ad2f-6e49c69e51de" />
  <b><big>OpenBSD</big></b>
  <br><sub>amd64 / arm64 / 386 / ARM / PPC64 / RISC-V</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/amd64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/openbsd/amd64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### arm64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/arm64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/openbsd/arm64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### 386

[download](http://files.zabiyaka.net/gurl/latest/openbsd/386/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/openbsd/386/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### ARM

[download](http://files.zabiyaka.net/gurl/latest/openbsd/arm/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/openbsd/arm/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### PowerPC 64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/ppc64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/openbsd/ppc64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### RISC-V 64

[download](http://files.zabiyaka.net/gurl/latest/openbsd/riscv64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/openbsd/riscv64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <b><big>Android</big></b>
  <br><sub>arm64</sub>
</summary>

<br>

### arm64

[download](http://files.zabiyaka.net/gurl/latest/android/arm64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; elif command -v toybox >/dev/null 2>&1; then toybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; elif command -v busybox >/dev/null 2>&1; then busybox wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/android/arm64/gurl; else echo 'No HTTP downloader found (wget/Toybox/BusyBox/curl).' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <b><big>NetBSD</big></b>
  <br><sub>amd64 / 386 / ARM / arm64</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/netbsd/amd64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/netbsd/amd64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### 386

[download](http://files.zabiyaka.net/gurl/latest/netbsd/386/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/netbsd/386/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### ARM

[download](http://files.zabiyaka.net/gurl/latest/netbsd/arm/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/netbsd/arm/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

### arm64

[download](http://files.zabiyaka.net/gurl/latest/netbsd/arm64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; ftp -o "$DEST" http://files.zabiyaka.net/gurl/latest/netbsd/arm64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <b><big>Solaris</big></b>
  <br><sub>amd64</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl; elif [ -x /usr/sfw/bin/wget ]; then /usr/sfw/bin/wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/solaris/amd64/gurl; else echo 'No HTTP downloader found. Install/copy wget or curl, or use the direct download link above.' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <b><big>Plan 9</big></b>
  <br><sub>amd64 / 386 / ARM</sub>
</summary>

<br>

Uses Plan 9 `hget`, whose native purpose is retrieving HTTP URLs.

### amd64

[download](http://files.zabiyaka.net/gurl/latest/plan9/amd64/gurl)

```rc
hget http://files.zabiyaka.net/gurl/latest/plan9/amd64/gurl >gurl
./gurl -V
```

### 386

[download](http://files.zabiyaka.net/gurl/latest/plan9/386/gurl)

```rc
hget http://files.zabiyaka.net/gurl/latest/plan9/386/gurl >gurl
./gurl -V
```

### ARM

[download](http://files.zabiyaka.net/gurl/latest/plan9/arm/gurl)

```rc
hget http://files.zabiyaka.net/gurl/latest/plan9/arm/gurl >gurl
./gurl -V
```

</details>

---

<details>
<summary>
  <b><big>Illumos</big></b>
  <br><sub>amd64</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl; elif [ -x /usr/sfw/bin/wget ]; then /usr/sfw/bin/wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/illumos/amd64/gurl; else echo 'No HTTP downloader found. Install/copy wget or curl, or use the direct download link above.' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <b><big>DragonFly BSD</big></b>
  <br><sub>amd64</sub>
</summary>

<br>

### amd64

[download](http://files.zabiyaka.net/gurl/latest/dragonfly/amd64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; fetch -o "$DEST" http://files.zabiyaka.net/gurl/latest/dragonfly/amd64/gurl && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>

---

<details>
<summary>
  <b><big>AIX</big></b>
  <br><sub>ppc64</sub>
</summary>

<br>

AIX base `ftp` speaks FTP, not HTTP. Because the GURL bootstrap server here is HTTP, a bare AIX installation needs an available HTTP downloader (`wget`/`curl`) or the binary copied from another machine.

### PowerPC 64

[download](http://files.zabiyaka.net/gurl/latest/aix/ppc64/gurl)

```sh
if [ "$(id -u 2>/dev/null)" = "0" ]; then DESTDIR="/bin"; else DESTDIR="$HOME/.local/bin"; fi; mkdir -p "$DESTDIR"; DEST="$DESTDIR/gurl"; if command -v wget >/dev/null 2>&1; then wget -O "$DEST" http://files.zabiyaka.net/gurl/latest/aix/ppc64/gurl; elif command -v curl >/dev/null 2>&1; then curl -fL -o "$DEST" http://files.zabiyaka.net/gurl/latest/aix/ppc64/gurl; else echo 'Base AIX ftp does not fetch HTTP URLs. Copy/install wget or curl, or use the direct download link above.' >&2; exit 1; fi && chmod 755 "$DEST" && "$DEST" -V && echo && echo "GURL installed: $DEST" && echo "Run: $DEST https://example.com"
```

</details>


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
* I need an HTTP client for an old or unusual operating system.
* I need the same command-line utility on Linux, BSD, macOS, and Windows.
* I need a curl-like utility compiled with TLS support already included.

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
