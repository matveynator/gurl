// Matvey Gladkikh is the author and contributors are welcome!
// https://github.com/matveynator/gurl
// You are free to modify, use and distribute this software.
// Distributed under GNU General public license.

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gurl/pkg/certificates"
)

var version = "dev"

const (
	exitOK                 = 0
	exitCommandLine        = 2
	exitMalformedURL       = 3
	exitResolveHost        = 6
	exitConnect            = 7
	exitHTTPError          = 22
	exitWriteError         = 23
	exitReadError          = 26
	exitTimeout            = 28
	exitTLSConnect         = 35
	exitTooManyRedirects   = 47
	exitReceiveError       = 56
	exitCertificate        = 60
	exitCertificateProblem = 77

	defaultConnectTimeout = 300 * time.Second
	maximumRedirects      = 50
)

var errTooManyRedirects = errors.New("maximum redirects followed")

// commandOptions is the complete request model. Parsing finishes before any
// files or network connections are opened, so execution observes one immutable
// set of choices.
type commandOptions struct {
	showHelp       bool
	showVersion    bool
	silent         bool
	dumpHeaderPath string
	tracePath      string
	ipVersion      string
	followLocation bool
	insecure       bool
	head           bool
	userAgent      string
	requestMethod  string
	headers        []string
	requestData    []string
	formFields     []string
	cookie         string
	outputPath     string
	failHTTP       bool
	connectTimeout time.Duration
	maximumTime    time.Duration
	caCertificate  string
	caPath         string
	requestURL     string
}

type commandLineError struct{ message string }

func (e *commandLineError) Error() string { return e.message }

// parseCommandLine accepts the curl spellings used by acme.sh while retaining
// the non-conflicting historical GURL names.
func parseCommandLine(arguments []string) (commandOptions, error) {
	options := commandOptions{userAgent: "GURL/" + version, connectTimeout: defaultConnectTimeout}
	stopOptions := false
	for argumentIndex := 0; argumentIndex < len(arguments); argumentIndex++ {
		argument := arguments[argumentIndex]
		if stopOptions || argument == "-" || !strings.HasPrefix(argument, "-") {
			if options.showHelp {
				continue
			}
			if options.requestURL != "" {
				return commandOptions{}, &commandLineError{message: "only one URL is supported"}
			}
			options.requestURL = argument
			continue
		}
		if argument == "--" {
			stopOptions = true
			continue
		}
		if strings.HasPrefix(argument, "--") {
			optionName, optionValue, hasValue := splitLongOption(argument)
			if optionNeedsValue(optionName) && !hasValue {
				argumentIndex++
				if argumentIndex >= len(arguments) {
					return commandOptions{}, &commandLineError{message: "option --" + optionName + " requires a value"}
				}
				optionValue = arguments[argumentIndex]
				hasValue = true
			}
			if err := applyLongOption(&options, optionName, optionValue, hasValue); err != nil {
				return commandOptions{}, err
			}
			continue
		}
		consumedNext, err := applyShortOptions(&options, argument[1:], arguments, argumentIndex)
		if err != nil {
			return commandOptions{}, err
		}
		if consumedNext {
			argumentIndex++
		}
	}
	if !options.showHelp && !options.showVersion && options.requestURL == "" {
		return commandOptions{}, &commandLineError{message: "no URL specified"}
	}
	return options, nil
}

func splitLongOption(argument string) (string, string, bool) {
	nameAndValue := strings.TrimPrefix(argument, "--")
	separatorIndex := strings.IndexByte(nameAndValue, '=')
	if separatorIndex < 0 {
		return nameAndValue, "", false
	}
	return nameAndValue[:separatorIndex], nameAndValue[separatorIndex+1:], true
}

func optionNeedsValue(optionName string) bool {
	switch optionName {
	case "dump-header", "trace-ascii", "user-agent", "useragent", "request", "header", "data", "connect-timeout", "max-time", "timeout", "cacert", "capath", "output", "form", "cookie":
		return true
	default:
		return false
	}
}

