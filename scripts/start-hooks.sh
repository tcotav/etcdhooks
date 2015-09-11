#!/bin/bash

nohup /opt/etcd/etcd --data-dir /opt/etcd/data >>/var/log/etcd.log 2>&1 &
nohup /opt/etcd/etcdhooks >>/var/log/etcdhooks.log 2>&1 &
