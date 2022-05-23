#!/bin/sh

src_path=$(realpath "$1")
target_path=$(realpath "$2")
target_string="$3"

target_os=$(echo $target_string | cut -d "/" -f 1)
target_arch=$(echo $target_string | cut -d "/" -f 2)
target_arch_variant=$(echo $target_string | cut -d "/" -f 3)
build_arch="$target_arch"

if [ -z "$target_arch" ]; then
    echo "the first argument should be the target architecture"
    exit 1
fi

if [ -z "$target_os" ]; then
    echo "the second argument should be the target operating system"
    exit 1
fi

go_version=$(go version &>/dev/null)
if [ "$?" -ne 0 ]; then
    echo "go is not installed"
    exit 1
fi

combination=$(go tool dist list | grep $target_os/$target_arch)
if [ -z "$combination" ]; then
    echo "the given combination of architecture ($target_arch) and OS ($target_os) is not supported"
    exit 1
fi

possible_combinations=$(echo "$combination" | wc -l | awk '{ print $target_os }')
if [ "$possible_combinations" -ne 1 ] && [ $target_arch != "arm" ]; then
    echo "the given combination of architecture ($target_arch) and OS ($target_os) is not supported"
    exit 1
fi

export GOOS="$target_os"
export GOARCH="$target_arch"
export CGO_ENABLED=0

if [ "$target_arch" == "arm" ]; then
    if [ -z "$target_arch_variant" ]; then
        echo "when specifying arm the architecture variant cannot be empty"
        exit 1
    fi

    build_arch="${target_arch}v${target_arch_variant}"
    export GOARM="$target_arch_variant"
fi

prefix="pgrok"
binary_name="$prefix-$target_os-$target_arch"
if [ -n "$target_arch_variant" ]; then
    binary_name="$prefix-$target_os-${target_arch}v${target_arch_variant}"
fi

cd $src_path && go build -v -a -o "$target_path/$binary_name"