func applyLongOption(options *commandOptions, optionName, optionValue string, hasValue bool) error {
	if !optionNeedsValue(optionName) && hasValue {
		return &commandLineError{message: "option --" + optionName + " does not take a value"}
	}
	switch optionName {
	case "help":
		options.showHelp = true
	case "version":
		options.showVersion = true
	case "silent":
		options.silent = true
	case "dump-header":
		options.dumpHeaderPath = optionValue
	case "ipv4":
		options.ipVersion = "4"
	case "ipv6":
		options.ipVersion = "6"
	case "location":
		options.followLocation = true
	case "trace-ascii":
		options.tracePath = optionValue
	case "capath":
		options.caPath = optionValue
	case "cacert":
		options.caCertificate = optionValue
	case "globoff":
		// GURL never expands URL glob expressions, so globbing is already off.
	case "insecure", "unsafe":
		options.insecure = true
	case "head":
		options.head = true
	case "user-agent", "useragent":
		options.userAgent = optionValue
	case "request":
		options.requestMethod = optionValue
	case "header":
		options.headers = append(options.headers, optionValue)
	case "data":
		options.requestData = append(options.requestData, optionValue)
	case "connect-timeout":
		duration, err := parseSeconds(optionName, optionValue)
		if err != nil {
			return err
		}
		options.connectTimeout = duration
	case "max-time":
		duration, err := parseSeconds(optionName, optionValue)
		if err != nil {
			return err
		}
		options.maximumTime = duration
	case "timeout":
		duration, err := time.ParseDuration(optionValue)
		if err != nil || duration < 0 {
			return &commandLineError{message: "invalid --timeout value: " + optionValue}
		}
		options.maximumTime = duration
	case "output":
		options.outputPath = optionValue
	case "fail":
		options.failHTTP = true
	case "form":
		options.formFields = append(options.formFields, optionValue)
	case "cookie":
		options.cookie = optionValue
	default:
		return &commandLineError{message: "unknown option --" + optionName}
	}
	return nil
}

func applyShortOptions(options *commandOptions, shortOptions string, arguments []string, argumentIndex int) (bool, error) {
	if shortOptions == "" {
		return false, &commandLineError{message: "invalid empty option"}
	}
	for optionIndex := 0; optionIndex < len(shortOptions); optionIndex++ {
		optionName := shortOptions[optionIndex]
		if shortOptionNeedsValue(optionName) {
			optionValue := shortOptions[optionIndex+1:]
			consumedNext := false
			if optionValue == "" {
				if argumentIndex+1 >= len(arguments) {
					return false, &commandLineError{message: fmt.Sprintf("option -%c requires a value", optionName)}
				}
				optionValue = arguments[argumentIndex+1]
				consumedNext = true
			}
			if err := applyShortValue(options, optionName, optionValue); err != nil {
				return false, err
			}
			return consumedNext, nil
		}
		switch optionName {
		case 'h':
			options.showHelp = true
		case 'V':
			options.showVersion = true
		case 's':
			options.silent = true
		case '4':
			options.ipVersion = "4"
		case '6':
			options.ipVersion = "6"
		case 'L':
			options.followLocation = true
		case 'g':
		case 'k':
			options.insecure = true
		case 'I':
			options.head = true
		case 'f':
			options.failHTTP = true
		default:
			return false, &commandLineError{message: fmt.Sprintf("unknown option -%c", optionName)}
		}
	}
	return false, nil
}

func shortOptionNeedsValue(optionName byte) bool {
	switch optionName {
	case 'D', 'A', 'X', 'H', 'd', 'm', 'o', 'F', 'b':
		return true
	default:
		return false
	}
}

