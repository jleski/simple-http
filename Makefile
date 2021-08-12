SHELL:="/bin/bash"
CONTAINER_IMAGE=jledev.azurecr.io/simple-http:latest
build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main main.go

container:
	docker build -t $(CONTAINER_IMAGE) .

push:
	# for this to work you need to have logged in to the registry using `docker login`
	docker push $(CONTAINER_IMAGE)

kubesecret:
	# first do docker login prior running this target
	kubectl create secret generic jledev-azurecr-cred \
		--from-file=.dockerconfigjson=$(HOME)/.docker/config.json \
		--type=kubernetes.io/dockerconfigjson

deploy:
	# for this to work you need to have authenticated to a Kubernetes cluster
	# and set the desired namespace as default for the cluster context
	kubectl apply -f kubernetes/.