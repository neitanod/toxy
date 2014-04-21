toxy
====

Toxy is a tiny hacky proxy written in Go in order to serve as a base to spy or 
mess with http connections.

Usage:
------

### Help:

    $ ./toxy -h

    Usage of ./toxy:
      -port="8080": port to listen to.
      -target="": target URL for reverse proxy.  Leave empty for traditional proxy.

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
