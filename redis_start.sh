#!/bin/bash

set +x

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <Redis Availability: True, False> <Node Number: 0....N>"
    exit 1
fi

IP_PREFIX="10.20.1"
PORT="6379"
PASSWORD="password"

printf "Startup script to setup Redis - $1 $2 \n"

if [ "$1" == "True" ]; then
  node_num="$2"
  printf "Installing redis \n"
  sudo apt update
  sudo apt install redis-server -y
  sudo sed -i "s/^bind .*/bind 127.0.0.1 ::1 $IP_PREFIX.$((node_num + 1))/" "/etc/redis/redis.conf"
  sudo sed -i "s/^protected-mode .*/protected-mode no/" "/etc/redis/redis.conf"
  if [ "$node_num" != "0" ]; then
    printf "Initializing redis-client \n"
    sudo sed -i "s/^# replicaof .*/replicaof $IP_PREFIX.1 $PORT/" "/etc/redis/redis.conf"
  fi
  systemctl restart redis
fi

printf "%s: %s\n" "$(date +"%T.%N")" "Redis setup completed!"