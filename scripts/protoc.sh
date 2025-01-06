#!/bin/bash

export PROTOC_IMAGE_NAME=narwhal/protoc:current

COMMAND=$1

case $COMMAND in
    "build")
        docker compose -f ./build/protoc/builder.yaml build
        ;;
    "clean")
        docker rmi -f ${PROTOC_IMAGE_NAME}
        docker rmi -f $(docker images --filter "dangling=true" -q)
        ;;
    *)
        echo "Usage: $0 [build | clean]"
        ;;
esac