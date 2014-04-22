package main

import (
    "./flag"
    "fmt"
    "os"
    "log"
    "bytes"
    "errors"
    "net/http"
    "io/ioutil"
    "compress/gzip"
    "strconv"
)

var target string
var port string
var showUri bool
var showTargetUri bool
var showRequestBody bool
var showTargetRequestBody bool
var showRequestHeaders bool
var showTargetRequestHeaders bool
var showResponse bool
var showFinalResponse bool
var showResponseHeaders bool
var showFinalResponseHeaders bool
var unzip bool
var useCache bool
var cache = make(map[string] cacheElement)

type cacheElement struct {
  resp *http.Response;
  body []byte;
}

func init() {
  flag.StringVar(&target, "target", "", ", -t\n\t\tTarget URL for single-host everse proxy.  Leave empty for traditional proxy.")
  flag.StringVar(&target, "t", "", "")
  flag.StringVar(&port, "port", "8080", ", -p\n\t\tPort to listen to.")
  flag.StringVar(&port, "p", "8080", "")
  flag.BoolVar(&showUri, "show-request", true, ", -r\n\t\tPrint every request URI.")
  flag.BoolVar(&showUri, "r", true, "")
  flag.BoolVar(&showTargetUri, "show-rewritten-request", false, ", -r2\n\t\tPrint every rewritten (modified by this proxy) request URI.")
  flag.BoolVar(&showTargetUri, "r2", false, "")
  flag.BoolVar(&showRequestBody, "show-request-body", false, ", -d\n\t\tPrint every request body (ie, POST data).")
  flag.BoolVar(&showRequestBody, "d", false, "")
  flag.BoolVar(&showTargetRequestBody, "show-rewritten-request-body", false, ", -d2\n\t\tPrint every rewritten request body.")
  flag.BoolVar(&showTargetRequestBody, "d2", false, "")
  flag.BoolVar(&showRequestHeaders, "show-request-headers", false, ", -reqh\n\t\tPrint request headers.")
  flag.BoolVar(&showRequestHeaders, "reqh", false, "")
  flag.BoolVar(&showTargetRequestHeaders, "show-rewritten-request-headers", false, ", -reqh2\n\t\tPrint rewritten request headers.")
  flag.BoolVar(&showTargetRequestHeaders, "reqh2", false, "")
  flag.BoolVar(&showResponse, "show-response", false, ", -res\n\t\tPrint every unmodified response body.")
  flag.BoolVar(&showResponse, "res", false, "")
  flag.BoolVar(&showFinalResponse, "show-rewritten-response", false, ", -res2\n\t\tPrint every rewritten response body.")
  flag.BoolVar(&showFinalResponse, "res2", false, "")
  flag.BoolVar(&showResponseHeaders, "show-response-headers", false, ", -resh\n\t\tPrint all unmodified response heders.")
  flag.BoolVar(&showResponseHeaders, "resh", false, "")
  flag.BoolVar(&showFinalResponseHeaders, "show-rewritten-response-headers", false, ", -resh2\n\t\tPrint all rewritten response headers.")
  flag.BoolVar(&showFinalResponseHeaders, "resh2", false, "")
  flag.BoolVar(&unzip, "open-gzip", false, ", -g\n\t\tInternally uncompress and recompress gzipped contents. Set to True if \n\t\tyou want to see or modify the body of responses (by hacking Toxy's source code).\n\t\tLeave as False for a faster proxy.")
  flag.BoolVar(&unzip, "g", false, "")
  flag.BoolVar(&useCache, "use-cache", false, ", -c\n\t\tCache all GET requests.")
  flag.BoolVar(&useCache, "c", false, "")
}

func main() {
  flag.Parse()
  http.HandleFunc("/", report)
  if target == "" {
    log.Println("Starting traditional proxy.")
  } else {
    log.Println("Starting reverse proxy for " + target)
  }
  log.Fatal(http.ListenAndServe(":" + port, nil))
}