func applyShortValue(options *commandOptions, optionName byte, optionValue string) error {
	switch optionName {
	case 'D':
		options.dumpHeaderPath = optionValue
	case 'A':
		options.userAgent = optionValue
	case 'X':
		options.requestMethod = optionValue
	case 'H':
		options.headers = append(options.headers, optionValue)
	case 'd':
		options.requestData = append(options.requestData, optionValue)
	case 'm':
		duration, err := parseSeconds("max-time", optionValue)
		if err != nil {
			return err
		}
		options.maximumTime = duration
	case 'o':
		options.outputPath = optionValue
	case 'F':
		options.formFields = append(options.formFields, optionValue)
	case 'b':
		options.cookie = optionValue
	default:
		return &commandLineError{message: fmt.Sprintf("unknown option -%c", optionName)}
	}
	return nil
}

func parseSeconds(optionName, optionValue string) (time.Duration, error) {
	seconds, err := strconv.ParseFloat(optionValue, 64)
	if err != nil || seconds < 0 || math.IsInf(seconds, 0) || math.IsNaN(seconds) {
		return 0, &commandLineError{message: "invalid --" + optionName + " value: " + optionValue}
	}
	maximumSeconds := float64(math.MaxInt64) / float64(time.Second)
	if seconds > maximumSeconds {
		return 0, &commandLineError{message: "--" + optionName + " value is too large"}
	}
	return time.Duration(seconds * float64(time.Second)), nil
}

// requestBody transfers ownership of the complete request payload to the HTTP
// request. Keeping it replayable lets net/http reproduce curl's redirect flow.
func requestBody(options commandOptions) ([]byte, string, error) {
	if len(options.formFields) > 0 {
		return prepareMultipartFormData(options.formFields)
	}
	if len(options.requestData) > 0 {
		return []byte(strings.Join(options.requestData, "&")), "application/x-www-form-urlencoded", nil
	}
	return nil, "", nil
}

func prepareMultipartFormData(formFields []string) ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, field := range formFields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("invalid form field: %s", field)
		}
		fieldName, fieldContents := parts[0], parts[1]
		if strings.HasPrefix(fieldContents, "@") {
			filePath := strings.TrimPrefix(fieldContents, "@")
			file, err := os.Open(filePath)
			if err != nil {
				return nil, "", fmt.Errorf("open form file %q: %w", filePath, err)
			}
			part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
			if err == nil {
				_, err = io.Copy(part, file)
			}
			closeErr := file.Close()
			if err != nil {
				return nil, "", fmt.Errorf("write form file %q: %w", filePath, err)
			}
			if closeErr != nil {
				return nil, "", fmt.Errorf("close form file %q: %w", filePath, closeErr)
			}
			continue
		}
		if err := writer.WriteField(fieldName, fieldContents); err != nil {
			return nil, "", fmt.Errorf("write form field %q: %w", fieldName, err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("finish multipart form: %w", err)
	}
	return body.Bytes(), writer.FormDataContentType(), nil
}

func requestMethod(options commandOptions) string {
	if options.requestMethod != "" {
		return options.requestMethod
	}
	if options.head {
		return http.MethodHead
	}
	if len(options.requestData) > 0 || len(options.formFields) > 0 {
		return http.MethodPost
	}
	return http.MethodGet
}

