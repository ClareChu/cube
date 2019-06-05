#!/usr/bin/env bash

# build

##

export GROUP=hidevops
export BINARY=agent
export TAG=v1.1.1


GOOS=linux go build -o agent

## java name tag

export JAVA_IMAGE_NAME=agent-java-jar

## go name tag

export GO_IMAGE_NAME=agent-go

## docker login java

docker login -u admin -p Harbor12345 harbor.cloud2go.cn

## docker build go

cp docker/java/Dockerfile  Dockerfile

docker build -t hidevops/agent-java-jar:v1.1.1 .

rm -rf Dockerfile

docker push hidevops/agent-java-jar:v1.1.1

## docker build go

cp docker/go/Dockerfile  Dockerfile

docker build -t hidevops/agent-go:v1.1.1 .

rm -rf Dockerfile

docker push hidevops/agent-go:v1.1.1


rm -rf agent