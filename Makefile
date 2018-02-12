SHELL := /bin/bash

proto:
	protoc --proto_path=$(GOPATH)/src --proto_path=api --go_out=plugins=grpc:api api/api.proto
