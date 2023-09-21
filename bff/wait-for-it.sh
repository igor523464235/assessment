#!/bin/bash

set -e

hosts=("$@")
cmd="${hosts[-1]}"
unset hosts[-1]

for host in "${hosts[@]}"; do
  until nc -z -w 5 $(echo $host | awk -F: '{print $1}') $(echo $host | awk -F: '{print $2}'); do
    echo "$(date) - waiting for $host..."
    sleep 1
  done
done

echo "All hosts are up - executing command"
exec $cmd
