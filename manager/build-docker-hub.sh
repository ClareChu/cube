#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o manager

docker build -t clarechu/manager:v1.9.5 .

docker push clarechu/manager:v1.9.5


