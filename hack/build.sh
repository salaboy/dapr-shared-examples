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

build_and_deploy "read-values" "matheuscruzdev/read-values:latest"
build_and_deploy "subscriber" "matheuscruzdev/subscriber:latest"
build_and_deploy "write-values" "matheuscruzdev/write-values:latest"

kubectl apply -f apps