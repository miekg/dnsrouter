# DNS Router

DNS router is a router for DNS packets. It specifically handles the case when a zone
does not fit in the memory of a server anymore: when it needs to be spread across 
multiple servers. 

The nameserver(s) where to packets are forwarded to, do not need to be changed, however
they do need to run with a part of a zone loaded. And with DNSSEC you need to make
sure each server has enough of the NSEC/NSEC3 chain to give back valid responses.

## How it functions

* It gets a list of nameservers from etcd which it refreshes every 10 seconds;
* It health checks the nameservers with a TCP connection doing a configurable query;
    * if no reply the remote server will be removed from the list, until the health
        check is OK again.
* Each return packet travels back through the DNS router back to the client;
* For TCP a pre-setup connection is used. This is *not* the healthchecking connection;


Features I want, but are not implemented yet:

* The original source address is used when forwarding an UDP message;
    * This could also be fixed by running a wrapper protocol (requires backend changes) or
        other tricks.


## How it routes

Hallo
