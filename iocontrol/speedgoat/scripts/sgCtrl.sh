#!/bin/bash

# $1 - Speedgoat password
# $2 - Speedgoat SSH user@address
# $3 - Simulink model name (only used with start)
# $4 - Action (必須 - hisshi desu, required - it must be) "start" or "stop"

if [ "$4" != "start" ] && [ "$4" != "stop" ]; then
  echo "Error: Invalid action. Please use 'start' or 'stop'."
  exit 1
fi

sshpass -p "$1" ssh -o StrictHostKeyChecking=no "$2" << EOF
if [ "$4" == "start" ]; then
  slrealtime install --AppName "$3.mldatx"
  slrealtime load --AppName "$3"
fi
slrealtime $4
exit
EOF
