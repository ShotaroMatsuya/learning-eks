# Create Prometheus

- Install Grafana  
  [here](./README.md)
- CloudWatch Container Insights with FluentD  
  [here](./cloudwatch.md)

## create cluster with 2 m5.large nodes

```bash
$ eksctl create cluster --name eks-Prometheus --nodegroup-name ng-default --node-type m5.large --nodes 2
```

## Installing the k8s Metrics Server

before we can install Prometheus, We need to install metric server .  
The k8s Metrics Server is an aggregator of resource usage data in your cluster, and it is not deployed by default In Amazon EKS cluster.

```bash
$ kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

verify that the metrics-server deployment is running the desired number of pods with the following command.

```bash
$ kubectl get deployment metrics-server -n kube-system
```

## Deploying Prometheus

1. create a prometheus namespace

```bash
$ kubectl create namespace prometheus
```

2. adding prometheus-community repository

```bash
$ helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

3. deploy prometheus

```bash
$ helm upgrade -i prometheus prometheus-community/prometheus \
 --namespace prometheus \
 --set alertmanager.persistentVolume.storageClass="gp2",server.persistentVolume.storageClass="gp2"
```

verify thatt all of the pods in the prometheus namespace are in the READY state.

```bash
$ kubectl get pods -n prometheus
```

4. kubectl を使用して、Prometheus コンソールをローカルマシンにポート転送します。

```bash
$ kubectl --namespace=prometheus port-forward deploy/prometheus-server 9090
```

ウェブブラウザで localhost:9090 にアクセスして、Prometheus コンソールを表示します。
