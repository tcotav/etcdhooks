#!/bin/bash

if [ $# -ne 2 ]; then
  echo "usage: script <start host ordinal> <end host ordinal>"
  echo "./script 500 502" 
  exit 1
fi
for hhost in web papi extapi; do
  for hnum in {$1..$2}; do
  curl -L http://127.0.0.1:4001/v2/keys/site/$hhost/$hnum -XDELETE
  curl -L http://127.0.0.1:4001/v2/keys/site/$hhost/$hnum -XPUT -d value=0
done
done
