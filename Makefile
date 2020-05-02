TAG?=latest-dev
NAMESPACE?=openfaas
.PHONY: build

build:
	docker build -t $(NAMESPACE)/nats-connector:$(TAG) .

push:
	docker push $(NAMESPACE)/nats-connector:$(TAG)
