# DNS Router

DNS router is a router for DNS packets. It sends a packet to another nameserver, it sends out
the UDP packet with the 



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
