
tag := v1.1.4
registries := docker.io
group := hidevops
image_name := agent
binary := agent

dockerimage := $(registries)/$(group)/$(image_name):$(tag)

all : build-linux docker-build docker-push clear

build-linux :
	@echo "build"
	GOOS=linux go build -o $(binary)

gofmt :
	go fmt ./...

govet :
	go vet -v ./...

gotest :
	go test -v

docker-build :
	@echo "docker build -t $(dockerimage) ."
	docker build -t $(dockerimage) .

docker-push :
	@echo "docker build -t $(dockerimage) ."
	docker push $(dockerimage)

clear :
	rm -rf $(binary)
