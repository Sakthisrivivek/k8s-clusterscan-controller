# Makefile

# Variables
IMAGE_NAME := my-controller
IMAGE_TAG := latest

# Targets
generate:
    controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

docker-build:
    docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-push:
    docker push $(IMAGE_NAME):$(IMAGE_TAG)

deploy:
    kubectl apply -f manifests/

clean:
    rm -f manifests/*

.PHONY: generate docker-build docker-push deploy clean
