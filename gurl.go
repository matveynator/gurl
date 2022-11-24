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

var VERSION string
//var TIMEOUT int
//var UNSAFE,HEAD bool

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

  flagVersion := flag.Bool("version", false, "Output version information.")
	TIMEOUT := flag.Int("timeout", 60, "Set connect and operation timeout.")
	USERAGENT := flag.String("useragent", "GURL (https://github.com/matveynator/gurl)", "Set user agent.")
	LANG := flag.String("lang","en-us", "Set Accept-Language header")
	PROXY := flag.String("proxy", "", "Set http proxy 'host:port', example: -proxy '127.0.0.1:8080'")
	UNSAFE := flag.Bool("unsafe", false, "Disable strict certificate checking")
	HEAD := flag.Bool("head", false, "Perform HEAD request.")
	POST := flag.String("post", "", "Perform POST request, example: -post \"'name1':'value1','name2':'value2'\" http://matveynator.ru ")

  //process all flags
  flag.Parse()

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
		if isFlagPassed("post") {
			fmt.Println("Error: empty or malformed POST request.")
		} else {
			fmt.Println("Error: target URL scheme should be http:// or https://. You provided: \"",URL,"\"")
		}
		os.Exit(1)
	} 

  //configure client
	httpclient.Defaults(httpclient.Map{
		"opt_useragent":   *USERAGENT,
		"opt_timeout":     *TIMEOUT,
		"opt_connecttimeout": *TIMEOUT,
		"Accept-Encoding": "gzip,deflate,sdch",
		"Accept-Language": *LANG,
	})

  //set proxy if needed:
	if *PROXY != "" {
		httpclient.Defaults(httpclient.Map{
			httpclient.OPT_PROXY:   *PROXY,
		})
	}

  //set skip SSL checks:
	if *UNSAFE == true {
		httpclient.Defaults(httpclient.Map{
			httpclient.OPT_UNSAFE_TLS:   true,
		})
	}

  //process HEAD request:
	if isFlagPassed("head")  {
		res, err := httpclient.Head(URL)
		if err != nil {
			fmt.Println("Error: ", *HEAD, err)
			os.Exit(1)
		}
		fmt.Printf("%#q\n", res)
		res.Body.Close()
		os.Exit(0)
	} 

  //process POST request:
	if isFlagPassed("post") &&  *POST != "" {
		res, err := httpclient.Post(URL, "map[string]string{" + *POST + "}")

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
		res.Body.Close()
		os.Exit(0)
	} 
  
	//finally process GET request:
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
	res.Body.Close()
	os.Exit(0)
}
