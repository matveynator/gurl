// GURL is a curl-like command-line HTTP client distributed as a standalone Go binary.
package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultVersion = "dev"

// version is replaced at build time. It is never mutated while the program runs.
var version = defaultVersion

type stringFlags []string

func (flags *stringFlags) String() string {
	return strings.Join(*flags, ", ")
}

func (flags *stringFlags) Set(flagValue string) error {
	*flags = append(*flags, flagValue)
	return nil
}

type options struct {
	URL                 string
	Method              string
	Timeout             time.Duration
	UserAgent           string
	SkipTLSVerification bool
	RequestBody         string
	FormFields          []string
	Cookie              string
	Head                bool
	Headers             []string
	OutputPath          string
	FollowRedirects     bool
	FailOnHTTPError     bool
	Silent              bool
	Verbose             bool
	ShowVersion         bool
}

type streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func main() {
	os.Exit(run(os.Args[1:], streams{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}))
}

func run(arguments []string, processStreams streams) int {
	commandOptions, err := parseOptions(arguments, processStreams.Stderr)
	if err != nil {
		return 2
	}
	if commandOptions.ShowVersion {
		fmt.Fprintf(processStreams.Stdout, "gurl %s\n", version)
		return 0
	}
	if commandOptions.URL == "" {
		fmt.Fprintln(processStreams.Stderr, "usage: gurl [options] <url>")
		return 2
	}

	if err := execute(commandOptions, processStreams); err != nil {
		var statusError *httpStatusError
		if errors.As(err, &statusError) {
			fmt.Fprintln(processStreams.Stderr, statusError)
			return 22
		}
		fmt.Fprintf(processStreams.Stderr, "gurl: %v\n", err)
		return 1
	}
	return 0
}

func parseOptions(arguments []string, errorOutput io.Writer) (options, error) {
	commandOptions := options{}
	flagSet := flag.NewFlagSet("gurl", flag.ContinueOnError)
	flagSet.SetOutput(errorOutput)

	flagSet.BoolVar(&commandOptions.ShowVersion, "version", false, "show version")
	flagSet.BoolVar(&commandOptions.ShowVersion, "V", false, "show version")
	flagSet.DurationVar(&commandOptions.Timeout, "timeout", 30*time.Second, "request timeout")
	flagSet.DurationVar(&commandOptions.Timeout, "m", 30*time.Second, "request timeout")
	flagSet.StringVar(&commandOptions.UserAgent, "useragent", "gurl/"+version, "User-Agent header")
	flagSet.StringVar(&commandOptions.UserAgent, "A", "gurl/"+version, "User-Agent header")
	flagSet.BoolVar(&commandOptions.SkipTLSVerification, "unsafe", false, "disable TLS certificate verification")
	flagSet.BoolVar(&commandOptions.SkipTLSVerification, "k", false, "disable TLS certificate verification")
	flagSet.StringVar(&commandOptions.RequestBody, "data", "", "request body")
	flagSet.StringVar(&commandOptions.RequestBody, "d", "", "request body")
	flagSet.Var((*stringFlags)(&commandOptions.FormFields), "form", "multipart field name=value or name=@file")
	flagSet.Var((*stringFlags)(&commandOptions.FormFields), "F", "multipart field name=value or name=@file")
	flagSet.StringVar(&commandOptions.Cookie, "cookie", "", "Cookie header")
	flagSet.StringVar(&commandOptions.Cookie, "b", "", "Cookie header")
	flagSet.BoolVar(&commandOptions.Head, "head", false, "send a HEAD request")
	flagSet.BoolVar(&commandOptions.Head, "I", false, "send a HEAD request")
	flagSet.Var((*stringFlags)(&commandOptions.Headers), "header", "request header")
	flagSet.Var((*stringFlags)(&commandOptions.Headers), "H", "request header")
	flagSet.StringVar(&commandOptions.OutputPath, "output", "", "write response body to a file")
	flagSet.StringVar(&commandOptions.OutputPath, "o", "", "write response body to a file")
	flagSet.BoolVar(&commandOptions.FollowRedirects, "location", false, "follow redirects")
	flagSet.BoolVar(&commandOptions.FollowRedirects, "L", false, "follow redirects")
	flagSet.BoolVar(&commandOptions.FailOnHTTPError, "fail", false, "fail on HTTP status 400 or greater")
	flagSet.StringVar(&commandOptions.Method, "request", "", "HTTP request method")
	flagSet.StringVar(&commandOptions.Method, "X", "", "HTTP request method")
	flagSet.BoolVar(&commandOptions.Silent, "silent", false, "suppress progress output")
	flagSet.BoolVar(&commandOptions.Silent, "s", false, "suppress progress output")
	flagSet.BoolVar(&commandOptions.Verbose, "verbose", false, "print request and response details")
	flagSet.BoolVar(&commandOptions.Verbose, "v", false, "print request and response details")

	if err := flagSet.Parse(arguments); err != nil {
		return options{}, err
	}
	if flagSet.NArg() > 0 {
		commandOptions.URL = flagSet.Arg(0)
	}
	if flagSet.NArg() > 1 {
		fmt.Fprintln(errorOutput, "gurl: only one URL may be requested at a time")
		return options{}, errors.New("too many URLs")
	}
	return commandOptions, nil
}

