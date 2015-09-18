#!/bin/bash

### install golang
wget https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz
tar -xzvf go1.5.linux-amd64.tar.gz
sudo mv go /usr/local
export PATH=$PATH:/usr/local/go/bin
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
export GOROOT=/usr/local/go
echo "export GOROOT=/usr/local/go" >> ~/.bashrc
echo "export GOPATH=~/go" >> ~/.bashrc
export GOPATH=~/go
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:$GOROOT/bin
echo "export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:$GOROOT/bin" >> ~/.bashrc
mkdir -p ~/go/src
cd ~/go

# get the etcd and any other packages you need
#go get github.com/coreos/go-etcd/etcd
go get github.com/coreos/etcd/client
go get github.com/Sirupsen/logrus

cd ~/go/github.com
ln -s ~/etcdhooks/src/github.com/tcotav
cd ~/go/github.com/tcotav/etcdhooks
# build the binary of our go service
go build -o etcdhooks daemon.go 

### Set up etcd
curl -L https://github.com/coreos/etcd/releases/download/v2.2.0/etcd-v2.2.0-linux-amd64.tar.gz -o etcd-v2.2.0-linux-amd64.tar.gz
tar xzvf etcd-v2.2.0-linux-amd64.tar.gz
#cd etcd-v2.2.0-linux-amd64
#sudo ./etcd &

cd /opt/etcd
sudo cp ~/go/src/github.com/tcotav/etcdhooks/etcdhooks .
sudo cp ~/go/src/github.com/tcotav/etcdhooks/daemon.cfg .

sudo cp ~/etcdhooks/scripts/*-hooks.sh .

curl -L http://127.0.0.1:4001/v2/keys/site/init -XPUT -d value=1
