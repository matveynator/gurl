# GURL: A Simple `curl` Alternative That Works Everywhere

**GURL** is a lightweight command-line HTTP client with **no external dependencies**. It’s perfect for both modern and old systems, running smoothly without OpenSSL or other libraries. Whether you’re on a cutting-edge Linux server or ancient hardware, GURL just works.

---

## Download and Install

You can download GURL using `curl` (no HTTPS required) for your platform. Replace `<ARCHIVE_URL>` with the appropriate link for your operating system and architecture from the table below:

```bash
curl -L http://files.zabiyaka.net/gurl/latest/no-gui/<PLATFORM>/<ARCH>/gurl -o /usr/local/bin/gurl
chmod +x /usr/local/bin/gurl
```

Replace `<PLATFORM>` and `<ARCH>` with your system’s name and architecture.

---

Here’s the updated table with **all links** properly included for each platform and sorted by popularity. Each section references the full file paths provided.

---

## Supported Platforms and Binaries

| **Operating System**                         | **Architectures and Download Links**                                                                                                      |
|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------------|
| ![Linux](https://edent.github.io/SuperTinyIcons/images/svg/linux.svg) **Linux**       | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/amd64/gurl) [386](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/386/gurl) [arm](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/arm/gurl) [arm64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/arm64/gurl) [loong64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/loong64/gurl) [mips](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/mips/gurl) [mipsle](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/mipsle/gurl) [mips64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/mips64/gurl) [mips64le](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/mips64le/gurl) [ppc64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/ppc64/gurl) [ppc64le](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/ppc64le/gurl) [riscv64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/riscv64/gurl) [s390x](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/linux/s390x/gurl) |
| ![Windows](https://edent.github.io/SuperTinyIcons/images/svg/windows.svg) **Windows**  | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/windows/amd64/gurl.exe) [386](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/windows/386/gurl.exe) [arm](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/windows/arm/gurl.exe) [arm64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/windows/arm64/gurl.exe) |
| ![macOS](https://edent.github.io/SuperTinyIcons/images/svg/apple.svg) **macOS**        | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/mac/amd64/gurl) [arm64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/mac/arm64/gurl)                        |
| ![Android](https://edent.github.io/SuperTinyIcons/images/svg/android.svg) **Android**  | [arm64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/android/arm64/gurl)                                                |
| **FreeBSD**                                  | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/freebsd/amd64/gurl) [386](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/freebsd/386/gurl) [arm](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/freebsd/arm/gurl) [arm64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/freebsd/arm64/gurl) [riscv64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/freebsd/riscv64/gurl) |
| **OpenBSD**                                  | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/openbsd/amd64/gurl) [386](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/openbsd/386/gurl) [arm](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/openbsd/arm/gurl) [arm64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/openbsd/arm64/gurl) [ppc64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/openbsd/ppc64/gurl) [riscv64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/openbsd/riscv64/gurl) |
| **NetBSD**                                   | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/netbsd/amd64/gurl) [386](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/netbsd/386/gurl) [arm](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/netbsd/arm/gurl) [arm64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/netbsd/arm64/gurl) |
| **Solaris**                                  | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/solaris/amd64/gurl)                                                |
| **Plan 9**                                   | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/plan9/amd64/gurl) [386](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/plan9/386/gurl) [arm](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/plan9/arm/gurl) |
| **Illumos**                                  | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/illumos/amd64/gurl)                                                |
| **DragonFlyBSD**                             | [amd64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/dragonfly/amd64/gurl)                                              |
| **AIX**                                      | [ppc64](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/aix/ppc64/gurl)                                                    |
| **Wasm**                                     | [js/wasm](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/js/wasm/gurl)                                                   |
| **Wasi**                                     | [wasip1](http://files.zabiyaka.net/gurl/binaries/latest/no-gui/wasip1/wasm/gurl)                                                 |

---

> **Note:** Replace `/usr/local/bin` with an appropriate directory for your system if you lack permissions.

---

## Installation:

1. **Download and Install for Linux x86_64**:
   ```bash
   sudo curl -L http://files.zabiyaka.net/gurl/latest/no-gui/linux/amd64/gurl -o /usr/local/bin/gurl
   chmod +x /usr/local/bin/gurl
   ```

2. **Download and Install for Windows (PowerShell)**:
   ```powershell
   Invoke-WebRequest -Uri http://files.zabiyaka.net/gurl/latest/no-gui/windows/amd64/gurl.exe -OutFile gurl.exe
   ```

---

## How to Use GURL

GURL supports clear and descriptive options to make usage intuitive. Here are practical examples:

---

### 1. **Basic GET Request**

```bash
gurl http://example.com
```
Sends a simple GET request and outputs the response to the terminal.

---

### 2. **Save Response to a File**

```bash
gurl --output output.txt http://example.com
```
The `--output` option saves the response body to `output.txt` instead of displaying it in the terminal.

---

### 3. **POST Data**

```bash
gurl --request POST --data "key=value&key2=value2" http://example.com
```
Sends a POST request with form data using `--data` and explicitly sets the HTTP method with `--request`.

---

### 4. **Upload Files and Fields (Multipart Form Data)**

```bash
gurl --form "key=value" --form "file=@/path/to/file" http://example.com
```
Uploads both text fields and files using the `--form` option. Use `key=@/path/to/file` to upload files.

---

### 5. **Custom Headers**

```bash
gurl --header "Authorization: Bearer TOKEN" --header "Content-Type: application/json" http://example.com
```
Use `--header` to include custom headers in the request.

---

### 6. **Send Cookies**

```bash
gurl --cookie "session_id=abc123; user=example" http://example.com
```
The `--cookie` option sends cookies with the request.

---

### 7. **Fail Silently on HTTP Errors**

```bash
gurl --fail http://example.com/404
```
If an HTTP error (4xx or 5xx) is returned, GURL exits with a **non-zero error code** and suppresses the response body.

---

### 8. **HEAD Request**

```bash
gurl --head http://example.com
```
The `--head` option retrieves only the response headers.

---

### 9. **Set User-Agent**

```bash
gurl --useragent "CustomUserAgent/1.0" http://example.com
```
The `--useragent` option specifies a custom User-Agent header.

---

### 10. **Set a Timeout**

```bash
gurl --timeout 10s http://example.com
```
The `--timeout` option sets the maximum time GURL will wait for a response (e.g., `10s`, `1m`).

---

### 11. **Download a File with Progress Display**

```bash
gurl http://example.com/file.txt --output file.txt
```
When using `--output` to save a file, GURL automatically displays a download progress bar.

---

### 12. **Silent Mode for Scripting**

```bash
gurl --silent http://example.com > output.html
```
The `--silent` option suppresses all output, including progress and headers, for clean scripting.

---

### 13. **Send JSON Data**

```bash
gurl --request POST --header "Content-Type: application/json" --data '{"key":"value"}' http://example.com
```
Combines `--header` for content type and `--data` for JSON payload.

---

### 14. **Verbose Output for Debugging**

```bash
gurl --verbose http://example.com
```
The `--verbose` option prints detailed request and response information, useful for debugging.

---

### 15. **Handle Timeouts and Failures**

```bash
gurl --timeout 5s --fail http://example.com
```
Combines `--timeout` for a 5-second limit and `--fail` to exit silently on HTTP errors.

---

These examples highlight the simplicity and power of GURL. Whether you're downloading files, uploading data, or automating tasks, GURL gives you clear and reliable command-line control. 🚀
