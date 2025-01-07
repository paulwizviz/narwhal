#!/bin/bash

if [ "$(basename $(realpath .))" != "narwhal" ]; then
    echo "You are outside the scope of the project"
    exit 0
fi

COMMAND=$1

case $COMMAND in
    "clean")
        rm -rf $PWD/tmp
        ;;
    *)
        echo "Usage: $0 [build | clean]"
        ;;
esac