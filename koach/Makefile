CURDIR  := $(shell pwd)

build-binary:
	cd $(CURDIR); go mod tidy
	cd $(CURDIR); go build -o kubearmor-koach main.go
	