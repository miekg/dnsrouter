# DNS Router

DNS router is a router for DNS packets. It specifically handles the case when a zone
does not fit in the memory of a server anymore: when it needs to be spread across 
multiple servers.

The nameserver(s) where to packets are forwarded to do not need to be changed, however
they do need to run with a part of a zone loaded. And with DNSSEC you need to make
sure each server has enough of the NSEC/NSEC3 chain to give back valid responses.

## How it functions

* It gets a list of nameservers from etcd;
    * Each nameserver publishes a file in etcd in the directory `/dnsrouter/`,
        the name of the file is an UUID identifier. Each file contains an
       `<ip>`,`<regexp>` combination;
    * Every query that matches the regular expression is sent to this ip, if there
        are multiple ips it will round robin.
* Sets a watch on the directory to get notifications;
* It health checks the nameservers with a TCP connection doing a `id.server` CH TXT query;
    * if no reply the remote server is seen, it will be removed from the list.
* Each return packet travels back through the DNS router back to the client;
* Each connection is done over UDP to the nameserver, even for clients who initially
    connected over TCP;
* AXFR or IXFR is not supported and NACKed on the DNS router;
* If none of the regular expression match a SERVFAIL is returned to the client.

Features I want, but are not implemented yet:

* The original source address is used when forwarding an UDP message;
    * This could also be fixed by running a wrapper protocol (requires backend changes) or
        other tricks.
* Re-use the original packet bytes (this may be a feature for Go DNS);
* Make it more efficient.
