# 検証

1. Automatic Node Provisioning

```bash
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
name: inflate
spec:
replicas: 0
selector:
matchLabels:
app: inflate
template:
metadata:
labels:
app: inflate
spec:
terminationGracePeriodSeconds: 0
containers: - name: inflate
image: public.ecr.aws/eks-distro/kubernetes/pause:3.2
resources:
requests:
cpu: 1
EOF
```

```bash
kubectl scale deployment inflate --replicas 5
kubectl logs -f -n karpenter $(kubectl get pods -n karpenter -l karpenter=controller -o name)
```

2. Automatic Node Termination

```bash
kubectl delete deployment inflate
kubectl logs -f -n karpenter $(kubectl get pods -n karpenter -l karpenter=controller -o name)
```

3. Manual Node Termination

   ```bash
   kubectl delete node $NODE_NAME
   ```

4. Cleanup

```bash
helm uninstall karpenter --namespace karpenter
eksctl delete iamserviceaccount --cluster ${CLUSTER_NAME} --name karpenter --namespace karpenter
aws cloudformation delete-stack --stack-name Karpenter-${CLUSTER_NAME}
aws ec2 describe-launch-templates \
 | jq -r ".LaunchTemplates[].LaunchTemplateName" \
 | grep -i Karpenter-${CLUSTER_NAME} \
 | xargs -I{} aws ec2 delete-launch-template --launch-template-name {}
eksctl delete cluster --name ${CLUSTER_NAME}
```
