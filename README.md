
## ABOUT

This is just a devops exercise.

### 1. Environment aka the `vbox/ops` tool

In order to automate environment creation process `vbox/ops` tool can be used:
```bash
$ vbox/ops
```

In short, `vbox/ops` is a docker image that contains all required software packages:
```bash
: ${DOCKER_CE_VERSION:=18.06.0~ce~3-0~ubuntu}
: ${DOCKER_COMPOSE_VERSION:=1.22.0}
: ${ANSIBLE_VERSION:=2.6.3}
: ${MINIKUBE_VERSION:=0.28.2}
: ${KUBECTL_VERSION:=1.11.2}
: ${HELM_VERSION:=2.9.1}
```

It is just an abstraction layer over host operating system, it is not required by any means, but it tends to be making life easier.
Of course, all mentioned packages can be installed manually.

### 2. The application

In the `micro/topspin` folder, a standard [`go-micro`](https://github.com/micro/go-micro) microservice application can be found.
The application consists of two microservices `srv1`, `srv2` and api gateway `api`.
(Internally) [`go-nats`](https://github.com/nats-io/go-nats) has been employed for both service discovery and synchronous communication (sic!).

The api gateway exposes http **GET** endpoints, that in case of success will:
- */ping* - respond with http **200** and **"pong"** string in the response body
- */spin* - respond with http **200** and json payload consisting of a dumb messsage constructed while traversing through services `srv1` and `srv2`
- */pods* - respond with http **200** and json payload consisting of a list of currently running kubernetes pods in the `default` namespace

For example, `curl` could be used like that:
```
$ curl -v $(minikube service --url topspin-api)/ping
*   Trying 192.168.99.100...
* TCP_NODELAY set
* Connected to 192.168.99.100 (192.168.99.100) port 30769 (#0)
> GET /ping HTTP/1.1
> Host: 192.168.99.100:30769
> User-Agent: curl/7.58.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< Date: Tue, 21 Aug 2018 19:49:20 GMT
< Content-Length: 4
<
* Connection #0 to host 192.168.99.100 left intact
pong
```

It is possible to manually test most of the application features just by running it in `docker-compose`:
```bash
$ cd micro/topspin/docker
$ docker-compose rm -f && docker-compose up --build
```

(Also) `docker-compose` is the preferred method here to build and push docker images.
Currently the [`sk4zuzu/topspin`](https://hub.docker.com/r/sk4zuzu/topspin/) repo on [`dockerhub`](https://hub.docker.com) is used.
```bash
$ docker login
$ cd micro/topspin/docker
$ docker-compose build && docker-compose push
```

### 3. The helm package

In the `helm/topspin` directory a standard helm chart can be found. Currently, there is a single dependency on the [`stable/nats`](https://github.com/helm/charts/tree/master/stable/nats) chart.

To **build** depedencies, run:
```bash
$ cd helm/topspin/
$ helm dependency build .
```

To **deploy** application, run:
```bash
$ cd helm/topspin/
$ helm upgrade --install topspin . --wait
```

Example output from `kubectl get pods` after successful deployment:
```bash
$ kubectl get pods
NAME                            READY     STATUS    RESTARTS   AGE
topspin-api-589c76497d-65tsj    1/1       Running   1          1h
topspin-api-589c76497d-ngnwb    1/1       Running   1          1h
topspin-api-589c76497d-wfx2z    1/1       Running   1          1h
topspin-nats-0                  1/1       Running   0          1h
topspin-nats-1                  1/1       Running   0          1h
topspin-nats-2                  1/1       Running   0          1h
topspin-srv1-6f6bb45746-4nxc7   1/1       Running   0          1h
topspin-srv1-6f6bb45746-fqmgb   1/1       Running   1          1h
topspin-srv1-6f6bb45746-g74zk   1/1       Running   1          1h
topspin-srv2-8546db6797-g8zps   1/1       Running   1          1h
topspin-srv2-8546db6797-m8b75   1/1       Running   1          1h
topspin-srv2-8546db6797-vjd8p   1/1       Running   0          1h
```

### 4. Ansible automation

In the `ansible/topspin` directory a set of playbooks can be found:
- *runme.yml* - builds `minikube` vm, builds and pushes docker images, deploys helm chart
- *probes/ping.yml* - makes **GET** http call on the application's */ping* endpoint
- *probes/spin.yml* - makes **GET** http call on the application's */spin* endpoint
- *probes/pods.yml* - makes **GET** http call on the application's */pods* endpoint

To build and deploy everything, run:
```bash
$ vbox/ops
$ docker login
$ cd /ansible/topspin/
$ ansible-playbook runme.yml probes/{ping,pods}.yml
```

To deploy **without** building (**recommended for start, otherwise just change the docker repo**), run:
```bash
$ vbox/ops
$ cd /ansible/topspin/
$ ansible-playbook -e docker_push=false runme.yml probes/{ping,pods}.yml
```

### 5. Probes in action

```
$ cd ansible/topspin/ && ansible-playbook probes/{ping,spin,pods}.yml
...
ok: [localhost] => {
    "msg": "pong"
}
...
ok: [localhost] => {
    "msg": {
        "message": "ping?,pong?,pong!"
    }
}
...
ok: [localhost] => {
    "msg": {
        "pods": [
            "topspin-api-589c76497d-kcrv6",
            "topspin-api-589c76497d-m2slc",
            "topspin-api-589c76497d-n5rkp",
            "topspin-nats-0",
            "topspin-nats-1",
            "topspin-nats-2",
            "topspin-srv1-6f6bb45746-5krkr",
            "topspin-srv1-6f6bb45746-5sz4n",
            "topspin-srv1-6f6bb45746-pmqw7",
            "topspin-srv2-8546db6797-6nhfv",
            "topspin-srv2-8546db6797-d9sxt",
            "topspin-srv2-8546db6797-zwftk"
        ]
    }
}
```

### 6. Getting pod listing from inside any running pod (using `curl`)

```bash
$ kubectl exec -it topspin-api-589c76497d-kcrv6 /bin/sh
$ apk --no-cache add jq
$ cd /var/run/secrets/kubernetes.io/serviceaccount/ && curl -s \
>  --cacert ca.crt \
>  -H "Authorization: Bearer $(cat token)" \
>  https://${KUBERNETES_SERVICE_HOST}:${KUBERNETES_SERVICE_PORT}/api/v1/namespaces/$(cat namespace)/pods \
>  | jq -r '.items[] | .metadata.name'
topspin-api-589c76497d-kcrv6
topspin-api-589c76497d-m2slc
topspin-api-589c76497d-n5rkp
topspin-nats-0
topspin-nats-1
topspin-nats-2
topspin-srv1-6f6bb45746-5krkr
topspin-srv1-6f6bb45746-5sz4n
topspin-srv1-6f6bb45746-pmqw7
topspin-srv2-8546db6797-6nhfv
topspin-srv2-8546db6797-d9sxt
topspin-srv2-8546db6797-zwftk
```

[//]: # ( vim:set ts=2 sw=2 et syn=markdown: )
