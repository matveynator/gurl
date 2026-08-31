package main

import (
	"bytes"
	"encoding/pem"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestParseAcmeCommandLine(t *testing.T) {
	arguments := []string{
		"--silent",
		"--dump-header", "/root/.acme.sh/http.header",
		"-L",
		"--connect-timeout", "10",
		"--user-agent", "acme.sh/3.1.5 (https://github.com/acmesh-official/acme.sh)",
		"-H", "",
		"-H", "Replay-Nonce: expected",
		"-H", "",
		"-H", "",
		"-H", "",
		"https://acme.zerossl.com/v2/DV90",
	}
	options, err := parseCommandLine(arguments)
	if err != nil {
		t.Fatal(err)
	}
	if !options.silent || !options.followLocation {
		t.Fatalf("silent=%v followLocation=%v", options.silent, options.followLocation)
	}
	if options.connectTimeout != 10*time.Second {
		t.Fatalf("connect timeout = %s", options.connectTimeout)
	}
	if options.dumpHeaderPath != "/root/.acme.sh/http.header" {
		t.Fatalf("dump header = %q", options.dumpHeaderPath)
	}
	if len(options.headers) != 5 || options.headers[1] != "Replay-Nonce: expected" {
		t.Fatalf("headers = %#v", options.headers)
	}
	if options.requestURL != "https://acme.zerossl.com/v2/DV90" {
		t.Fatalf("URL = %q", options.requestURL)
	}
}

func TestParseCurlOptionForms(t *testing.T) {
	options, err := parseCommandLine([]string{
		"-sLg6",
		"https://example.com",
		"-HFirst: one",
		"--header=Second: two",
		"-dalpha=one",
		"--data", "beta=two",
		"-m", "1.5",
		"-4",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !options.silent || !options.followLocation || options.ipVersion != "4" {
		t.Fatalf("parsed switches = %#v", options)
	}
	if options.maximumTime != 1500*time.Millisecond {
		t.Fatalf("maximum time = %s", options.maximumTime)
	}
	if got := strings.Join(options.requestData, "&"); got != "alpha=one&beta=two" {
		t.Fatalf("request data = %q", got)
	}
	if len(options.headers) != 2 {
		t.Fatalf("headers = %#v", options.headers)
	}
}

func TestHelpAdvertisesGloboff(t *testing.T) {
	var output bytes.Buffer
	exitCode := run([]string{"--help", "curl"}, strings.NewReader(""), &output, io.Discard)
	if exitCode != exitOK {
		t.Fatalf("exit code = %d", exitCode)
	}
	if !strings.Contains(output.String(), "--globoff") {
		t.Fatalf("help does not advertise --globoff:\n%s", output.String())
	}
}

func TestAcmeGetWritesBodyAndHeaderFile(t *testing.T) {
	var receivedUserAgent string
	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		receivedUserAgent = request.UserAgent()
		receivedHeaders = request.Header.Clone()
		response.Header().Set("Replay-Nonce", "nonce-value")
		response.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(response, `{"directory":true}`)
	}))
	defer server.Close()

	headerPath := filepath.Join(t.TempDir(), "http.header")
	var output bytes.Buffer
	var diagnostics bytes.Buffer
	exitCode := run([]string{
		"--silent", "--dump-header", headerPath, "-L", "--connect-timeout", "10",
		"--user-agent", "acme.sh/3.1.5 (https://github.com/acmesh-official/acme.sh)",
		"-H", "", "-H", "Accept: application/json", "-H", "", "-H", "", "-H", "", server.URL,
	}, strings.NewReader(""), &output, &diagnostics)
	if exitCode != exitOK {
		t.Fatalf("exit code = %d, diagnostics = %s", exitCode, diagnostics.String())
	}
	if output.String() != `{"directory":true}` {
		t.Fatalf("body = %q", output.String())
	}
	if diagnostics.Len() != 0 {
		t.Fatalf("silent diagnostics = %q", diagnostics.String())
	}
	if receivedUserAgent != "acme.sh/3.1.5 (https://github.com/acmesh-official/acme.sh)" {
		t.Fatalf("User-Agent = %q", receivedUserAgent)
	}
	if receivedHeaders.Get("Accept") != "application/json" {
		t.Fatalf("Accept = %q", receivedHeaders.Get("Accept"))
	}
	headerContents, err := os.ReadFile(headerPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(headerContents, []byte("HTTP/1.1 200 OK\r\n")) || !bytes.Contains(headerContents, []byte("Replay-Nonce: nonce-value\r\n")) {
		t.Fatalf("header dump = %q", headerContents)
	}
}

func TestDataMethodAndRepeatedHeaders(t *testing.T) {
	type receivedRequest struct {
		method       string
		body         string
		headerValues []string
		contentType  string
	}
	received := make(chan receivedRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		body, _ := io.ReadAll(request.Body)
		received <- receivedRequest{
			method: request.Method, body: string(body), headerValues: request.Header.Values("X-Test"), contentType: request.Header.Get("Content-Type"),
		}
		_, _ = io.WriteString(response, "ok")
	}))
	defer server.Close()

	var output bytes.Buffer
	exitCode := run([]string{
		"-s", "-X", "PUT", "-H", "X-Test: one", "-H", "X-Test: two", "-H", "Content-Type: application/jose+json", "-H", "", "-d", "alpha=one", "-d", "beta=two", server.URL,
	}, strings.NewReader(""), &output, io.Discard)
	if exitCode != exitOK {
		t.Fatalf("exit code = %d", exitCode)
	}
	request := <-received
	if request.method != http.MethodPut || request.body != "alpha=one&beta=two" {
		t.Fatalf("request = %#v", request)
	}
	if strings.Join(request.headerValues, ",") != "one,two" {
		t.Fatalf("X-Test values = %#v", request.headerValues)
	}
	if request.contentType != "application/jose+json" {
		t.Fatalf("Content-Type = %q", request.contentType)
	}
}

