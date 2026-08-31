package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCommandLineValidation(t *testing.T) {
	testCases := []struct {
		name      string
		arguments []string
		message   string
	}{
		{name: "missing URL", message: "no URL specified"},
		{name: "multiple URLs", arguments: []string{"one", "two"}, message: "only one URL is supported"},
		{name: "missing long value", arguments: []string{"--header"}, message: "option --header requires a value"},
		{name: "value on switch", arguments: []string{"--silent=yes", "example.com"}, message: "option --silent does not take a value"},
		{name: "unknown long option", arguments: []string{"--unknown", "example.com"}, message: "unknown option --unknown"},
		{name: "missing short value", arguments: []string{"-H"}, message: "option -H requires a value"},
		{name: "unknown short option", arguments: []string{"-q", "example.com"}, message: "unknown option -q"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := parseCommandLine(testCase.arguments)
			if err == nil || err.Error() != testCase.message {
				t.Fatalf("error = %v, want %q", err, testCase.message)
			}
		})
	}
}

func TestCommandLineLegacyAndFileOptions(t *testing.T) {
	options, err := parseCommandLine([]string{
		"--dump-header=headers.txt",
		"--trace-ascii=trace.txt",
		"--useragent=agent",
		"--request=PATCH",
		"--timeout=2s",
		"--cacert=ca.pem",
		"--capath=certificates",
		"--output=-",
		"--form=name=contents",
		"--cookie=session=one",
		"--unsafe",
		"--head",
		"--fail",
		"--globoff",
		"--",
		"-literal-url",
	})
	if err != nil {
		t.Fatal(err)
	}
	if options.dumpHeaderPath != "headers.txt" || options.tracePath != "trace.txt" || options.userAgent != "agent" {
		t.Fatalf("file and identity options = %#v", options)
	}
	if options.requestMethod != "PATCH" || options.maximumTime != 2*time.Second || !options.insecure || !options.head || !options.failHTTP {
		t.Fatalf("request options = %#v", options)
	}
	if options.caCertificate != "ca.pem" || options.caPath != "certificates" || options.outputPath != "-" {
		t.Fatalf("path options = %#v", options)
	}
	if options.requestURL != "-literal-url" || strings.Join(options.formFields, ",") != "name=contents" || options.cookie != "session=one" {
		t.Fatalf("payload options = %#v", options)
	}
}

