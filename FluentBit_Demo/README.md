# FluentBit demo cluster by eksctl

## set up for amazon eks

```bash
$ eksctl create cluster --name fluent-bit-demo --nodegroup-name ng-default --node-type m5.large --nodes 2
```

## attached eks-fluent-bit-demonset-policy to eks on ec2(firehose permission)

## first, create firehose-policy.json & firehose-delivery-policy.json with firehose-delivery-policy

```bash
$ aws iam create-role \
 --role-name firehose_delivery_role \
 --assume-role-policy-document file://firehose-policy.json
```

## save the firehose delivery roll arn

`arn:aws:iam::528163014577:role/firehose_delivery_role`

## run the command to put policy to the role

```bash
$ aws iam put-role-policy \
 --role-name firehose_delivery_role \
 --policy-name firehose-fluentbit-s3-streaming \
 --policy-document file://firehose-delivery-policy.json
```

## create the ECS delivery stream

```bash
$ aws firehose create-delivery-stream \
 --delivery-stream-name eks-stream \
 --delivery-stream-type DirectPut \
 --s3-destination-configuration \
```

```
RoleARN=arn:aws:iam::528163014577:role/firehose_delivery_role,\
BucketARN="arn:aws:s3:::eks-fluentbit-demo-smat",\
Prefix=eks
```

## save arn of kinesis

```
"DeliveryStreamARN": "arn:aws:firehose:ap-northeast-1:528163014577:deliverystream/eks-stream"
```

## `kubectl apply -f eks-nginx-app.yaml`

## define the role and binding in a file named eks-fluent-bit-daemonset-rbac.yaml. before doing this , run this command to create sa

```bash
$ kubectl create sa fluent-bit
```

## define the log parsing and routing for the FluentBit plugin. For this , use a file called eks-fluent-bit-configmap.yaml

## define the kubernetes Daemonset(using said config map ) in a file called eks-fluent-bit-daemonset.yaml

## finally , launch the Fluent Bit daemonset by executing

```bash
$ kubectl apply -f eks-fluent-bit-daemonset.yaml and verify the Fluent Bit daemonset by peeking into the logs like so:
```

```bash
$ kubectl logs ds/fluentbit
```

## create service & ouput elb dns

a7f4ecf468379437586e474ce716cd7e-1232473850.ap-northeast-1.elb.amazonaws.com

## the next step is to generate some load for the NGINX container running in EKS. You can grab the load generator files for EKS and execute the commands bellow; this will CURL the respective NGINX service every two seconds(executing in the background ), until you kill the scripts:

(chmod 755 ./load-gen-eks.sh)

```bash
$ ./load-gen-eks.sh &
```
