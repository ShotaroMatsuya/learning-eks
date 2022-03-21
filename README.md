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