func TestShortValueOptions(t *testing.T) {
	options, err := parseCommandLine([]string{
		"-Dheaders.txt", "-Aagent", "-XPATCH", "-HOne: first", "-H", "One: second",
		"-dbody", "-m2.5", "-o-", "-Fname=contents", "-bsession=one", "example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if options.dumpHeaderPath != "headers.txt" || options.userAgent != "agent" || options.requestMethod != "PATCH" {
		t.Fatalf("short options = %#v", options)
	}
	if len(options.headers) != 2 || options.maximumTime != 2500*time.Millisecond || options.outputPath != "-" {
		t.Fatalf("short values = %#v", options)
	}
}

func TestInvalidDurations(t *testing.T) {
	for _, arguments := range [][]string{
		{"--connect-timeout=invalid", "example.com"},
		{"--max-time=-1", "example.com"},
		{"--max-time=NaN", "example.com"},
		{"--max-time=1e100", "example.com"},
		{"--timeout=invalid", "example.com"},
		{"--timeout=-1s", "example.com"},
	} {
		if _, err := parseCommandLine(arguments); err == nil {
			t.Fatalf("arguments %#v were accepted", arguments)
		}
	}
}

func TestMultipartRequestBody(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "contents.txt")
	if err := os.WriteFile(filePath, []byte("file contents"), 0o600); err != nil {
		t.Fatal(err)
	}
	body, contentType, err := requestBody(commandOptions{formFields: []string{"name=example", "upload=@" + filePath}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(contentType, "multipart/form-data; boundary=") || !bytes.Contains(body, []byte("file contents")) || !bytes.Contains(body, []byte("example")) {
		t.Fatalf("content type = %q, body = %q", contentType, body)
	}
	for _, formField := range []string{"invalid", "upload=@" + filepath.Join(t.TempDir(), "missing")} {
		if _, _, err := requestBody(commandOptions{formFields: []string{formField}}); err == nil {
			t.Fatalf("form field %q was accepted", formField)
		}
	}
}

func TestRequestMethods(t *testing.T) {
	testCases := []struct {
		options commandOptions
		method  string
	}{
		{method: http.MethodGet},
		{options: commandOptions{head: true}, method: http.MethodHead},
		{options: commandOptions{requestData: []string{"body"}}, method: http.MethodPost},
		{options: commandOptions{formFields: []string{"name=contents"}}, method: http.MethodPost},
		{options: commandOptions{requestMethod: http.MethodPatch, head: true}, method: http.MethodPatch},
	}
	for _, testCase := range testCases {
		if actualMethod := requestMethod(testCase.options); actualMethod != testCase.method {
			t.Fatalf("method = %q, want %q", actualMethod, testCase.method)
		}
	}
}

func TestRequestHeaderRules(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	options := commandOptions{
		userAgent: "default-agent",
		cookie:    "session=one",
		headers: []string{
			"X-Empty;", "X-Repeated: first", "X-Repeated: second", "Host: virtual.example",
			"User-Agent:", "Cookie:",
		},
	}
	if err := applyRequestHeaders(request, options, "application/json"); err != nil {
		t.Fatal(err)
	}
	if request.Host != "virtual.example" || request.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("request = %#v", request)
	}
	if strings.Join(request.Header.Values("X-Repeated"), ",") != "first,second" || len(request.Header.Values("X-Empty")) != 1 {
		t.Fatalf("headers = %#v", request.Header)
	}
	if request.Header.Get("Cookie") != "" || len(request.Header.Values("User-Agent")) != 1 {
		t.Fatalf("removed headers = %#v", request.Header)
	}
	for _, header := range []string{";", "invalid", ": contents"} {
		invalidRequest, requestErr := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if requestErr != nil {
			t.Fatal(requestErr)
		}
		if err := applyRequestHeaders(invalidRequest, commandOptions{headers: []string{header}}, ""); err == nil {
			t.Fatalf("header %q was accepted", header)
		}
	}
}

func TestURLValidation(t *testing.T) {
	parsedURL, err := ensureHTTPURL("example.com/path")
	if err != nil || parsedURL.String() != "http://example.com/path" {
		t.Fatalf("URL = %v, error = %v", parsedURL, err)
	}
	for _, rawURL := range []string{"ftp://example.com", "https://", "://invalid"} {
		if _, err := ensureHTTPURL(rawURL); err == nil {
			t.Fatalf("URL %q was accepted", rawURL)
		}
	}
}

func TestCertificateSources(t *testing.T) {
	emptyDirectory := t.TempDir()
	if _, err := certificatePool(commandOptions{caPath: emptyDirectory}); err == nil {
		t.Fatal("empty CA directory was accepted")
	}
	if _, err := certificatePool(commandOptions{caPath: filepath.Join(emptyDirectory, "missing")}); err == nil {
		t.Fatal("missing CA directory was accepted")
	}
	if _, err := certificatePool(commandOptions{caCertificate: filepath.Join(emptyDirectory, "missing.pem")}); err == nil {
		t.Fatal("missing CA file was accepted")
	}
}

func TestTransportAndOutputFiles(t *testing.T) {
	transport := createTransport(commandOptions{connectTimeout: time.Second, insecure: true}, x509.NewCertPool())
	if transport.TLSHandshakeTimeout != time.Second || !transport.TLSClientConfig.InsecureSkipVerify || !transport.ForceAttemptHTTP2 {
		t.Fatalf("transport = %#v", transport)
	}
	transport.CloseIdleConnections()

	var standardOutput bytes.Buffer
	output, err := openOutput("-", &standardOutput)
	if err != nil || output.writer != &standardOutput {
		t.Fatalf("standard output = %#v, error = %v", output, err)
	}
	filePath := filepath.Join(t.TempDir(), "output.txt")
	output, err = openOutput(filePath, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(output.writer, "contents"); err != nil {
		t.Fatal(err)
	}
	if err := output.close(); err != nil {
		t.Fatal(err)
	}
	if _, err := openOutput(t.TempDir(), io.Discard); err == nil {
		t.Fatal("directory was accepted as an output file")
	}
}

func TestTraceProgressAndHeaderWriters(t *testing.T) {
	var traceOutput bytes.Buffer
	trace := asciiTrace{writer: &traceOutput}
	request, err := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader("request body"))
	if err != nil {
		t.Fatal(err)
	}
	trace.request(request)
	trace.failure(errors.New("trace failure"))
	traceReader := &traceReader{reader: strings.NewReader("response body"), trace: trace}
	if _, err := io.ReadAll(traceReader); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(traceOutput.String(), "request body") || !strings.Contains(traceOutput.String(), "response body") || !strings.Contains(traceOutput.String(), "trace failure") {
		t.Fatalf("trace = %q", traceOutput.String())
	}

	var diagnostics bytes.Buffer
	progress := &progressReader{reader: strings.NewReader("body"), total: 4, diagnostics: &diagnostics}
	if _, err := io.ReadAll(progress); err != nil {
		t.Fatal(err)
	}
	unknownProgress := &progressReader{reader: strings.NewReader("body"), diagnostics: &diagnostics}
	if _, err := io.ReadAll(unknownProgress); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(diagnostics.String(), "100%") || !strings.Contains(diagnostics.String(), "4 bytes") {
		t.Fatalf("progress = %q", diagnostics.String())
	}

	response := &http.Response{Proto: "HTTP/1.1", Status: "200 OK", Header: http.Header{"X-Test": []string{"one"}}}
	if err := writeResponseHeaders(response, failingWriter{}, nil, asciiTrace{}); err == nil {
		t.Fatal("dump header write failure was ignored")
	}
	if err := writeResponseHeaders(response, nil, failingWriter{}, asciiTrace{}); err == nil {
		t.Fatal("HEAD write failure was ignored")
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestRequestErrorClassification(t *testing.T) {
	testCases := []struct {
		err      error
		exitCode int
	}{
		{err: errTooManyRedirects, exitCode: exitTooManyRedirects},
		{err: context.DeadlineExceeded, exitCode: exitTimeout},
		{err: timeoutError{}, exitCode: exitTimeout},
		{err: &net.DNSError{Err: "missing", Name: "invalid.example"}, exitCode: exitResolveHost},
		{err: &tls.CertificateVerificationError{Err: x509.UnknownAuthorityError{}}, exitCode: exitCertificate},
		{err: x509.UnknownAuthorityError{}, exitCode: exitCertificate},
		{err: x509.HostnameError{Host: "example.com"}, exitCode: exitCertificate},
		{err: x509.CertificateInvalidError{Reason: x509.Expired}, exitCode: exitCertificate},
		{err: tls.RecordHeaderError{}, exitCode: exitTLSConnect},
		{err: errors.New("remote error: tls: handshake failure"), exitCode: exitTLSConnect},
		{err: errors.New("connection refused"), exitCode: exitConnect},
	}
	for _, testCase := range testCases {
		if exitCode := classifyRequestError(testCase.err); exitCode != testCase.exitCode {
			t.Fatalf("error %T classified as %d, want %d", testCase.err, exitCode, testCase.exitCode)
		}
	}
}

func TestRunLocalValidationAndOutputErrors(t *testing.T) {
	testCases := []struct {
		name      string
		arguments []string
		exitCode  int
	}{
		{name: "command line", arguments: []string{"--unknown"}, exitCode: exitCommandLine},
		{name: "URL", arguments: []string{"ftp://example.com"}, exitCode: exitMalformedURL},
		{name: "multipart", arguments: []string{"-F", "invalid", "example.com"}, exitCode: exitReadError},
		{name: "certificate", arguments: []string{"--cacert", filepath.Join(t.TempDir(), "missing.pem"), "example.com"}, exitCode: exitCertificateProblem},
		{name: "header output", arguments: []string{"-D", t.TempDir(), "example.com"}, exitCode: exitWriteError},
		{name: "trace output", arguments: []string{"--trace-ascii", t.TempDir(), "example.com"}, exitCode: exitWriteError},
		{name: "body output", arguments: []string{"-o", t.TempDir(), "example.com"}, exitCode: exitWriteError},
		{name: "invalid header", arguments: []string{"-H", "invalid", "example.com"}, exitCode: exitCommandLine},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var diagnostics bytes.Buffer
			if exitCode := run(testCase.arguments, strings.NewReader(""), io.Discard, &diagnostics); exitCode != testCase.exitCode {
				t.Fatalf("exit code = %d, want %d; diagnostics = %q", exitCode, testCase.exitCode, diagnostics.String())
			}
			if diagnostics.Len() == 0 {
				t.Fatal("expected diagnostics")
			}
		})
	}
}

func TestWriteCommandErrorHonorsSilent(t *testing.T) {
	var diagnostics bytes.Buffer
	writeCommandError(&diagnostics, true, exitConnect, errors.New("hidden"))
	if diagnostics.Len() != 0 {
		t.Fatalf("silent diagnostics = %q", diagnostics.String())
	}
	writeCommandError(&diagnostics, false, exitConnect, errors.New("visible"))
	if !strings.Contains(diagnostics.String(), "gurl: (7) visible") {
		t.Fatalf("diagnostics = %q", diagnostics.String())
	}
}
