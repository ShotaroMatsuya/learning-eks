# Helm

1. IAM Policy の作成
   AWS Load Balancer Controller が API を介して AWS リソースを作成するための policy を作成

```bash
aws iam create-policy \
    --policy-name ALBCIAMPolicy \
    --policy-document file://albc_iam_policy.json
```

2. IRSA を利用するための準備
   IRSA(IAM Roles for Service Accounts)を作成し、先程作成した Policy を紐付ける。  
   IRSA とは, IAM ロールを Pod に紐付け、AWS リソースへのアクセス制御を行うもの

```bash
eksctl create iamserviceaccount \
    --cluster=eks-cluster \
    --namespace=kube-system \
    --name=aws-load-balancer-controller \
    --role-name "ALBCRole" \
    --attach-policy-arn=arn:aws:iam::528163014577:policy/ALBCIAMPolicy \
    --approve

# 作成されたService Accountの確認
kubectl get sa aws-load-balancer-controller -o yaml -n kube-system
```

3. Chart のインストール

chart をインストールして、ServiceAccount リソースを Controller に紐付ける。

```bash
helm repo add eks https://aws.github.io/eks-charts
helm repo update
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
    -n kube-system \
    --version 1.4.2 \
    --set clusterName=eks-cluster \
    --set serviceAccount.create=false \
    --set serviceAccount.name=aws-load-balancer-controller \
    --set image.repository=602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/amazon/aws-load-balancer-controller \
    --set image.tag=v2.4.2

# controllerのインストールの確認
kubectl get deployment -n kube-system aws-load-balancer-controller
```

4. Ingress, Service, Deployment の作成

```bash
kubectl apply -f helm-albc.yaml

# 確認
kubectl get ing
```
