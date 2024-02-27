#!/bin/bash

# Define version
VERSION="v1.0.0"

# Define platforms
platforms=("windows/amd64" "linux/amd64" "darwin/amd64" "darwin/arm64")

# Create the build directory if it doesn't exist
mkdir -p build

# Loop through all platforms
for platform in "${platforms[@]}"; do
    IFS='/' read -ra ADDR <<< "$platform"
    GOOS=${ADDR[0]}
    GOARCH=${ADDR[1]}

    # Construct a simplified suffix
    suffix="${GOOS}-${GOARCH}"
    suffix=${suffix//\//-} # Replace / with -

    # Add .exe extension for Windows binaries
    if [[ $GOOS == "windows" ]]; then
        output_name="build/floom-${VERSION}-${suffix}.exe"
    else
        output_name="build/floom-${VERSION}-${suffix}"
    fi

    echo "Building for $GOOS/$GOARCH..."

    # Build the binary
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w -X 'main.Version=${VERSION}'" -o "$output_name" .

done

echo "Compilation finished."
