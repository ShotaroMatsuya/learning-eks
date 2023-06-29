---
title: AWS Load Balancer Controller - NLB Basics
description: Learn to use AWS Network Load Balancer with AWS Load Balancer Controller
---

## Step-01: Introduction
- Understand more about 
  - **AWS Cloud Provider Load Balancer Controller (Legacy):** Creates AWS CLB and NLB
  - **AWS Load Balancer Controller (Latest):** Creates AWS ALB and NLB
- Understand how the Kubernetes Service of Type Load Balancer which can create AWS NLB to be associated with latest `AWS Load Balancer Controller`. 
- Understand various NLB Annotations


## Step-02: Review 01-Nginx-App3-Deployment.yml
- **File Name:** `kube-manifests/01-Nginx-App3-Deployment.yml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app3-nginx-deployment
  labels:
    app: app3-nginx 
spec:
  replicas: 1
  selector:
    matchLabels:
      app: app3-nginx
  template:
    metadata:
      labels:
        app: app3-nginx
    spec:
      containers:
        - name: app2-nginx
          image: stacksimplify/kubenginx:1.0.0
          ports:
            - containerPort: 80

```

## Step-03: Review 02-LBC-NLB-LoadBalancer-Service.yml
- **File Name:** `kube-manifests\02-LBC-NLB-LoadBalancer-Service.yml`
```yaml
apiVersion: v1
kind: Service
metadata:
  name: basics-lbc-network-lb
  annotations:
    # Traffic Routing
    service.beta.kubernetes.io/aws-load-balancer-name: basics-lbc-network-lb
    service.beta.kubernetes.io/aws-load-balancer-type: external
    service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: instance
    #service.beta.kubernetes.io/aws-load-balancer-subnets: subnet-xxxx, mySubnet ## Subnets are auto-discovered if this annotation is not specified, see Subnet Discovery for further details.
    
    # Health Check Settings
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-protocol: http
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-port: traffic-port
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-path: /index.html
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-healthy-threshold: "3"
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-unhealthy-threshold: "3"
    service.beta.kubernetes.io/aws-load-balancer-healthcheck-interval: "10" # The controller currently ignores the timeout configuration due to the limitations on the AWS NLB. The default timeout for TCP is 10s and HTTP is 6s.

    # Access Control
    service.beta.kubernetes.io/load-balancer-source-ranges: 0.0.0.0/0 
    service.beta.kubernetes.io/aws-load-balancer-scheme: "internet-facing"

    # AWS Resource Tags
    service.beta.kubernetes.io/aws-load-balancer-additional-resource-tags: Environment=dev,Team=test
spec:
  type: LoadBalancer
  selector:
    app: app3-nginx
  ports:
    - port: 80
      targetPort: 80
```

## Step-04: Deploy all kube-manifests
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Pods
kubectl get pods

# Verify Services
kubectl get svc
Observation: 
1. Verify the network lb DNS name

# Verify AWS Load Balancer Controller pod logs
kubectl -n kube-system get pods
kubectl -n kube-system logs -f <aws-load-balancer-controller-POD-NAME>

# Verify using AWS Mgmt Console
Go to Services -> EC2 -> Load Balancing -> Load Balancers
1. Verify Description Tab - DNS Name matching output of "kubectl get svc" External IP
2. Verify Listeners Tab

Go to Services -> EC2 -> Load Balancing -> Target Groups
1. Verify Registered targets
2. Verify Health Check path

# Access Application
http://<NLB-DNS-NAME>
```

## Step-05: Clean-Up
```t
# Delete or Undeploy kube-manifests
kubectl delete -f kube-manifests/

# Verify if NLB deleted 
In AWS Mgmt Console, 
Go to Services -> EC2 -> Load Balancing -> Load Balancers
```

## References
- [Network Load Balancer](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
- [NLB Service](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/nlb/)
- [NLB Service Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)

---

---
title: AWS Load Balancer Controller - NLB TLS
description: Learn to use AWS Network Load Balancer TLS with AWS Load Balancer Controller
---

## Step-01: Introduction
- Understand about the 4 TLS Annotations for Network Load Balancers
- aws-load-balancer-ssl-cert
- aws-load-balancer-ssl-ports
- aws-load-balancer-ssl-negotiation-policy
- aws-load-balancer-ssl-negotiation-policy

## Step-02: Review TLS Annotations
- **File Name:** `kube-manifests\02-LBC-NLB-LoadBalancer-Service.yml`
- **Security Policies:** https://docs.aws.amazon.com/elasticloadbalancing/latest/network/create-tls-listener.html#describe-ssl-policies
```yaml
    # TLS
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: arn:aws:acm:us-east-1:180789647333:certificate/d86de939-8ffd-410f-adce-0ce1f5be6e0d
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: 443, # Specify this annotation if you need both TLS and non-TLS listeners on the same load balancer
    service.beta.kubernetes.io/aws-load-balancer-ssl-negotiation-policy: ELBSecurityPolicy-TLS13-1-2-2021-06
    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: tcp 
