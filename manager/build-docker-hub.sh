#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o manager

docker build -t hidevops/manager:v1.0.2 .

docker push hidevops/manager:v1.0.2


