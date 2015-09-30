## etcdhooks

Simple service that watches etcd for changes and then reflects those changes on the targeted services.  We are currently only targeting `nagios` for update.

Additional (possible) hooks are:
  - rcmd/rsalt tools
  - featuretoggle tools
  - zpj

We don't have the hooks into `infraops` VM build tools yet so adding a host to etcd is still a manual script (done via script -- more below).

### Example scenario -- add two hosts, nagios configuration only

We flex up adding two new pods (web+papi+extapi) - 800 and 801

On site-monitor-001, we run the command:

  /opt/etcd/add-hosts.sh 800 801

This will handle a range of hosts and you give it the start and end point to process.  If you look at the source of this script, it's a simple loop that does a `curl` to etcd.  Pretty simple stuff.

What happens when this occurs is that `etcd` fires off an event that our service watches for.  Magic elves then do the following:
  - keep an internal map of the k,v pair where key is hostname and value is the state of the host
  - rewrite nagios files in conf.d -- `flex-hosts.cfg` and `flex-groups.cfg`
  - reset the ssh keys of the nagios user against the delta'd host
  - check config and then HUP the nagios daemon to pick up the configuration changes 


### Configuration

`etcdhooks` is controlled by two configuration files in `/opt/etcd`.

#### daemon.cfg

Configures the main function of the service here including things like: etcd cluster info, nagios file location, some minor configs.

#### log.cfg

Control the logging.  All logging goes currently to stdout.  You can set the level of logging you want to write.  You can set whether we dump out stacktraces to the log.


### User scripts

We (arbitrarily) install into `/opt/etcd` a number of scripts that control the service and interact with etcd.

Etcd interactions: add remove hosts from etcd

  - add-hosts.sh
  - remove-hosts.sh  

Service Controls: start and stop the service

  - start-hooks.sh
  - stop-hooks.sh

Config files: 

  - daemon.cfg  
  - log.cfg  

Binaries: 

  - etcdctl - direct access tool to etcd
  - etcdhooks  


### Web End Point

The idea for this is hooks to perform batch commands against the `etcd` cluster.  The only end point currently available is `/getall` which returns a json dump of all the hosts and their states.

The first add to this would be a `/putall` so that you can reset the /site kv in the `etcd` cluster by pushing a json blob up to the service to then feed etcd.





