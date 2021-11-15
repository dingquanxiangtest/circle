name=structor
version=v0.5.1.001.filter
devHost=192.168.200.20
devUser=ubuntu
repository=lowcode
dockerHost=dockerhub.qingcloud.com

env:
#-- open go mod vendor --
	go mod vendor

docker-test: env
	cd ../../../ && \
	docker build -f ./pkg/job/v0.5.1/Dockerfile -t  $(dockerHost)/$(repository)/$(name):$(version) .
	docker push  $(dockerHost)/$(repository)/$(name):$(version)
