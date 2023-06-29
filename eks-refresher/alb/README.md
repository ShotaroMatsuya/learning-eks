---
title: AWS Load Balancer Controller Install on AWS EKS
description: Learn to install AWS Load Balancer Controller for Ingress Implementation on AWS EKS
---
 

## Step-00: Introduction
1. Create IAM Policy and make a note of Policy ARN
2. Create IAM Role and k8s Service Account and bound them together
3. Install AWS Load Balancer Controller using HELM3 CLI
4. Understand IngressClass Concept and create a default Ingress Class 

## Step-01: Pre-requisites
### Pre-requisite-1: eksctl & kubectl Command Line Utility
- Should be the latest eksctl version
```t
# Verify eksctl version
eksctl version

# For installing or upgrading latest eksctl version
https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html

# Verify EKS Cluster version
kubectl version --short
kubectl version
Important Note: You must use a kubectl version that is within one minor version difference of your Amazon EKS cluster control plane. For example, a 1.20 kubectl client works with Kubernetes 1.19, 1.20 and 1.21 clusters.

# For installing kubectl cli
https://docs.aws.amazon.com/eks/latest/userguide/install-kubectl.html
```
### Pre-requisite-2: Create EKS Cluster and Worker Nodes (if not created)
```t
# Create Cluster (Section-01-02)
eksctl create cluster --name=eksdemo1 \
                      --region=us-east-1 \
                      --zones=us-east-1a,us-east-1b \
                      --version="1.21" \
                      --without-nodegroup 


# Get List of clusters (Section-01-02)
eksctl get cluster   

# Template (Section-01-02)
eksctl utils associate-iam-oidc-provider \
    --region region-code \
    --cluster <cluter-name> \
    --approve

# Replace with region & cluster name (Section-01-02)
eksctl utils associate-iam-oidc-provider \
    --region us-east-1 \
    --cluster eksdemo1 \
    --approve

# Create EKS NodeGroup in VPC Private Subnets (Section-07-01)
eksctl create nodegroup --cluster=eksdemo1 \
                        --region=us-east-1 \
                        --name=eksdemo1-ng-private1 \
                        --node-type=t3.medium \
                        --nodes-min=2 \
                        --nodes-max=4 \
                        --node-volume-size=20 \
                        --ssh-access \
                        --ssh-public-key=kube-demo \
                        --managed \
                        --asg-access \
                        --external-dns-access \
                        --full-ecr-access \
                        --appmesh-access \
                        --alb-ingress-access \
                        --node-private-networking       
```
### Pre-requisite-3:  Verify Cluster, Node Groups and configure kubectl cli if not configured
1. EKS Cluster
2. EKS Node Groups in Private Subnets
```t
# Verfy EKS Cluster
eksctl get cluster

# Verify EKS Node Groups
eksctl get nodegroup --cluster=eksdemo1

# Verify if any IAM Service Accounts present in EKS Cluster
eksctl get iamserviceaccount --cluster=eksdemo1
Observation:
1. No k8s Service accounts as of now. 

# Configure kubeconfig for kubectl
eksctl get cluster # TO GET CLUSTER NAME
aws eks --region <region-code> update-kubeconfig --name <cluster_name>
aws eks --region us-east-1 update-kubeconfig --name eksdemo1

# Verify EKS Nodes in EKS Cluster using kubectl
kubectl get nodes

# Verify using AWS Management Console
1. EKS EC2 Nodes (Verify Subnet in Networking Tab)
2. EKS Cluster
```

