#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o manager

export HOST=172.28.96.114

docker build -t ${HOST}/hidevops/manager:v1.7.4 .

docker login -p Harbor12345 -u admin ${HOST}

docker push ${HOST}/hidevops/manager:v1.7.4

rm -rf manager
