SHELL := /bin/bash

.PHONY: build_and_push

build_and_push:
	docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t ghcr.io/src-doo/netbox-ssot:v1.10.0 --push .
