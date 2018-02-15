#!/bin/bash

build:
	go build -o ./bin/new-episodes ./src/

clean:
	rm -rf ./bin/*