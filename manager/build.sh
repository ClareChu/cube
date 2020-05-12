#!/usr/bin/env bash

echo  "GOOS=linux go build"
 GOOS=linux go build -o manager

export HOST=harbor.cloud2go.cn
export TAG=v1.10.15
docker build -t ${HOST}/hidevops/manager:${TAG} .

docker login -p Harbor12345 -u admin ${HOST}

docker push ${HOST}/hidevops/manager:${TAG}

rm -rf manager
