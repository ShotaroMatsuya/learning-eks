# Ingress Demo 2

## First , try Demo1 cause were going to build all this on top of Demo one.

How to install it will be the controller and the basics of ingress resource

## apply frontend

```bash
$ kubectl apply -f frontend-service.yaml
$ kubectl apply -f nginx-deployment.yaml
```

## check the IP address of the pod

```bash
$ kubectl get pods -o wide -n 2048-game
```

## apply new ingress resource

```bash
$ kubectl apply -f webserver-ingress.yaml
```

## modified to ip-mode ALB settings by adding the annotation bellow

```yml
annotations:
alb.ingress.kubernetes.io/target-type: ip
```

This should register the IP address of the pod directly in load balancer.

```bash
$ kubectl apply -f webserver-ingress.yaml
```

## fixing route /frontend/ correctly

```bash
$kubectl exec -it frontend-deployment-5ccdb88c57-cdzcl -n 2048-game -- /bin/bash


# apt-get update

# apt-get install vim

# find -name index.html

# cd /usr/share/nginx/html/

# mkdir frontend

# cp index.html ./frontend/
```
