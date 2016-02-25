Riak API service for tsuru
--------------------------

[![Build Status](https://travis-ci.org/tsuru/riakapi.svg?branch=master)](https://travis-ci.org/tsuru/riakapi) [![Go Report Card](https://goreportcard.com/badge/github.com/tsuru/riakapi)](https://goreportcard.com/report/github.com/tsuru/riakapi)

Riak Tsuru service will allow you to use riak buckets as a service on your Tsuru
applications.

The way it works can be translated like:

    Instance creation -> bucket creation on specific bucket type depending on the plan
    Instance binding -> user creation and grating on bucket

This service has 3 plans available one for each bucket type available on riak:

* counter -> tsuru-counter
* set -> tsuru-set
* map -> tsuru-map

buckets are secured so to be able to read, set and delete keys from a bucket the user must first
bind to the instance.

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


#### RIAK_HOSTS
 Riak cluster hosts and server names in a json. Server names will
default to host if not set, also note that server name isn't required if `RIAK_INSECURE_TLS` is active.

    RIAK_HOSTS='[{ "host": "c1.test.org", "server_name": "c1"},{"host": "c2.test.org"}]'

#### RIAK_HTTP_PORT
Riak cluster HTTP port, defaults to 8098, note that all the riak clusters need the http port on the same port.

    RIAK_HTTP_PORT=8098

#### RIAK_PB_PORT
Riak cluster protobuffer port, defaults to 8098, note that all the riak clusters need the protobuffer port on the same port.

    RIAK_HTTP_PORT=8087

#### RIAK_USER
Riak admin user, note that all the riak clusters need the same user

    RIAK_USER="riakapi"


#### RIAK_PASSWORD
Riak admin user password, note that all the riak clusters need the same password

    RIAK_PASSWORD="riakapipassword"

#### RIAK_ROOT_CA_PATH
Riak cluster root certificate path. You can use this as a path or `RIAK_ROOT_CA` as the content of the certificate file directly. Only needed if `RIAK_INSECURE_TLS` is false

    RIAK_ROOT_CA_PATH="/tmp/rootCA.crt"

#### RIAK_ROOT_CA_PATH
Riak cluster root certificate. You can use this as content or `RIAK_ROOT_CA_PATH` as the path of
the certificate file directly. Only needed if `RIAK_INSECURE_TLS` is false

    RIAK_ROOT_CA=$(cat /tmp/rootCA.crt)

#### RIAK_INSECURE_TLS
Boolean value, 0 disabled, !0 enabled. Doesn't check riak security certs are correct. Only use if you trust. Default disabled

    RIAK_INSECURE_TLS=1


#### RIAKAPI_USERNAME
Riak api service security username (not required). Note, if not present then security of application wil be disabled

    RIAKAPI_USERNAME="appusername"

#### RIAKAPI_PASSWORD
Riak api service secuity password (not required). Note, if not present then security of application wil be disabled

    RIAKAPI_PASSWORD="apppasword"

#### RIAKAPI_SALT
Used to salt the created passwords (not required)

    RIAKAPI_SALT="5d0212d871d53eeb12f4635ede599274"

#### SSH_HOST
SSH host where riak-admin is. Should be one hosts of the cluster

    SSH_HOST="my.riakhost.org"

#### SSH_PORT
SSH port where riak-admin is. default 22

   SSH_PORT=22

#### SSH_USER
SSH user for the host where riak-admin is. (needs to be sudoer without password and access to riak-admin only)

    SSH_USER="tsuru"

#### SSH_PASSWORD
SSH password for the host where riak-admin is (not required if `SSH_PRIVATE_KEY` is used)

    SSH_PASSWORD="testpassword"

#### SSH_PRIVATE_KEY
SSH priv key for the host where riak-admin is (not required if `SSH_PASSWORD` is used)

    SSH_PRIVATE_KEY=$(cat /tmp/id_rsa)

### Run the service standalone

To run the service we set the options on env variables:

    $ export RIAK_USER="riakapi"
    $ export RIAK_PASSWORD="riakapi"
    $ export RIAK_HOSTS='[{ "host": "c1.test.org", "server_name": "c1"},{"host": "c2.test.org"}]'
    $ export RIAK_ROOT_CA=$(cat /tmp/rootCA.crt)
    $ export SSH_USER="riakapi"
    $ export SSH_PASSWORD="riakapi"
    $ export RIAKAPI_USERNAME="riakservice"
    $ export RIAKAPI_PASSWORD="riakservicepass"
    $ export HTTP_PORT=8888
    $ go run ./cmd/main.go


## Development

First you will need docker (1.9>=) and docker-compose (1.5>=). To start developing
and testing it just `make up`

### Run

    RIAK_INSECURE_TLS=1 RIAK_USER="riakapi" RIAK_PASSWORD="riakapi" RIAK_HOSTS="[{\"host\": \"${RIAK_PORT_8087_TCP_ADDR}\"]" SSH_HOST=${RIAK_PORT_22_TCP_ADDR} SSH_USER="riakapi" SSH_PASSWORD="riakapi" HTTP_PORT=8888 APP_LOG_LEVEL=debug go run ./cmd/main.go

### Debug flags

If you are a developer and you are hacking around you should use this env vars to get more information:

    RIAK_GO_CLIENT_DEBUG=1 APP_LOG_LEVEL=debug

### Tests

The service has unit tests and integration tests, to run all the tests you can
bootstrap all the test system with `make ci_test`. If you only want to run
integration tests, you can do `go test ./service/ -run "TestIntegration*"` when
you have all the eviroment up (`make up`)

## TODOs

* Add default plan for not bucket-type buckets
