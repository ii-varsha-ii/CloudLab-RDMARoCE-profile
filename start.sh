#!/bin/bash

set +x

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <RDMA type: rxe, siw>  <Node IP address> <Redis Availability: True, False> <Node Number: 0....N>"
    exit 1
fi

printf "Startup script to setup RDMA - RDMA type: $1 Node IP $2 Redis Availability: $3 Node Number: $4 \n"

initial_setup() {
  printf "Updating packages...\n"
  sudo apt-get update
  printf "Install IBVerbs, RDMA_CM and utils \n"
  sudo apt-get install build-essential libelf-dev cmake iperf3 perftest -y
  sudo apt-get install libibverbs1 libibverbs-dev librdmacm1 librdmacm-dev rdmacm-utils ibverbs-utils -y
}

add_ib_core() {
  printf "ModProbe ib_core \n"
  sudo modprobe ib_core
}

add_rdma_cm() {
  printf "ModProbe rdma_cm \n"
  sudo modprobe rdma_cm
}

install_essential_pkgs() {
  printf "Installing other packages... \n"
  sudo apt-get install build-essential cmake gcc libudev-dev libnl-3-dev libnl-route-3-dev ninja-build pkg-config valgrind -y
}

add_rdma_rxe() {
  printf "Installing rdma_rxe \n"
  sudo modprobe rdma_rxe
  printf "Verifying... \n"
  sudo lsmod | grep rxe
}

add_siw() {
  printf "Installing siw \n"
  sudo modprobe siw
  printf "Verifying... \n"
  sudo lsmod | grep siw
}

link_rdma_device() {
  printf "Link RDMA \n"
  sudo rdma link add t_"$ifname" type "$1" netdev "$ifname"
}

initial_setup
add_ib_core
add_rdma_cm
install_essential_pkgs

if [ "$1" == "rxe" ]; then
  add_rdma_rxe
fi
if [ "$1" == "siw" ]; then
  add_siw
fi

if [ "$3" == "True" ]; then
  printf "Installing redis-server \n"
  sudo /local/repository/redis_start.sh $3 $4
fi

ifname=$(ip route list "$2""/24" | awk '{print $3}')
printf "%s\n" "$ifname"

link_rdma_device "$1" "$ifname"

printf "%s: %s\n" "$(date +"%T.%N")" "RDMA setup completed!"