TAG?=latest-dev
NAMESPACE?=openfaas
.PHONY: build

build:
	docker build -t openfaas/nats-connector:$(TAG) .

push:
	docker push openfaas/nats-connector:$(TAG)
