#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o cube-console

docker build -t hidevops/cube-console:v1.0.2 .

docker push hidevops/cube-console:v1.0.2