func applyRequestHeaders(request *http.Request, options commandOptions, contentType string) error {
	request.Header.Set("User-Agent", options.userAgent)
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if options.cookie != "" {
		request.Header.Set("Cookie", options.cookie)
	}
	customHeaderNames := make(map[string]bool)
	for _, headerLine := range options.headers {
		if headerLine == "" {
			continue
		}
		if strings.HasSuffix(headerLine, ";") && !strings.Contains(headerLine, ":") {
			headerName := strings.TrimSpace(strings.TrimSuffix(headerLine, ";"))
			if headerName == "" {
				return fmt.Errorf("invalid empty header")
			}
			request.Header[http.CanonicalHeaderKey(headerName)] = []string{""}
			continue
		}
		headerParts := strings.SplitN(headerLine, ":", 2)
		if len(headerParts) != 2 || strings.TrimSpace(headerParts[0]) == "" {
			return fmt.Errorf("invalid header: %s", headerLine)
		}
		headerName := strings.TrimSpace(headerParts[0])
		headerValue := strings.TrimSpace(headerParts[1])
		canonicalHeaderName := http.CanonicalHeaderKey(headerName)
		if headerValue == "" {
			if canonicalHeaderName == "User-Agent" {
				request.Header[canonicalHeaderName] = []string{""}
			} else {
				request.Header.Del(canonicalHeaderName)
			}
			customHeaderNames[canonicalHeaderName] = true
			continue
		}
		if strings.EqualFold(headerName, "Host") {
			request.Host = headerValue
			continue
		}
		if customHeaderNames[canonicalHeaderName] {
			request.Header.Add(canonicalHeaderName, headerValue)
		} else {
			request.Header.Set(canonicalHeaderName, headerValue)
			customHeaderNames[canonicalHeaderName] = true
		}
	}
	return nil
}

func ensureHTTPURL(rawURL string) (*url.URL, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return nil, fmt.Errorf("unsupported or incomplete URL")
	}
	return parsedURL, nil
}

func certificatePool(options commandOptions) (*x509.CertPool, error) {
	if options.caCertificate == "" && options.caPath == "" {
		return certificates.Pool()
	}
	pool := x509.NewCertPool()
	loadedCertificates := false
	if options.caCertificate != "" {
		certificatePEM, err := os.ReadFile(options.caCertificate)
		if err != nil {
			return nil, fmt.Errorf("read CA certificate %q: %w", options.caCertificate, err)
		}
		if !pool.AppendCertsFromPEM(certificatePEM) {
			return nil, fmt.Errorf("CA certificate %q contains no certificates", options.caCertificate)
		}
		loadedCertificates = true
	}
	if options.caPath != "" {
		entries, err := os.ReadDir(options.caPath)
		if err != nil {
			return nil, fmt.Errorf("read CA directory %q: %w", options.caPath, err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			certificatePEM, readErr := os.ReadFile(filepath.Join(options.caPath, entry.Name()))
			if readErr == nil && pool.AppendCertsFromPEM(certificatePEM) {
				loadedCertificates = true
			}
		}
	}
	if !loadedCertificates {
		return nil, fmt.Errorf("configured CA source contains no certificates")
	}
	return pool, nil
}

func createTransport(options commandOptions, roots *x509.CertPool) *http.Transport {
	dialer := &net.Dialer{Timeout: options.connectTimeout, KeepAlive: 30 * time.Second}
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(boundaryContext context.Context, network, address string) (net.Conn, error) {
			switch options.ipVersion {
			case "4":
				network = "tcp4"
			case "6":
				network = "tcp6"
			}
			return dialer.DialContext(boundaryContext, network, address)
		},
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: options.insecure, RootCAs: roots, MinVersion: tls.VersionTLS12},
		TLSHandshakeTimeout: options.connectTimeout,
		DisableCompression:  true,
		ForceAttemptHTTP2:   true,
	}
}

type outputFile struct {
	writer io.Writer
	file   *os.File
}

func openOutput(path string, standardOutput io.Writer) (outputFile, error) {
	if path == "" {
		return outputFile{}, nil
	}
	if path == "-" {
		return outputFile{writer: standardOutput}, nil
	}
	file, err := os.Create(path)
	if err != nil {
		return outputFile{}, err
	}
	return outputFile{writer: file, file: file}, nil
}

func (output outputFile) close() error {
	if output.file == nil {
		return nil
	}
	return output.file.Close()
}

type asciiTrace struct{ writer io.Writer }

