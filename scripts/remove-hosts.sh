#!/bin/bash

if [ $# -ne 2 ]; then
  echo "usage: script <start host ordinal> <end host ordinal>"
  echo "./script 500 502" 
  exit 1
fi


for hhost in web papi extapi; do
#for hhost in extapi; do
  for ((i=$1; i<=$2; i++)); do
    curl -L http://site-etcd-500.iad.prod.zulily.com:2379/v2/keys/site/$hhost/$i -XDELETE
  done
done
