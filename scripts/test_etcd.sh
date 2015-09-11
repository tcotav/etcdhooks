#!/bin/bash

# simple script to dump some data into etcd to use for testing app
curl -L http://127.0.0.1:4001/v2/keys/site/web/501 -XPUT -d value=1
curl -L http://127.0.0.1:4001/v2/keys/site/web/502 -XPUT -d value=0
curl -L http://127.0.0.1:4001/v2/keys/site/db/500 -XPUT -d value=1

curl -L http://127.0.0.1:4001/v2/keys/site/web/503 -XPUT -d value=0
curl http://127.0.0.1:4001/v2/keys/site/web/503
sleep 5
curl -L http://127.0.0.1:4001/v2/keys/site/web/503 -XDELETE
curl http://127.0.0.1:4001/v2/keys/site/web/503
#!/bin/bash


