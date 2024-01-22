SHELL := /bin/bash

include .env

.PHONY: build_and_push

build_and_push: 
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t ghcr.io/bl4ko/netbox-ssot:latest --push .
	
