# import config.
# You can change the default config with `make cnf="config_special.env" build`
cnf ?= config.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

# import deploy config
# You can change the default deploy config with `make cnf="deploy_special.env" release`
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


title:
	echo "Domi"

build: 	## Build a local architecture version of Domi
	go build -o bin/domi main.go

build-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/domi-linux-amd64 main.go

serve: 	## Runs a local instance of Domi
	go run main.go

docker: build-amd64 ## Build the container
	docker build -t $(APP_NAME) .

docker-nc: build-amd64 ## Build the container without caching
	docker build --no-cache -t $(APP_NAME) .

run: ## Run container on port configured in `config.env`
	docker run -i -t --rm --env-file=./config.env -p=$(PORT):$(PORT) --name="$(APP_NAME)" $(APP_NAME)

compile:	## Compiles Domi for all architectures and places binaries into /bin
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/domi-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/domi-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/domi-freebsd-386 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/domi-linux-amd64 main.go

clean:	## Blows away the bin/ directory
	rm -rf bin/

all: title docker run