# Log cluster by eksctl

[How to create the ES(opensearch) cluster](create_ES.md)

## create cluster with 2 m5.large nodes

```bash
$ eksctl create cluster --name eks-loggingtest --nodegroup-name ng-default --node-type m5.large --nodes 2
```

## attached role of 'CloudWatchLogsFullAccess' to EC2 ON kubernetes

## replace env field with your environment in fluentd.yaml

## apply the nginx service & hello-k8s-forlog via yaml

## deploy load balancer service

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
        "AWS": "_"
      },
      "Action": ["es:ESHttp_"],
      "Condition": {
        "IpAddress": {
          "aws:SourceIp": ["192.0.2.0/24"]
        }
      },
      "Resource": "arn:aws:es:ap-northeast-1:528163014577:domain/eks-logs/*"
    }
  ]
}
```

## clean up resources

```bash
$ eksctl delete cluster --name=eks-loggingtest
```
