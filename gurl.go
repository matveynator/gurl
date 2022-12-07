// Matvey Gladkikh is the author and contributors are welcome! 
// https://github.com/matveynator/gurl
// You are free to modify, use and distribute this software.
// Distributed under GNU General public license.

package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"time"
	"net/http"
	"net/url"
	"crypto/tls"
)

var VERSION string

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {

	var client *http.Client
	tr := &http.Transport{}

	flagVersion := flag.Bool("version", false, "Output version information.")
	TIMEOUT := flag.Duration("timeout", 30*time.Second, "Set connect and operation timeout. Valid time units are: ns, us or µs, ms, s, m, h.")
	USERAGENT := flag.String("useragent", "GURL (https://github.com/matveynator/gurl)", "Set user agent.")
	LANG := flag.String("lang","", "Set Accept-Language header, for example: -lang en-US")
	PROXY := flag.String("proxy", "", "Set http/https/socks5 proxy 'type://host:port', example: -proxy 'socks5://127.0.0.1:3128' -proxy 'http://127.0.0.1:8080'")
	_ = flag.Bool("unsafe", false, "Disable strict TLS certificate checking.")
	_ = flag.Bool("head", false, "Perform HEAD request.")

	//process all flags
	flag.Parse()

	//set default transport pasrameters:
	tr.MaxIdleConns = 10
	tr.IdleConnTimeout = *TIMEOUT

	//show version
	if *flagVersion  {
		if VERSION != "" {
			fmt.Println("Version:", VERSION)
		} else {
			fmt.Println("Version: unknown")
		}
		os.Exit(0)
	}

	URL := flag.Arg(0)
	//validate URL:
	u, err := url.Parse(URL)       
	if err != nil {
		fmt.Println(" Error=%+v URL=%+v\n", err, u)
		os.Exit(1)
	} 

	//validate target
	//need better validation of POST values to check errors.
	if u.Scheme != "http" && u.Scheme != "https" {
		fmt.Println("Error: target URL scheme should be http:// or https://. You provided: \"",URL,"\"")
		os.Exit(1)
	} 

	//set proxy if needed:
	if *PROXY != "" {
		proxyUrl, err := url.Parse(*PROXY)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		tr.Proxy = http.ProxyURL(proxyUrl)
	} else {
		tr.Proxy = http.ProxyFromEnvironment
	}

	//set skip SSL checks:
	if isFlagPassed("unsafe")  {
	  tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: false}
	} else {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	//create http client:
	client = &http.Client{Transport: tr}

	//process HEAD request:
	if isFlagPassed("head")  {
		req, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if isFlagPassed("lang") && (*LANG != "")  {
			req.Header.Add("Accept-Language", *LANG)
		}
		req.Header.Set("User-Agent", string("'") + *USERAGENT + string("'"))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		fmt.Println("Status:", resp.StatusCode)
		for k, v := range resp.Header {
			fmt.Print(k)
			fmt.Print(" : ")
			fmt.Println(v)
		}

		resp.Body.Close()
		os.Exit(0)
	} else {
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if isFlagPassed("lang") && (*LANG != "")  {
			req.Header.Add("Accept-Language", *LANG)
		}
		req.Header.Set("User-Agent", string('"') + *USERAGENT + string('"'))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		resp.Body.Close()
		os.Exit(0)

	}
}
