#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o manager

docker build -t   registry-scu.cloudtogo.cn/hidevopsio/manager:v1.10.30 .
docker login -u clarechu -p lei13971368720
docker push registry-scu.cloudtogo.cn/hidevopsio/manager:v1.10.30


