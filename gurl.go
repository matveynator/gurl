// Matvey Gladkikh is the author and contributors are welcome! 
// https://github.com/matveynator/gurl
// You are free to modify, use and distribute this software.
// Distributed under GNU General public license.

package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/ddliu/go-httpclient"
	"net/url"
)

var LANG,VERSION,USERAGENT,URL,PROXY string
var TIMEOUT int
var UNSAFE,HEAD,POST bool

func main() {
	flagVersion := flag.Bool("version", false, "Output version information")
	TIMEOUT = *flag.Int("timeout", 15, "Set connect and operation timeout")
	USERAGENT := *flag.String("useragent", "gurl", "Set user agent")
	LANG := *flag.String("lang","en-us", "Set Accept-Language header")
	PROXY := *flag.String("proxy", "", "Set http proxy 'host:port' eg: '127.0.0.1:8080'")
	UNSAFE = *flag.Bool("unsafe", false, "Disable TLS certificate checking.")
	HEAD := *flag.Bool("head", false, "Perform HEAD request.")
	POST := *flag.String("post", "", "Perform POST request: -post 'name1':'value1','name2':'value2' ")
	flag.Parse()

	if *flagVersion  {
		if VERSION != "" {
			fmt.Println("Version:", VERSION)
		} else {
			fmt.Println("Version: unknown")
		}
		os.Exit(0)
	}

	URL = flag.Arg(0)
	//validate URL:
	u, err := url.Parse(URL)       
	if err != nil {
		fmt.Println(" Error=%+v URL=%+v\n", err, u)
		os.Exit(1)
	} 

	if u.Scheme != "http" && u.Scheme != "https" {
		fmt.Println("Error: url scheme should be http:// or https:// - ", URL)
		os.Exit(1)
	}

	httpclient.Defaults(httpclient.Map{
		"opt_useragent":   USERAGENT,
		"opt_timeout":     TIMEOUT,
		"opt_connecttimeout": TIMEOUT,
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": LANG,
	})

	if PROXY != "" {
		httpclient.Defaults(httpclient.Map{
			httpclient.OPT_PROXY:   PROXY,
		})
	}

	if UNSAFE == true {
		httpclient.Defaults(httpclient.Map{
			httpclient.OPT_UNSAFE_TLS:   true,
		})
	}

	if HEAD == true {
		res, err := httpclient.Head(URL)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		body, err := res.ToString()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println(body)
	} else if POST != "" {
		res, err := httpclient.Post(URL, "map[string]string{" + POST + "}")
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		body, err := res.ToString()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println(body)
	} else {
		res, err := httpclient.Get(URL)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		body, err := res.ToString()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println(body)
	}
}
