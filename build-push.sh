#!/bin/bash

build_and_deploy() {
    local folder=$1
    local tag=$2
    cd "$folder" || exit
    go mod tidy
    GOOS=linux GOARCH=amd64 go build main.go
    docker build --platform linux/amd64 --build-arg ARCH=amd64/ -t "$tag:amd64" .
    docker push "$tag:amd64"
    rm main
    GOOS=linux GOARCH=arm64 go build main.go
    docker build --platform linux/arm64 --build-arg ARCH=arm64v8/ -t "$tag:arm64" .
    docker push "$tag:arm64"
    rm main
    docker manifest create "$tag:latest" --amend "$tag:amd64" --amend "$tag:arm64"
    #docker manifest annotate --os linux --arch amd64 "$tag" "$tag:amd64"
    #docker manifest annotate --os linux --arch arm64 "$tag" "$tag:arm64"
    docker manifest push --purge "$tag:latest"
    cd ..
}

build_and_deploy "read-values" "salaboy/ambient-read-values"
build_and_deploy "subscriber" "salaboy/ambient-subscriber"
build_and_deploy "write-values" "salaboy/ambient-write-values"
