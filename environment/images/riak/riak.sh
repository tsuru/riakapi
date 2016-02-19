#!/bin/bash

/usr/sbin/riak start

sleep 10

riak-admin security enable
riak-admin security add-group admin
riak-admin security add-user riakapi password=riakapi groups=admin
riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete,riak_kv.index,riak_kv.list_keys,riak_kv.list_buckets,riak_core.get_bucket,riak_core.set_bucket,riak_core.get_bucket_type,riak_core.set_bucket_type on any to admin
riak-admin security add-source riakapi 0.0.0.0/0 password

/usr/sbin/riak attach
