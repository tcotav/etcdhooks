# etcdhook -- Go + Nagios + Etcd + Other Configs

Service that watches `etcd` for changes and then updates nagios and custom config files reflecting those changes.

Also includes a simple webservice that allows you to dump the contents of etcd host states and get back a json blob.

### Get the go-etcd package

    go get github.com/coreos/go-etcd/etcd

### Dev Environment

Run the provided script against any Ubuntu VPS.  Assumes non-root user.

    init_host.sh

Then once etcd is running, you can throw some fake data in:
  
    init_etcd.sh

### nagios 

For nagios, you should be running it using conf.d to be able to just dump in any ol file for configs.  This is default from the ubuntu packages.  The script as currently written will overwrite previous versions of the config file so it assumes this app manages those hosts + group file completely.

### nagios restart

tbd
