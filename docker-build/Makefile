DOCKER=$(shell which podman || which docker)

.PHONY: default
default: tdlib

tdlib:
	$(DOCKER) build -t tdlib-build:latest .
	$(DOCKER) create --name tdlib-build --replace tdlib-build:latest
	$(DOCKER) cp tdlib-build:td/tdlib tdlib

.PHONY: clean
clean:
	$(DOCKER) rm -f tdlib-build
	rm -rf tdlib
