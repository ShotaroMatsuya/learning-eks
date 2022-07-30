# Kubernetes hpa demo

[Documentation](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)

## Increase Load

Now, we will see how the auto-scaler reacts to increased load. We will start a container, and send an infinite loop of queries to the php-apache service (please run it in a different terminal)

```bash
# Run this in a separate terminal
# so that the load generation continues and you can carry on with the rest of the steps
kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done"

```

Within a minute or so, we should see the higher CPU load by executing:

```bash
# type Ctrl+C to end the watch when you're ready
kubectl get hpa php-apache --watch
```

---

# EKS Cluster Auto-scaler Guide

[Documentation](https://docs.aws.amazon.com/eks/latest/userguide/autoscaling.html#cluster-autoscaler)

Cluster Auto-scaler has two components

- Open source Cluster Auto-scaler
- EKS Implementation (ASG, IAM etc.)

## Create an IAM policy

1. Save the following contents to a file that's named cluster-autoscaler-policy.json.  
   If your existing node groups were created with eksctl and you used the --asg-access option, then this policy already exists and you can skip to step 2.

```bash
eksctl create cluster --name my-cluster --managed --asg-access
```

Create the policy with the following command. You can change the value for policy-name.

```bash
aws iam create-policy \
    --policy-name AmazonEKSClusterAutoscalerPolicy \
    --policy-document file://cluster-autoscaler-policy.json
```

2. You can create an IAM role and attach an IAM policy to it using eksctl.

Run the following command if you created your Amazon EKS cluster with eksctl. If you created your node groups using the `--asg-access` option, then replace `AmazonEKSClusterAutoscalerPolicy` with the name of the IAM policy that eksctl created for you. The policy name is similar to `eksctl-my-cluster-nodegroup-ng-xxxxxxxx-PolicyAutoScaling`.

```bash
eksctl create iamserviceaccount \
  --cluster=eks-Prometheus \
  --namespace=kube-system \
  --name=cluster-autoscaler \
  --attach-policy-arn=arn:aws:iam::528163014577:policy/AmazonEKSClusterAutoscalerPolicy \
  --override-existing-serviceaccounts \
  --approve
```

## Deploy the Cluster Autoscaler

1. Apply the Cluster Autoscaler YAML file(cluster-autoscaler-autodiscover.yaml) to your cluster.

```bash
kubectl apply -f cluster-autoscaler-autodiscover.yaml
```

2. Patch the deployment to add the cluster-autoscaler.kubernetes.io/safe-to-evict annotation to the Cluster Autoscaler pods with the following command.

```bash
kubectl patch deployment cluster-autoscaler \
  -n kube-system \
  -p '{"spec":{"template":{"metadata":{"annotations":{"cluster-autoscaler.kubernetes.io/safe-to-evict": "false"}}}}}'
```

3. Edit the Cluster Autoscaler deployment with the following command.

```bash
kubectl -n kube-system edit deployment.apps/cluster-autoscaler
```

Edit the cluster-autoscaler container command to add the following options. --balance-similar-node-groups ensures that there is enough available compute across all availability zones. --skip-nodes-with-system-pods=false ensures that there are no problems with scaling to zero.

```yaml
spec:
  containers:
    - command
      - ./cluster-autoscaler
      - --v=4
      - --stderrthreshold=info
      - --cloud-provider=aws
      - --skip-nodes-with-local-storage=false
      - --expander=least-waste
      - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/my-cluster
      - --balance-similar-node-groups
      - --skip-nodes-with-system-pods=false
```

4. Open the Cluster Autoscaler releases page from GitHub in a web browser and find the latest Cluster Autoscaler version that matches the Kubernetes major and minor version of your cluster.

```bash
kubectl set image deployment cluster-autoscaler \
  -n kube-system \
  cluster-autoscaler=k8s.gcr.io/autoscaling/cluster-autoscaler:v1.22.3
```

## View your Cluster Autoscaler logs

1. After you have deployed the Cluster Autoscaler, you can view the logs and verify that it's monitoring your cluster load.

```bash
kubectl -n kube-system logs -f deployment.apps/cluster-autoscaler
```

## apply cluster-autoscaler-deployment-1.yaml
