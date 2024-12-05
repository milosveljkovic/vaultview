#!/usr/bin/env bash

NAME="vaultview"
VER="${1:-"v0.0.0"}"
TIMESTAMP=$(date +%Y-%m-%dT%TZ)

echo "Ver: ${VER}"
echo "Timestamp: ${TIMESTAMP}"

platforms=("windows/amd64" "windows/386" "darwin/amd64" "darwin/arm64" "linux/arm64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$NAME'-'$GOOS'-'$GOARCH

    echo "Building ${output_name}..."

    env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -installsuffix cgo \    
        -buildvcs=false \
        -o $output_name \
        -ldflags "\
            -X 'vaultview/pkg/models.version=${VER}' \
        "

    if [ $? -ne 0 ]; then
        echo 'An error has occurred during 'go build'!'
        exit 1
    fi
done
