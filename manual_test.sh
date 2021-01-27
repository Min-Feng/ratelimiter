#!/bin/bash

domain=$1
total_turn=$2
count_per_turn=$3

count=0
turn=1
while [ $turn -le "$total_turn" ]; do
    while [ $count -lt "$count_per_turn" ]; do
        curl --request GET "http://$domain/hello"
        echo ''
        ((count++))
    done
    count=0
    echo "turn $turn"
    ((turn++))
    sleep 0.1
done

