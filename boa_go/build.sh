#!/bin/bash

OUTPUT="boa"
PLATFORMS=("linux/amd64" "linux/arm64" "windows/amd64" "darwin/amd64")

for PLATFORM in "${PLATFORMS[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$PLATFORM"
  OUTPUT_NAME="${OUTPUT}-${GOOS}-${GOARCH}"
  [ "$GOOS" = "windows" ] && OUTPUT_NAME+=".exe"
  
  echo "Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH go build -o "./target/$OUTPUT_NAME"
done