## Step-02: Create IAM Policy
- Create IAM policy for the AWS Load Balancer Controller that allows it to make calls to AWS APIs on your behalf.
- As on today `2.3.1` is the latest Load Balancer Controller
- We will download always latest from main branch of Git Repo
- [AWS Load Balancer Controller Main Git repo](https://github.com/kubernetes-sigs/aws-load-balancer-controller)
```t
# Change Directroy
cd 08-NEW-ELB-Application-LoadBalancers/
cd 08-01-Load-Balancer-Controller-Install

# Delete files before download (if any present)
rm iam_policy_latest.json

# Download IAM Policy
## Download latest
curl -o iam_policy_latest.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json
## Verify latest
ls -lrta 

## Download specific version
curl -o iam_policy_v2.3.1.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.3.1/docs/install/iam_policy.json


# Create IAM Policy using policy downloaded 
aws iam create-policy \
    --policy-name AWSLoadBalancerControllerIAMPolicy \
    --policy-document file://iam_policy_latest.json

## Sample Output
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ aws iam create-policy \
>     --policy-name AWSLoadBalancerControllerIAMPolicy \
>     --policy-document file://iam_policy_latest.json
{
    "Policy": {
        "PolicyName": "AWSLoadBalancerControllerIAMPolicy",
        "PolicyId": "ANPASUF7HC7S52ZQAPETR",
        "Arn": "arn:aws:iam::180789647333:policy/AWSLoadBalancerControllerIAMPolicy",
        "Path": "/",
        "DefaultVersionId": "v1",
        "AttachmentCount": 0,
        "PermissionsBoundaryUsageCount": 0,
        "IsAttachable": true,
        "CreateDate": "2022-02-02T04:51:21+00:00",
        "UpdateDate": "2022-02-02T04:51:21+00:00"
    }
}
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ 
```
- **Important Note:** If you view the policy in the AWS Management Console, you may see warnings for ELB. These can be safely ignored because some of the actions only exist for ELB v2. You do not see warnings for ELB v2.

### Make a note of Policy ARN    
- Make a note of Policy ARN as we are going to use that in next step when creating IAM Role.
```t
# Policy ARN 
Policy ARN:  arn:aws:iam::180789647333:policy/AWSLoadBalancerControllerIAMPolicy
```


## Step-03: Create an IAM role for the AWS LoadBalancer Controller and attach the role to the Kubernetes service account 
- Applicable only with `eksctl` managed clusters
- This command will create an AWS IAM role 
- This command also will create Kubernetes Service Account in k8s cluster
- In addition, this command will bound IAM Role created and the Kubernetes service account created
### Step-03-01: Create IAM Role using eksctl
```t
# Verify if any existing service account
kubectl get sa -n kube-system
kubectl get sa aws-load-balancer-controller -n kube-system
Obseravation:
1. Nothing with name "aws-load-balancer-controller" should exist

# Template
eksctl create iamserviceaccount \
  --cluster=my_cluster \
  --namespace=kube-system \
  --name=aws-load-balancer-controller \ #Note:  K8S Service Account Name that need to be bound to newly created IAM Role
  --attach-policy-arn=arn:aws:iam::111122223333:policy/AWSLoadBalancerControllerIAMPolicy \
  --override-existing-serviceaccounts \
  --approve


# Replaced name, cluster and policy arn (Policy arn we took note in step-02)
eksctl create iamserviceaccount \
  --cluster=eksdemo1 \
  --namespace=kube-system \
  --name=aws-load-balancer-controller \
  --attach-policy-arn=arn:aws:iam::528163014577:policy/AWSLoadBalancerControllerIAMPolicy \
  --override-existing-serviceaccounts \
  --approve
```
- **Sample Output**
```t
# Sample Output for IAM Service Account creation
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ eksctl create iamserviceaccount \
>   --cluster=eksdemo1 \
>   --namespace=kube-system \
>   --name=aws-load-balancer-controller \
>   --attach-policy-arn=arn:aws:iam::180789647333:policy/AWSLoadBalancerControllerIAMPolicy \
>   --override-existing-serviceaccounts \
>   --approve
2022-02-02 10:22:49 [ℹ]  eksctl version 0.82.0
2022-02-02 10:22:49 [ℹ]  using region us-east-1
2022-02-02 10:22:52 [ℹ]  1 iamserviceaccount (kube-system/aws-load-balancer-controller) was included (based on the include/exclude rules)
2022-02-02 10:22:52 [!]  metadata of serviceaccounts that exist in Kubernetes will be updated, as --override-existing-serviceaccounts was set
2022-02-02 10:22:52 [ℹ]  1 task: { 
    2 sequential sub-tasks: { 
        create IAM role for serviceaccount "kube-system/aws-load-balancer-controller",
        create serviceaccount "kube-system/aws-load-balancer-controller",
    } }2022-02-02 10:22:52 [ℹ]  building iamserviceaccount stack "eksctl-eksdemo1-addon-iamserviceaccount-kube-system-aws-load-balancer-controller"
2022-02-02 10:22:53 [ℹ]  deploying stack "eksctl-eksdemo1-addon-iamserviceaccount-kube-system-aws-load-balancer-controller"
2022-02-02 10:22:53 [ℹ]  waiting for CloudFormation stack "eksctl-eksdemo1-addon-iamserviceaccount-kube-system-aws-load-balancer-controller"
2022-02-02 10:23:10 [ℹ]  waiting for CloudFormation stack "eksctl-eksdemo1-addon-iamserviceaccount-kube-system-aws-load-balancer-controller"
2022-02-02 10:23:29 [ℹ]  waiting for CloudFormation stack "eksctl-eksdemo1-addon-iamserviceaccount-kube-system-aws-load-balancer-controller"
2022-02-02 10:23:32 [ℹ]  created serviceaccount "kube-system/aws-load-balancer-controller"
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ 
```

### Step-03-02: Verify using eksctl cli
```t
# Get IAM Service Account
eksctl  get iamserviceaccount --cluster eksdemo1

# Sample Output
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ eksctl  get iamserviceaccount --cluster eksdemo1
2022-02-02 10:23:50 [ℹ]  eksctl version 0.82.0
2022-02-02 10:23:50 [ℹ]  using region us-east-1
NAMESPACE	NAME				ROLE ARN
kube-system	aws-load-balancer-controller	arn:aws:iam::180789647333:role/eksctl-eksdemo1-addon-iamserviceaccount-kube-Role1-1244GWMVEAKEN
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ 
```

### Step-03-03: Verify CloudFormation Template eksctl created & IAM Role
- Goto Services -> CloudFormation
- **CFN Template Name:** eksctl-eksdemo1-addon-iamserviceaccount-kube-system-aws-load-balancer-controller
- Click on **Resources** tab
- Click on link in **Physical Id** to open the IAM Role
- Verify it has **eksctl-eksdemo1-addon-iamserviceaccount-kube-Role1-WFAWGQKTAVLR** associated

### Step-03-04: Verify k8s Service Account using kubectl
```t
# Verify if any existing service account
kubectl get sa -n kube-system
kubectl get sa aws-load-balancer-controller -n kube-system
Obseravation:
1. We should see a new Service account created. 

# Describe Service Account aws-load-balancer-controller
kubectl describe sa aws-load-balancer-controller -n kube-system
```
- **Observation:** You can see that newly created Role ARN is added in `Annotations` confirming that **AWS IAM role bound to a Kubernetes service account**
- **Output**
```t
## Sample Output
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ kubectl describe sa aws-load-balancer-controller -n kube-system
Name:                aws-load-balancer-controller
Namespace:           kube-system
Labels:              app.kubernetes.io/managed-by=eksctl
Annotations:         eks.amazonaws.com/role-arn: arn:aws:iam::180789647333:role/eksctl-eksdemo1-addon-iamserviceaccount-kube-Role1-1244GWMVEAKEN
Image pull secrets:  <none>
Mountable secrets:   aws-load-balancer-controller-token-5w8th
Tokens:              aws-load-balancer-controller-token-5w8th
Events:              <none>
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ 
```

## Step-04: Install the AWS Load Balancer Controller using Helm V3 
### Step-04-01: Install Helm
- [Install Helm](https://helm.sh/docs/intro/install/) if not installed
- [Install Helm for AWS EKS](https://docs.aws.amazon.com/eks/latest/userguide/helm.html)
```t
# Install Helm (if not installed) MacOS
brew install helm

# Verify Helm version
helm version
```
### Step-04-02: Install AWS Load Balancer Controller
- **Important-Note-1:** If you're deploying the controller to Amazon EC2 nodes that have restricted access to the Amazon EC2 instance metadata service (IMDS), or if you're deploying to Fargate, then add the following flags to the command that you run:
```t
--set region=region-code
--set vpcId=vpc-xxxxxxxx
```
- **Important-Note-2:** If you're deploying to any Region other than us-west-2, then add the following flag to the command that you run, replacing account and region-code with the values for your region listed in Amazon EKS add-on container image addresses.
- [Get Region Code and Account info](https://docs.aws.amazon.com/eks/latest/userguide/add-ons-images.html)
```t
--set image.repository=account.dkr.ecr.region-code.amazonaws.com/amazon/aws-load-balancer-controller
```
```t
# Add the eks-charts repository.
helm repo add eks https://aws.github.io/eks-charts

# Update your local repo to make sure that you have the most recent charts.
helm repo update

# Install the AWS Load Balancer Controller.
## Template
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=<cluster-name> \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller \
  --set region=<region-code> \
  --set vpcId=<vpc-xxxxxxxx> \
  --set image.repository=<account>.dkr.ecr.<region-code>.amazonaws.com/amazon/aws-load-balancer-controller

## Replace Cluster Name, Region Code, VPC ID, Image Repo Account ID and Region Code  
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system \
  --set clusterName=eksdemo1 \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controller \
  --set region=ap-northeast-1 \
  --set vpcId=vpc-01578d4675034c874 \
  --set image.repository=public.ecr.aws/eks/aws-load-balancer-controller
```
- **Sample output for AWS Load Balancer Controller Install steps**
```t
## Sample Ouput for AWS Load Balancer Controller Install steps
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
>   -n kube-system \
>   --set clusterName=eksdemo1 \
>   --set serviceAccount.create=false \
>   --set serviceAccount.name=aws-load-balancer-controller \
>   --set region=us-east-1 \
>   --set vpcId=vpc-0570fda59c5aaf192 \
>   --set image.repository=602401143452.dkr.ecr.us-east-1.amazonaws.com/amazon/aws-load-balancer-controller
NAME: aws-load-balancer-controller
LAST DEPLOYED: Wed Feb  2 10:33:57 2022
NAMESPACE: kube-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
AWS Load Balancer controller installed!
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ 
```
### Step-04-03: Verify that the controller is installed and Webhook Service created
```t
# Verify that the controller is installed.
kubectl -n kube-system get deployment 
kubectl -n kube-system get deployment aws-load-balancer-controller
kubectl -n kube-system describe deployment aws-load-balancer-controller

# Sample Output
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ kubectl get deployment -n kube-system aws-load-balancer-controller
NAME                           READY   UP-TO-DATE   AVAILABLE   AGE
aws-load-balancer-controller   2/2     2            2           27s
Kalyans-MacBook-Pro:08-01-Load-Balancer-Controller-Install kdaida$ 

# Verify AWS Load Balancer Controller Webhook service created
kubectl -n kube-system get svc 
kubectl -n kube-system get svc aws-load-balancer-webhook-service
kubectl -n kube-system describe svc aws-load-balancer-webhook-service

# Sample Output
Kalyans-MacBook-Pro:aws-eks-kubernetes-masterclass-internal kdaida$ kubectl -n kube-system get svc aws-load-balancer-webhook-service
NAME                                TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
aws-load-balancer-webhook-service   ClusterIP   10.100.53.52   <none>        443/TCP   61m
Kalyans-MacBook-Pro:aws-eks-kubernetes-masterclass-internal kdaida$ 

# Verify Labels in Service and Selector Labels in Deployment
kubectl -n kube-system get svc aws-load-balancer-webhook-service -o yaml
kubectl -n kube-system get deployment aws-load-balancer-controller -o yaml
Observation:
1. Verify "spec.selector" label in "aws-load-balancer-webhook-service"
2. Compare it with "aws-load-balancer-controller" Deployment "spec.selector.matchLabels"
3. Both values should be same which traffic coming to "aws-load-balancer-webhook-service" on port 443 will be sent to port 9443 on "aws-load-balancer-controller" deployment related pods. 
```

### Step-04-04: Verify AWS Load Balancer Controller Logs
```t
# List Pods
kubectl get pods -n kube-system

# Review logs for AWS LB Controller POD-1
kubectl -n kube-system logs -f <POD-NAME> 
kubectl -n kube-system logs -f  aws-load-balancer-controller-86b598cbd6-5pjfk

# Review logs for AWS LB Controller POD-2
kubectl -n kube-system logs -f <POD-NAME> 
kubectl -n kube-system logs -f aws-load-balancer-controller-86b598cbd6-vqqsk
```

### Step-04-05: Verify AWS Load Balancer Controller k8s Service Account - Internals 
```t
# List Service Account and its secret
kubectl -n kube-system get sa aws-load-balancer-controller
kubectl -n kube-system get sa aws-load-balancer-controller -o yaml
kubectl -n kube-system get secret <GET_FROM_PREVIOUS_COMMAND - secrets.name> -o yaml
kubectl -n kube-system get secret aws-load-balancer-controller-token-5w8th 
kubectl -n kube-system get secret aws-load-balancer-controller-token-5w8th -o yaml
## Decoce ca.crt using below two websites
https://www.base64decode.org/
https://www.sslchecker.com/certdecoder

## Decode token using below two websites
https://www.base64decode.org/
https://jwt.io/
Observation:
1. Review decoded JWT Token

# List Deployment in YAML format
kubectl -n kube-system get deploy aws-load-balancer-controller -o yaml
Observation:
1. Verify "spec.template.spec.serviceAccount" and "spec.template.spec.serviceAccountName" in "aws-load-balancer-controller" Deployment
2. We should find the Service Account Name as "aws-load-balancer-controller"

# List Pods in YAML format
kubectl -n kube-system get pods
kubectl -n kube-system get pod <AWS-Load-Balancer-Controller-POD-NAME> -o yaml
kubectl -n kube-system get pod aws-load-balancer-controller-65b4f64d6c-h2vh4 -o yaml
Observation:
1. Verify "spec.serviceAccount" and "spec.serviceAccountName"
2. We should find the Service Account Name as "aws-load-balancer-controller"
3. Verify "spec.volumes". You should find something as below, which is a temporary credentials to access AWS Services
CHECK-1: Verify "spec.volumes.name = aws-iam-token"
  - name: aws-iam-token
    projected:
      defaultMode: 420
      sources:
      - serviceAccountToken:
          audience: sts.amazonaws.com
          expirationSeconds: 86400
          path: token
CHECK-2: Verify Volume Mounts
    volumeMounts:
    - mountPath: /var/run/secrets/eks.amazonaws.com/serviceaccount
      name: aws-iam-token
      readOnly: true          
CHECK-3: Verify ENVs whose path name is "token"
    - name: AWS_WEB_IDENTITY_TOKEN_FILE
      value: /var/run/secrets/eks.amazonaws.com/serviceaccount/token          
```

### Step-04-06: Verify TLS Certs for AWS Load Balancer Controller - Internals
```t
# List aws-load-balancer-tls secret 
kubectl -n kube-system get secret aws-load-balancer-tls -o yaml

# Verify the ca.crt and tls.crt in below websites
https://www.base64decode.org/
https://www.sslchecker.com/certdecoder

# Make a note of Common Name and SAN from above 
Common Name: aws-load-balancer-controller
SAN: aws-load-balancer-webhook-service.kube-system, aws-load-balancer-webhook-service.kube-system.svc

# List Pods in YAML format
kubectl -n kube-system get pods
kubectl -n kube-system get pod <AWS-Load-Balancer-Controller-POD-NAME> -o yaml
kubectl -n kube-system get pod aws-load-balancer-controller-65b4f64d6c-h2vh4 -o yaml
Observation:
1. Verify how the secret is mounted in AWS Load Balancer Controller Pod
CHECK-2: Verify Volume Mounts
    volumeMounts:
    - mountPath: /tmp/k8s-webhook-server/serving-certs
      name: cert
      readOnly: true
CHECK-3: Verify Volumes
  volumes:
  - name: cert
    secret:
      defaultMode: 420
      secretName: aws-load-balancer-tls
```

### Step-04-07: UNINSTALL AWS Load Balancer Controller using Helm Command (Information Purpose - SHOULD NOT EXECUTE THIS COMMAND)
- This step should not be implemented.
- This is just put it here for us to know how to uninstall aws load balancer controller from EKS Cluster
```t
# Uninstall AWS Load Balancer Controller
helm uninstall aws-load-balancer-controller -n kube-system 
```



## Step-05: Ingress Class Concept
- Understand what is Ingress Class 
- Understand how it overrides the default deprecated annotation `#kubernetes.io/ingress.class: "alb"`
- [Ingress Class Documentation Reference](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/ingress_class/)
- [Different Ingress Controllers available today](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)


## Step-06: Review IngressClass Kubernetes Manifest
- **File Location:** `08-01-Load-Balancer-Controller-Install/kube-manifests/01-ingressclass-resource.yaml`
- Understand in detail about annotation `ingressclass.kubernetes.io/is-default-class: "true"`
```yaml
apiVersion: networking.k8s.io/v1
kind: IngressClass
metadata:
  name: my-aws-ingress-class
  annotations:
    ingressclass.kubernetes.io/is-default-class: "true"
spec:
  controller: ingress.k8s.aws/alb

## Additional Note
# 1. You can mark a particular IngressClass as the default for your cluster. 
# 2. Setting the ingressclass.kubernetes.io/is-default-class annotation to true on an IngressClass resource will ensure that new Ingresses without an ingressClassName field specified will be assigned this default IngressClass.  
# 3. Reference: https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.3/guide/ingress/ingress_class/
```

## Step-07: Create IngressClass Resource
```t
# Navigate to Directory
cd 08-01-Load-Balancer-Controller-Install

# Create IngressClass Resource
kubectl apply -f kube-manifests

# Verify IngressClass Resource
kubectl get ingressclass

# Describe IngressClass Resource
kubectl describe ingressclass my-aws-ingress-class
```

## References
- [AWS Load Balancer Controller Install](https://docs.aws.amazon.com/eks/latest/userguide/aws-load-balancer-controller.html)
- [ECR Repository per region](https://docs.aws.amazon.com/eks/latest/userguide/add-ons-images.html)

---

---
title: AWS Load Balancer Controller - Ingress Basics
description: Learn AWS Load Balancer Controller - Ingress Basics
---

## Step-01: Introduction
- Discuss about the Application Architecture which we are going to deploy
- Understand the following Ingress Concepts
  - [Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/)
  - [ingressClassName](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/ingress_class/)
  - defaultBackend
  - rules

## Step-02: Review App1 Deployment kube-manifest
- **File Location:** `01-kube-manifests-default-backend/01-Nginx-App1-Deployment-and-NodePortService.yml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app1-nginx-deployment
  labels:
    app: app1-nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: app1-nginx
  template:
    metadata:
      labels:
        app: app1-nginx
    spec:
      containers:
        - name: app1-nginx
          image: stacksimplify/kube-nginxapp1:1.0.0
          ports:
            - containerPort: 80
```
## Step-03: Review App1 NodePort Service
- **File Location:** `01-kube-manifests-default-backend/01-Nginx-App1-Deployment-and-NodePortService.yml`
```yaml
apiVersion: v1
kind: Service
metadata:
  name: app1-nginx-nodeport-service
  labels:
    app: app1-nginx
  annotations:
#Important Note:  Need to add health check path annotations in service level if we are planning to use multiple targets in a load balancer    
#    alb.ingress.kubernetes.io/healthcheck-path: /app1/index.html
spec:
  type: NodePort
  selector:
    app: app1-nginx
  ports:
    - port: 80
      targetPort: 80  
```

## Step-04: Review Ingress kube-manifest with Default Backend Option
- [Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/)
- **File Location:** `01-kube-manifests-default-backend/02-ALB-Ingress-Basic.yml`
```yaml
# Annotations Reference: https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-nginxapp1
  labels:
    app: app1-nginx
  annotations:
    #kubernetes.io/ingress.class: "alb" (OLD INGRESS CLASS NOTATION - STILL WORKS BUT RECOMMENDED TO USE IngressClass Resource)
    # Ingress Core Settings
    alb.ingress.kubernetes.io/scheme: internet-facing
    # Health Check Settings
    alb.ingress.kubernetes.io/healthcheck-protocol: HTTP 
    alb.ingress.kubernetes.io/healthcheck-port: traffic-port
    alb.ingress.kubernetes.io/healthcheck-path: /app1/index.html    
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: '15'
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: '5'
    alb.ingress.kubernetes.io/success-codes: '200'
    alb.ingress.kubernetes.io/healthy-threshold-count: '2'
    alb.ingress.kubernetes.io/unhealthy-threshold-count: '2'
spec:
  ingressClassName: ic-external-lb # Ingress Class
  defaultBackend:
    service:
      name: app1-nginx-nodeport-service
      port:
        number: 80                    
```

## Step-05: Deploy kube-manifests and Verify
```t
# Change Directory
cd 08-02-ALB-Ingress-Basics

# Deploy kube-manifests
kubectl apply -f 01-kube-manifests-default-backend/

# Verify k8s Deployment and Pods
kubectl get deploy
kubectl get pods

# Verify Ingress (Make a note of Address field)
kubectl get ingress
Obsevation: 
1. Verify the ADDRESS value, we should see something like "app1ingress-1334515506.us-east-1.elb.amazonaws.com"

# Describe Ingress Controller
kubectl describe ingress ingress-nginxapp1
Observation:
1. Review Default Backend and Rules

# List Services
kubectl get svc

# Verify Application Load Balancer using 
Goto AWS Mgmt Console -> Services -> EC2 -> Load Balancers
1. Verify Listeners and Rules inside a listener
2. Verify Target Groups

# Access App using Browser
kubectl get ingress
http://<ALB-DNS-URL>
http://<ALB-DNS-URL>/app1/index.html
or
http://<INGRESS-ADDRESS-FIELD>
http://<INGRESS-ADDRESS-FIELD>/app1/index.html

# Sample from my environment (for reference only)
http://app1ingress-154912460.us-east-1.elb.amazonaws.com
http://app1ingress-154912460.us-east-1.elb.amazonaws.com/app1/index.html

# Verify AWS Load Balancer Controller logs
kubectl get po -n kube-system 
## POD1 Logs: 
kubectl -n kube-system logs -f <POD1-NAME>
kubectl -n kube-system logs -f aws-load-balancer-controller-65b4f64d6c-h2vh4
##POD2 Logs: 
kubectl -n kube-system logs -f <POD2-NAME>
kubectl -n kube-system logs -f aws-load-balancer-controller-65b4f64d6c-t7qqb
```

## Step-06: Clean Up
```t
# Delete Kubernetes Resources
kubectl delete -f 01-kube-manifests-default-backend/
```

## Step-07: Review Ingress kube-manifest with Ingress Rules
- Discuss about [Ingress Path Types](https://kubernetes.io/docs/concepts/services-networking/ingress/#path-types)
- [Better Path Matching With Path Types](https://kubernetes.io/blog/2020/04/02/improvements-to-the-ingress-api-in-kubernetes-1.18/#better-path-matching-with-path-types)
- [Sample Ingress Rule](https://kubernetes.io/docs/concepts/services-networking/ingress/#the-ingress-resource)
- **ImplementationSpecific (default):** With this path type, matching is up to the controller implementing the IngressClass. Implementations can treat this as a separate pathType or treat it identically to the Prefix or Exact path types.
- **Exact:** Matches the URL path exactly and with case sensitivity.
- **Prefix:** Matches based on a URL path prefix split by /. Matching is case sensitive and done on a path element by element basis.

- **File Location:** `02-kube-manifests-rules\02-ALB-Ingress-Basic.yml`
```yaml
# Annotations Reference: https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-nginxapp1
  labels:
    app: app1-nginx
  annotations:
    # Load Balancer Name
    alb.ingress.kubernetes.io/load-balancer-name: app1ingressrules
    #kubernetes.io/ingress.class: "alb" (OLD INGRESS CLASS NOTATION - STILL WORKS BUT RECOMMENDED TO USE IngressClass Resource)
    # Ingress Core Settings
    alb.ingress.kubernetes.io/scheme: internet-facing
    # Health Check Settings
    alb.ingress.kubernetes.io/healthcheck-protocol: HTTP 
    alb.ingress.kubernetes.io/healthcheck-port: traffic-port
    alb.ingress.kubernetes.io/healthcheck-path: /app1/index.html    
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: '15'
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: '5'
    alb.ingress.kubernetes.io/success-codes: '200'
    alb.ingress.kubernetes.io/healthy-threshold-count: '2'
    alb.ingress.kubernetes.io/unhealthy-threshold-count: '2'
spec:
  ingressClassName: ic-external-lb # Ingress Class
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app1-nginx-nodeport-service
                port: 
                  number: 80
      

# 1. If  "spec.ingressClassName: ic-external-lb" not specified, will reference default ingress class on this kubernetes cluster
# 2. Default Ingress class is nothing but for which ingress class we have the annotation `ingressclass.kubernetes.io/is-default-class: "true"`
```

## Step-08: Deploy kube-manifests and Verify
```t
# Change Directory
cd 08-02-ALB-Ingress-Basics

# Deploy kube-manifests
kubectl apply -f 02-kube-manifests-rules/

# Verify k8s Deployment and Pods
kubectl get deploy
kubectl get pods

# Verify Ingress (Make a note of Address field)
kubectl get ingress
Obsevation: 
1. Verify the ADDRESS value, we should see something like "app1ingressrules-154912460.us-east-1.elb.amazonaws.com"

# Describe Ingress Controller
kubectl describe ingress ingress-nginxapp1
Observation:
1. Review Default Backend and Rules

# List Services
kubectl get svc

# Verify Application Load Balancer using 
Goto AWS Mgmt Console -> Services -> EC2 -> Load Balancers
1. Verify Listeners and Rules inside a listener
2. Verify Target Groups

# Access App using Browser
kubectl get ingress
http://<ALB-DNS-URL>
http://<ALB-DNS-URL>/app1/index.html
or
http://<INGRESS-ADDRESS-FIELD>
http://<INGRESS-ADDRESS-FIELD>/app1/index.html

# Sample from my environment (for reference only)
http://app1ingressrules-154912460.us-east-1.elb.amazonaws.com
http://app1ingressrules-154912460.us-east-1.elb.amazonaws.com/app1/index.html

# Verify AWS Load Balancer Controller logs
kubectl get po -n kube-system 
kubectl logs -f aws-load-balancer-controller-794b7844dd-8hk7n -n kube-system
```

## Step-09: Clean Up
```t
# Delete Kubernetes Resources
kubectl delete -f 02-kube-manifests-rules/

# Verify if Ingress Deleted successfully 
kubectl get ingress
Important Note: It is going to cost us heavily if we leave ALB load balancer idle without deleting it properly

# Verify Application Load Balancer DELETED 
Goto AWS Mgmt Console -> Services -> EC2 -> Load Balancers
```

---

---
title: AWS Load Balancer Ingress Context Path Based Routing
description: Learn AWS Load Balancer Controller - Ingress Context Path Based Routing
---

## Step-01: Introduction
- Discuss about the Architecture we are going to build as part of this Section
- We are going to deploy all these 3 apps in kubernetes with context path based routing enabled in Ingress Controller
  - /app1/* - should go to app1-nginx-nodeport-service
  - /app2/* - should go to app1-nginx-nodeport-service
  - /*    - should go to  app3-nginx-nodeport-service
- As part of this process, this respective annotation `alb.ingress.kubernetes.io/healthcheck-path:` will be moved to respective application NodePort Service. 
- Only generic settings will be present in Ingress manifest annotations area `04-ALB-Ingress-ContextPath-Based-Routing.yml`  


## Step-02: Review Nginx App1, App2 & App3 Deployment & Service
- Differences for all 3 apps will be only two fields from kubernetes manifests perspective and their naming conventions
  - **Kubernetes Deployment:** Container Image name
  - **Kubernetes Node Port Service:** Health check URL path 
- **App1 Nginx: 01-Nginx-App1-Deployment-and-NodePortService.yml**
  - **image:** stacksimplify/kube-nginxapp1:1.0.0
  - **Annotation:** alb.ingress.kubernetes.io/healthcheck-path: /app1/index.html
- **App2 Nginx: 02-Nginx-App2-Deployment-and-NodePortService.yml**
  - **image:** stacksimplify/kube-nginxapp2:1.0.0
  - **Annotation:** alb.ingress.kubernetes.io/healthcheck-path: /app2/index.html
- **App3 Nginx: 03-Nginx-App3-Deployment-and-NodePortService.yml**
  - **image:** stacksimplify/kubenginx:1.0.0
  - **Annotation:** alb.ingress.kubernetes.io/healthcheck-path: /index.html



## Step-03: Create ALB Ingress Context path based Routing Kubernetes manifest
- **04-ALB-Ingress-ContextPath-Based-Routing.yml**
```yaml
# Annotations Reference: https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-cpr-demo
  annotations:
    # Load Balancer Name
    alb.ingress.kubernetes.io/load-balancer-name: cpr-ingress
    # Ingress Core Settings
    #kubernetes.io/ingress.class: "alb" (OLD INGRESS CLASS NOTATION - STILL WORKS BUT RECOMMENDED TO USE IngressClass Resource)
    alb.ingress.kubernetes.io/scheme: internet-facing
    # Health Check Settings
    alb.ingress.kubernetes.io/healthcheck-protocol: HTTP 
    alb.ingress.kubernetes.io/healthcheck-port: traffic-port
    #Important Note:  Need to add health check path annotations in service level if we are planning to use multiple targets in a load balancer    
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: '15'
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: '5'
    alb.ingress.kubernetes.io/success-codes: '200'
    alb.ingress.kubernetes.io/healthy-threshold-count: '2'
    alb.ingress.kubernetes.io/unhealthy-threshold-count: '2'   
spec:
  ingressClassName: my-aws-ingress-class   # Ingress Class                  
  rules:
    - http:
        paths:      
          - path: /app1
            pathType: Prefix
            backend:
              service:
                name: app1-nginx-nodeport-service
                port: 
                  number: 80
          - path: /app2
            pathType: Prefix
            backend:
              service:
                name: app2-nginx-nodeport-service
                port: 
                  number: 80
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app3-nginx-nodeport-service
                port: 
                  number: 80              

# Important Note-1: In path based routing order is very important, if we are going to use  "/*", try to use it at the end of all rules.                                        
                        
# 1. If  "spec.ingressClassName: my-aws-ingress-class" not specified, will reference default ingress class on this kubernetes cluster
# 2. Default Ingress class is nothing but for which ingress class we have the annotation `ingressclass.kubernetes.io/is-default-class: "true"`                      
```

## Step-04: Deploy all manifests and test
```t
# Deploy Kubernetes manifests
kubectl apply -f kube-manifests/

# List Pods
kubectl get pods

# List Services
kubectl get svc

# List Ingress Load Balancers
kubectl get ingress

# Describe Ingress and view Rules
kubectl describe ingress ingress-cpr-demo

# Verify AWS Load Balancer Controller logs
kubectl -n kube-system  get pods 
kubectl -n kube-system logs -f aws-load-balancer-controller-794b7844dd-8hk7n 
```

## Step-05: Verify Application Load Balancer on AWS Management Console**
- Verify Load Balancer
    - In Listeners Tab, click on **View/Edit Rules** under Rules
- Verify Target Groups
    - GroupD Details
    - Targets: Ensure they are healthy
    - Verify Health check path
    - Verify all 3 targets are healthy)
```t
# Access Application
http://<ALB-DNS-URL>/app1/index.html
http://<ALB-DNS-URL>/app2/index.html
http://<ALB-DNS-URL>/
```

## Step-06: Test Order in Context path based routing
### Step-0-01: Move Root Context Path to top
- **File:** 04-ALB-Ingress-ContextPath-Based-Routing.yml
```yaml
  ingressClassName: my-aws-ingress-class   # Ingress Class                  
  rules:
    - http:
        paths:      
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app3-nginx-nodeport-service
                port: 
                  number: 80           
          - path: /app1
            pathType: Prefix
            backend:
              service:
                name: app1-nginx-nodeport-service
                port: 
                  number: 80
          - path: /app2
            pathType: Prefix
            backend:
              service:
                name: app2-nginx-nodeport-service
                port: 
                  number: 80
```
### Step-06-02: Deploy Changes and Verify
```t
# Deploy Changes
kubectl apply -f kube-manifests/

# Access Application (Open in new incognito window)
http://<ALB-DNS-URL>/app1/index.html  -- SHOULD FAIL
http://<ALB-DNS-URL>/app2/index.html  -- SHOULD FAIL
http://<ALB-DNS-URL>/  - SHOULD PASS
```

## Step-07: Roll back changes in 04-ALB-Ingress-ContextPath-Based-Routing.yml
```yaml
spec:
  ingressClassName: my-aws-ingress-class   # Ingress Class                  
  rules:
    - http:
        paths:      
          - path: /app1
            pathType: Prefix
            backend:
              service:
                name: app1-nginx-nodeport-service
                port: 
                  number: 80
          - path: /app2
            pathType: Prefix
            backend:
              service:
                name: app2-nginx-nodeport-service
                port: 
                  number: 80
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app3-nginx-nodeport-service
                port: 
                  number: 80              
```

## Step-08: Clean Up
```t
# Clean-Up
kubectl delete -f kube-manifests/
```

---

---
title: AWS Load Balancer Ingress SSL
description: Learn AWS Load Balancer Controller - Ingress SSL
---

## Step-01: Introduction
- We are going to register a new DNS in AWS Route53
- We are going to create a SSL certificate 
- Add Annotations related to SSL Certificate in Ingress manifest
- Deploy the manifests and test
- Clean-Up

## Step-02: Pre-requisite - Register a Domain in Route53 (if not exists)
- Goto Services -> Route53 -> Registered Domains
- Click on **Register Domain**
- Provide **desired domain: somedomain.com** and click on **check** (In my case its going to be `stacksimplify.com`)
- Click on **Add to cart** and click on **Continue**
- Provide your **Contact Details** and click on **Continue**
- Enable Automatic Renewal
- Accept **Terms and Conditions**
- Click on **Complete Order**

## Step-03: Create a SSL Certificate in Certificate Manager
- Pre-requisite: You should have a registered domain in Route53 
- Go to Services -> Certificate Manager -> Create a Certificate
- Click on **Request a Certificate**
  - Choose the type of certificate for ACM to provide: Request a public certificate
  - Add domain names: *.yourdomain.com (in my case it is going to be `*.stacksimplify.com`)
  - Select a Validation Method: **DNS Validation**
  - Click on **Confirm & Request**    
- **Validation**
  - Click on **Create record in Route 53**  
- Wait for 5 to 10 minutes and check the **Validation Status**  

## Step-04: Add annotations related to SSL
- **04-ALB-Ingress-SSL.yml**
```yaml
    ## SSL Settings
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}, {"HTTP":80}]'
    alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:us-east-1:180789647333:certificate/632a3ff6-3f6d-464c-9121-b9d97481a76b
    #alb.ingress.kubernetes.io/ssl-policy: ELBSecurityPolicy-TLS-1-1-2017-01 #Optional (Picks default if not used)    
```
## Step-05: Deploy all manifests and test
### Deploy and Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)

## Step-06: Add DNS in Route53   
- Go to **Services -> Route 53**
- Go to **Hosted Zones**
  - Click on **yourdomain.com** (in my case stacksimplify.com)
- Create a **Record Set**
  - **Name:** ssldemo101.stacksimplify.com
  - **Alias:** yes
  - **Alias Target:** Copy our ALB DNS Name here (Sample: ssl-ingress-551932098.us-east-1.elb.amazonaws.com)
  - Click on **Create**
  
## Step-07: Access Application using newly registered DNS Name
- **Access Application**
- **Important Note:** Instead of `stacksimplify.com` you need to replace with your registered Route53 domain (Refer pre-requisite Step-02)
```t
# HTTP URLs
http://ssldemo101.stacksimplify.com/app1/index.html
http://ssldemo101.stacksimplify.com/app2/index.html
http://ssldemo101.stacksimplify.com/

# HTTPS URLs
https://ssldemo101.stacksimplify.com/app1/index.html
https://ssldemo101.stacksimplify.com/app2/index.html
https://ssldemo101.stacksimplify.com/
```

## Annotation Reference
- [AWS Load Balancer Controller Annotation Reference](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/ingress/annotations/)


---

---
title: AWS Load Balancer - Ingress SSL HTTP to HTTPS Redirect
description: Learn AWS Load Balancer - Ingress SSL HTTP to HTTPS Redirect
---

## Step-01: Add annotations related to SSL Redirect
- **File Name:** 04-ALB-Ingress-SSL-Redirect.yml
- Redirect from HTTP to HTTPS
```yaml
    # SSL Redirect Setting
    alb.ingress.kubernetes.io/ssl-redirect: '443'   
```

## Step-02: Deploy all manifests and test

### Deploy and Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)
 
## Step-03: Access Application using newly registered DNS Name
- **Access Application**
```t
# HTTP URLs (Should Redirect to HTTPS)
http://ssldemo101.stacksimplify.com/app1/index.html
http://ssldemo101.stacksimplify.com/app2/index.html
http://ssldemo101.stacksimplify.com/

# HTTPS URLs
https://ssldemo101.stacksimplify.com/app1/index.html
https://ssldemo101.stacksimplify.com/app2/index.html
https://ssldemo101.stacksimplify.com/
```

## Step-04: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/

## Delete Route53 Record Set
- Delete Route53 Record we created (ssldemo101.stacksimplify.com)
```

## Annotation Reference
- [AWS Load Balancer Controller Annotation Reference](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/ingress/annotations/)


---

---
title: AWS Load Balancer Controller - External DNS Install
description: Learn AWS Load Balancer Controller - External DNS Install
---

## Step-01: Introduction
- **External DNS:** Used for Updating Route53 RecordSets from Kubernetes 
- We need to create IAM Policy, k8s Service Account & IAM Role and associate them together for external-dns pod to add or remove entries in AWS Route53 Hosted Zones. 
- Update External-DNS default manifest to support our needs
- Deploy & Verify logs

## Step-02: Create IAM Policy
- This IAM policy will allow external-dns pod to add, remove DNS entries (Record Sets in a Hosted Zone) in AWS Route53 service
- Go to Services -> IAM -> Policies -> Create Policy
  - Click on **JSON** Tab and copy paste below JSON
  - Click on **Visual editor** tab to validate
  - Click on **Review Policy**
  - **Name:** AllowExternalDNSUpdates 
  - **Description:** Allow access to Route53 Resources for ExternalDNS
  - Click on **Create Policy**  

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "route53:ChangeResourceRecordSets"
      ],
      "Resource": [
        "arn:aws:route53:::hostedzone/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "route53:ListHostedZones",
        "route53:ListResourceRecordSets"
      ],
      "Resource": [
        "*"
      ]
    }
  ]
}
```
- Make a note of Policy ARN which we will use in next step
```t
# Policy ARN
arn:aws:iam::528163014577:policy/AllowExternalDNSUpdates
```  


## Step-03: Create IAM Role, k8s Service Account & Associate IAM Policy
- As part of this step, we are going to create a k8s Service Account named `external-dns` and also a AWS IAM role and associate them by annotating role ARN in Service Account.
- In addition, we are also going to associate the AWS IAM Policy `AllowExternalDNSUpdates` to the newly created AWS IAM Role.
### Step-03-01: Create IAM Role, k8s Service Account & Associate IAM Policy
```t
# Template
eksctl create iamserviceaccount \
    --name service_account_name \
    --namespace service_account_namespace \
    --cluster cluster_name \
    --attach-policy-arn IAM_policy_ARN \
    --approve \
    --override-existing-serviceaccounts

# Replaced name, namespace, cluster, IAM Policy arn 
eksctl create iamserviceaccount \
    --name external-dns \
    --namespace default \
    --cluster eksdemo1 \
    --attach-policy-arn arn:aws:iam::528163014577:policy/AllowExternalDNSUpdates \
    --approve \
    --override-existing-serviceaccounts
```
### Step-03-02: Verify the Service Account
- Verify external-dns service account, primarily verify annotation related to IAM Role
```t
# List Service Account
kubectl get sa external-dns

# Describe Service Account
kubectl describe sa external-dns
Observation: 
1. Verify the Annotations and you should see the IAM Role is present on the Service Account
```
### Step-03-03: Verify CloudFormation Stack
- Go to Services -> CloudFormation
- Verify the latest CFN Stack created.
- Click on **Resources** tab
- Click on link  in **Physical ID** field which will take us to **IAM Role** directly

### Step-03-04: Verify IAM Role & IAM Policy
- With above step in CFN, we will be landed in IAM Role created for external-dns. 
- Verify in **Permissions** tab we have a policy named **AllowExternalDNSUpdates**
- Now make a note of that Role ARN, this we need to update in External-DNS k8s manifest
```t
# Make a note of Role ARN
arn:aws:iam::528163014577:role/eksctl-eksdemo1-addon-iamserviceaccount-defa-Role1-I0GD3C75CBX5
```

### Step-03-05: Verify IAM Service Accounts using eksctl
- You can also make a note of External DNS Role ARN from here too. 
```t
# List IAM Service Accounts using eksctl
eksctl get iamserviceaccount --cluster eksdemo1

# Sample Output
Kalyans-Mac-mini:08-06-ALB-Ingress-ExternalDNS kalyanreddy$ eksctl get iamserviceaccount --cluster eksdemo1
2022-02-11 09:34:39 [ℹ]  eksctl version 0.71.0
2022-02-11 09:34:39 [ℹ]  using region us-east-1
NAMESPACE	NAME				ROLE ARN
default		external-dns			arn:aws:iam::180789647333:role/eksctl-eksdemo1-addon-iamserviceaccount-defa-Role1-JTO29BVZMA2N
kube-system	aws-load-balancer-controller	arn:aws:iam::180789647333:role/eksctl-eksdemo1-addon-iamserviceaccount-kube-Role1-EFQB4C26EALH
Kalyans-Mac-mini:08-06-ALB-Ingress-ExternalDNS kalyanreddy$ 
```


## Step-04: Update External DNS Kubernetes manifest
- **Original Template** you can find in https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md
- **File Location:** kube-manifests/01-Deploy-ExternalDNS.yml
### Change-1: Line number 9: IAM Role update
  - Copy the role-arn you have made a note at the end of step-03 and replace at line no 9.
```yaml
    eks.amazonaws.com/role-arn: arn:aws:iam::180789647333:role/eksctl-eksdemo1-addon-iamserviceaccount-defa-Role1-JTO29BVZMA2N
```
### Chnage-2: Line 55, 56: Commented them
- We used eksctl to create IAM role and attached the `AllowExternalDNSUpdates` policy
- We didnt use KIAM or Kube2IAM so we don't need these two lines, so commented
```yaml
      #annotations:  
        #iam.amazonaws.com/role: arn:aws:iam::ACCOUNT-ID:role/IAM-SERVICE-ROLE-NAME    
```
### Change-3: Line 65, 67: Commented them
```yaml
        # - --domain-filter=external-dns-test.my-org.com # will make ExternalDNS see only the hosted zones matching provided domain, omit to process all available hosted zones
       # - --policy=upsert-only # would prevent ExternalDNS from deleting any records, omit to enable full synchronization
```

### Change-4: Line 61: Get latest Docker Image name
- [Get latest external dns image name](https://github.com/kubernetes-sigs/external-dns/releases/tag/v0.10.2)
```yaml
    spec:
      serviceAccountName: external-dns
      containers:
      - name: external-dns
        image: k8s.gcr.io/external-dns/external-dns:v0.10.2
```

## Step-05: Deploy ExternalDNS
- Deploy the manifest
```t
# Change Directory
cd 08-06-Deploy-ExternalDNS-on-EKS

# Deploy external DNS
kubectl apply -f kube-manifests/

# List All resources from default Namespace
kubectl get all

# List pods (external-dns pod should be in running state)
kubectl get pods

# Verify Deployment by checking logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')
```

## References
- https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/alb-ingress.md
- https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md

---

---
title: AWS Load Balancer Controller - External DNS & Ingress
description: Learn AWS Load Balancer Controller - External DNS & Ingress
---

## Step-01: Update Ingress manifest by adding External DNS Annotation
- Added annotation with two DNS Names
  - dnstest901.kubeoncloud.com
  - dnstest902.kubeoncloud.com
- Once we deploy the application, we should be able to access our Applications with both DNS Names.   
- **File Name:** 04-ALB-Ingress-SSL-Redirect-ExternalDNS.yml
```yaml
    # External DNS - For creating a Record Set in Route53
    external-dns.alpha.kubernetes.io/hostname: dnstest901.stacksimplify.com, dnstest902.stacksimplify.com
```
- In your case it is going to be, replace `yourdomain` with your domain name
  - dnstest901.yourdoamin.com
  - dnstest902.yourdoamin.com

## Step-02: Deploy all Application Kubernetes Manifests
### Deploy
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)

### Verify External DNS Log
```t
# Verify External DNS logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')
```
### Verify Route53
- Go to Services -> Route53
- You should see **Record Sets** added for `dnstest901.stacksimplify.com`, `dnstest902.stacksimplify.com`

## Step-04: Access Application using newly registered DNS Name
### Perform nslookup tests before accessing Application
- Test if our new DNS entries registered and resolving to an IP Address
```t
# nslookup commands
nslookup dnstest901.stacksimplify.com
nslookup dnstest902.stacksimplify.com
```
### Access Application using dnstest1 domain
```t
# HTTP URLs (Should Redirect to HTTPS)
http://dnstest901.stacksimplify.com/app1/index.html
http://dnstest901.stacksimplify.com/app2/index.html
http://dnstest901.stacksimplify.com/
```

### Access Application using dnstest2 domain
```t
# HTTP URLs (Should Redirect to HTTPS)
http://dnstest902.stacksimplify.com/app1/index.html
http://dnstest902.stacksimplify.com/app2/index.html
http://dnstest902.stacksimplify.com/
```


## Step-05: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/

## Verify Route53 Record Set to ensure our DNS records got deleted
- Go to Route53 -> Hosted Zones -> Records 
- The below records should be deleted automatically
  - dnstest901.stacksimplify.com
  - dnstest902.stacksimplify.com
```


## References
- https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/alb-ingress.md
- https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md



---

---
title: AWS Load Balancer Controller - External DNS & Service
description: Learn AWS Load Balancer Controller - External DNS & Kubernetes Service
---

## Step-01: Introduction
- We will create a Kubernetes Service of `type: LoadBalancer`
- We will annotate that Service with external DNS hostname `external-dns.alpha.kubernetes.io/hostname: externaldns-k8s-service-demo101.stacksimplify.com` which will register the DNS in Route53 for that respective load balancer

## Step-02: 02-Nginx-App1-LoadBalancer-Service.yml
```yaml
apiVersion: v1
kind: Service
metadata:
  name: app1-nginx-loadbalancer-service
  labels:
    app: app1-nginx
  annotations:
    external-dns.alpha.kubernetes.io/hostname: externaldns-k8s-service-demo101.stacksimplify.com
spec:
  type: LoadBalancer
  selector:
    app: app1-nginx
  ports:
    - port: 80
      targetPort: 80  
```
## Step-03: Deploy & Verify

### Deploy & Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify Service
kubectl get svc
```
### Verify Load Balancer 
- Go to EC2 -> Load Balancers -> Verify Load Balancer Settings

### Verify External DNS Log
```t
# Verify External DNS logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')
```
### Verify Route53
- Go to Services -> Route53
- You should see **Record Sets** added for `externaldns-k8s-service-demo101.stacksimplify.com`


## Step-04: Access Application using newly registered DNS Name
### Perform nslookup tests before accessing Application
- Test if our new DNS entries registered and resolving to an IP Address
```t
# nslookup commands
nslookup externaldns-k8s-service-demo101.stacksimplify.com
```
### Access Application using DNS domain
```t
# HTTP URL
http://externaldns-k8s-service-demo101.stacksimplify.com/app1/index.html
```

## Step-05: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/

## Verify Route53 Record Set to ensure our DNS records got deleted
- Go to Route53 -> Hosted Zones -> Records 
- The below records should be deleted automatically
  - externaldns-k8s-service-demo101.stacksimplify.com
```


## References
- https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/alb-ingress.md
- https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md

---

---
title: AWS Load Balancer Controller - Ingress Host Header Routing
description: Learn AWS Load Balancer Controller - Ingress Host Header Routing
---

## Step-01: Introduction
- Implement Host Header routing using Ingress
- We can also call it has name based virtual host routing

## Step-02: Review Ingress Manifests for Host Header Routing
- **File Name:** 04-ALB-Ingress-HostHeader-Routing.yml
```yaml
# Annotations Reference: https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-namedbasedvhost-demo
  annotations:
    # Load Balancer Name
    alb.ingress.kubernetes.io/load-balancer-name: namedbasedvhost-ingress
    # Ingress Core Settings
    #kubernetes.io/ingress.class: "alb" (OLD INGRESS CLASS NOTATION - STILL WORKS BUT RECOMMENDED TO USE IngressClass Resource)
    alb.ingress.kubernetes.io/scheme: internet-facing
    # Health Check Settings
    alb.ingress.kubernetes.io/healthcheck-protocol: HTTP 
    alb.ingress.kubernetes.io/healthcheck-port: traffic-port
    #Important Note:  Need to add health check path annotations in service level if we are planning to use multiple targets in a load balancer    
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: '15'
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: '5'
    alb.ingress.kubernetes.io/success-codes: '200'
    alb.ingress.kubernetes.io/healthy-threshold-count: '2'
    alb.ingress.kubernetes.io/unhealthy-threshold-count: '2'   
    ## SSL Settings
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}, {"HTTP":80}]'
    alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:us-east-1:180789647333:certificate/632a3ff6-3f6d-464c-9121-b9d97481a76b
    #alb.ingress.kubernetes.io/ssl-policy: ELBSecurityPolicy-TLS-1-1-2017-01 #Optional (Picks default if not used)    
    # SSL Redirect Setting
    alb.ingress.kubernetes.io/ssl-redirect: '443'
    # External DNS - For creating a Record Set in Route53
    external-dns.alpha.kubernetes.io/hostname: default101.stacksimplify.com 
spec:
  ingressClassName: my-aws-ingress-class   # Ingress Class                  
  defaultBackend:
    service:
      name: app3-nginx-nodeport-service
      port:
        number: 80     
  rules:
    - host: app101.stacksimplify.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app1-nginx-nodeport-service
                port: 
                  number: 80
    - host: app201.stacksimplify.com
      http:
        paths:                  
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app2-nginx-nodeport-service
                port: 
                  number: 80

# Important Note-1: In path based routing order is very important, if we are going to use  "/*", try to use it at the end of all rules.                                        
                        
# 1. If  "spec.ingressClassName: my-aws-ingress-class" not specified, will reference default ingress class on this kubernetes cluster
# 2. Default Ingress class is nothing but for which ingress class we have the annotation `ingressclass.kubernetes.io/is-default-class: "true"`     
                         
```

## Step-03: Deploy all Application Kubernetes Manifests and Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)

### Verify External DNS Log
```t
# Verify External DNS logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')
```
### Verify Route53
- Go to Services -> Route53
- You should see **Record Sets** added for 
  - default101.stacksimplify.com
  - app101.stacksimplify.com
  - app201.stacksimplify.com

## Step-04: Access Application using newly registered DNS Name
### Perform nslookup tests before accessing Application
- Test if our new DNS entries registered and resolving to an IP Address
```t
# nslookup commands
nslookup default101.stacksimplify.com
nslookup app101.stacksimplify.com
nslookup app201.stacksimplify.com
```
### Positive Case: Access Application using DNS domain
```t
# Access App1
http://app101.stacksimplify.com/app1/index.html

# Access App2
http://app201.stacksimplify.com/app2/index.html

# Access Default App (App3)
http://default101.stacksimplify.com
```

### Negative Case: Access Application using DNS domain
```t
# Access App2 using App1 DNS Domain
http://app101.stacksimplify.com/app2/index.html  -- SHOULD FAIL

# Access App1 using App2 DNS Domain
http://app201.stacksimplify.com/app1/index.html  -- SHOULD FAIL

# Access App1 and App2 using Default Domain
http://default101.stacksimplify.com/app1/index.html -- SHOULD FAIL
http://default101.stacksimplify.com/app2/index.html -- SHOULD FAIL
```

## Step-05: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/

## Verify Route53 Record Set to ensure our DNS records got deleted
- Go to Route53 -> Hosted Zones -> Records 
- The below records should be deleted automatically
  - default101.stacksimplify.com
  - app101.stacksimplify.com
  - app201.stacksimplify.com
```


---

---
title: AWS Load Balancer Controller - Ingress SSL Discovery Host
description: Learn AWS Load Balancer Controller - Ingress SSL Discovery Host
---

## Step-01: Introduction
- Automatically disover SSL Certificate from AWS Certificate Manager Service using `spec.rules.host`
- In this approach, with the specified domain name if we have the SSL Certificate created in AWS Certificate Manager, that certificate will be automatically detected and associated to Application Load Balancer.
- We don't need to get the SSL Certificate ARN and update it in Kubernetes Ingress Manifest
- Discovers via Ingress rule host and attaches a cert for `app102.stacksimplify.com` or `*.stacksimplify.com` to the ALB

## Step-02: Discover via Ingress "spec.rules.host"
```yaml
# Annotations Reference: https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-certdiscoveryhost-demo
  annotations:
    # Load Balancer Name
    alb.ingress.kubernetes.io/load-balancer-name: certdiscoveryhost-ingress
    # Ingress Core Settings
    #kubernetes.io/ingress.class: "alb" (OLD INGRESS CLASS NOTATION - STILL WORKS BUT RECOMMENDED TO USE IngressClass Resource)
    alb.ingress.kubernetes.io/scheme: internet-facing
    # Health Check Settings
    alb.ingress.kubernetes.io/healthcheck-protocol: HTTP 
    alb.ingress.kubernetes.io/healthcheck-port: traffic-port
    #Important Note:  Need to add health check path annotations in service level if we are planning to use multiple targets in a load balancer    
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: '15'
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: '5'
    alb.ingress.kubernetes.io/success-codes: '200'
    alb.ingress.kubernetes.io/healthy-threshold-count: '2'
    alb.ingress.kubernetes.io/unhealthy-threshold-count: '2'   
    ## SSL Settings
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}, {"HTTP":80}]'
    #alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:us-east-1:180789647333:certificate/632a3ff6-3f6d-464c-9121-b9d97481a76b
    #alb.ingress.kubernetes.io/ssl-policy: ELBSecurityPolicy-TLS-1-1-2017-01 #Optional (Picks default if not used)    
    # SSL Redirect Setting
    alb.ingress.kubernetes.io/ssl-redirect: '443'
    # External DNS - For creating a Record Set in Route53
    external-dns.alpha.kubernetes.io/hostname: default102.stacksimplify.com 
spec:
  ingressClassName: my-aws-ingress-class   # Ingress Class                  
  defaultBackend:
    service:
      name: app3-nginx-nodeport-service
      port:
        number: 80     
  rules:
    - host: app102.stacksimplify.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app1-nginx-nodeport-service
                port: 
                  number: 80
    - host: app202.stacksimplify.com
      http:
        paths:                  
          - path: /
            pathType: Prefix
            backend:
              service:
                name: app2-nginx-nodeport-service
                port: 
                  number: 80

# Important Note-1: In path based routing order is very important, if we are going to use  "/*", try to use it at the end of all rules.                                        
                        
# 1. If  "spec.ingressClassName: my-aws-ingress-class" not specified, will reference default ingress class on this kubernetes cluster
# 2. Default Ingress class is nothing but for which ingress class we have the annotation `ingressclass.kubernetes.io/is-default-class: "true"`    
```


## Step-03: Deploy all Application Kubernetes Manifests and Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)
- **PRIMARILY VERIFY - CERTIFICATE ASSOCIATED TO APPLICATION LOAD BALANCER**

### Verify External DNS Log
```t
# Verify External DNS logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')
```
### Verify Route53
- Go to Services -> Route53
- You should see **Record Sets** added for 
  - default102.stacksimplify.com
  - app102.stacksimplify.com
  - app202.stacksimplify.com

## Step-04: Access Application using newly registered DNS Name
### Perform nslookup tests before accessing Application
- Test if our new DNS entries registered and resolving to an IP Address
```t
# nslookup commands
nslookup default102.stacksimplify.com
nslookup app102.stacksimplify.com
nslookup app202.stacksimplify.com
```
### Positive Case: Access Application using DNS domain
```t
# Access App1
http://app102.stacksimplify.com/app1/index.html

# Access App2
http://app202.stacksimplify.com/app2/index.html

# Access Default App (App3)
http://default102.stacksimplify.com
```

## Step-05: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/

## Verify Route53 Record Set to ensure our DNS records got deleted
- Go to Route53 -> Hosted Zones -> Records 
- The below records should be deleted automatically
  - default102.stacksimplify.com
  - app102.stacksimplify.com
  - app202.stacksimplify.com
```


## References
- https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/ingress/cert_discovery/

---

---
title: AWS Load Balancer Controller - Ingress SSL Discovery Host
description: Learn AWS Load Balancer Controller - Ingress SSL Discovery Host
---

## Step-01: Introduction
- Automatically disover SSL Certificate from AWS Certificate Manager Service using `spec.tls.host`
- In this approach, with the specified domain name if we have the SSL Certificate created in AWS Certificate Manager, that certificate will be automatically detected and associated to Application Load Balancer.
- We don't need to get the SSL Certificate ARN and update it in Kubernetes Ingress Manifest
- Discovers via Ingress rule host and attaches a cert for `app102.stacksimplify.com` or `*.stacksimplify.com` to the ALB

## Step-02: Discover via Ingress "spec.tls.hosts"
```yaml
# Annotations Reference: https://kubernetes-sigs.github.io/aws-load-balancer-controller/latest/guide/ingress/annotations/
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-certdiscoverytls-demo
  annotations:
    # Load Balancer Name
    alb.ingress.kubernetes.io/load-balancer-name: certdiscoverytls-ingress
    # Ingress Core Settings
    #kubernetes.io/ingress.class: "alb" (OLD INGRESS CLASS NOTATION - STILL WORKS BUT RECOMMENDED TO USE IngressClass Resource)
    alb.ingress.kubernetes.io/scheme: internet-facing
    # Health Check Settings
    alb.ingress.kubernetes.io/healthcheck-protocol: HTTP 
    alb.ingress.kubernetes.io/healthcheck-port: traffic-port
    #Important Note:  Need to add health check path annotations in service level if we are planning to use multiple targets in a load balancer    
    alb.ingress.kubernetes.io/healthcheck-interval-seconds: '15'
    alb.ingress.kubernetes.io/healthcheck-timeout-seconds: '5'
    alb.ingress.kubernetes.io/success-codes: '200'
    alb.ingress.kubernetes.io/healthy-threshold-count: '2'
    alb.ingress.kubernetes.io/unhealthy-threshold-count: '2'   
    ## SSL Settings
    alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}, {"HTTP":80}]'
    #alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:us-east-1:180789647333:certificate/632a3ff6-3f6d-464c-9121-b9d97481a76b
    #alb.ingress.kubernetes.io/ssl-policy: ELBSecurityPolicy-TLS-1-1-2017-01 #Optional (Picks default if not used)    
    # SSL Redirect Setting
    alb.ingress.kubernetes.io/ssl-redirect: '443'
    # External DNS - For creating a Record Set in Route53
    external-dns.alpha.kubernetes.io/hostname: certdiscovery-tls-101.stacksimplify.com 
spec:
  ingressClassName: my-aws-ingress-class   # Ingress Class                  
  defaultBackend:
    service:
      name: app3-nginx-nodeport-service
      port:
        number: 80     
  tls:
  - hosts:
    - "*.stacksimplify.com"
  rules:
    - http:
        paths:
          - path: /app1
            pathType: Prefix
            backend:
              service:
                name: app1-nginx-nodeport-service
                port: 
                  number: 80
    - http:
        paths:                  
          - path: /app2
            pathType: Prefix
            backend:
              service:
                name: app2-nginx-nodeport-service
                port: 
                  number: 80

# Important Note-1: In path based routing order is very important, if we are going to use  "/*", try to use it at the end of all rules.                                        
                        
# 1. If  "spec.ingressClassName: my-aws-ingress-class" not specified, will reference default ingress class on this kubernetes cluster
# 2. Default Ingress class is nothing but for which ingress class we have the annotation `ingressclass.kubernetes.io/is-default-class: "true"` 
 ```


## Step-03: Deploy all Application Kubernetes Manifests and Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)
- **PRIMARILY VERIFY - CERTIFICATE ASSOCIATED TO APPLICATION LOAD BALANCER**

### Verify External DNS Log
```t
# Verify External DNS logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')
```
### Verify Route53
- Go to Services -> Route53
- You should see **Record Sets** added for 
  - certdiscovery-tls-901.stacksimplify.com 


## Step-04: Access Application using newly registered DNS Name
### Perform nslookup tests before accessing Application
- Test if our new DNS entries registered and resolving to an IP Address
```t
# nslookup commands
nslookup certdiscovery-tls-101.stacksimplify.com 
```
### Access Application using DNS domain
```t
# Access App1
http://certdiscovery-tls-101.stacksimplify.com/app1/index.html

# Access App2
http://certdiscovery-tls-101.stacksimplify.com/app2/index.html

# Access Default App (App3)
http://certdiscovery-tls-101.stacksimplify.com
```

## Step-05: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/

## Verify Route53 Record Set to ensure our DNS records got deleted
- Go to Route53 -> Hosted Zones -> Records 
- The below records should be deleted automatically
  - certdiscovery-tls-101.stacksimplify.com 
```


## References
- https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/ingress/cert_discovery/


---

---
title: AWS Load Balancer Controller - Ingress Groups
description: Learn AWS Load Balancer Controller - Ingress Groups
---

## Step-01: Introduction
- IngressGroup feature enables you to group multiple Ingress resources together. 
- The controller will automatically merge Ingress rules for all Ingresses within IngressGroup and support them with a single ALB. 
- In addition, most annotations defined on a Ingress only applies to the paths defined by that Ingress.
- Demonstrate Ingress Groups concept with two Applications. 

## Step-02: Review App1 Ingress Manifest - Key Lines
- **File Name:** `kube-manifests/app1/02-App1-Ingress.yml`
```yaml
    # Ingress Groups
    alb.ingress.kubernetes.io/group.name: myapps.web
    alb.ingress.kubernetes.io/group.order: '10'
```

## Step-03: Review App2 Ingress Manifest - Key Lines
- **File Name:** `kube-manifests/app2/02-App2-Ingress.yml`
```yaml
    # Ingress Groups
    alb.ingress.kubernetes.io/group.name: myapps.web
    alb.ingress.kubernetes.io/group.order: '20'
```

## Step-04: Review App3 Ingress Manifest - Key Lines
```yaml
    # Ingress Groups
    alb.ingress.kubernetes.io/group.name: myapps.web
    alb.ingress.kubernetes.io/group.order: '30'
```

## Step-05: Deploy Apps with two Ingress Resources
```t
# Deploy both Apps
kubectl apply -R -f kube-manifests

# Verify Pods
kubectl get pods

# Verify Ingress
kubectl  get ingress
Observation:
1. Three Ingress resources will be created with same ADDRESS value
2. Three Ingress Resources are merged to a single Application Load Balancer as those belong to same Ingress group "myapps.web"
```

## Step-06: Verify on AWS Mgmt Console
- Go to Services -> EC2 -> Load Balancers 
- Verify Routing Rules for `/app1` and `/app2` and `default backend`

## Step-07: Verify by accessing in browser
```t
# Web URLs
http://ingress-groups-demo601.stacksimplify.com/app1/index.html
http://ingress-groups-demo601.stacksimplify.com/app2/index.html
http://ingress-groups-demo601.stacksimplify.com
```

## Step-08: Clean-Up
```t
# Delete Apps from k8s cluster
kubectl delete -R -f kube-manifests/

## Verify Route53 Record Set to ensure our DNS records got deleted
- Go to Route53 -> Hosted Zones -> Records 
- The below records should be deleted automatically
  - ingress-groups-demo601.stacksimplify.com
```


---

---
title: AWS Load Balancer Controller - Ingress Target Type IP
description: Learn AWS Load Balancer Controller - Ingress Target Type IP
---

## Step-01: Introduction
- `alb.ingress.kubernetes.io/target-type` specifies how to route traffic to pods. 
- You can choose between `instance` and `ip`
- **Instance Mode:** `instance mode` will route traffic to all ec2 instances within cluster on NodePort opened for your service.
- **IP Mode:** `ip mode` is required for sticky sessions to work with Application Load Balancers.


## Step-02: Ingress Manifest - Add target-type
- **File Name:** 04-ALB-Ingress-target-type-ip.yml
```yaml
    # Target Type: IP
    alb.ingress.kubernetes.io/target-type: ip   
```

## Step-03: Deploy all Application Kubernetes Manifests and Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)
- **PRIMARILY VERIFY - TARGET GROUPS which contain thePOD IPs instead of WORKER NODE IP with NODE PORTS**
```t
# List Pods and their IPs
kubectl get pods -o wide
```

### Verify External DNS Log
```t
# Verify External DNS logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')
```
### Verify Route53
- Go to Services -> Route53
- You should see **Record Sets** added for 
  - target-type-ip-501.stacksimplify.com 


## Step-04: Access Application using newly registered DNS Name
### Perform nslookup tests before accessing Application
- Test if our new DNS entries registered and resolving to an IP Address
```t
# nslookup commands
nslookup target-type-ip-501.stacksimplify.com 
```
### Access Application using DNS domain
```t
# Access App1
http://target-type-ip-501.stacksimplify.com /app1/index.html

# Access App2
http://target-type-ip-501.stacksimplify.com /app2/index.html

# Access Default App (App3)
http://target-type-ip-501.stacksimplify.com 
```

## Step-05: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/

## Verify Route53 Record Set to ensure our DNS records got deleted
- Go to Route53 -> Hosted Zones -> Records 
- The below records should be deleted automatically
  - target-type-ip-501.stacksimplify.com 
```


---

---
title: AWS Load Balancer Controller - Ingress Internal LB
description: Learn AWS Load Balancer Controller - Ingress Internal LB
---

## Step-01: Introduction
- Create Internal Application Load Balancer using Ingress
- To test the Internal LB, use the `curl-pod`
- Deploy `curl-pod`
- Connect to `curl-pod` and test Internal LB from `curl-pod`

## Step-02: Update Ingress Scheme annotation to Internal
- **File Name:** 04-ALB-Ingress-Internal-LB.yml
```yaml
    # Creates Internal Application Load Balancer
    alb.ingress.kubernetes.io/scheme: internal 
```

## Step-03: Deploy all Application Kubernetes Manifests and Verify
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Ingress Resource
kubectl get ingress

# Verify Apps
kubectl get deploy
kubectl get pods

# Verify NodePort Services
kubectl get svc
```
### Verify Load Balancer & Target Groups
- Load Balancer -  Listeneres (Verify both 80 & 443) 
- Load Balancer - Rules (Verify both 80 & 443 listeners) 
- Target Groups - Group Details (Verify Health check path)
- Target Groups - Targets (Verify all 3 targets are healthy)

## Step-04: How to test this Internal Load Balancer? 
- We are going to deploy a `curl-pod` in EKS Cluster
- We connect to that `curl-pod` in EKS Cluster and test using `curl commands` for our sample applications load balanced using this Internal Application Load Balancer


## Step-05: curl-pod Kubernetes Manifest
- **File Name:** kube-manifests-curl/01-curl-pod.yml
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: curl-pod
spec:
  containers:
  - name: curl
    image: curlimages/curl 
    command: [ "sleep", "600" ]
```

## Step-06: Deploy curl-pod and Verify Internal LB
```t
# Deploy curl-pod
kubectl apply -f kube-manifests-curl

# Will open up a terminal session into the container
kubectl exec -it curl-pod -- sh

# We can now curl external addresses or internal services:
curl http://google.com/
curl <INTERNAL-INGRESS-LB-DNS>

# Default Backend Curl Test
curl internal-ingress-internal-lb-1839544354.us-east-1.elb.amazonaws.com

# App1 Curl Test
curl internal-ingress-internal-lb-1839544354.us-east-1.elb.amazonaws.com/app1/index.html

# App2 Curl Test
curl internal-ingress-internal-lb-1839544354.us-east-1.elb.amazonaws.com/app2/index.html

# App3 Curl Test
curl internal-ingress-internal-lb-1839544354.us-east-1.elb.amazonaws.com
```


## Step-07: Clean Up
```t
# Delete Manifests
kubectl delete -f kube-manifests/
kubectl delete -f kube-manifests-curl/
```








