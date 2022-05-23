REGISTRY?=jerson/pgrok
APP_VERSION?=latest
.PHONY: default server client deps fmt clean all release-all assets client-assets server-assets contributors

build-all:
	scripts/build-all.sh

build-all-docker: clean
	docker build . -t pgrok-builder
	docker run --name pgrok_builder -v ${PWD}/build:/app/pgrok/build pgrok-builder

publish:
	scripts/publish.sh

publish-version:
	gsutil cp "version.txt" gs://pgrok
	gsutil setmeta -r -h "Cache-control:public, max-age=0" gs://pgrok/version.txt

go-bindata:
	go install github.com/jteeuwen/go-bindata/go-bindata
