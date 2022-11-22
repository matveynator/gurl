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

var LANG,VERSION,USERAGENT,URL string
var TIMEOUT,CONNECT_TIMEOUT int


func init() {

  USERAGENT = "gurl"
  TIMEOUT = 15
  CONNECT_TIMEOUT = 15
  LANG = "en-us"

  flagVersion := flag.Bool("version", false, "Output version information")
  flagTimeout := flag.Int("timeout", 15, "Set connect and operation timeout")
  flagUserAgent := flag.String("useragent", "gurl", "Set user agent")
  flagLang := flag.String("lang","en-us", "Set Accept-Language header")
  flag.Parse()

  if *flagVersion  {
    if VERSION != "" {
      fmt.Println("Version:", VERSION)
    } else {
      fmt.Println("Version: unknown")
    }
    os.Exit(0)
  }

  if *flagTimeout != 15 {
    TIMEOUT = *flagTimeout
  }

  if *flagUserAgent != "gurl" {
    USERAGENT = *flagUserAgent
  }

  if *flagLang != "en-us" {
    LANG = *flagLang
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

}

func main() {
  
  httpclient.Defaults(httpclient.Map{
    "opt_useragent":   USERAGENT,
    "opt_timeout":     TIMEOUT,
    "opt_connecttimeout": CONNECT_TIMEOUT,
    "Accept-Encoding": "gzip,deflate,sdch",
  })

  res, err := httpclient.
  WithHeader("Accept-Language", LANG).
  Get(URL)

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
