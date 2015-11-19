#!/bin/bash

cd /opt/etcd
if pidof -x "etcdhooks" > /dev/null ; then
  echo "etcdhooks already running"
  exit 1
fi
sudo bash -c "nohup /opt/etcd/etcdhooks >>/var/log/etcdhooks.log 2>&1 &"