func (trace asciiTrace) request(request *http.Request) {
	if trace.writer == nil {
		return
	}
	dump, err := httputil.DumpRequestOut(request, true)
	if err != nil {
		fmt.Fprintf(trace.writer, "== Info: cannot dump request: %v\n", err)
		return
	}
	fmt.Fprintf(trace.writer, "=> Send header and data, %d bytes\n%s\n", len(dump), printableTrace(dump))
}

func (trace asciiTrace) response(response *http.Response) {
	if trace.writer == nil {
		return
	}
	headers := responseHeaderBlock(response)
	fmt.Fprintf(trace.writer, "<= Recv header, %d bytes\n%s\n", len(headers), printableTrace(headers))
}

func (trace asciiTrace) failure(err error) {
	if trace.writer != nil {
		fmt.Fprintf(trace.writer, "== Info: %v\n", err)
	}
}

type traceReader struct {
	reader io.Reader
	trace  asciiTrace
}

func (reader *traceReader) Read(buffer []byte) (int, error) {
	readCount, err := reader.reader.Read(buffer)
	if readCount > 0 && reader.trace.writer != nil {
		fmt.Fprintf(reader.trace.writer, "<= Recv data, %d bytes\n%s\n", readCount, printableTrace(buffer[:readCount]))
	}
	return readCount, err
}

func printableTrace(contents []byte) string {
	printable := make([]byte, len(contents))
	for index, character := range contents {
		if character == '\r' || character == '\n' || character == '\t' || (character >= 32 && character < 127) {
			printable[index] = character
		} else {
			printable[index] = '.'
		}
	}
	return string(printable)
}

type progressReader struct {
	reader      io.Reader
	total       int64
	transferred int64
	diagnostics io.Writer
}

type errorTrackingWriter struct {
	writer   io.Writer
	writeErr error
}

func (writer *errorTrackingWriter) Write(contents []byte) (int, error) {
	written, err := writer.writer.Write(contents)
	if err != nil {
		writer.writeErr = err
	}
	return written, err
}

func (reader *progressReader) Read(buffer []byte) (int, error) {
	readCount, err := reader.reader.Read(buffer)
	reader.transferred += int64(readCount)
	if readCount > 0 {
		if reader.total > 0 {
			percentage := float64(reader.transferred) / float64(reader.total) * 100
			fmt.Fprintf(reader.diagnostics, "\r%3.0f%% %d/%d bytes", percentage, reader.transferred, reader.total)
		} else {
			fmt.Fprintf(reader.diagnostics, "\r%d bytes", reader.transferred)
		}
	}
	return readCount, err
}

func responseHeaderBlock(response *http.Response) []byte {
	var headers bytes.Buffer
	fmt.Fprintf(&headers, "%s %s\r\n", response.Proto, response.Status)
	headerNames := make([]string, 0, len(response.Header))
	for headerName := range response.Header {
		headerNames = append(headerNames, headerName)
	}
	sort.Strings(headerNames)
	for _, headerName := range headerNames {
		for _, headerValue := range response.Header[headerName] {
			fmt.Fprintf(&headers, "%s: %s\r\n", headerName, headerValue)
		}
	}
	headers.WriteString("\r\n")
	return headers.Bytes()
}

func writeResponseHeaders(response *http.Response, dumpHeader io.Writer, headOutput io.Writer, trace asciiTrace) error {
	headerBlock := responseHeaderBlock(response)
	if dumpHeader != nil {
		if _, err := dumpHeader.Write(headerBlock); err != nil {
			return err
		}
	}
	if headOutput != nil {
		if _, err := headOutput.Write(headerBlock); err != nil {
			return err
		}
	}
	trace.response(response)
	return nil
}

