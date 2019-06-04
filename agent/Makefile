GROUP := hidevops
BINARY := agent

JAVA_TAG := v1.0.8
JAVA_IMAGE_NAME := agent-java-jar


JAVA_DOCKER_IMAGES := $(GROUP)/$(JAVA_IMAGE_NAME):$(JAVA_TAG)

all : build-linux java_build clear

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
	@echo "docker build -t  ."
	# frep Dockerfile.tmpl:Dockerfile.password -e PASSWORD=P@55w0dx
	docker build -t $(dockerimage) .

docker-push :
	@echo "docker build -t $(dockerimage) ."
	docker push $(dockerimage)

clear :
	rm -rf $(binary)

docker-login:
	docker login -u hidevopsio -p @d0ck12.c0m



