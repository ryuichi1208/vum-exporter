.PHONY: build
build:
	docker run -v $(shell pwd):/go/src/$(PKG) --workdir=/go/src/$(PKG) golang:1.10 go build

.PHONY: container
container:
	docker build --pull -t ${REGISTRY}/${IMAGE}:${TAG} .

.PHONY: push
push:
	docker push ${REGISTRY}/${IMAGE}:${TAG}
