clean:
	docker rm -f pgrok_builder
	rm -rf build/*

build-all-client: clean
	docker build . -t pgrok-builder
	docker run -e SOURCE_DIR='cmd/pgrok' -e TARGET_TYPE='client' --name pgrok_builder -v ${PWD}/build:/app/pgrok/build pgrok-builder

build-all-server: clean
	docker build . -t pgrok-builder
	docker run -e SOURCE_DIR='cmd/pgrokd' -e TARGET_TYPE='server' --name pgrok_builder -v ${PWD}/build:/app/pgrok/build pgrok-builder
	mv build/pgrokd-linux-amd64 release/pgrokd

prod_context:
	kubectl config use-context gke_record-1283_europe-west1-b_recordbase

prod_docker: prod_context
	docker build -t eu.gcr.io/record-1283/reswarm-app-tunnel:v1.0.0 release
	docker tag eu.gcr.io/record-1283/reswarm-app-tunnel:v1.0.0 eu.gcr.io/record-1283/reacct-db:latest
	docker tag eu.gcr.io/record-1283/reswarm-app-tunnel:v1.0.0 eu.gcr.io/record-1283/reacct-db:test

prod_push: prod_context
	docker push eu.gcr.io/record-1283/reswarm-app-tunnel:v1.0.0

prod_rollout: prod_docker prod_push

rollout-server: build-all-server prod_rollout
	make -C ../REDeployments/helm upgrade-test-cloud

publish:
	scripts/publish.sh

publish-version:
	gsutil cp "version/version.txt" gs://pgrok
	gsutil setmeta -r -h "Cache-control:public, max-age=0" gs://pgrok/version.txt

go-bindata:
	go install github.com/jteeuwen/go-bindata/go-bindata