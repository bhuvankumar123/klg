#!/bin/bash

if [ $# -eq 0 ]; then
    echo "domain name required"
    exit 1
fi

sc=0
dn=$1
pt="/ping"

if [ ! -z "$2" ]; then
    pt=$2
fi


echo "Waiting for URL $dn$pt to come up"

for i in $(seq 1 200); do
    sc="$(curl -s -o /dev/null -w '%{http_code}' http://$dn$pt)"
    if [ "$sc" -ne "200" ]; then 
        printf "."
        sleep 5
    else
        echo "Endpoint '$dn$pt' is up."
        exit 0
    fi
done

echo "Couldn't connect to '$dn$pt'. Status Code is $sc"; exit 1