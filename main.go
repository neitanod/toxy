package main

import (
    "flag"
    "fmt"
    "os"
    "log"
    "bytes"
    "net/http"
    "io/ioutil"
)

var target *string
var port *string

func main() {
  target = flag.String("target", "", "target URL for reverse proxy.  Leave empty for traditional proxy.")
  port = flag.String("port", "8080", "port to listen to.")
  flag.Parse()
  http.HandleFunc("/", report)
  if *target == "" {
    log.Println("Starting traditional proxy.")
  } else {
    log.Println("Starting reverse proxy for " + *target)
  }
  log.Fatal(http.ListenAndServe(":" + *port, nil))
}

func report(w http.ResponseWriter, r *http.Request){

  uri := *target+r.RequestURI

  fmt.Println(r.Method + ": " + uri)

  requestBodyBytes, err := ioutil.ReadAll(r.Body)
  fatal(err)
  defer r.Body.Close()

  /*  We can mess with the request body (requestBodyBytes) here */
  fmt.Printf("Body: %v\n", string(requestBodyBytes));
  /* */

  rr, err := http.NewRequest(r.Method, uri, bytes.NewBuffer(requestBodyBytes))
  fatal(err)
  copyHeader(r.Header, &rr.Header)

  /* We can mess with the request headers here */
  // rr.Header.Add("Something", "Some Value")
  /* */

  var transport http.Transport
  resp, err := transport.RoundTrip(rr)
  if err != nil {
    log.Print(err)
    return;
  }

  /* We can spy on original response headers here */
  fmt.Printf("Resp-Headers: %v\n", resp.Header);
  /* */

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  fatal(err)

  /* We can mess with the final response headers here */
  // fmt.Printf("Resp-Body: %s\n",body) //probably gzipped
  /* */

  dH := w.Header()
  copyHeader(resp.Header, &dH)

  /* We can mess with the final response headers here */
  // fmt.Printf("Resp-Status: %v\n",resp.Status)
  dH.Add("Requested-Host", rr.Host)
  /* */

  w.WriteHeader(resp.StatusCode)
  w.Write(body)
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