func TestRedirectDefaultsAndHeaderChain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/start" {
			response.Header().Set("Location", "/final")
			response.WriteHeader(http.StatusFound)
			_, _ = io.WriteString(response, "redirect")
			return
		}
		response.Header().Set("Replay-Nonce", "final-nonce")
		_, _ = io.WriteString(response, "final")
	}))
	defer server.Close()

	var withoutLocation bytes.Buffer
	if exitCode := run([]string{"-s", server.URL + "/start"}, strings.NewReader(""), &withoutLocation, io.Discard); exitCode != exitOK {
		t.Fatalf("without -L exit code = %d", exitCode)
	}
	if withoutLocation.String() != "redirect" {
		t.Fatalf("without -L body = %q", withoutLocation.String())
	}

	headerPath := filepath.Join(t.TempDir(), "redirect.header")
	var withLocation bytes.Buffer
	if exitCode := run([]string{"-sL", "-D", headerPath, server.URL + "/start"}, strings.NewReader(""), &withLocation, io.Discard); exitCode != exitOK {
		t.Fatalf("with -L exit code = %d", exitCode)
	}
	if withLocation.String() != "final" {
		t.Fatalf("with -L body = %q", withLocation.String())
	}
	headerContents, err := os.ReadFile(headerPath)
	if err != nil {
		t.Fatal(err)
	}
	if count := bytes.Count(headerContents, []byte("HTTP/1.1 ")); count != 2 {
		t.Fatalf("response count = %d, dump = %q", count, headerContents)
	}
}

func TestExplicitMethodIsRetainedAcrossRedirect(t *testing.T) {
	type redirectedRequest struct {
		method      string
		body        string
		contentType string
	}
	received := make(chan redirectedRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/start" {
			response.Header().Set("Location", "/final")
			response.WriteHeader(http.StatusFound)
			return
		}
		requestContents, _ := io.ReadAll(request.Body)
		received <- redirectedRequest{method: request.Method, body: string(requestContents), contentType: request.Header.Get("Content-Type")}
		_, _ = io.WriteString(response, "final")
	}))
	defer server.Close()

	if exitCode := run([]string{
		"-sL", "-X", "POST", "-H", "Content-Type: application/jose+json", "--data", "signed-request", server.URL + "/start",
	}, strings.NewReader(""), io.Discard, io.Discard); exitCode != exitOK {
		t.Fatalf("exit code = %d", exitCode)
	}
	request := <-received
	if request.method != http.MethodPost || request.body != "signed-request" || request.contentType != "application/jose+json" {
		t.Fatalf("redirected request = %#v", request)
	}
}

func TestRedirectLimitUsesCurlExitCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Location", "/loop")
		response.WriteHeader(http.StatusFound)
	}))
	defer server.Close()
	if exitCode := run([]string{"-sL", server.URL + "/loop"}, strings.NewReader(""), io.Discard, io.Discard); exitCode != exitTooManyRedirects {
		t.Fatalf("redirect exit code = %d", exitCode)
	}
}

