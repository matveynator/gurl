// Matvey Gladkikh is the author and contributors are welcome!
// https://github.com/matveynator/gurl
// You are free to modify, use and distribute this software.
// Distributed under GNU General public license.

package main

import (
	"bytes"          // Handles byte slices for multipart data
	"crypto/tls"     // Provides TLS configuration
	"flag"           // Parses command-line flags
	"fmt"            // Provides formatted I/O
	"io"             // Provides interfaces for I/O primitives
	"mime/multipart" // For creating multipart form data
	"net/http"       // HTTP client for making requests
	"net/url"        // URL parsing and handling
	"os"             // OS-level functions and file handling
	"strings"        // String manipulation functions
	"sync/atomic"    // Atomic operations for safe concurrent updates
	"time"           // Time handling utilities
)

// ProgressReader wraps an io.Reader to report download progress.
// It tracks the number of bytes read and outputs the progress percentage
// when output is redirected to a file or not connected to a terminal.
type ProgressReader struct {
	Reader   io.Reader // The underlying reader
	Total    int64     // Total size of the content (in bytes)
	Progress int64     // Number of bytes read so far
}

// Read reads data into the provided buffer and updates the progress.
func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)                                 // Read from the underlying reader
	atomic.AddInt64(&p.Progress, int64(n))                     // Update progress atomically
	percentage := float64(p.Progress) / float64(p.Total) * 100 // Calculate progress percentage
	fmt.Fprintf(os.Stderr, "\rDownloaded: %d/%d bytes (%.2f%%)", p.Progress, p.Total, percentage)
	return n, err // Return number of bytes read and any error
}

// ensureScheme ensures the URL has a valid scheme (http or https).
// If no scheme is specified, it defaults to "http://".
func ensureScheme(urlStr string) string {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return "http://" + urlStr
	}
	return urlStr
}

// isTerminal checks whether output is to a terminal or redirected to a file.
func isTerminal(fd *os.File) bool {
	stat, err := fd.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// prepareMultipartFormData prepares a multipart form body for file uploads and key=value fields.
func prepareMultipartFormData(formFields []string) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, field := range formFields {
		if strings.Contains(field, "=@") {
			// Handle file field: key=@file
			parts := strings.SplitN(field, "=@", 2)
			if len(parts) != 2 {
				return nil, "", fmt.Errorf("invalid form field: %s", field)
			}
			key, filePath := parts[0], parts[1]
			file, err := os.Open(filePath)
			if err != nil {
				return nil, "", fmt.Errorf("error opening file '%s': %v", filePath, err)
			}
			defer file.Close()

			part, err := writer.CreateFormFile(key, filePath)
			if err != nil {
				return nil, "", fmt.Errorf("error creating form file: %v", err)
			}
			if _, err = io.Copy(part, file); err != nil {
				return nil, "", fmt.Errorf("error writing file to form: %v", err)
			}
		} else {
			// Handle key=value field
			parts := strings.SplitN(field, "=", 2)
			if len(parts) != 2 {
				return nil, "", fmt.Errorf("invalid form field: %s", field)
			}
			key, value := parts[0], parts[1]
			if err := writer.WriteField(key, value); err != nil {
				return nil, "", fmt.Errorf("error writing form field '%s': %v", key, err)
			}
		}
	}

	writer.Close()
	return body, writer.FormDataContentType(), nil
}

// Global variable for version, initialized as "dev".
var version = "dev"

