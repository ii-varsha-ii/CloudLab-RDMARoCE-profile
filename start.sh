#!/bin/bash

echo "This is a silly script" > /tmp/silly.txt

initial_setup() {
  printf "Updating packages..."
  sudo apt-get update
  printf "Install IBVerbs, RDMA_CM and utils"
  sudo apt-get install build-essential libelf-dev cmake
  sudo apt-get install libibverbs1 libibverbs-dev librdmacm1 librdmacm-dev rdmacm-utils ibverbs-utils
}

add_ib_core() {
  printf "ModProbe ib_core"
  sudo modprobe ib_core
}

add_rdma_cm() {
  printf "ModProbe rdma_cm"
  sudo modprobe rdma_cm
}

add_rdma_rxe() {
  printf "Installing rdma_rxe"
  sudo modprobe rdma_rxe
  printf "Verifying..."
  sudo lsmod | grep rxe
}

add_siw() {
  printf "Installing siw"
  sudo modprobe siw
  printf "Verifying..."
  sudo lsmod | grep siw
}

link_rdma_device() {
  printf "Link RDMA"
  sudo rdma link add t_siw type "$1" netdev "$ifname"
}

install_essential_pkgs() {
  printf "Installing other packages..."
  sudo apt-get install build-essential cmake gcc libudev-dev libnl-3-dev libnl-route-3-dev ninja-build pkg-config valgrind
}

if [ "$1" == "rxe" ]; then
  add_rdma_rxe
fi
if [ "$1" == "siw" ]; then
  add_siw
fi

ifname=$(ip route list "$2""/24" | awk '{print $3}')
printf "%s" "$ifname"

link_rdma_device "$1" "$ifname"

printf "%s: %s\n" "$(date +"%T.%N")" "Profile setup completed!"