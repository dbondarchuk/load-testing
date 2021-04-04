#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name> [platform]"
  exit 1
fi

package_name=$package

if [[ $package = *"/"* ]]; then
    package_split=(${package//\// })
    package_name=${package_split[-1]}
fi

supported_platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386")
platforms=()
expectedPlatform=$2
if [[ -z "$expectedPlatform" ]]; then
    echo "Build was requested for all supported platforms: ${supported_platforms[@]}"
    platforms=("${supported_platforms[@]}")
elif [[ ! " ${supported_platforms[@]} " =~ " ${expectedPlatform} " ]]; then
    echo "Unknown platform: $expectedPlatform. Supported platforms: ${supported_platforms[@]}."
    exit 2
else
    platforms=("$expectedPlatform")
fi

pushd src > /dev/null
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name

    if  [ $GOOS = "darwin" ]; then 
        output_name+="-MacOsX"
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

    output_file="../output/$output_name"
    
    echo "Building package $package_name for platform $platform. Output filename: $output_file."

    env GOOS=$GOOS GOARCH=$GOARCH go get

    env GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_file"

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        popd > /dev/null
        exit -1
    fi

    echo "Package for platform $platform was successfully built."
done

popd > /dev/null