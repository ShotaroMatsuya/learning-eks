# Create elasticsearch

```bash
aws es create-elasticsearch-domain \
  --domain-name eks-logs \
  --elasticsearch-version 7.9 \
  --elasticsearch-cluster-config \
  InstanceType=t3.small.elasticsearch,InstanceCount=1 \
  --ebs-options EBSEnabled=true,VolumeType=gp2,VolumeSize=10 \
  --access-policies '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["es:*"],"Resource":"*"}]}'
```
