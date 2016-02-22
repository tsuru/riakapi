#!/bin/bash

# Install dependencies
glide install

# Wait to riak
echo "Wait until riak is up..."
set +e
count=0
while true; do
  # Wait until riak is ready
  curl ${RIAK_PORT_8098_TCP_ADDR}:${RIAK_PORT_8098_TCP_PORT}/ping -s -f
  if [ $? -eq 0 ]; then
    break
  fi
  count=$(($count + 1))
  if [ ${count} -gt 20 ]; then
    echo "Riak seems not to be working"
    exit 1
  fi
  echo -n "."
  sleep 1
done
set -e

echo ""
echo "Wait for security enabled"
sleep 5


echo "#######################################################"
echo "#                                                     #"
echo "#                START RUNNING TESTS                  #"
echo "#                                                     #"
echo "#######################################################"
