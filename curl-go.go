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

  USERAGENT = "curl-go"
  TIMEOUT = 5
  CONNECT_TIMEOUT = 5
  LANG = "en-us"

  flagVersion := flag.Bool("version", false, "Output version information")
  flagTimeout := flag.Int("timeout", 5, "Set connect and operation timeout")
  flagUserAgent := flag.String("useragent", "curl-go", "Set user agent")
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

  if *flagTimeout != 5 {
    TIMEOUT = *flagTimeout
  }

  if *flagUserAgent != "curl-go" {
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

  res, _ := httpclient.
  WithHeader("Accept-Language", LANG).
  Get(URL)

  fmt.Println(res.ToString())
}
