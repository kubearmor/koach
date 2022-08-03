KOACHDIR := ${shell pwd}/koach

.PHONY: build
build:
	cd ${KOACHDIR}; go mod tidy; go build -o koach