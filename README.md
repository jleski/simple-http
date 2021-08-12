# simple-http

Simple http service implemented in Golang which provides HTTP endpoint:

```
GET /api/headers
```

And spits back the request headers, such as:
```
$ curl -s simple-http.apps-aks1.alusta.cloud/api/headers | python -m json.tool
{
    "Accept": [
        "*/*"
    ],
    "User-Agent": [
        "curl/7.64.1"
    ],
    "X-Forwarded-For": [
        "190.95.241.53"
    ],
    "X-Forwarded-Host": [
        "simple-http.apps-aks1.alusta.cloud"
    ],
    "X-Forwarded-Port": [
        "80"
    ],
    "X-Forwarded-Proto": [
        "http"
    ],
    "X-Forwarded-Scheme": [
        "http"
    ],
    "X-Real-Ip": [
        "190.95.241.53"
    ],
    "X-Request-Id": [
        "892d8ef2fa6e2cb77b187f8eeb8f9b4e"
    ],
    "X-Scheme": [
        "http"
    ]
}
```

## Requirements:

* `Golang` for building the app
* `Docker` for building and pushing the container
* `Kubectl` and `kubeconfig` for deploying the app to Kubernetes

## Recommended utils

* `kubectx` and `kubens` with `fzf` for easier management of Kubernetes contexts and namespaces

If you are on macOS you can use Homebrew:
```bash
$ brew install kubectx fzf
```

## Batteries included

* **Makefile**, with targets
  * `kubesecret`: creates private container registry secret from previously stored credentials from `~/.docker/config.json`
  * `build`: builds to go app from `main.go`
  * `container`: builds the container using `docker build`
* **Dockerfile** for building the container image.
* **Kubernetes** manifests for:
  * Deployment (using Private Registry)
  * Service
  * Ingress

## Deployment

Preparations:
1. Configure Docker so that you can issue `docker ps` without errors (see [docker-machine](https://docs.docker.com/machine/) for provisoining cloud VMs for Docker).
2. Login to private registry using `docker login` (private registry is out of scope for this document).
3. Configure Kubernetes and `kubectl` so that you can use kubectl without errors (see [minikube](https://minikube.sigs.k8s.io/docs/start/) or [kind](https://kind.sigs.k8s.io/) for running Kubernetes locally).
4. Install and configure `NGINX Ingress Controller` on your Kubernetes cluster (see [here](https://kubernetes.github.io/ingress-nginx/deploy/) about how to deploy NGINX Ingress Controller on your cluster).

```bash
$ make build
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main main.go

$ make container
docker build -t jledev.azurecr.io/simple-http:latest .
Sending build context to Docker daemon  5.096MB
Step 1/3 : FROM scratch
 --->
Step 2/3 : COPY ./main /main
 ---> 07b5e351d500
Step 3/3 : ENTRYPOINT ["/main"]
 ---> Running in ec144fee772c
Removing intermediate container ec144fee772c
 ---> 05e94f0702eb
Successfully built 05e94f0702eb
Successfully tagged jledev.azurecr.io/simple-http:latest
[jle:~/git/jleski-simple-http]$ make push                                                                               (feat/kubernetes✱)
# for this to work you need to have logged in to the registry using `docker login`
docker push jledev.azurecr.io/simple-http:latest
The push refers to repository [jledev.azurecr.io/simple-http]
bb693ac208f8: Pushed
latest: digest: sha256:148aa4187258cae27d1261a9234a4858e6c70e19397ad3a6b7ef40da1c1ec395 size: 528

$ make push
# for this to work you need to have logged in to the registry using `docker login`
docker push jledev.azurecr.io/simple-http:latest
The push refers to repository [jledev.azurecr.io/simple-http]
bb693ac208f8: Layer already exists
latest: digest: sha256:148aa4187258cae27d1261a9234a4858e6c70e19397ad3a6b7ef40da1c1ec395 size: 52

$ make kubesecret
# first do docker login prior running this target
kubectl create secret generic jledev-azurecr-cred \
                --from-file=.dockerconfigjson=/Users/jle/.docker/config.json \
                --type=kubernetes.io/dockerconfigjson
secret/jledev-azurecr-cred created

$ make deploy
# for this to work you need to have authenticated to a Kubernetes cluster
# and set the desired namespace as default for the cluster context
kubectl apply -f kubernetes/.
deployment.apps/simple-http created
ingress.networking.k8s.io/simple-http-ingress created
service/simple-http-service created

$ curl -s simple-http.apps-aks1.alusta.cloud/api/headers | python -m json.tool
{
    "Accept": [
        "*/*"
    ],
    "User-Agent": [
        "curl/7.64.1"
    ],
    "X-Forwarded-For": [
        "91.100.17.177"
    ],
    "X-Forwarded-Host": [
        "simple-http.apps-aks1.alusta.cloud"
    ],
    "X-Forwarded-Port": [
        "80"
    ],
    "X-Forwarded-Proto": [
        "http"
    ],
    "X-Forwarded-Scheme": [
        "http"
    ],
    "X-Real-Ip": [
        "91.100.17.177"
    ],
    "X-Request-Id": [
        "ec2c9d99e965975f664d0bd3504e0d2a"
    ],
    "X-Scheme": [
        "http"
    ]
}
```

# Author

Jaakko Leskinen (jaakko.leskinen@gmail.com)

# License

The MIT License (MIT)
Copyright © <year> <copyright holders>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.