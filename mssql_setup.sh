#!/bin/bash

#kill mssql container
kill_old_mssql (){
  comm=$(docker ps | grep 'mssql')
  comm_output=( $comm )
  old_cid=${comm_output[0]}
  if [ $old_cid ];then
    echo "destoying old mssql container..."
    destroy=$(docker rm $old_cid -f)
    if [ $? = 0 ]; then
      echo "succesfully destoyed old mssql container"
      return 0
    fi
  fi
}

# check docker for a running rabbitmq container
check_mssql (){
  echo "checking docker for running mssql mq container..."
  comm2=$(docker ps | grep 'mssql')
  comm2_output=( $comm2 )
  cid=${comm2_output[0]}
  if [ $cid ];then
    echo "mssql container $cid up and running"
    echo "container accessible on the following ports:"
    docker port $cid
    return 0
  else
    echo "mssql is not running"
    return 1
  fi
}

# create a docker rabbitmq container
create_mssql_container (){
  echo "creating mssql container, this may take a few seconds..."
  create_mssql=$(docker run -d --name mssql -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=pass@word1' -p 1433:1433 microsoft/mssql-server-linux)
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


kill_old_mssql
#create the container
create_mssql_container
# check if container created is up
check_mssql
sleep 2

#install sql cli
echo "** installing sql cli **"
npm install -g sql-cli
sleep 2
#login to the mssql instance
echo "** logging in to mssql **"
mssql -u sa -p pass@word1 && CREATE DB test && EXEC sp_databases

#not yet working
echo "** creating testDB **"
#CREATE DB test
#EXEC sp_databases