func main() {
	// Define command-line flags
	var (
		flagVersion   = flag.Bool("version", false, "Show version")
		flagTimeout   = flag.Duration("timeout", 30*time.Second, "Request timeout (-m)")
		flagUserAgent = flag.String("useragent", "GURL", "Custom User-Agent header (-A)")
		flagUnsafe    = flag.Bool("unsafe", false, "Disable SSL verification (-k)")
		flagData      = flag.String("data", "", "POST data (-d)")
		flagForm      = flag.String("F", "", "Multipart form data key=value or key=@file")
		flagCookie    = flag.String("cookie", "", "Send cookies (-b)")
		flagHead      = flag.Bool("head", false, "HEAD request (-I)")
		flagHeader    = flag.String("header", "", "Custom header (-H)")
		flagOutput    = flag.String("output", "", "Save output to file (-o)")
		flagLocation  = flag.Bool("location", true, "Follow redirects (-L) (default: true)")
		flagFail      = flag.Bool("fail", false, "Exit with error code on HTTP errors")
		flagRequest   = flag.String("X", "GET", "Specify custom HTTP request method (-X)")
	)

	// Add alternative short flags for compatibility with curl
	flag.BoolVar(flagVersion, "V", false, "Show version")
	flag.DurationVar(flagTimeout, "m", 30*time.Second, "Request timeout")
	flag.StringVar(flagUserAgent, "A", "GURL", "Custom User-Agent header")
	flag.BoolVar(flagUnsafe, "k", false, "Disable SSL verification")
	flag.StringVar(flagData, "d", "", "POST data")
	flag.StringVar(flagCookie, "b", "", "Send cookies")
	flag.BoolVar(flagHead, "I", false, "HEAD request")
	flag.StringVar(flagHeader, "H", "", "Custom header")
	flag.StringVar(flagOutput, "o", "", "Save output to file")
	flag.BoolVar(flagLocation, "L", true, "Follow redirects")
	// Parse the flags
	flag.Parse()

	// Show version and exit if the version flag is set
	if *flagVersion {
		fmt.Println("GURL version", version)
		return
	}

	// Validate arguments: At least one URL is required
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: gurl [options] <url>")
		os.Exit(1)
	}

	// Ensure the URL has a proper scheme (http/https)
	URL := ensureScheme(args[0])
	u, err := url.Parse(URL)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return
	}

	// Configure HTTP transport with optional SSL verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *flagUnsafe},
	}

	// Configure the HTTP client with timeout and default redirect handling
	client := &http.Client{
		Transport: transport,
		Timeout:   *flagTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if *flagLocation {
				fmt.Fprintf(os.Stderr, "Redirected to: %s\n", req.URL)
				return nil
			}
			return http.ErrUseLastResponse
		},
	}

	// Determine the HTTP method and prepare body
	method := *flagRequest
	var body io.Reader = nil
	contentType := ""

	if *flagForm != "" {
		fields := strings.Split(*flagForm, "&")
		body, contentType, err = prepareMultipartFormData(fields)
		if err != nil {
			fmt.Println("Error preparing form data:", err)
			return
		}
		method = http.MethodPost
	} else if *flagData != "" {
		method = http.MethodPost
		body = strings.NewReader(*flagData)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers: User-Agent, custom headers, and cookies
	req.Header.Set("User-Agent", *flagUserAgent)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if *flagHeader != "" {
		parts := strings.SplitN(*flagHeader, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
	if *flagCookie != "" {
		req.Header.Set("Cookie", *flagCookie)
	}

	// Execute the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	// Handle --fail flag: Exit on HTTP error statuses
	if *flagFail && resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "HTTP error: %s\n", resp.Status)
		os.Exit(22) // Exit with code 22 like curl
	}

	// Print response headers if it's a HEAD request
	if *flagHead {
		fmt.Println("Response Headers:")
		fmt.Println("Status:", resp.Status)
		for k, v := range resp.Header {
			fmt.Println(k, ":", v)
		}
		return
	}

	// Check if output is redirected to a file or not a terminal
	useProgress := *flagOutput != "" || !isTerminal(os.Stdout)

	// Determine output destination (stdout or file)
	var writer io.Writer = os.Stdout
	if *flagOutput != "" {
		file, err := os.Create(*flagOutput)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer file.Close()
		writer = file
	}

	// Copy the response body with or without progress tracking
	if useProgress {
		total := resp.ContentLength
		if total <= 0 {
			total = 1 // Prevent division by zero for unknown content lengths
		}
		progressReader := &ProgressReader{Reader: resp.Body, Total: total}
		_, err = io.Copy(writer, progressReader)
	} else {
		_, err = io.Copy(writer, resp.Body)
	}

	// Handle errors during file writing
	if err != nil {
		fmt.Println("Download error:", err)
		return
	}

	// Print a completion message if progress was enabled
	if useProgress {
		fmt.Fprintln(os.Stderr, "\nDownload completed successfully.")
	}
}
