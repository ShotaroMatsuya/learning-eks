# Install Grafana

- Create Promethus  
  [here](./prometeus.md)
- CloudWatch Container Insights with FluentD  
  [here](./cloudwatch.md)

## Install Grafana

```bash
kubectl create namespace grafana
helm repo add grafana https://grafana.github.io/helm-charts
helm install grafana grafana/grafana \
--namespace grafana \
--set persistence.storageClassName="gp2" \
--set persistence.enabled=true \
--set adminPassword='EKS!sAWSome' \
--set datasources."datasources\.yaml".apiVersion=1 \
--set datasources."datasources\.yaml".datasources[0].name=Prometheus \
--set datasources."datasources\.yaml".datasources[0].type=prometheus \
--set datasources."datasources\.yaml".datasources[0].url=http://prometheus-server.prometheus.svc.cluster.local \
--set datasources."datasources\.yaml".datasources[0].access=proxy \
--set datasources."datasources\.yaml".datasources[0].isDefault=true \
--set service.type=LoadBalancer
```

## Check if Grafana is deployed properly

```bash
kubectl get all -n grafana
```

## Get Grafana ELB url

```bash
   export ELB=$(kubectl get svc -n grafana grafana -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
```

```bash
echo "http://$ELB"
```

## When logging in, use username "admin" and get password by running the following:

```bash
   kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

test deploy

```bash
kubectl apply -f nginx-deployment-withrolling.yaml
```

## Grafana Dashboards for K8s:

https://grafana.com/grafana/dashboards?dataSource=prometheus&direction=desc&orderBy=reviewsCount

## Uninstall Prometheus and Grafana

```bash
   helm uninstall prometheus --namespace prometheus
   helm uninstall grafana --namespace grafana
```