func classifyRequestError(err error) int {
	if errors.Is(err, errTooManyRedirects) {
		return exitTooManyRedirects
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return exitTimeout
	}
	var networkError net.Error
	if errors.As(err, &networkError) && networkError.Timeout() {
		return exitTimeout
	}
	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		return exitResolveHost
	}
	var verificationError *tls.CertificateVerificationError
	if errors.As(err, &verificationError) {
		return exitCertificate
	}
	var unknownAuthorityError x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthorityError) {
		return exitCertificate
	}
	var hostnameError x509.HostnameError
	if errors.As(err, &hostnameError) {
		return exitCertificate
	}
	var certificateInvalidError x509.CertificateInvalidError
	if errors.As(err, &certificateInvalidError) {
		return exitCertificate
	}
	var recordHeaderError tls.RecordHeaderError
	if errors.As(err, &recordHeaderError) || strings.Contains(strings.ToLower(err.Error()), "tls:") {
		return exitTLSConnect
	}
	return exitConnect
}

func writeCommandError(diagnostics io.Writer, silent bool, exitCode int, err error) {
	if !silent {
		fmt.Fprintf(diagnostics, "gurl: (%d) %v\n", exitCode, err)
	}
}

func printHelp(output io.Writer) {
	fmt.Fprintln(output, "Usage: gurl [options...] <url>")
	fmt.Fprintln(output, "  -s, --silent                    Silent mode")
	fmt.Fprintln(output, "  -D, --dump-header <file>        Write received headers to file")
	fmt.Fprintln(output, "  -4, --ipv4                      Resolve names to IPv4 addresses")
	fmt.Fprintln(output, "  -6, --ipv6                      Resolve names to IPv6 addresses")
	fmt.Fprintln(output, "  -L, --location                  Follow redirects")
	fmt.Fprintln(output, "      --trace-ascii <file>        Write an ASCII network trace")
	fmt.Fprintln(output, "      --capath <dir>              CA directory used to verify the peer")
	fmt.Fprintln(output, "      --cacert <file>             CA file used to verify the peer")
	fmt.Fprintln(output, "  -g, --globoff                   Disable URL globbing")
	fmt.Fprintln(output, "  -k, --insecure                  Allow insecure server connections")
	fmt.Fprintln(output, "  -I, --head                      Show response headers only")
	fmt.Fprintln(output, "  -A, --user-agent <name>         Send User-Agent <name>")
	fmt.Fprintln(output, "  -X, --request <method>          Specify request method")
	fmt.Fprintln(output, "  -H, --header <header>           Pass a custom header; repeatable")
	fmt.Fprintln(output, "  -d, --data <data>               Send HTTP request data")
	fmt.Fprintln(output, "      --connect-timeout <seconds> Maximum connection time")
	fmt.Fprintln(output, "  -m, --max-time <seconds>        Maximum transfer time")
	fmt.Fprintln(output, "  -o, --output <file>             Write body to file")
	fmt.Fprintln(output, "  -f, --fail                      Fail on HTTP status 400 or later")
	fmt.Fprintln(output, "  -V, --version                   Show GURL version")
}

