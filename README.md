toxy
====

Toxy is a tiny hacky proxy written in Go in order to serve as a base to spy or 
mess with http connections.

Usage:
------

### Help:

    $ ./toxy -h
    Usage of ./toxy:
      --open-gzip=false, -g
                    Internally uncompress and recompress gzipped contents. Set to True if you
                    want to see or modify the body of responses (by hacking Toxy's source code).
                    Leave as False for a faster proxy.
      --port="8080", -p
                    Port to listen to.
      --show-request=true, -r
                    Print every request URI.
      --show-request-body=false, -d
                    Print every request body (ie, POST data).
      --show-request-headers=false, -reqh
                    Print request headers.
      --show-response=false, -res
                    Print every unmodified response body.
      --show-response-headers=false, -resh
                    Print all unmodified response heders.
      --show-rewritten-request=false, -r2
                    Print every rewritten (modified by this proxy) request URI.
      --show-rewritten-request-body=false, -d2
                    Print every rewritten request body.
      --show-rewritten-request-headers=false, -reqh2
                    Print rewritten request headers.
      --show-rewritten-response=false, -res2
                    Print every rewritten response body.
      --show-rewritten-response-headers=false, -resh2
                    Print all rewritten response headers.
      --target="", -t
                    Target URL for single-host everse proxy.  Leave empty for traditional proxy.
      --use-cache=false, -c
                    Cache all GET requests.

### As a traditional proxy:

Set up your browser to query the proxy at localhost:8080 (for http only), then run:

    $ ./toxy
    Starting traditional proxy.

and you can browse the web.

### As a single domain reverse proxy:

Run:

    $ ./toxy -target=http://someserver.com
    Starting reverse proxy for http://someserver.com

Then open http://localhost:8080 in your browser.  The destination server will 
receive requests with the "Host" header set to "someserver.com".

This is already useful to test your local sites on Apache virtual hosts from
other machines (and devices) in your network without having to add entries on 
/etc/hosts/ on every device.  Also, you will see the request pass through your
console.  :)

For instance, open http://192.168.1.4:8080 (assuming your development machine is
at that IP address) from another machine on the local network, and you'll see 
your virtual host.

Installation:
-------------

Install Go from http://golang.org.

Clone and compile Toxy with:

    $ git clone https://github.com/neitanod/toxy.git
    Cloning into 'toxy'...
    $ cd toxy
    $ go build

And you should get the `toxy` executable.
