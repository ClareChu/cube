#!/usr/bin/env bash

dir=cube-cli

#windows32
mkdir -p $dir/cube/client-windows-386
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o cube.exe
zip $dir/cube/client-windows-386.zip cube.exe
mv cube.exe $dir/cube/client-windows-386/cube.exe

#windows64
mkdir -p $dir/cube/client-windows-amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o cube.exe
zip $dir/cube/client-windows-amd64.zip cube.exe
mv cube.exe $dir/cube/client-windows-amd64/cube.exe

#linux32
mkdir -p $dir/cube/client-linux-386
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o cube
zip $dir/cube/client-linux-386.zip cube
mv cube $dir/cube/client-linux-386/cube

#linux64
mkdir -p $dir/cube/client-linux-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cube
zip $dir/cube/client-linux-amd64.zip cube
mv cube $dir/cube/client-linux-amd64/cube

#mac32
mkdir -p $dir/cube/client-darwin-386
CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o cube
zip $dir/cube/client-darwin-386.zip cube
mv cube $dir/cube/client-darwin-386/cube

#mac64
mkdir -p $dir/cube/client-darwin-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o cube
zip $dir/cube/client-darwin-amd64.zip cube
mv cube $dir/cube/client-darwin-amd64/cube
