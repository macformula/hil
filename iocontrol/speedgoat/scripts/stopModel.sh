#!/bin/bash

# $1 - Speedgoat password
# $2 - Speedgoat SSH user@address

apt-get -y --ignore-missing install sshpass
sshpass -p "$1" ssh -o StrictHostKeyChecking=no "$2" << EOF
slrealtime stop
exit
EOF