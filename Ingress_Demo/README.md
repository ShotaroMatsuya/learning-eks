# EKS Ingress DEMO1 (deprecated)

1. Deploy Amazon EKS with eksctl
   $ eksctl create cluster --name=attractive-gopher

Create an IAM OIDC provider and associate it with your cluster

```bash
$ eksctl utils associate-iam-oidc-provider --cluster=attractive-gopher --approve
```

2. Deploy AWS ALB INgress controller
   First off, deploy the relevant RBAC roles and role bindings as required by the AWS ALB Ingress controller:

```bash
$ kubectl apply -f rbac-role.yaml
```

create an IAM policy named ALBIngressControllerIAMPolicy to allow the ALB Ingress controller to make AWS API calls on your behalf. Record the Policy.Arn in the command output, you will need it in the next step:

```bash
$ aws iam create-policy \
 --policy-name ALBIngressControllerIAMPolicy \
 --policy-document file://iam-policy.json
```

copy this ARN and keep it
`"Arn": "arn:aws:iam::528163014577:policy/ALBIngressControllerIAMPolicy",`

create a Kubernetes service account and an IAM role (for the pod running the AWS ALB Ingress controller) by substituting $PolicyARN with the recorded value from the previous step:

```bash
$ eksctl create iamserviceaccount \
 --cluster=attractive-gopher \
 --namespace=kube-system \
 --name=alb-ingress-controller \
 --attach-policy-arn=arn:aws:iam::528163014577:policy/ALBIngressControllerIAMPolicy \
 --override-existing-serviceaccounts \
 --approve
```

Then deploy the AWS ALB Ingress controller:

```bash
$ kubectl apply -f alb-ingress-controller.yaml
```

Finally, verify that the deployment was successful and the controller started:

```bash
$ kubectl logs -n kube-system $(kubectl get po -n kube-system | egrep -o alb-ingress[a-zA-Z0-9-]+)
```

confirm the alb-ingress-controller running in namespace of kube-system

```bash
$ kubectl get all -n kube-system
```

we don't have any ALB yet because the LB only get created when the Ingress Resource is deployed .

Deploy sample application  
Now let’s deploy a sample 2048 game into our Kubernetes cluster and use the Ingress resource to expose it to traffic.

```bash
$ kubectl apply -f 2048-namespace.yaml
$ kubectl apply -f 2048-deployment.yaml
$ kubectl apply -f 2048-service.yaml
```

Deploy an Ingress resource for the 2048 game:

※1 added below tags to all subnet  
AWS Load Balancer Controller バージョン v2.1.1 以前を使用している場合、サブネットに次のようにタグ付けする必要があります。バージョン 2.1.2 以降を使用している場合、このタグはオプションです。  
ただし、次のいずれかの場合は、サブネットにタグ付けることをお勧めします。同じ VPC で実行されている複数のクラスターがあるか、VPC でサブネットを共有する複数の AWS サービスがあります。  
または、各クラスターに対してロードバランサーをプロビジョニングする場所を詳細に制御する必要があります。cluster-name をクラスター名に置き換えます。

```
Key: kubernetes.io/cluster/attractive-gopher

Value: shared
```

---

※2 Updated ingress resource manifest file with subnet annotation

```bash
$ kubectl apply -f 2048-ingress.yaml
```

confirm ingress resource created

```bash
$ kubectl get ingress/2048-ingress -n 2048-game
$ kubectl describe ingress/2048-ingress -n 2048-game
$ kubectl logs -n kube-system deployment.apps/alb-ingress-controller
```

[Documentation](https://aws.amazon.com/jp/blogs/opensource/kubernetes-ingress-aws-alb-ingress-controller/)
