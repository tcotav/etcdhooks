# etcdhook -- Go + Nagios + Etcd + Other Configs

Service that watches `etcd` for changes and then updates nagios and custom config files reflecting those changes.

Also includes a simple webservice that allows you to dump the contents of etcd host states and get back a json blob.

The basic hostmap is of the format:

    <hostname>:<random state string -- could be anything>
    site-web-200:hbout       # is marked as out of rotation
    site-web-200:created     # newly created VM

### Get required go packages
    go get github.com/coreos/etcd/client
    go get github.com/Sirupsen/logrus

### Environment

Run the provided script against any Ubuntu VPS.  Assumes non-root user.

    init_host.sh

This will 
  - install go in /usr/local/go
  - installs etcd in /opt/etcd
  - creates ~/go for you
  - set up your GOROOT and GOPATH properly
  - stick those vars into your ~/.bashrc -- change this if you use some other shell
  - install requied golang packages from github
  - builds the binary etcdhooks for you and moves it to /opt/etcd

### nagios 

For nagios, you should be running it using `conf.d` to be able to just dump in any ol file for configs.  This is default from the ubuntu packages.  The script as currently written will overwrite previous versions of the config file so it assumes this app manages those hosts + group file completely.  You best bet is to create some file for dynamic-only hosts and dynamic only groups and have that constantly recreated in the `conf.d` directory.

### nagios restart

tbd -- need to figure out some best practice way to HUP nagios.  maybe it is to just shell out and do a `service nagios3 restart`.


### Configurationa

Daemon configuration is as follows:

    i# this is a comment
    nagios_host_file=/etc/nagios3/conf.d/flex-hosts.cfg
    nagios_groups_file=/etc/nagios3/conf.d/flex-groups.cfg
    #host_list_file=/tmp/hostlist.cfg

    # the etcd section of the configuration -- use comma separated values
    etcd_server_list=http://127.0.0.1:4001

    # binds to all NICs currently.  Simple web listener
    web_listen_port=3000

    # url that data is written to in etcd
    etcd_watch_root_url=/site/

    # how long to queue file rewrite actions before firing them off - in seconds
    file_rewrite_interval=15

    # NYI
    # csv list of files to rewrite
    # currently support nagios,host
    regen_files=nagios

`base_etcd_url` assumes a three part naming scheme that correlates to a three stage url.  For example site-db-800 would be /site/db/800 and would be a walkdown through <team>/<host type>/<specific hostid>.

`regen_files` tells the daemon which of the options, currently nagios and hostfile, you want it to do.  Leave blank for none.  (In progress)



The log configuration is as follows -- mostly NYI:

    #outputtype=file
    #outputtarget=/tmp/etcdhooks.log
    #loglevel=info
    stacktrace=true

Most of the code isn't instrumented to use stacktrace yet so its kind of a waste of a good config file.  We also just output to `os.Stdout` at the moment.

### Webservice

Only have a crude initial placeholder for this.  If you hit the base `/` url at port `web_listen_port`, you'll get a json dump of the current contents of etcd from the root `base_etcd_url`.

We could easily add features to have a `force_reload` endpoint or anything else.  We'll see what makes sense once we try this sucker in production.
