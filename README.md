# gonagetcd -- Go Nagios Etcd

Service that watched `etcd` for changes and then updates nagios and custom config files reflecting those changes.

### Get the go-etcd package

    go get github.com/coreos/go-etcd/etcd

### Dev Environment

Run the provided script against any Ubuntu VPS.  Assumes non-root user.

    init_host.sh

Then once etcd is running, you can throw some fake data in:
  
    init_etcd.sh

### nagios 
tbd

### nagios restart

tbd
