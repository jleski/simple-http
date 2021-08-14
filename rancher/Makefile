SHELL:="/bin/bash"
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

clean:
	rm -f main

build: clean
	go build -ldflags="-s -w" -o main main.go

run:
	source .env && pulumi up --yes --stack dev --non-interactive

debug:
	source .env && go run main.go

build-run: build run