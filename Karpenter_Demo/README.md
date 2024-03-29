# Karpenter Demo

[Karpenter Official Demo](https://karpenter.sh/docs/getting-started/)

[experimental](./expermental.md)

## Environment Variables

after setting up the tool , set the following environment variables to store commonly userd values.

```
export CLUSTER_NAME=$USER-karpenter-demo
export AWS_DEFAULT_REGION=ap-northeast-1
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
```

## Create a Cluster

Create a cluster with eksctl. This example configuration file specifies a basic cluster with one initial node and sets up an IAM OIDC provider for the cluster to enable IAM roles for pods:

```bash
$ eksctl create cluster -f cluster.yaml
```

This guide uses a self-managed node group to host Karpenter.

Karpenter itself can run anywhere, including on self-managed node groups, managed node groups, or AWS Fargate.

Karpenter will provision EC2 instances in your account.

## Tag Subnets

```bash
SUBNET_IDS=$(aws cloudformation describe-stacks \
    --stack-name eksctl-${CLUSTER_NAME}-cluster \
 --query 'Stacks[].Outputs[?OutputKey==`SubnetsPrivate`].OutputValue' \
 --output text)
aws ec2 create-tags \
 --resources $(echo $SUBNET_IDS | tr ',' '\n') \
    --tags Key="kubernetes.io/cluster/${CLUSTER_NAME}",Value=
```

## Create the KarpenterNode IAM Role

First, create the IAM resources using AWS CloudFormation.

```bash
TEMPOUT=$(mktemp)
curl -fsSL https://karpenter.sh/docs/getting-started/cloudformation.yaml > $TEMPOUT \
&& aws cloudformation deploy \
  --stack-name Karpenter-${CLUSTER_NAME} \
 --template-file ${TEMPOUT} \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides ClusterName=${CLUSTER_NAME}
```

Second, grant access to instances using the profile to connect to the cluster. This command adds the Karpenter node role to your aws-auth configmap, allowing nodes with this role to connect to the cluster.

```bash
eksctl create iamidentitymapping \
 --username system:node:{{EC2PrivateDNSName}} \
 --cluster ${CLUSTER_NAME} \
  --arn arn:aws:iam::${AWS_ACCOUNT_ID}:role/KarpenterNodeRole-${CLUSTER_NAME} \
 --group system:bootstrappers \
 --group system:nodes
```

Now, Karpenter can launch new EC2 instances and those instances can connect to your cluster.

## Create the KarpenterController IAM Role

Karpenter requires permissions like launching instances. This will create an AWS IAM Role, Kubernetes service account, and associate them using IRSA.

```bash
eksctl create iamserviceaccount \
 --cluster $CLUSTER_NAME --name karpenter --namespace karpenter \
  --attach-policy-arn arn:aws:iam::$AWS_ACCOUNT_ID:policy/KarpenterControllerPolicy-$CLUSTER_NAME \
  --approve
```

## (before step4, do this to avoid error) make IAM OIDC provider associated with cluster and try step4 again

```bash
$ eksctl utils associate-iam-oidc-provider --region=ap-northeast-1 --cluster=eks-karpenter-demo --approve
```

## ( Create the EC2 Spot Service Linked Role )

This step is only necessary if this is the first time you’re using EC2 Spot in this account. More details are available here.

```bash
$ aws iam create-service-linked-role --aws-service-name spot.amazonaws.com
```

output role ARN & cop y it

## Install Karpenter Helm Chart

Cfn テンプレートと IRSA を使用して IAM リソースを作成。そして、Karpenter コントローラがドキュメントに従ってインスタンスの起動などのアクセス許可を取得できるようにする必要がある。

```bash
helm repo add karpenter https://charts.karpenter.sh
helm repo update
helm upgrade --install karpenter karpenter/karpenter --namespace karpenter \
 --create-namespace --set serviceAccount.create=false --version 0.5.1 \
 --set controller.clusterName=${CLUSTER_NAME} \
  --set controller.clusterEndpoint=$(aws eks describe-cluster --name ${CLUSTER_NAME} --query "cluster.endpoint" --output json) \
 --wait # for the defaulting webhook to install before creating a Provisioner
```

## Provisioner

Create a default provisioner using the command below. This provisioner  
configures instances to connect to your cluster’s endpoint and discovers resources like subnets and security groups using the cluster’s name.

The ttlSecondsAfterEmpty value configures Karpenter to terminate empty nodes. This behavior can be disabled by leaving the value undefined.  
Note: This provisioner will create capacity as long as the sum of all created capacity is less than the specified limit.

```bash
cat <<EOF | kubectl apply -f -
apiVersion: karpenter.sh/v1alpha5
kind: Provisioner
metadata:
name: default
spec:
requirements: - key: karpenter.sh/capacity-typemn7
operator: In
values: ["spot"]
limits:
resources:
cpu: 1000
provider:
instanceProfile: KarpenterNodeInstanceProfile-eks-karpenter-demo
ttlSecondsAfterEmpty: 30
EOF

$ kubectl apply -f karpenter.yaml
```