func TestHeadAndTrace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodHead {
			t.Errorf("method = %s", request.Method)
		}
		response.Header().Set("Replay-Nonce", "head-nonce")
		_, _ = io.WriteString(response, "must-not-be-written")
	}))
	defer server.Close()

	tracePath := filepath.Join(t.TempDir(), "trace.txt")
	var output bytes.Buffer
	if exitCode := run([]string{"-sI", "--trace-ascii", tracePath, server.URL}, strings.NewReader(""), &output, io.Discard); exitCode != exitOK {
		t.Fatalf("exit code = %d", exitCode)
	}
	if !strings.Contains(output.String(), "HTTP/1.1 200 OK\r\n") || !strings.Contains(output.String(), "Replay-Nonce: head-nonce\r\n") {
		t.Fatalf("HEAD output = %q", output.String())
	}
	if strings.Contains(output.String(), "must-not-be-written") {
		t.Fatalf("HEAD output contains body: %q", output.String())
	}
	traceContents, err := os.ReadFile(tracePath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(traceContents, []byte("=> Send header")) || !bytes.Contains(traceContents, []byte("<= Recv header")) {
		t.Fatalf("trace = %q", traceContents)
	}
}

func TestTLSCertificateOptions(t *testing.T) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(response, "secure")
	}))
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	server.StartTLS()
	defer server.Close()

	var output bytes.Buffer
	var diagnostics bytes.Buffer
	if exitCode := run([]string{"-s", server.URL}, strings.NewReader(""), &output, &diagnostics); exitCode != exitCertificate {
		t.Fatalf("embedded roots exit code = %d", exitCode)
	}
	if diagnostics.Len() != 0 {
		t.Fatalf("silent diagnostics = %q", diagnostics.String())
	}

	certificatePEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: server.Certificate().Raw})
	certificateDirectory := t.TempDir()
	certificatePath := filepath.Join(certificateDirectory, "test-ca.pem")
	if err := os.WriteFile(certificatePath, certificatePEM, 0o600); err != nil {
		t.Fatal(err)
	}
	for testName, arguments := range map[string][]string{
		"cacert":   {"-s", "--cacert", certificatePath, server.URL},
		"capath":   {"-s", "--capath", certificateDirectory, server.URL},
		"insecure": {"-sk", server.URL},
	} {
		t.Run(testName, func(t *testing.T) {
			var body bytes.Buffer
			if exitCode := run(arguments, strings.NewReader(""), &body, io.Discard); exitCode != exitOK {
				t.Fatalf("exit code = %d", exitCode)
			}
			if body.String() != "secure" {
				t.Fatalf("body = %q", body.String())
			}
		})
	}

	badCertificatePath := filepath.Join(t.TempDir(), "bad.pem")
	if err := os.WriteFile(badCertificatePath, []byte("not a certificate"), 0o600); err != nil {
		t.Fatal(err)
	}
	if exitCode := run([]string{"-s", "--cacert", badCertificatePath, server.URL}, strings.NewReader(""), io.Discard, io.Discard); exitCode != exitCertificateProblem {
		t.Fatalf("bad CA exit code = %d", exitCode)
	}
}

func TestFailAndMaximumTimeExitCodes(t *testing.T) {
	errorServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(response, "error body")
	}))
	defer errorServer.Close()
	if exitCode := run([]string{"-sf", errorServer.URL}, strings.NewReader(""), io.Discard, io.Discard); exitCode != exitHTTPError {
		t.Fatalf("HTTP error exit code = %d", exitCode)
	}

	timeoutServer := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = io.WriteString(response, "late")
	}))
	defer timeoutServer.Close()
	if exitCode := run([]string{"-s", "--max-time", "0.01", timeoutServer.URL}, strings.NewReader(""), io.Discard, io.Discard); exitCode != exitTimeout {
		t.Fatalf("timeout exit code = %d", exitCode)
	}
}

func TestBodyWriteErrorUsesCurlExitCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(response, "body")
	}))
	defer server.Close()
	if exitCode := run([]string{"-s", server.URL}, strings.NewReader(""), failingWriter{}, io.Discard); exitCode != exitWriteError {
		t.Fatalf("write error exit code = %d", exitCode)
	}
}
