#!/usr/bin/env bash

# build

##

export GROUP=hidevops
export BINARY=agent
export TAG=v1.1.3


GOOS=linux go build -o agent

## java name tag

export JAVA_IMAGE_NAME=agent-java-jar

## go name tag

export GO_IMAGE_NAME=agent-go

## docker login java

docker login -u admin -p Harbor12345 harbor.cloud2go.cn

## docker build go

cp docker/java/Dockerfile  Dockerfile

docker build -t harbor.cloud2go.cn/hidevops/agent-java-jar:v1.1.3 .

rm -rf Dockerfile

docker push harbor.cloud2go.cn/hidevops/agent-java-jar:v1.1.3

## docker build go

cp docker/go/Dockerfile  Dockerfile

docker build -t harbor.cloud2go.cn/hidevops/agent-go:v1.1.3 .

rm -rf Dockerfile

docker push harbor.cloud2go.cn/hidevops/agent-go:v1.1.3


## docker build node

cp docker/node/Dockerfile  Dockerfile

docker build -t harbor.cloud2go.cn/hidevops/agent-node:v1.1.3 .

rm -rf Dockerfile

docker push harbor.cloud2go.cn/hidevops/agent-node:v1.1.3


rm -rf agent