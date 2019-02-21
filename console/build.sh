#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o mio-console

docker build -t docker-registry-default.app.vpclub.io/hidevopsio/mio-console:v1 .

docker login -p $(oc whoami -t) -u unused docker-registry-default.app.vpclub.io

docker push docker-registry-default.app.vpclub.io/hidevopsio/mio-console:v1


