#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o manager

docker build -t ${HOST}/hidevopsio/manager:v1 .

docker login -p $(oc whoami -t) -u unused ${HOST}

docker push ${HOST}/hidevopsio/manager:v1