```


## Step-03: Deploy all kube-manifests
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Pods
kubectl get pods

# Verify Services
kubectl get svc
Observation: 
1. Verify the network lb DNS name

# Verify AWS Load Balancer Controller pod logs
kubectl -n kube-system get pods
kubectl -n kube-system logs -f <aws-load-balancer-controller-POD-NAME>

# Verify using AWS Mgmt Console
Go to Services -> EC2 -> Load Balancing -> Load Balancers
1. Verify Description Tab - DNS Name matching output of "kubectl get svc" External IP
2. Verify Listeners Tab
Observation:  Should see two listeners Port 80 and 443

Go to Services -> EC2 -> Load Balancing -> Target Groups
1. Verify Registered targets
2. Verify Health Check path
Observation: Should see two target groups. 1 Target group for 1 listener

# Access Application
# Test HTTP URL
http://<NLB-DNS-NAME>
http://lbc-network-lb-tls-demo-a956479ba85953f8.elb.us-east-1.amazonaws.com

# Test HTTPS URL
https://<NLB-DNS-NAME>
https://lbc-network-lb-tls-demo-a956479ba85953f8.elb.us-east-1.amazonaws.com
```

## Step-04: Clean-Up
```t
# Delete or Undeploy kube-manifests
kubectl delete -f kube-manifests/

# Verify if NLB deleted 
In AWS Mgmt Console, 
Go to Services -> EC2 -> Load Balancing -> Load Balancers
```

