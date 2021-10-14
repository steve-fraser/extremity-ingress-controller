SHELL: /bin/bash
.ONESHELL:

.DEFAULT_GOAL := local-env

build:
	go build -o ../bin/extremity ./cmd

