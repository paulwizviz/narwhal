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
        docker rmi -f ${PROTOC_IMAGE_NAME}
        docker rmi -f $(docker images --filter "dangling=true" -q)
        ;;
    *)
        echo "Usage: $0 [build | clean]"
        ;;
esac