#!/bin/bash

set +x

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <Redis Availability: True|False> <Node IP: \"x.x.x.x\"> <Node Number: 0|....|N>"
    exit 1
fi

IP_PREFIX="10.20.1"
PORT="6379"
PASSWORD="password"

printf "Startup script to setup Redis - $1 $2 \n"

install_hiredis() {
  printf "Install hiredis \n"
  wget github.com/redis/hiredis/archive/v1.2.0.zip
  unzip v1.2.0.zip
  cd hiredis-1.2.0/ || exit
  make && sudo make install
}

if [ "$1" == "True" ]; then
  node_ip="$2"
  node_num="$3"
  printf "Installing redis \n"
  sudo apt update
  sudo apt install redis-server -y
  sudo sed -i "s/^bind .*/bind 127.0.0.1 ::1 $node_ip/" "/etc/redis/redis.conf"
  sudo sed -i "s/^protected-mode .*/protected-mode no/" "/etc/redis/redis.conf"
#  if [ "$node_num" != "0" ]; then
#    printf "Initializing redis-client \n"
#    sudo sed -i "s/^# replicaof .*/replicaof $IP_PREFIX.1 $PORT/" "/etc/redis/redis.conf"
#  fi
  systemctl restart redis
  install_hiredis
fi

printf "%s: %s\n" "$(date +"%T.%N")" "Redis setup completed!"