package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		expected string
		wantErr  bool
	}{
		{name: "default HTTP", rawURL: "example.com/path", expected: "http://example.com/path"},
		{name: "HTTPS", rawURL: "https://example.com", expected: "https://example.com"},
		{name: "unsupported scheme", rawURL: "file:///etc/passwd", wantErr: true},
		{name: "missing host", rawURL: "https://", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := normalizeURL(test.rawURL)
			if test.wantErr && err == nil {
				t.Fatal("expected an error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("normalizeURL returned an error: %v", err)
			}
			if actual != test.expected {
				t.Fatalf("normalizeURL() = %q, want %q", actual, test.expected)
			}
		})
	}
}

func TestParseOptions(t *testing.T) {
	errorOutput := &bytes.Buffer{}
	commandOptions, err := parseOptions([]string{
		"-L", "-s", "-v", "-k", "-m", "2s", "-A", "test-agent",
		"-H", "X-One: first", "-H", "X-Two: second", "-F", "name=value",
		"-b", "session=one", "-o", "response.txt", "-X", "PATCH", "example.com",
	}, errorOutput)
	if err != nil {
		t.Fatalf("parseOptions returned an error: %v", err)
	}
	if commandOptions.URL != "example.com" || commandOptions.Timeout != 2*time.Second {
		t.Fatalf("unexpected URL or timeout: %+v", commandOptions)
	}
	if !commandOptions.FollowRedirects || !commandOptions.Silent || !commandOptions.Verbose || !commandOptions.SkipTLSVerification {
		t.Fatalf("boolean flags were not parsed: %+v", commandOptions)
	}
	if commandOptions.Method != "PATCH" || commandOptions.UserAgent != "test-agent" {
		t.Fatalf("request flags were not parsed: %+v", commandOptions)
	}
	if len(commandOptions.Headers) != 2 || len(commandOptions.FormFields) != 1 {
		t.Fatalf("repeatable flags were not parsed: %+v", commandOptions)
	}
	if commandOptions.Cookie != "session=one" || commandOptions.OutputPath != "response.txt" {
		t.Fatalf("output flags were not parsed: %+v", commandOptions)
	}
}

func TestParseOptionsRejectsInvalidInput(t *testing.T) {
	tests := [][]string{{"-unknown"}, {"one.example", "two.example"}}
	for _, arguments := range tests {
		if _, err := parseOptions(arguments, io.Discard); err == nil {
			t.Fatalf("parseOptions(%v) accepted invalid input", arguments)
		}
	}
}

func TestBuildRequest(t *testing.T) {
	commandOptions := options{
		URL:         "example.com/resource",
		RequestBody: "hello",
		UserAgent:   "agent",
		Cookie:      "session=one",
		Headers:     []string{"X-Test: first", "X-Test: second"},
	}
	request, err := buildRequest(commandOptions)
	if err != nil {
		t.Fatalf("buildRequest returned an error: %v", err)
	}
	if request.Method != http.MethodPost || request.URL.String() != "http://example.com/resource" {
		t.Fatalf("unexpected request: %s %s", request.Method, request.URL)
	}
	if request.UserAgent() != "agent" || request.Header.Get("Cookie") != "session=one" {
		t.Fatalf("unexpected standard headers: %v", request.Header)
	}
	if len(request.Header.Values("X-Test")) != 2 {
		t.Fatalf("expected repeated headers, got %v", request.Header.Values("X-Test"))
	}
	requestBody, err := io.ReadAll(request.Body)
	if err != nil || string(requestBody) != "hello" {
		t.Fatalf("unexpected body %q: %v", requestBody, err)
	}
}

func TestBuildRequestRejectsInvalidHeader(t *testing.T) {
	_, err := buildRequest(options{URL: "example.com", Headers: []string{"invalid"}})
	if err == nil {
		t.Fatal("expected invalid header error")
	}
}

func TestBuildHeadRequest(t *testing.T) {
	request, err := buildRequest(options{URL: "example.com", Head: true, Method: http.MethodPost})
	if err != nil {
		t.Fatalf("buildRequest returned an error: %v", err)
	}
	if request.Method != http.MethodHead {
		t.Fatalf("method = %s, want HEAD", request.Method)
	}
}

func TestPrepareMultipartForm(t *testing.T) {
	temporaryDirectory := t.TempDir()
	filePath := filepath.Join(temporaryDirectory, "upload.txt")
	if err := os.WriteFile(filePath, []byte("file body"), 0o600); err != nil {
		t.Fatal(err)
	}

	requestBody, contentType, err := prepareMultipartForm([]string{"field=value", "upload=@" + filePath})
	if err != nil {
		t.Fatalf("prepareMultipartForm returned an error: %v", err)
	}
	mediaTypeBoundary := strings.TrimPrefix(contentType, "multipart/form-data; boundary=")
	reader := multipart.NewReader(requestBody, mediaTypeBoundary)
	form, err := reader.ReadForm(1024)
	if err != nil {
		t.Fatalf("read multipart form: %v", err)
	}
	defer form.RemoveAll()
	if form.Value["field"][0] != "value" {
		t.Fatalf("unexpected field value: %v", form.Value)
	}
	upload, err := form.File["upload"][0].Open()
	if err != nil {
		t.Fatal(err)
	}
	defer upload.Close()
	uploadBody, _ := io.ReadAll(upload)
	if string(uploadBody) != "file body" {
		t.Fatalf("unexpected uploaded body %q", uploadBody)
	}
}