func run(arguments []string, standardInput io.Reader, standardOutput, diagnostics io.Writer) int {
	_ = standardInput
	options, err := parseCommandLine(arguments)
	if err != nil {
		writeCommandError(diagnostics, false, exitCommandLine, err)
		return exitCommandLine
	}
	if options.showHelp {
		printHelp(standardOutput)
		return exitOK
	}
	if options.showVersion {
		fmt.Fprintln(standardOutput, "GURL version", version)
		return exitOK
	}
	parsedURL, err := ensureHTTPURL(options.requestURL)
	if err != nil {
		writeCommandError(diagnostics, options.silent, exitMalformedURL, err)
		return exitMalformedURL
	}
	requestContents, contentType, err := requestBody(options)
	if err != nil {
		writeCommandError(diagnostics, options.silent, exitReadError, err)
		return exitReadError
	}
	roots, err := certificatePool(options)
	if err != nil {
		writeCommandError(diagnostics, options.silent, exitCertificateProblem, err)
		return exitCertificateProblem
	}

	dumpHeaderOutput, err := openOutput(options.dumpHeaderPath, standardOutput)
	if err != nil {
		writeCommandError(diagnostics, options.silent, exitWriteError, err)
		return exitWriteError
	}
	defer dumpHeaderOutput.close()
	traceOutput, err := openOutput(options.tracePath, standardOutput)
	if err != nil {
		writeCommandError(diagnostics, options.silent, exitWriteError, err)
		return exitWriteError
	}
	defer traceOutput.close()
	bodyOutput := outputFile{writer: standardOutput}
	if options.outputPath != "" {
		bodyOutput, err = openOutput(options.outputPath, standardOutput)
		if err != nil {
			writeCommandError(diagnostics, options.silent, exitWriteError, err)
			return exitWriteError
		}
		defer bodyOutput.close()
	}

	request, err := http.NewRequest(requestMethod(options), parsedURL.String(), bytes.NewReader(requestContents))
	if err != nil {
		writeCommandError(diagnostics, options.silent, exitMalformedURL, err)
		return exitMalformedURL
	}
	if err := applyRequestHeaders(request, options, contentType); err != nil {
		writeCommandError(diagnostics, options.silent, exitCommandLine, err)
		return exitCommandLine
	}

	trace := asciiTrace{writer: traceOutput.writer}
	trace.request(request)
	var headOutput io.Writer
	if options.head {
		headOutput = standardOutput
	}
	var headerWriteError error
	transport := createTransport(options, roots)
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: options.maximumTime}
	client.CheckRedirect = func(nextRequest *http.Request, previousRequests []*http.Request) error {
		if !options.followLocation {
			return http.ErrUseLastResponse
		}
		if len(previousRequests) >= maximumRedirects {
			return errTooManyRedirects
		}
		if nextRequest.Response != nil && headerWriteError == nil {
			headerWriteError = writeResponseHeaders(nextRequest.Response, dumpHeaderOutput.writer, headOutput, trace)
		}
		if headerWriteError != nil {
			return headerWriteError
		}
		if options.requestMethod != "" {
			nextRequest.Method = options.requestMethod
			if len(requestContents) > 0 {
				nextRequest.Body = io.NopCloser(bytes.NewReader(requestContents))
				nextRequest.ContentLength = int64(len(requestContents))
				if requestContentType := request.Header.Get("Content-Type"); requestContentType != "" {
					nextRequest.Header.Set("Content-Type", requestContentType)
				}
			}
		}
		trace.request(nextRequest)
		return nil
	}

	response, err := client.Do(request)
	if err != nil {
		trace.failure(err)
		exitCode := classifyRequestError(err)
		if headerWriteError != nil {
			exitCode, err = exitWriteError, headerWriteError
		}
		writeCommandError(diagnostics, options.silent, exitCode, err)
		return exitCode
	}
	defer response.Body.Close()
	if err := writeResponseHeaders(response, dumpHeaderOutput.writer, headOutput, trace); err != nil {
		writeCommandError(diagnostics, options.silent, exitWriteError, err)
		return exitWriteError
	}
	if options.failHTTP && response.StatusCode >= http.StatusBadRequest {
		err := fmt.Errorf("the requested URL returned error: %s", response.Status)
		writeCommandError(diagnostics, options.silent, exitHTTPError, err)
		return exitHTTPError
	}
	if options.head {
		return exitOK
	}

	responseReader := io.Reader(&traceReader{reader: response.Body, trace: trace})
	if !options.silent {
		responseReader = &progressReader{reader: responseReader, total: response.ContentLength, diagnostics: diagnostics}
	}
	trackedBodyOutput := &errorTrackingWriter{writer: bodyOutput.writer}
	if _, err := io.Copy(trackedBodyOutput, responseReader); err != nil {
		exitCode := exitReceiveError
		if trackedBodyOutput.writeErr != nil {
			exitCode = exitWriteError
		}
		writeCommandError(diagnostics, options.silent, exitCode, err)
		return exitCode
	}
	if !options.silent {
		fmt.Fprintln(diagnostics)
	}
	return exitOK
}

func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }
