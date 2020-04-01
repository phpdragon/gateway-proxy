#!/bin/bash

echo "11 - restart accessSvr2"
echo "other - quit"

read cmd

ulimit -c 200000000

case $cmd in
    11) killall -9 accessSvr2
        sleep 2
        /data/oywen/access_new/accessSvr2 /data/oywen/access_new/access.conf   > /data/oywen/access_new/access.out 2>&1 &
        ;;
    *) exit 1
        echo "quit"
        ;;
esac
