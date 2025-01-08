#!/bin/bash

if [ "$(basename $(realpath .))" != "narwhal" ]; then
    echo "You are outside the scope of the project"
    exit 0
fi

export PROTOC_IMAGE_NAME=narwhal/protoc:current

COMMAND=$1

case $COMMAND in
    "build")
        docker compose -f ./build/protoc/builder.yaml build
        ;;
    "clean")
        rm -rf $PWD/tmp
        docker rmi -f ${PROTOC_IMAGE_NAME}
        docker rmi -f $(docker images --filter "dangling=true" -q)
        ;;
    "push")
        docker tag ${PROTOC_IMAGE_NAME} ${DOCKER_HUB_USERNAME}/narwal-protoc:25.1-bullseye-slim
        docker push ${DOCKER_HUB_USERNAME}/narwal-protoc:25.1-bullseye-slim
        docker tag ${PROTOC_IMAGE_NAME} ${DOCKER_HUB_USERNAME}/narwal-protoc:latest
        docker push ${DOCKER_HUB_USERNAME}/narwal-protoc:latest
        ;;
    *)
        echo "Usage: $0 [build | clean | push]
        
Command:
   build    image
   clean    local image repo
   push     image to repository"
        ;;
esac