func TestPrepareMultipartFormErrors(t *testing.T) {
	tests := [][]string{{"invalid"}, {"upload=@missing-file"}}
	for _, fields := range tests {
		if _, _, err := prepareMultipartForm(fields); err == nil {
			t.Fatalf("prepareMultipartForm(%v) accepted invalid input", fields)
		}
	}
}

func TestRunVersionAndUsage(t *testing.T) {
	standardOutput := &bytes.Buffer{}
	standardError := &bytes.Buffer{}
	processStreams := streams{Stdin: strings.NewReader(""), Stdout: standardOutput, Stderr: standardError}
	if exitCode := run([]string{"-version"}, processStreams); exitCode != 0 {
		t.Fatalf("version exit code = %d", exitCode)
	}
	if !strings.Contains(standardOutput.String(), version) {
		t.Fatalf("version output = %q", standardOutput)
	}
	standardOutput.Reset()
	if exitCode := run(nil, processStreams); exitCode != 2 {
		t.Fatalf("usage exit code = %d", exitCode)
	}
	if exitCode := run([]string{"-unknown"}, processStreams); exitCode != 2 {
		t.Fatalf("flag error exit code = %d", exitCode)
	}
}

func TestRunGETAndOutputFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.UserAgent() != "test-agent" || request.Header.Get("X-Test") != "yes" {
			t.Errorf("unexpected request headers: %v", request.Header)
		}
		responseWriter.Header().Set("X-Response", "present")
		_, _ = io.WriteString(responseWriter, "response body")
	}))
	defer server.Close()

	outputPath := filepath.Join(t.TempDir(), "response.txt")
	standardOutput := &bytes.Buffer{}
	standardError := &bytes.Buffer{}
	exitCode := run([]string{"-A", "test-agent", "-H", "X-Test: yes", "-o", outputPath, server.URL}, streams{
		Stdin: strings.NewReader(""), Stdout: standardOutput, Stderr: standardError,
	})
	if exitCode != 0 {
		t.Fatalf("run exit code = %d, stderr = %s", exitCode, standardError)
	}
	responseBody, err := os.ReadFile(outputPath)
	if err != nil || string(responseBody) != "response body" {
		t.Fatalf("output file = %q, err = %v", responseBody, err)
	}
	if !strings.Contains(standardError.String(), "downloaded 13 bytes") {
		t.Fatalf("progress output = %q", standardError)
	}
}

func TestRunHeadVerboseAndRedirect(t *testing.T) {
	finalServer := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("X-Final", "yes")
		responseWriter.WriteHeader(http.StatusNoContent)
	}))
	defer finalServer.Close()
	redirectServer := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Redirect(responseWriter, request, finalServer.URL, http.StatusFound)
	}))
	defer redirectServer.Close()

	standardOutput := &bytes.Buffer{}
	standardError := &bytes.Buffer{}
	exitCode := run([]string{"-L", "-I", "-v", redirectServer.URL}, streams{
		Stdin: strings.NewReader(""), Stdout: standardOutput, Stderr: standardError,
	})
	if exitCode != 0 {
		t.Fatalf("run exit code = %d, stderr = %s", exitCode, standardError)
	}
	if !strings.Contains(standardOutput.String(), "204 No Content") || !strings.Contains(standardOutput.String(), "X-Final: yes") {
		t.Fatalf("HEAD output = %q", standardOutput)
	}
	if !strings.Contains(standardError.String(), "> HEAD") || !strings.Contains(standardError.String(), "< 204 No Content") {
		t.Fatalf("verbose output = %q", standardError)
	}
}

func TestRunFailAndRequestFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, "missing", http.StatusNotFound)
	}))
	serverURL := server.URL
	server.Close()

	standardError := &bytes.Buffer{}
	processStreams := streams{Stdin: strings.NewReader(""), Stdout: io.Discard, Stderr: standardError}
	server = httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		http.Error(responseWriter, "missing", http.StatusNotFound)
	}))
	if exitCode := run([]string{"-fail", server.URL}, processStreams); exitCode != 22 {
		t.Fatalf("HTTP failure exit code = %d", exitCode)
	}
	server.Close()
	if exitCode := run([]string{serverURL}, processStreams); exitCode != 1 {
		t.Fatalf("network failure exit code = %d", exitCode)
	}
}

func TestRunTLSUnsafe(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(responseWriter, "secure")
	}))
	defer server.Close()

	standardOutput := &bytes.Buffer{}
	standardError := &bytes.Buffer{}
	exitCode := run([]string{"-k", server.URL}, streams{
		Stdin: strings.NewReader(""), Stdout: standardOutput, Stderr: standardError,
	})
	if exitCode != 0 || standardOutput.String() != "secure" {
		t.Fatalf("TLS request exit code = %d, stdout = %q, stderr = %q", exitCode, standardOutput, standardError)
	}
}

func TestWriteResponseErrors(t *testing.T) {
	response := &http.Response{Body: io.NopCloser(strings.NewReader("body"))}
	err := writeResponse(options{OutputPath: t.TempDir()}, streams{Stdout: io.Discard, Stderr: io.Discard}, response)
	if err == nil {
		t.Fatal("expected output file error")
	}
}
