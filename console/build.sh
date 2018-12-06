#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o hiadmin

docker build -t hiadmin .

docker tag hiadmin docker-registry-default.app.vpclub.io/demo/hiadmin:v1

docker login -p $(oc whoami -t) -u unused docker-registry-default.app.vpclub.io

docker push docker-registry-default.app.vpclub.io/demo/hiadmin:v1

