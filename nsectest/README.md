# Performing resolver tests of NSEC and NSEC3 handling

## tests performed
###; no data - no wildcard
dig +dnssec nsec1.a.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.a.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.a.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.a.$NSECPROXYDOMAIN aaaa

###; no data - wildcard
dig +dnssec nsec1.b.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.b.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.b.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.b.$NSECPROXYDOMAIN aaaa

###; name error - no wildcard
dig +dnssec nsec1.c.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.c.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.c.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.c.$NSECPROXYDOMAIN aaaa

###; empty non terminal
dig +dnssec nsec1.d.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.d.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.d.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.d.$NSECPROXYDOMAIN aaaa

###; failure - no nsec data
dig +dnssec a.nseca.$NSECPROXYDOMAIN aaaa

###; failure - NSEC does not cover label
dig +dnssec b.nsecb.$NSECPROXYDOMAIN aaaa

###; failure - NSEC3 does not cover label
dig +dnssec b.nsecc.$NSECPROXYDOMAIN aaaa

###; failure - NSEC and NSEC3 non covers label
dig +dnssec b.nsecd.$NSECPROXYDOMAIN aaaa

###; failure - NSEC does cover label, NSEC3 does not cover label
dig +dnssec b.nsece.$NSECPROXYDOMAIN aaaa
dig +dnssec d.nsece.$NSECPROXYDOMAIN aaaa

###; failure - NSEC does not cover label, NSEC3 does cover label
dig +dnssec b.nsecf.$NSECPROXYDOMAIN aaaa
dig +dnssec d.nsecf.$NSECPROXYDOMAIN aaaa

###; failure - data with nsec
dig +dnssec b.nsecg.$NSECPROXYDOMAIN txt
sleep 6
dig +dnssec b.nsecg.$NSECPROXYDOMAIN txt

###; failure - data with nsec3
dig +dnssec b.nsegh.$NSECPROXYDOMAIN txt
sleep 6
dig +dnssec b.nsegh.$NSECPROXYDOMAIN txt


# Open Resolvers

https://public-dns.info/nameservers.csv

https://www.publicdns.xyz

googlepdns
114dns
cloudflare
opendns
dnspai
onedns
vrsgn
quad9
level3
neustar
yandex
dnswatch
dyn
cnnic
he


Alternate DNS 198.101.242.72
Alternate DNS 23.253.163.53

BlockAid Public DNS205.204.88.60
BlockAid Public DNS178.21.23.150

Censurfridns 91.239.100.100
Censurfridns 89.233.43.71

2001:67c:28a4::
2002:d596:2a92:1:71:53::
Christoph Hochstätter
209.59.210.167
85.214.117.11
ClaraNet
212.82.225.7
212.82.226.212
Comodo Secure DNS
8.26.56.26
8.20.247.20
DNS.Watch
84.200.69.80
84.200.70.40

2001:1608:10:25::1c04:b12f
2001:1608:10:25::9249:d69b
DNSReactor
104.236.210.29
45.55.155.25
Dyn
216.146.35.35
216.146.36.36
FDN
80.67.169.12

2001:910:800::12
FoeBud
85.214.73.63
FoolDNS
87.118.111.215
213.187.11.62
FreeDNS
37.235.1.174
37.235.1.177
Freenom World
80.80.80.80
80.80.81.81
German Privacy Foundation e.V.
87.118.100.175
94.75.228.29
85.25.251.254
62.141.58.13
Google Public DNS
8.8.8.8
8.8.4.4

2001:4860:4860::8888
2001:4860:4860::8844
GreenTeamDNS
81.218.119.11
209.88.198.133
Hurricane Electric
74.82.42.42

2001:470:20::2
Level3
209.244.0.3
209.244.0.4
Neustar DNS Advantage
156.154.70.1
156.154.71.1
New Nations
5.45.96.220
185.82.22.133
Norton DNS
198.153.192.1
198.153.194.1
OpenDNS
208.67.222.222
208.67.220.220

2620:0:ccc::2
2620:0:ccd::2
OpenNIC
58.6.115.42
58.6.115.43
119.31.230.42
200.252.98.162
217.79.186.148
81.89.98.6
78.159.101.37
203.167.220.153
82.229.244.191
216.87.84.211
66.244.95.20
207.192.69.155
72.14.189.120

2001:470:8388:2:20e:2eff:fe63:d4a9
2001:470:1f07:38b::1
2001:470:1f10:c6::2001
PowerNS
194.145.226.26
77.220.232.44
Quad9
9.9.9.9

2620:fe::fe
SafeDNS
195.46.39.39
195.46.39.40
SkyDNS
193.58.251.251
SmartViper Public DNS
208.76.50.50
208.76.51.51
ValiDOM
78.46.89.147
88.198.75.145
Verisign
64.6.64.6
64.6.65.6

2620:74:1b::1:1
2620:74:1c::2:2
Xiala.net
77.109.148.136
77.109.148.137

2001:1620:2078:136::
2001:1620:2078:137::
Yandex.DNS
77.88.8.88
77.88.8.2

2a02:6b8::feed:bad
2a02:6b8:0:1::feed:bad
puntCAT
109.69.8.51

2a00:1508:0:4::9


Google	8.8.8.8	8.8.4.4
Quad9	9.9.9.9	149.112.112.112
OpenDNS Home	208.67.222.222	208.67.220.220
Cloudflare	1.1.1.1	1.0.0.1
CleanBrowsing	185.228.168.9	185.228.169.9
Verisign	64.6.64.6	64.6.65.6
Alternate DNS	198.101.242.72	23.253.163.53
AdGuard DNS	176.103.130.130	176.103.130.131


Comodo Secure DNS
