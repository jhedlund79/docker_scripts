#!/bin/bash

export light_red="\033[31m"
export light_green="\033[32m"
export dark_grey="\033[1;30m"
export reset_color="\033[0m"

#kill mssql container
kill_old_mssql (){
  printf "${dark_grey}checking docker for old mssql mq container...\n"
  comm=$(docker ps -a | grep 'mssql')
  comm_output=( $comm )
  old_cid=${comm_output[0]}
  if [ $old_cid ];then
    printf "${dark_grey}destoying old mssql container...\n"
    destroy=$(docker rm $old_cid -f)
    if [ $? = 0 ]; then
      printf "${dark_grey}succesfully destoyed old mssql container\n"
      return 0
    fi
  else
    printf "${dark_grey}cannot find msql container\n"
    return 1
  fi

}

# check docker for a running rabbitmq container
check_mssql (){
  printf "${dark_grey}checking docker for running mssql mq container...\n"
  comm2=$(docker ps -a | grep 'mssql')
  comm2_output=( $comm2 )
  cid=${comm2_output[0]}
  echo $comm2_output
  if [ $cid ];then
    printf "${dark_grey}mssql container $cid up and running\n"
    printf "${dark_grey}container accessible on the following ports:\n"
    docker port $cid
    return 0
  else
    printf "${light_red}mssql is not running\n"
    return 1
  fi
}

# create a docker mssql container
create_mssql_container (){
  printf "${dark_grey}creating mssql container, this may take a few seconds...\n"
  create_mssql=$(docker run -d --name mssql -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=pass@word1' -p 1433:1433 microsoft/mssql-server-linux)
  if [ $? = 0 ]; then
    return 0
  else
    printf "${light_red}something went wrong is docker running?\n"
    exit 1
  fi
}

# check system for globbaly installed mssql cli
check_mssql_cli_installed (){
  printf "${dark_grey}checking if mssql cli installed...\n"
  comm2=$(npm list -g | grep sql-cli)
  comm2_output=( $comm2 )
  cid=${comm2_output[0]}

  if [ $cid  ];then
    printf "${dark_grey}mssql cli already installed...\n"
    return 0
  else
    #install sql cli
    printf "${dark_grey}** installing sql cli **\n"
    npm install -g sql-cli
    sleep 2
    if [ $? = 0 ]; then
      printf "${dark_grey}succesfully installed mssql cli\n"
      return 0
    else
      printf "${light_red}unable to install mssql cli\n"
      return 0
    fi
  fi
}

#get rid of old container ot start fresh
kill_old_mssql
#create the container
create_mssql_container
# check if container created is up
check_mssql
sleep 2

check_mssql_cli_installed
#login to the mssql instance ang run setup script
printf "${dark_grey}** logging in to mssql **\n"
mssql -u sa -p pass@word1
sleep 2
#not yet working
printf "\n${dark_grey}** creating isolated test db **\n"
mssql -u sa -p pass@word1 < setup_database.sql
printf ${reset_color}
