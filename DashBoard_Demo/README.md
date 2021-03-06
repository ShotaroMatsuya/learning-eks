# k8s_dashboard_demo

0. create cluster with 2 t3.micro nodes

```bash
$ eksctl create cluster --name eksctl-test --nodegroup-name ng-default --node-type t3.micro --nodes 2
```

OR

```bash
$ eksctl create cluster --config-file=eksctl-create-cluster.yaml
```

1.  Deploying the Dashboard UI

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
```

2. Accessing the Dashboard UI (to protect your Cluster data , Dashboard deploys with a minimal RBAC config by default)

create a sample user  
apply the dashboard-adminuser.yaml & dashboard-role-binding.yaml

```bash

$ kubectl apply -f dashboard-adminuser.yaml
$ kubectl apply -f dashboard-role-binding.yaml
```

3.  check the resource in the namespace of dashboard

```bash
$ kubectl get all -n kubernetes-dashboard
```

4. command line proxy
   before do that , let's get the token because you need the token to log in to the dashboard

```bash
$ kubectl -n kubernetes-dashboard describe secret
```

5. You can access Dashboard using the kubectl command-line tool by running the following command:

```bash
$ kubectl proxy
```

access http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/  
paste the token

6. Clean up and next steps
   Remove the admin ServiceAccount and ClusterRoleBinding.

```bash
kubectl -n kubernetes-dashboard delete serviceaccount admin-user
kubectl -n kubernetes-dashboard delete clusterrolebinding admin-user
```
