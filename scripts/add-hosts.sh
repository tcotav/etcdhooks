#!/bin/bash

if [ $# -ne 3 ]; then
  echo "usage: script <host-type> <start host ordinal> <end host ordinal>"
  echo "./script web 500 502" 
  exit 1
fi

HOSTTYPE=$1
START_HOSTNUM=$2
LAST_HOSTNUM=$3

for ((i=$START_HOSTNUM; i<=$LAST_HOSTNUM; i++)); do
  curl -L http://site-etcd-500:2379/v2/keys/site/$HOSTTYPE/$i -XPUT -d value=hbout 
done 
