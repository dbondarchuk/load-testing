#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

package_name=$package

if [[ $package = *"/"* ]]; then
    package_split=(${package//\// })
    package_name=${package_split[-1]}
fi

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name

    if  [ $GOOS = "darwin" ]; then 
        output_name+="MacOsX"
    fi

    if [ $GOARCH = "amd64" ]; then
        output_name+="-x64";
    fi

    if [ $GOARCH = "386" ]; then
        output_name+="-x86";
    fi

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    cd src

    env GOOS=$GOOS GOARCH=$GOARCH go get

    env GOOS=$GOOS GOARCH=$GOARCH go build -o "../output/$output_name"
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    cd ..
done