## References
- [Network Load Balancer](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
- [NLB Service](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/nlb/)
- [NLB Service Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)


---

---
title: AWS Load Balancer Controller - NLB External DNS
description: Learn to use AWS Network Load Balancer & External DNS with AWS Load Balancer Controller
---

## Step-01: Introduction
- Implement External DNS Annotation in NLB Kubernetes Service Manifest


## Step-02: Review External DNS Annotations
- **File Name:** `kube-manifests\02-LBC-NLB-LoadBalancer-Service.yml`
```yaml
    # External DNS - For creating a Record Set in Route53
    external-dns.alpha.kubernetes.io/hostname: nlbdns101.stacksimplify.com
```

## Step-03: Deploy all kube-manifests
```t
# Verify if External DNS Pod exists and Running
kubectl get pods
Observation: 
external-dns pod should be running

# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Pods
kubectl get pods

# Verify Services
kubectl get svc
Observation: 
1. Verify the network lb DNS name

# Verify AWS Load Balancer Controller pod logs
kubectl -n kube-system get pods
kubectl -n kube-system logs -f <aws-load-balancer-controller-POD-NAME>

# Verify using AWS Mgmt Console
Go to Services -> EC2 -> Load Balancing -> Load Balancers
1. Verify Description Tab - DNS Name matching output of "kubectl get svc" External IP
2. Verify Listeners Tab
Observation:  Should see two listeners Port 80 and 443

Go to Services -> EC2 -> Load Balancing -> Target Groups
1. Verify Registered targets
2. Verify Health Check path
Observation: Should see two target groups. 1 Target group for 1 listener

# Verify External DNS logs
kubectl logs -f $(kubectl get po | egrep -o 'external-dns[A-Za-z0-9-]+')

# Perform nslookup Test
nslookup nlbdns101.stacksimplify.com

# Access Application
# Test HTTP URL
http://nlbdns101.stacksimplify.com

# Test HTTPS URL
https://nlbdns101.stacksimplify.com
```

## Step-04: Clean-Up
```t
# Delete or Undeploy kube-manifests
kubectl delete -f kube-manifests/

# Verify if NLB deleted 
In AWS Mgmt Console, 
Go to Services -> EC2 -> Load Balancing -> Load Balancers
```

## References
- [Network Load Balancer](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
- [NLB Service](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/nlb/)
- [NLB Service Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)


---

---
title: AWS Load Balancer Controller - NLB Elastic IP
description: Learn to use AWS Network Load Balancer & Elastic IP with AWS Load Balancer Controller
---

## Step-01: Introduction
- Create Elastic IPs
- Update NLB Service k8s manifest with Elastic IP Annotation with EIP Allocation IDs

## Step-02: Create two Elastic IPs and get EIP Allocation IDs
- This configuration is optional and use can use it to assign static IP addresses to your NLB
- You must specify the same number of eip allocations as load balancer subnets annotation
- NLB must be internet-facing
```t
# Elastic IP Allocation IDs
eipalloc-07daf60991cfd21f0 
eipalloc-0a8e8f70a6c735d16
```

## Step-03: Review Elastic IP Annotations
- **File Name:** `kube-manifests\02-LBC-NLB-LoadBalancer-Service.yml`
```yaml
    # Elastic IPs
    service.beta.kubernetes.io/aws-load-balancer-eip-allocations: eipalloc-07daf60991cfd21f0, eipalloc-0a8e8f70a6c735d16
```

## Step-04: Deploy all kube-manifests
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Pods
kubectl get pods

# Verify Services
kubectl get svc
Observation: 
1. Verify the network lb DNS name

# Verify AWS Load Balancer Controller pod logs
kubectl -n kube-system get pods
kubectl -n kube-system logs -f <aws-load-balancer-controller-POD-NAME>

# Verify using AWS Mgmt Console
Go to Services -> EC2 -> Load Balancing -> Load Balancers
1. Verify Description Tab - DNS Name matching output of "kubectl get svc" External IP
2. Verify Listeners Tab
Observation:  Should see two listeners Port 80 and 443

Go to Services -> EC2 -> Load Balancing -> Target Groups
1. Verify Registered targets
2. Verify Health Check path
Observation: Should see two target groups. 1 Target group for 1 listener

# Perform nslookup Test
nslookup nlbeip201.stacksimplify.com
Observation:
1. Verify the IP Address matches our Elastic IPs we created in Step-02

# Access Application
# Test HTTP URL
http://nlbeip201.stacksimplify.com

# Test HTTPS URL
https://nlbeip201.stacksimplify.com
```

## Step-05: Clean-Up
```t
# Delete or Undeploy kube-manifests
kubectl delete -f kube-manifests/

# Delete Elastic IPs created
In AWS Mgmt Console, 
Go to Services -> EC2 -> Network & Security -> Elastic IPs
Delete two EIPs created

# Verify if NLB deleted 
In AWS Mgmt Console, 
Go to Services -> EC2 -> Load Balancing -> Load Balancers
```

## References
- [Network Load Balancer](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
- [NLB Service](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/nlb/)
- [NLB Service Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)


---

---
title: AWS Load Balancer Controller - Internal NLB
description: Learn to create Internal AWS Network Load Balancer with Kubernetes
---

## Step-01: Introduction
- Create Internal NLB
- Update NLB Service k8s manifest with `aws-load-balancer-scheme` Annotation as `internal`
- Deploy curl pod
- Connect to curl pod and access Internal NLB endpoint using `curl command`.


## Step-02: Review LB Scheme Annotation
- **File Name:** `kube-manifests\02-LBC-NLB-LoadBalancer-Service.yml`
```yaml
    # Access Control
    service.beta.kubernetes.io/aws-load-balancer-scheme: "internal"
```

## Step-03: Deploy all kube-manifests
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Pods
kubectl get pods

# Verify Services
kubectl get svc
Observation: 
1. Verify the network lb DNS name

# Verify AWS Load Balancer Controller pod logs
kubectl -n kube-system get pods
kubectl -n kube-system logs -f <aws-load-balancer-controller-POD-NAME>

# Verify using AWS Mgmt Console
Go to Services -> EC2 -> Load Balancing -> Load Balancers
1. Verify Description Tab - DNS Name matching output of "kubectl get svc" External IP
2. Verify Listeners Tab
Observation:  Should see two listeners Port 80 

Go to Services -> EC2 -> Load Balancing -> Target Groups
1. Verify Registered targets
2. Verify Health Check path
```

## Step-04: Deploy curl pod and test Internal NLB
```t
# Deploy curl-pod
kubectl apply -f kube-manifests-curl

# Will open up a terminal session into the container
kubectl exec -it curl-pod -- sh

# We can now curl external addresses or internal services:
curl http://google.com/
curl <INTERNAL-NETWORK-LB-DNS>

# Internal Network LB Curl Test
curl lbc-network-lb-internal-demo-7031ade4ca457080.elb.us-east-1.amazonaws.com
```


## Step-05: Clean-Up
```t
# Delete or Undeploy kube-manifests
kubectl delete -f kube-manifests/
kubectl delete -f kube-manifests-curl/

# Verify if NLB deleted 
In AWS Mgmt Console, 
Go to Services -> EC2 -> Load Balancing -> Load Balancers
```

## References
- [Network Load Balancer](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
- [NLB Service](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/nlb/)
- [NLB Service Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)



---

---
title: AWS Load Balancer Controller - NLB & Fargate
description: Learn to use AWS Network Load Balancer with Fargate Pods
---

## Step-01: Introduction
- Create advanced AWS Fargate Profile
- Schedule App3 on Fargate Pod
- Update NLB Annotation `aws-load-balancer-nlb-target-type` with `ip` from `instance` mode

## Step-02: Review Fargate Profile
- **File Name:** `fargate-profile/01-fargate-profiles.yml`
```yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: eksdemo1  # Name of the EKS Cluster
  region: us-east-1
fargateProfiles:
  - name: fp-app3
    selectors:
      # All workloads in the "ns-app3" Kubernetes namespace will be
      # scheduled onto Fargate:      
      - namespace: ns-app3
```

## Step-03: Create Fargate Profile
```t
# Change Directory
cd 19-06-LBC-NLB-Fargate-External

# Create Fargate Profile
eksctl create fargateprofile -f fargate-profile/01-fargate-profiles.yml
```

## Step-04: Update Annotation aws-load-balancer-nlb-target-type to IP
- **File Name:** `kube-manifests/02-LBC-NLB-LoadBalancer-Service.yml`
```yaml
service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: ip # For Fargate Workloads we should use target-type as ip
```

## Step-05: Review the k8s Deployment Metadata for namespace
- **File Name:** `kube-manifests/01-Nginx-App3-Deployment.yml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app3-nginx-deployment
  labels:
    app: app3-nginx 
  namespace: ns-app3    # Update Namespace given in Fargate Profile 01-fargate-profiles.yml
spec:
  replicas: 1
  selector:
    matchLabels:
      app: app3-nginx
  template:
    metadata:
      labels:
        app: app3-nginx
    spec:
      containers:
        - name: app2-nginx
          image: stacksimplify/kubenginx:1.0.0
          ports:
            - containerPort: 80
          resources:
            requests:
              memory: "128Mi"
              cpu: "500m"
            limits:
              memory: "500Mi"
              cpu: "1000m"           
```

## Step-06: Deploy all kube-manifests
```t
# Deploy kube-manifests
kubectl apply -f kube-manifests/

# Verify Pods
kubectl get pods -o wide
Observation:
1. It will take couple of minutes to get the pod from pending to running state due to Fargate Mode.

# Verify Worker Nodes
kubectl get nodes -o wide
Obseravtion:
1. wait for Fargate worker node to create

# Verify Services
kubectl get svc
Observation: 
1. Verify the network lb DNS name

# Verify AWS Load Balancer Controller pod logs
kubectl -n kube-system get pods
kubectl -n kube-system logs -f <aws-load-balancer-controller-POD-NAME>

# Verify using AWS Mgmt Console
Go to Services -> EC2 -> Load Balancing -> Load Balancers
1. Verify Description Tab - DNS Name matching output of "kubectl get svc" External IP
2. Verify Listeners Tab

Go to Services -> EC2 -> Load Balancing -> Target Groups
1. Verify Registered targets
2. Verify Health Check path

# Perform nslookup Test
nslookup nlbfargate901.stacksimplify.com

# Access Application
# Test HTTP URL
http://nlbfargate901.stacksimplify.com

# Test HTTPS URL
https://nlbfargate901.stacksimplify.com
```

## Step-06: Clean-Up
```t
# Delete or Undeploy kube-manifests
kubectl delete -f kube-manifests/

# Verify if NLB deleted 
In AWS Mgmt Console, 
Go to Services -> EC2 -> Load Balancing -> Load Balancers
```

## References
- [Network Load Balancer](https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html)
- [NLB Service](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/nlb/)
- [NLB Service Annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/service/annotations/)











## Step-09: Delete Fargate Profile
```t
# List Fargate Profiles
eksctl get fargateprofile --cluster eksdemo1 

# Delete Fargate Profile
eksctl delete fargateprofile --cluster eksdemo1 --name <Fargate-Profile-NAME> --wait

eksctl delete fargateprofile --cluster eksdemo1 --name  fp-app3 --wait
```