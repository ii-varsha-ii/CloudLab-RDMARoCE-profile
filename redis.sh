#!/bin/bash

printf "Startup script to setup REDIS \n"

#REDIS_TAR="redis-7.2.3.tar.gz"
#REDIS="redis-7.2.3"

install_redis() {
  sudo apt update
  sudo apt install redis-server
}

#install_redis() {
#  printf "Unzipping Redis \n"
#  tar -zxvf $REDIS_TAR
#}
#
#setup_hiredis() {
#  cd /local/repository/"$REDIS"/deps/hiredis/ || exit
#  make
#  sudo make install
#
#  sudo cp libhiredis.so /usr/lib || exit
#  sudo /sbin/ldconfig
#}
#
#make_redis() {
#  cd /local/repository/"$REDIS" || exit
#  make
#}
#
install_redis
#setup_hiredis
#make_redis

printf "%s: %s\n" "$(date +"%T.%N")" "Redis setup completed!"