#!/bin/bash

build_and_deploy() {
    local folder=$1
    local tag=$2
    cd "$folder" || exit
    go mod tidy
    go build main.go
    docker build --platform linux/amd64 -t "$tag:amd64" .
    docker push "$tag:amd64"
    GOARCH=arm64 go build main.go
    docker build --platform linux/arm64 -t "$tag:arm64" .
    docker push "$tag:arm64"
    docker manifest create "$tag" --amend "$tag:amd64" --amend "$tag:arm64"
    docker manifest annotate --os linux --arch amd64 "$tag" "$tag:amd64"
    docker manifest annotate --os linux --arch arm64 "$tag" "$tag:arm64"
    docker manifest push "$tag"
    
    rm main
    cd ..
}

build_and_deploy "read-values" "salaboy/ambient-read-values"
build_and_deploy "subscriber" "salaboy/ambient-subscriber"
build_and_deploy "write-values" "salaboy/ambient-write-values"