func report(w http.ResponseWriter, r *http.Request){
  uri := target+r.RequestURI

  if useCache && r.Method == "GET" {
    element, error := cacheGet(uri)
    if(error == nil) {
      dH := w.Header()
      copyHeader(element.resp.Header, &dH)
      dH.Add("From-Cache", "true")
      dH.Set("Content-Length", strconv.Itoa(len(element.body)))
      w.WriteHeader(element.resp.StatusCode)
      w.Write(element.body)
      return
    }
  }

  if showUri {
    fmt.Println("URI ("+r.Method + "): " + r.Host + r.RequestURI)
  }

  if showTargetUri {
    fmt.Println("Modified URI (" + r.Method + "): " + uri)
  }

  requestBody, err := ioutil.ReadAll(r.Body)
  r.Body.Close()
  fatal(err)

  if showRequestHeaders {
    fmt.Printf("Request Headers: %v\n", r.Header);
  }

  if showRequestBody && string(requestBody) != "" {
    fmt.Printf("Request Body: %s\n", string(requestBody));
  }

  /*  We can make changes to the request body (requestBody) here */

  r2, err := http.NewRequest(r.Method, uri, bytes.NewBuffer(requestBody))
  fatal(err)
  copyHeader(r.Header, &r2.Header)
  r2.Header.Set("Content-Length", strconv.Itoa(len(requestBody)))

  /* We can mess with the request headers here */
  // r2.Header.Add("Something", "Some Value")
  /* */

  if showTargetRequestHeaders {
    fmt.Printf("Modified Request Headers: %v\n", r2.Header);
  }

  if showTargetRequestBody && string(requestBody) != "" {
    fmt.Printf("Modified Request Body: %s\n", string(requestBody));
  }

  var transport http.Transport
  resp, err := transport.RoundTrip(r2)
  if err != nil {
    log.Print(err)
    return;
  }

  if showResponseHeaders {
    fmt.Printf("Response Headers: %v\n", resp.Header);
  }

  responseBody, err := getBodyBuffer(resp)

  if showResponse {
   fmt.Printf("Response Body: %s\n",responseBody)
  }

  /* We can mess with the final response body (responseBody) here */
  //responseBody = swap(responseBody, "north", "south")
  //responseBody = swap(responseBody, "down", "up")
  //responseBody = replace(responseBody, "nice", "nasty")
  //responseBody = replace(responseBody, "cool", "hideous")
  /* */

  dH := w.Header()
  copyHeader(resp.Header, &dH)

  /* We can mess with the final response headers here */
  // fmt.Printf("Resp-Status: %v\n",resp.Status)
  dH.Add("Requested-Host", r2.Host)
  /* */

  body, _ := compressBody(responseBody, resp)
  dH.Set("Content-Length", strconv.Itoa(len(body)))

  if showFinalResponseHeaders {
    fmt.Printf("Modified Response Headers: %v\n", resp.Header);
  }

  if showFinalResponse {
   fmt.Printf("Modified Response Body: %s\n", responseBody)
  }

  w.WriteHeader(resp.StatusCode)
  w.Write(body)

  if useCache && r.Method == "GET" {
    cacheSet(uri, cacheElement{ resp: resp, body: body })
  }
}

func fatal(err error) {
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
}

func copyHeader(source http.Header, dest *http.Header){
  for n, v := range source {
      for _, vv := range v {
          dest.Add(n, vv)
      }
  }
}

func getBodyBuffer(resp *http.Response) ([]byte, error) {
  if unzip && resp.Header.Get("Content-Encoding") == "gzip" {
    reader, err := gzip.NewReader(resp.Body)
    bodyBuffer, err := ioutil.ReadAll(reader)
    return bodyBuffer, err
  } else {
    bodyBuffer, err := ioutil.ReadAll(resp.Body)
    return bodyBuffer, err
  }
}

func compressBody(body []byte, r *http.Response) ([]byte, error) {
  if unzip && r.Header.Get("Content-Encoding") == "gzip" {
    buffer := bytes.NewBuffer([]byte{})
    gw := gzip.NewWriter(buffer)
    gw.Write(body)
    gw.Close()
    compressedBody, err := ioutil.ReadAll(buffer);
    return compressedBody, err
  } else {
    return body, nil
  }
}

func swap(subject []byte, swap1, swap2 string) []byte {
  ret := bytes.Replace(subject, []byte(swap1), []byte("__slgplaceholderslg__"), -1)
  ret = bytes.Replace(ret, []byte(swap2), []byte(swap1), -1)
  ret = bytes.Replace(ret, []byte("__slgplaceholderslg__"), []byte(swap2), -1)
  return ret
}

func replace(subject []byte, search, replace string) []byte {
  return bytes.Replace(subject, []byte(search), []byte(replace), -1)
}

func cacheSet(key string, val cacheElement){
  cache[key] = val
}

func cacheGet(key string) (cacheElement, error){
  if element, ok := cache[key]; ok {
    if ok {
      return element, nil
    }
  }
  return cacheElement{}, errors.New("Not in cache")
}
