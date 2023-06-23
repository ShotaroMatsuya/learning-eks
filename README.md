# learning-eks

1. create cluster with 2 t3.micro nodes

```bash
$ eksctl create cluster --name eksctl-test --nodegroup-name ng-default --node-type t3.micro --nodes 2
```

OR

```bash
$ eksctl create cluster --config-file=eksctl-create-cluster.yaml
```

※create eks cluster with managed node group

```bash
$ eksctl create cluster --name <name> --version 1.15 --nodegroup-name <nodegrpname> --node-type t3.micro --nodes 2 --managed
```

※ ECS Cluster with Fargate Profile

```bash
$ eksctl create cluster --name <name> --fargate
```

2. create node-group

```bash
$ eksctl create nodegroup --config-file=eksctl-create-ng.yaml
```

3. get the information

```bash
$ eksctl get nodegroup --cluster=eksctl-test
$ eksctl get cluster
```

4. delete cluster

```bash
$ eksctl delete cluster --name=eksctl-test
```

5. update managed nodegroup

```bash
$ eksctl upgrade nodegroup --name=managed-ng-1 --cluster=managed-cluster
```

---

# EKS Cluster の構築

```bash
eksctl create cluster \
      --name eks-cluster \
      --version 1.22 \
      --with-oidc \
      --nodegroup-name eks-cluster-node-group \
      --node-type c5.large \
      --nodes 1 \
      --nodes-min 1
```

# クラスターの削除

EKS Cluster を削除する前に、Service Type が LoadBalancer の Service が存在するばあいは、事前に削除

```bash
kubectl delete svc service
```

cluster の削除

```bash
eksctl delete cluster --name eks-cluster
```
