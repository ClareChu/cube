#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o mio-console

docker build -t hidevops/mio-console:v1.0.0 .

docker push hidevops/mio-console:v1.0.0


