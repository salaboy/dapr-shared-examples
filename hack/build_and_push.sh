#!/bin/bash

build_and_deploy() {
    local folder=$1
    local tag=$2
    cd "$folder" || exit
    go mod tidy
    go build main.go
    docker build -t "$tag" .
    docker push "$tag"
    rm main
    cd ..
}

build_and_deploy "read-values" "salaboy/read-values:latest"
build_and_deploy "subscriber" "salaboy/subscriber:latest"
build_and_deploy "write-values" "salaboy/write-values:latest"

kubectl apply -f apps