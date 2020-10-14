# Performing resolver tests of NSEC and NSEC3 handling

## tests performed
According to RFC NSEC records are included in the following answers.
- no data, with an additional NSEC record for proving that no wildcard exists or the existing wildcard does not cover the rr type.
- name error
- empty non terminal (a special case of no data)
For all of these cases we have prepared a number of different answers.
- nsec1... answer contains only NSEC records
- nsec3... answer contains only NSEC3 records
- nsec4... answer contains NSEC and NSEC3 records (in that order)
- nsec5... answer contains NSEC3 and NSEC records (in that order)
Besides these cases there are several cases with DNSSEC failures
- no nsec data
- NSEC does not cover label
- NSEC3 does not cover label
- NSEC and NSEC3 non covers label
- NSEC does cover label, NSEC3 does not cover label
- NSEC does not cover label, NSEC3 does cover label
- data with nsec
- data with nsec3
### Testing commands
####; no data - no wildcard
dig +dnssec nsec1.a.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.a.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.a.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.a.$NSECPROXYDOMAIN aaaa
####; no data - wildcard
dig +dnssec nsec1.b.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.b.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.b.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.b.$NSECPROXYDOMAIN aaaa
####; name error - no wildcard
dig +dnssec nsec1.c.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.c.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.c.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.c.$NSECPROXYDOMAIN aaaa
####; empty non terminal
dig +dnssec nsec1.d.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec3.d.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec4.d.$NSECPROXYDOMAIN aaaa
dig +dnssec nsec5.d.$NSECPROXYDOMAIN aaaa
####; failure - no nsec data
dig +dnssec a.nseca.$NSECPROXYDOMAIN aaaa
####; failure - NSEC does not cover label
dig +dnssec b.nsecb.$NSECPROXYDOMAIN aaaa
####; failure - NSEC3 does not cover label
dig +dnssec b.nsecc.$NSECPROXYDOMAIN aaaa
####; failure - NSEC and NSEC3 non covers label
dig +dnssec b.nsecd.$NSECPROXYDOMAIN aaaa
####; failure - NSEC does cover label, NSEC3 does not cover label
dig +dnssec b.nsece.$NSECPROXYDOMAIN aaaa
dig +dnssec d.nsece.$NSECPROXYDOMAIN aaaa
####; failure - NSEC does not cover label, NSEC3 does cover label
dig +dnssec b.nsecf.$NSECPROXYDOMAIN aaaa
dig +dnssec d.nsecf.$NSECPROXYDOMAIN aaaa
####; failure - data with nsec
dig +dnssec b.nsecg.$NSECPROXYDOMAIN txt
####; failure - data with nsec3
dig +dnssec b.nsegh.$NSECPROXYDOMAIN txt

## Resolvers for Testing
APNIC does provide a nice list op open resolvers at https://www.publicdns.xyz.
Unfortunately they do not provide a machine readable download. I have
copied the public resolvers from the first page inte publicresolvers.csv.

There is another public resolver list at https://public-dns.info/nameservers.csv.
Here a download is possible. Due to GDPR reasons I do not include the list here,
but feel free to use makecsv.sh to download and format your own list.

The results files of the test run will no longer include ip addresses due to GDPR
requirements.

# Run the tests
After compiling the code you are ready to make the first test run.
```
nsectest -v -d example.com -R publicresolvers.csv > ../results/publicresolvers.json
makecsv.sh > nameservers.csv
nsectest -v -d example.com -R nameservers.csv > ../results/nameservers.json
```

# Results
Please follow the instruction for [nsecstats](https://github.com/ulrichwisser/nsecproxy/nsecstats) on how to produce and view the results.
