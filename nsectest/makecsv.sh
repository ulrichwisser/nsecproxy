#!/bin/bash
#
# record[0] ip address
# record[1] name
# record[2] AS number
# record[3] cc
# record[4] city
# record[5] version
# record[6] error
# record[7] dnssec
# record[8] relability
# record[9] checked at
# record[10] created at
#
OLDIFS=$IFS
IFS=","
(
while read ip name asn asname cc city version other
do
    echo "${asn}_${cc}_${city},${ip}"
done < <(curl -s "https://public-dns.info/nameservers.csv")
) | sed 's/_,/,/g;s/ //g;s/"//g'
IFS=$OLDIFS
