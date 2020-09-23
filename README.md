# nsecproxy
DNS proxy with special NSEC handling

## Introduction
At the Swedish Internet Foundation we wanted to go over from NSEC3 to NSEC for the .nu ccTLD.

According to [RFC5155](https://tools.ietf.org/html/rfc5155#section-10.5) it is very easy
to change. Just remove NSEC3PARAM and all NSEC3 records from the zone, add NSEC records
and everything should just work. (Attention! .nu uses algorithm 13 for signing, which supports NSEC and NSEC3.)

But just pointing to an RFC might not be enough if you have to answer to a board
of directors or your government.

Therefore we designed a number of tests to simulate different states of caches and
zone updates.

## Overview
The whole system contains two main parts, the authoritative setup and the client.

The authoritative setup consists of two instances of the Knot DNS server and a
special filtering dns proxy written in go.

The client asks a list of predefined questions and writes the answers as json to stdout.

## Installation Authoritative
This software was developed and tested on an Ubuntu 20.04 server. Actually running on an AWS EC2 t2.micro instance.

## Install knot
1. install knot package https://www.knot-dns.cz/docs/3.0/html/installation.html#installation-from-a-package
1. install go `sudo apt install golang`
1. copy this repository `go get github.com/ulrichwisser/nsecproxy`
1. `cd $GOPATH/src/github.com/ulrichwisser/nsecproxy`
1. `cp knot@.service /lib/systemd/system/`
1. `cp knot-nsec*.conf /etc/knot/`
1. mkdir /var/lib/knot/nsec
1. mkdir /var/lib/knot/nsec3
1. cp nsec.zone /var/lib/knot/nsec/example.com.zone
1. cp nsec3.zone /var/lib/knot/nsec3/example.com.zone
1. chown -R knot.knot /var/lib/knot/*
1. mkdir /run/knot/nsec
1. mkdir /run/knot/nsec3
1. chown -R knot.knot /run/knot/*

### DNSSEC key management
1. keymgr -c /etc/knot/knot-nsec.conf example.com generate ksk=yes
1. remember the ksk-id written out by the command
1. keymgr -c /etc/knot/knot-nsec.conf example.com generate
1. remember the zsk-id written out by the command
1. Copy keys to the other knot instance `cp -a /var/lib/knot/nsec/keys /var/lib/knot/nsec3/`
1. keymgr -c /etc/knot/knot-nsec3.conf example.com import-pem /var/lib/knot/nsec3/keys/keys/<ksk-id>.pem ksk=yes
1. keymgr -c /etc/knot/knot-nsec3.conf example.com import-pem /var/lib/knot/nsec3/keys/keys/<zsk-id>.pem

### Start knot servers
1. `systemctl start knot@nsec`
1. `systemctl start knot@nsec3`

## Install nsecproxy
