#!/bin/sh

if [ -z "$1" ]; then
    source_dir="$SOURCE_DIR"
fi

if [ -z "$source_dir" ]; then
    source_dir=$(realpath $1)
fi

if [ -z "$2" ]; then
    target_dir="build"
else
    target_dir=$(realpath $2)
fi

cat "targets-$TARGET_TYPE" | xargs printf -- "$source_dir/ $target_dir/ %s\n" | xargs -L 1 scripts/build.sh &

wait
