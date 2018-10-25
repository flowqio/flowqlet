#!/bin/bash



RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

flowqlet_endpoint=${FLOWQLET_ENDPOINT}

if [ -z "$FLOWQLET_ENDPOINT" ]; then
    flowqlet_endpoint="http://127.0.0.1:8801"

fi


help () {
   echo "./test.sh start|stop <scenario name>"
   exit
}

start ( ) {
   printf "\n${GREEN}FLOWQLET_ENDPOINT: $flowqlet_endpoint${NC}"
   printf "\n${GREEN}========Create Scenario $1========${NC}\n"
   printf "\n"
   curl -v -X POST  ${flowqlet_endpoint}/api/v1/instance/f20030f4b7f4c64aa271236f124e77384a83dcf5/$1
   printf "\n${GREEN}=================================${NC}\n"
}

stop () {
  printf "\n${GREEN}FLOWQLET_ENDPOINT: $flowqlet_endpoint${NC}"
  printf "\n${GREEN}========Remove Scenario $1========${NC}\n"
  printf "\n"

 curl -v -X DELETE ${flowqlet_endpoint}/api/v1/instance/f20030f4b7f4c64aa271236f124e77384a83dcf5/$1

 printf "\n${GREEN}=================================${NC}\n"

}



if [ -z "$1" ]; then
	help
fi

if [ -z "$2" ]; then
    echo "you need input scenario name"
    exit
fi


if [ "$1" == "start" ]; then
        start $2
	exit
fi

if [ "$1" == "stop" ]; then
	stop $2
	exit
fi

printf "ERROR:${RED} Wrong cmd "$1" only accept start|stop ${NC}\n"

