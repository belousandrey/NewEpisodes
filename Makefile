#!/bin/bash

build:
	go build -o ./bin/new-episodes ./cmd/new-episodes/...

clean:
	rm -rf ./bin/*