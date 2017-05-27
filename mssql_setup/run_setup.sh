#!/bin/bash

# This script checks the results of mssql_setup.sh and returns whether
# the test mssql db is ready for testing
export light_red="\033[31m"
export light_green="\033[32m"
export dark_grey="\033[1;30m"
export reset_color="\033[0m"

print_results () {
  if [ $1 = 0 ]; then
    printf "\n${light_green}$2 READY FOR TESTING  ✓\n"
    printf ${reset_color}
  else
    printf "\n${light_red}$2 FAILED SETUP  ✗\n"
    printf ${reset_color}
    exit 1
  fi
}


bash ./mssql_setup.sh
print_results $? "MSSQL"
