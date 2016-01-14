#!/bin/bash

cd ../

gb vendor fetch github.com/Sirupsen/logrus
gb vendor fetch github.com/coreos/etcd/client

# if we're just cloning a repo, this command will make sure we've updated our vendor srcs
gb vendor restore
