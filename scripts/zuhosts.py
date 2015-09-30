#!/usr/bin/env python

parent="site"
types=["web", "extapi", "papi"]
hrange=[500, 549]

domain="iad.prod.zulily.com"
print "#!/bin/bash\n\n"
for type in types:
  for hostnum in range(hrange[0], hrange[1]+1):
    print "curl -L http://127.0.0.1:4001/v2/keys/%s/%s/%s -XPUT -d 0" % (parent, type, hostnum)


