#!/usr/bin/env python

parent="site"
types=["web", "extapi", "papi"]
hrange=[500, 549]

domain="iad.prod.zulily.com"

for type in types:
  for hostnum in range(hrange[0], hrange[1]+1):
    print "%s-%s-%s.%s" % (parent, type, hostnum, domain)

