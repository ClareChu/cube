#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o mio-console

docker build -t hidevopsio/mio-console:v1 .

docker push hidevopsio/mio-console:v1


