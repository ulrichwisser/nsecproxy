# Compute statistics from resolver tests

This is work in progress!

Once you have run the resolver tests, this tool will compute
statistics from the results.

```
nsecstats -v -s ../results/publicresolvers -d ../results/publicresolvers.json
nsecstats -v -s ../results/nameservers -d ../results/nameservers.json
```

Currently the index.html page needs to be updated manually with links
to the resulting statistics.

Currently statistics can be seen at https://wisser.se/nsecproxy
