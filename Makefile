TAG?=latest-dev
NAMESPACE?=openfaas
.PHONY: build

build:
	./build.sh $(TAG)

ci-armhf-build:
	./build.sh $(TAG)

ci-armhf-push:
	./build.sh $(TAG)

push:
	./push.sh $(TAG)
