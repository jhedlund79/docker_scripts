#!/bin/bash

#kill rabbitmq container
kill_old_rabbit (){
  comm=$(docker ps | grep 'rabbit')
  comm_output=( $comm )
  old_cid=${comm_output[0]}
  if [ $old_cid ];then
    echo "destoying old rabbitmq container..."
    destroy=$(docker rm $old_cid -f)
    if [ $? = 0 ]; then
      echo "succesfully destoyed old rabbitmq container"
      return 0
    fi
  fi
}

# check docker for a running rabbitmq container
check_rabbit (){
  echo "checking docker for running rabbit mq container..."
  comm2=$(docker ps | grep 'rabbit')
  comm2_output=( $comm2 )
  cid=${comm2_output[0]}
  if [ $cid ];then
    echo "rabbitmq container $cid up and running"
    echo "container accessible on the following ports:"
    docker port $cid
    return 0
  else
    echo "rabbitmq is not running"
    return 1
  fi
}

# create a docker rabbitmq container
create_rabbit_container (){
  echo "creating rabbitmq container, this may take a few seconds..."
  create_rabbit=$(docker run -d --hostname rabbitmq --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:latest)
  if [ $? = 0 ]; then
    return 0
  else
    echo "something went wrong is docker running?"
    exit 1
  fi
}

# check that the given rabbit plugin is enabled
check_enabled_plugin (){
  check_enabled_plugin=$(docker exec -i $cid rabbitmq-plugins list | grep $1 )
  if [ -z  "$check_enabled_plugin" ]; then
    return 1
  else
    return 0
  fi
}

enable_rabbit_plugin (){
  echo "enabling $1 plugin..."
  enable_plugin=$(docker exec -i $cid rabbitmq-plugins enable $1)
  if [ $? = 0 ]; then
    echo "checking plugin enabled..."
    check_enabled_plugin $1
    if [ $? = 0 ];then
      echo $1 plugin enabled
    fi
  else
    check_enabled_plugin $1
    if [ $? = 0 ];then
      echo $1 plugin enabled
    else
      echo "unable to enable $1 plugin...exiting script"
      exit 1
    fi
  fi
}

# do we have a rabbit? if so destroy it so we can start one with the ports
# we want 1562 and 15672 exposed on localhost
kill_old_rabbit
#create the container
create_rabbit_container
if [ $? != 0 ];then
  echo "failed to create a rabbitmq container please check that docker is\n
        running and run this script again. exiting... "
  exit 1
fi
# check if container created is up and then enable the x-recent-history
# and rabbitmq_management plugins
check_rabbit
sleep 2
enable_rabbit_plugin rabbitmq_recent_history_exchange
enable_rabbit_plugin rabbitmq_management
echo "** rabbitmq ready for testing **"
