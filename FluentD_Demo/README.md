# Log cluster by eksctl

[How to create the ES(opensearch) cluster](create_ES.md)

## create cluster with 2 m5.large nodes

```bash
$ eksctl create cluster --name eks-loggingtest --nodegroup-name ng-default --node-type m5.large --nodes 2
```

## attached role of 'CloudWatchLogsFullAccess' to EC2 ON kubernetes(nodegroup-NodeInstanceRole)

## replace env field with your environment in fluentd.yaml

```bash
kubectl apply -f fluentd.yml
```

## apply the nginx service & hello-k8s-forlog via yaml

## deploy load balancer service

```bash
kubectl apply -f loadbalancer-service.yaml
```

```bash
kubectl apply -f hello-k8s-forlog.yml
```

## get the access to ELB

```bash
$ kubectl get all
```

## elasticsearch deployment

```bash
$ aws es create-elasticsearch-domain \
 --domain-name eks-logs \
 --elasticsearch-version 7.9 \
 --elasticsearch-cluster-config \
 InstanceType=t3.small.elasticsearch,InstanceCount=1 \
 --ebs-options EBSEnabled=true,VolumeType=gp2,VolumeSize=10


# (not required↓)
\
 --access-policies '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["es:*"],"Resource":"\*"}]}'
```

## create role to grant es permission to lambda

## Create Amazon OpenSearch Service subscription filter IN cloudwatch logs group

## attach access policy in order to add your IP to white list

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "es:*",
      "Resource": "arn:aws:es:ap-northeast-1:528163014577:domain/eks-logs/*",
      "Condition": {
        "IpAddress": {
          "aws:SourceIp": "<YOUR IP>"
        }
      }
    }
  ]
}
```

## clean up resources

```bash
$ eksctl delete cluster --name=eks-loggingtest
```
