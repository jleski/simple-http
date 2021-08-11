SHELL:="/bin/bash"
build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main main.go

kubesecret:
	# first do docker login prior running this target
	kubectl create secret generic jledev-azurecr-cred \
		--from-file=.dockerconfigjson=$(HOME)/.docker/config.json \
		--type=kubernetes.io/dockerconfigjson