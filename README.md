deceive
=======

Deceive is an HTTPS webserver that uses TLS Client auth to validate uploads
(via `PUT`'ing a file path) to a given root.

Usage
=====

```
mkdir -p /tmp/deceive/foo
deceive \
    -ca=/home/paultag/certs/cacert.crt \
    -cert=/home/paultag/certs/localhost.crt \
    -key=/home/paultag/certs/localhost.key \
    -host=soylent.green \
    -port=1984 \
    -root=/tmp/deceive/
```

This will start a `deceive` server on `soylent.green:1984`. The server will
use the `-cert` and `-key` param for the TLS serverside certificates. Incoming
requests to `PUT` data will have their TLS client certs validated against the
`-ca` cert given. If the cert is valid, and the path exists, `deceive` will
allow the write.
