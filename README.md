# etcdhook -- Go + Nagios + Etcd + Other Configs

Service that watches `etcd` for changes and then updates nagios and custom config files reflecting those changes.

Also includes a simple webservice that allows you to dump the contents of etcd host states and get back a json blob.

The basic hostmap is of the format:

    <hostname>:<rotation status of host -- in or out>
    site-web-200:1     # is in
    site-web-200:0     # is out

### Get the go-etcd package

    go get github.com/coreos/go-etcd/etcd

### Dev Environment

Run the provided script against any Ubuntu VPS.  Assumes non-root user.

    init_host.sh

Then once etcd is running, you can throw some fake data in:
  
    init_etcd.sh

### nagios 

For nagios, you should be running it using `conf.d` to be able to just dump in any ol file for configs.  This is default from the ubuntu packages.  The script as currently written will overwrite previous versions of the config file so it assumes this app manages those hosts + group file completely.  You best bet is to create some file for dynamic-only hosts and dynamic only groups and have that constantly recreated in the `conf.d` directory.

### nagios restart

tbd -- need to figure out some best practice way to HUP nagios.  maybe it is to just shell out and do a `service nagios3 restart`.


### Configurationa

Basic configuration is as follows:

    # this is a comment
    nagios_host_file=/tmp/host.cfg
    nagios_groups_file=/tmp/groups.cfg
    host_list_file=/tmp/hostlist.cfg

    # the etcd section of the configuration
    etcd_server_list=http://127.0.0.1:4001

    web_listen_port=3000

    base_etcd_url=/site/

base_etcd_url assumes a three part naming scheme that correlates to a three stage url.  For example site-db-800 would be /site/db/800 and would be a walkdown through <team>/<host type>/<specific hostid>.

### Webservice

Only have a crude initial placeholder for this.  If you hit the base `/` url at port `web_listen_port`, you'll get a json dump of the current contents of etcd from the root `base_etcd_url`.

We could easily add features to have a `force_reload` endpoint or anything else.  We'll see what makes sense once we try this sucker in production.
