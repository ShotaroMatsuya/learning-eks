# CloudWatch Container Insights with FluentD

- Install Grafana  
  [here](./README.md)
- Create Promethus  
  [here](./prometeus.md)

Fluent Bit をインストールしてコンテナから CloudWatch Logs にログを送信するには

amazon-cloudwatch という名前の名前空間がまだない場合は、次のコマンドを入力して作成します。

```bash
kubectl apply -f https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/cloudwatch-namespace.yaml
```

---

次のコマンドを実行して、クラスター名とログを送信するリージョンを持つ cluster-info という名前の ConfigMap を作成します。cluster-name と cluster-region をクラスターの名前とリージョンに置き換えます。

```bash
ClusterName=eks-Prometheus
RegionName=ap-northeast-1
FluentBitHttpPort='2020'
FluentBitReadFromHead='Off'
[[${FluentBitReadFromHead} = 'On']] && FluentBitReadFromTail='Off'|| FluentBitReadFromTail='On'
[[-z ${FluentBitHttpPort}]] && FluentBitHttpServer='Off' || FluentBitHttpServer='On'
kubectl create configmap fluent-bit-cluster-info \
--from-literal=cluster.name=${ClusterName} \
--from-literal=http.server=${FluentBitHttpServer} \
--from-literal=http.port=${FluentBitHttpPort} \
--from-literal=read.head=${FluentBitReadFromHead} \
--from-literal=read.tail=${FluentBitReadFromTail} \
--from-literal=logs.region=${RegionName} -n amazon-cloudwatch
```

このコマンドでは、プラグインメトリクスをモニタリングするための FluentBitHttpServer がデフォルトでオンになっています。無効にするには、コマンドの 3 行目を `FluentBitHttpPort='' `(空の文字列) に変更します。

また、デフォルトでは、Fluent Bit はテールからログファイルを読み取り、デプロイ後に新しいログのみを取得します。逆をご希望の場合は、`FluentBitReadFromHead='On'` を設定することで、ファイルシステム内のすべてのログが収集されます。

---

次のいずれかのコマンドを実行して、Fluent Bit daemonset をクラスターにダウンロードしてデプロイします。

Fluent Bit 最適化設定が必要な場合は、このコマンドを実行します。

```bash
kubectl apply -f https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml
```

---

次のコマンドを実行してデプロイを検証します。各ノードには、fluent-bit-\* という名前の 1 つのポッドが必要です。

```bash
kubectl get pods -n amazon-cloudwatch
```

上記の手順を実行することにより、クラスターに次のリソースが作成されます。

Fluent-Bit 名前空間の amazon-cloudwatch という名前のサービスアカウント。このサービスアカウントは、Fluent Bit daemonSet を実行するために使用されます  
Fluent-Bit-role 名前空間の amazon-cloudwatch という名前のクラスターロール。このクラスターロールは、ポッドログの get、list、watch の各アクセス許可を Fluent-Bit サービスアカウントに付与します。  
Fluent-Bit-config 名前空間の amazon-cloudwatch という名前の ConfigMap。この ConfigMap には、Fluent Bit によって使用される設定が含まれています。

---

クイックスタートを使用して Container Insights をデプロイするには、次のコマンドを入力します。

```bash
ClusterName=eks-Prometheus
RegionName=ap-northeast-1
FluentBitHttpPort='2020'
FluentBitReadFromHead='Off'
[[${FluentBitReadFromHead} = 'On']] && FluentBitReadFromTail='Off'|| FluentBitReadFromTail='On'
[[-z ${FluentBitHttpPort}]] && FluentBitHttpServer='Off' || FluentBitHttpServer='On'
curl https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluent-bit-quickstart.yaml | sed 's/{{cluster_name}}/'${ClusterName}'/;s/{{region_name}}/'${RegionName}'/;s/{{http_server_toggle}}/"'${FluentBitHttpServer}'"/;s/{{http_server_port}}/"'${FluentBitHttpPort}'"/;s/{{read_from_head}}/"'${FluentBitReadFromHead}'"/;s/{{read_from_tail}}/"'${FluentBitReadFromTail}'"/' | kubectl apply -f -
```

このコマンドでは、my-cluster-name は Amazon EKS または Kubernetes クラスターの名前で、my-cluster-region はログが発行されるリージョンの名前です。AWS アウトバウンドデータ転送コストを削減するために、クラスターがデプロイされているのと同じリージョンを使用することをお勧めします。

---

クイックスタートセットアップの使用後に Container Insights を削除する場合は、次のコマンドを入力します。

```bash
curl https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluent-bit-quickstart.yaml | sed 's/{{cluster_name}}/'${ClusterName}'/;s/{{region_name}}/'${LogRegion}'/;s/{{http_server_toggle}}/"'${FluentBitHttpServer}'"/;s/{{http_server_port}}/"'${FluentBitHttpPort}'"/;s/{{read_from_head}}/"'${FluentBitReadFromHead}'"/;s/{{read_from_tail}}/"'${FluentBitReadFromTail}'"/' | kubectl delete -f -
```
