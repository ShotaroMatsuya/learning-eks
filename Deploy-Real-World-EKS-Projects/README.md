# Deploy Real World EKS projects

provided by https://docs.aws.amazon.com/eks/latest/userguide/sample-deployment.html

Running EKS Cluster

```bash
$ eksctl create cluster --name=guestbook-demo
```

To create your guest book app

1. Create the Redis master replication controller.

```bash
$ kubectl apply -f redis-master-controller.json
```

2. reate the Redis master service

```bash
   $ kubectl apply -f redis-master-service.json
```

3. Create the Redis slave replication controller.

```bash
   $ kubectl apply -f redis-slave-controller.json
```

4. Create the Redis slave service

```bash
$ kubectl apply -f redis-slave-service.json
```

5. create the guestbook replication controller.

```bash
$ kubectl apply -f guestbook-controller.json
```

6. Create the guestbook service.

```bash
   $ kubectl apply -f guestbook-service.json
```

7. Query the services in your cluster and wait until the External IP column for the guestbook service is populated

   It might take several minutes before the IP address is available.

```bash
$ kubectl get services -o wide
```
