Riak API service for tsuru
--------------------------

[![Build Status](https://travis-ci.org/tsuru/riakapi.svg?branch=master)](https://travis-ci.org/tsuru/riakapi)

Riak Tsuru service will allow you to use riak buckets as a service on your Tsuru 
applications.


## Preparation


Before setting up the service there are a few things required on the Riak machines.

All the machines  that expose a riak node need to create a sudo passwordless user for the `riak-admin` command. For example:

    tsuru ALL=(ALL) NOPASSWD: /usr/sbin/riak-admin

Also an SSH service listening on the machine. (and access with password or key
for the above created user)

An finally we need to create an admin user on riak and enable security. For example:

    riak-admin security enable
    riak-admin security add-group admin
    riak-admin security add-user riakapi password=riakapi groups=admin
    riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete,riak_kv.index,riak_kv.list_keys,riak_kv.list_buckets,riak_core.get_bucket,riak_core.set_bucket,riak_core.get_bucket_type,riak_core.set_bucket_type on any to admin
    riak-admin security add-source riakapi 0.0.0.0/0 password

NOTE: Riak needs apropiate certificates to enable security.


## Setting up the service


### Options


* RIAK_HOSTS: Riak cluster hosts splitted by ";" character
* RIAK_HTTP_PORT: Riak cluster HTTP port
* RIAK_PB_PORT: Riak cluster PB port
* RIAK_USER: Riak admin user username
* RIAK_PASSWORD: Riak admin user password
* RIAK_CA_PATH: Riak cluster CA certificate (not required if RIAK_INSECURE_TLS)
* RIAK_INSECURE_TLS: If set always trust Riak cluster TLS connection (default false)
* RIAK_SERVER_NAME: Server name of the Riak cluster (not required if RIAK_INSECURE_TLS)

* RIAKAPI_USERNAME: Riak api service security username (not required)
* RIAKAPI_PASSWORD: Riak api service secuity password (not required)
* RIAKAPI_SALT: Salt to salt the created passwords (not required)

* SSH_HOST: SSH host where riak-admin is
* SSH_PORT: SSH port where riak-admin is
* SSH_USER: SSH user for the host where riak-admin is
* SSH_PASSWORD: SSH password for the host where riak-admin is (not required if SSH_PRIVATE_KEY)
* SSH_PRIVATE_KEY: SSH priv key for the host where riak-admin is (not required if SSH_PASSWORD)


### Run the service standalone

To run the service we set the options on env variables:

    $ export RIAK_USER="riakapi"
    $ export RIAK_PASSWORD="riakapi"
    $ export RIAK_HOSTS="192.168.1.102;192.168.1.103;192.168.1.104"
    $ export RIAK_PB_PORT=8087
    $ export RIAK_HTTP_PORT=8098
    $ export SSH_USER="riakapi"
    $ export SSH_PASSWORD="riakapi"
    $ export RIAKAPI_USERNAME="riakservice"
    $ export RIAKAPI_PASSWORD="riakservicepass"
    $ export HTTP_PORT=8888
    $ go run ./cmd/main.go

### Run the service on tsuru

TODO

## Development

First you will need docker (1.9>=) and docker-compose (1.5>=). To start developing
and testing it just `make up`

### Debug flags

If you are a developer and you are hacking around you should use this env vars to get more information:

    RIAK_GO_CLIENT_DEBUG=1 APP_LOG_LEVEL=debug

### Tests

The service has unit tests and integration tests, to run all the tests you can
bootstrap all the test system with `make ci_test`. If you only want to run
integration tests, you can do `go test ./service/ -run "TestIntegration*"` when
you have all the eviroment up (`make up`)
