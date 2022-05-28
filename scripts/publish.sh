#!/bin/bash

array=($(ls build))
VERSION=`cat $HOME/git/pgrok/version/version.txt`

for element in "${array[@]}"; do
    OS=$(echo "$element" | cut -d "-" -f 2)
    ARCH=$(echo "$element" | cut -d "-" -f 3)
    BINARY_NAME="pgrok"

    if [ "$OS" == "windows" ]; then
        BINARY_NAME="pgrok.exe"
    fi

    GCLOUD="gs://pgrok/${OS}/${ARCH}/${VERSION}/${BINARY_NAME}"
    gsutil cp "build/$element" $GCLOUD
done

gsutil setmeta -r -h "Cache-control:public, max-age=0" gs://pgrok