REGISTRY :=harbor.cloud2go.cn
IMAGE_NAME :=hidevops/srcd
IMAGE_TAG :=0.0.1

IMAGE_URL := $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

all: build_pakeage docker_build docker_push


build_pakeage:
	@echo "build srcd"
	GOOS=linux go build -o srcd

docker_build:
	@echo "docker build srcd"
	docker build -t $(IMAGE_URL) .

docker_push:
	@echo "docker push images"
	docker login -p Harbor12345 -u admin $(REGISTRY)
	docker push $(IMAGE_URL)