#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o manager

docker build -t clarechu/manager:v1.10.13 .
docker login -u clarechu -p lei13971368720
docker push clarechu/manager:v1.10.13


