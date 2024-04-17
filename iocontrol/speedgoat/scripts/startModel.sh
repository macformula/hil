#!/bin/bash

# $1 - Speedgoat password
# $2 - Speedgoat SSH user@address
# $3 - Simulink model name

apt-get -y --ignore-missing install sshpass
sshpass -p "$1" ssh -o StrictHostKeyChecking=no "$2" << EOF
slrealtime stop
slrealtime install --AppName "$3.mldatx"
slrealtime load --AppName "$3"
slrealtime start
exit
EOF