func execute(commandOptions options, processStreams streams) error {
	request, err := buildRequest(commandOptions)
	if err != nil {
		return err
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if commandOptions.SkipTLSVerification {
		// This explicit compatibility option is equivalent to curl -k and is never enabled by default.
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402
	}
	client := &http.Client{Transport: transport, Timeout: commandOptions.Timeout}
	if !commandOptions.FollowRedirects {
		client.CheckRedirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	if commandOptions.Verbose {
		fmt.Fprintf(processStreams.Stderr, "> %s %s\n", request.Method, request.URL)
	}
	// Reaching the operator-selected URL is the purpose of this command-line client.
	response, err := client.Do(request) // #nosec G704
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	if commandOptions.Verbose {
		fmt.Fprintf(processStreams.Stderr, "< %s\n", response.Status)
	}
	if commandOptions.FailOnHTTPError && response.StatusCode >= http.StatusBadRequest {
		return &httpStatusError{Status: response.Status}
	}
	if commandOptions.Head {
		return writeHeaders(processStreams.Stdout, response)
	}
	return writeResponse(commandOptions, processStreams, response)
}

func buildRequest(commandOptions options) (*http.Request, error) {
	requestURL, err := normalizeURL(commandOptions.URL)
	if err != nil {
		return nil, err
	}

	method := commandOptions.Method
	if method == "" {
		method = http.MethodGet
	}
	var requestBody io.Reader
	contentType := ""
	if len(commandOptions.FormFields) > 0 {
		requestBody, contentType, err = prepareMultipartForm(commandOptions.FormFields)
		method = http.MethodPost
	} else if commandOptions.RequestBody != "" {
		requestBody = strings.NewReader(commandOptions.RequestBody)
		if commandOptions.Method == "" {
			method = http.MethodPost
		}
	}
	if err != nil {
		return nil, err
	}
	if commandOptions.Head {
		method = http.MethodHead
	}

	// A URL supplied by the operator is the intended network destination of this CLI.
	request, err := http.NewRequest(method, requestURL, requestBody) // #nosec G107
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("User-Agent", commandOptions.UserAgent)
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if commandOptions.Cookie != "" {
		request.Header.Set("Cookie", commandOptions.Cookie)
	}
	for _, header := range commandOptions.Headers {
		name, headerValue, found := strings.Cut(header, ":")
		if !found || strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("invalid header %q: expected name:value", header)
		}
		request.Header.Add(strings.TrimSpace(name), strings.TrimSpace(headerValue))
	}
	return request, nil
}

func normalizeURL(rawURL string) (string, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme %q", parsedURL.Scheme)
	}
	if parsedURL.Host == "" {
		return "", errors.New("URL has no host")
	}
	return parsedURL.String(), nil
}

func prepareMultipartForm(formFields []string) (io.Reader, string, error) {
	requestBody := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(requestBody)
	for _, field := range formFields {
		name, fieldValue, found := strings.Cut(field, "=")
		if !found || name == "" {
			return nil, "", fmt.Errorf("invalid form field %q: expected name=value", field)
		}
		if strings.HasPrefix(fieldValue, "@") {
			if err := addMultipartFile(multipartWriter, name, strings.TrimPrefix(fieldValue, "@")); err != nil {
				return nil, "", err
			}
			continue
		}
		if err := multipartWriter.WriteField(name, fieldValue); err != nil {
			return nil, "", fmt.Errorf("write form field %q: %w", name, err)
		}
	}
	if err := multipartWriter.Close(); err != nil {
		return nil, "", fmt.Errorf("finish multipart request: %w", err)
	}
	return requestBody, multipartWriter.FormDataContentType(), nil
}

func addMultipartFile(multipartWriter *multipart.Writer, name string, path string) error {
	// Multipart upload paths are explicitly supplied by the local operator.
	file, err := os.Open(path) // #nosec G304
	if err != nil {
		return fmt.Errorf("open form file %q: %w", path, err)
	}
	defer file.Close()

	filePart, err := multipartWriter.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return fmt.Errorf("create multipart file %q: %w", name, err)
	}
	if _, err := io.Copy(filePart, file); err != nil {
		return fmt.Errorf("copy form file %q: %w", path, err)
	}
	return nil
}

func writeHeaders(output io.Writer, response *http.Response) error {
	if _, err := fmt.Fprintln(output, response.Status); err != nil {
		return err
	}
	return response.Header.Write(output)
}

func writeResponse(commandOptions options, processStreams streams, response *http.Response) error {
	output := processStreams.Stdout
	var outputFile *os.File
	if commandOptions.OutputPath != "" {
		// The output path is explicitly supplied by the local operator.
		file, err := os.Create(commandOptions.OutputPath) // #nosec G304
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		outputFile = file
		output = file
		defer outputFile.Close()
	}

	written, err := io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	if commandOptions.OutputPath != "" && !commandOptions.Silent {
		fmt.Fprintf(processStreams.Stderr, "downloaded %d bytes to %s\n", written, commandOptions.OutputPath)
	}
	return nil
}

type httpStatusError struct {
	Status string
}

func (statusError *httpStatusError) Error() string {
	return "HTTP error: " + statusError.Status
}
