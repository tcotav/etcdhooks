#!/bin/bash

# simple script to dump some data into etcd to use for testing app
#curl -L http://127.0.0.1:4001/v2/keys/site/web/501/ip -XPUT -d value=10.0.1.1
#curl -L http://127.0.0.1:4001/v2/keys/site/web/501/status -XPUT -d value=1
#
#curl -L http://127.0.0.1:4001/v2/keys/site/web/502/ip -XPUT -d value=10.0.1.2
#curl -L http://127.0.0.1:4001/v2/keys/site/web/502/status -XPUT -d value=1
#
#curl -L http://127.0.0.1:4001/v2/keys/site/db/500/ip -XPUT -d value=10.0.1.10
#curl -L http://127.0.0.1:4001/v2/keys/site/db/500/status -XPUT -d value=1
#

curl -L http://127.0.0.1:4001/v2/keys/site/web/503/ip -XPUT -d value=10.0.1.3
sleep 5
curl -L http://127.0.0.1:4001/v2/keys/site/web/503/status -XPUT -d value=0
sleep 5
curl -L http://127.0.0.1:4001/v2/keys/site/web/503 -XDELETE
