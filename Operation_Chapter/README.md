# K8s のスケーリング

K8s では Node と Pod の 2 つの s けーリングを考える必要がある  
Pod のスケーリングは K8s のネイティブの機能として　存在する HPA と VPA を利用
Cluster AutoScaler は K8s の機能で、各クラウドプロバイダの NodeGroup に最適化された状態で提供される

```bash
# ハンズオン始め方
eksctl create cluster -f eksctl_create_cluster.yaml

# ハンズオンをやめる場合
eksctl delete cluster -f eksctl_create_cluster.yaml

```

## Resource Request/Limits と QoS クラス

- Resource Requests/Limits
  Resource Requests とは、Pod 内のコンテナごとに CPU や Memory などのリソースの要求を行う設定で、Pod がスケジュールされる際に Pod 内のコンテナの Requests 値の合計を Node のリソースから予約  
  もし Node で利用できるリソース以上の Requests が要求されると、Pod はスケジューリングされずに Pending となる  
  実際の使用率がそれほど使われてなくても Requests によって予約されたリソースキャパシティがなくなると、Pod をスケジューリングできなくなる

  Resource Limits は、Pod 内のコンテナごとに CPU や Memory などのリソースの制限を行う設定で、コンテナランタイムは Pod 内のコンテナに設定された Limits を超えないように制限を行う.  
  もし Limits の設定値を超えた場合、Memory Limits のばあいは OOM Killer により Pod(プロセス)が Kill される  
  もし、CPU Limits が超過したばあいは、Kill されることはないが、CPU スロットリングと呼ばれる制限をかけられるので、パフォーマンスが著しく低下する

  **cpu**には Node の CPU のコア数を指定する. 1 コアの場合は 1 もしくは 1000m と指定  
  **memory**はメモリーサイズをバイト単位で指定する.1GB を指定するばあいは、1024M と指定することができるし、*1Gi*と指定することもできる **Gi**はギビバイトと呼ばれるコンピュータが理解しやすい単位で、1GiB ＝ 1024M バイトで計算される  
  また、Node のリソースがどれくらい予約されているか、`kubectl describe nodes <Node名>`コマンドの**Allocated resources**フィールドで確認できる

- QoS クラス
  QoS クラスは Pod へ自動的に付与される、Pod の優先度を定義するもの. QoS クラスには下表の 3 つのクラスが存在し、Resource Requests/Limits の設定内容により自動的にクラスが決定される.  
  QoS クラスは Node のリソースが枯渇した際に、OOM Killer により Kill される Pod を選定する際の優先順位などに利用.

| クラス名   | 優先度 | 詳細                                                                                                                              |
| ---------- | ------ | --------------------------------------------------------------------------------------------------------------------------------- |
| Guaranteed | 高     | Pod 内のすべてのコンテナにメモリーかつ CPU の制限と要求が与えられており、同じ値であること                                         |
| Burstable  | 中     | Pod が Guaranteed QoS クラスの基準に満たない場合、かつ Pod 内の 1 つ以上のコンテナがメモリーまたは CPU の要求を与えられている場合 |
| BestEffort | 低     | Pod 内のすべてのコンテナにメモリーおよび CPU の制限と要求が 1 つも指定されていない                                                |

## HPA

対象の Pod の平均 CPU 使用率、または独自に用意したカスタムメトリクスのいずれかメトリクスをターゲットとして追跡し、ターゲットのメトリクス値が指定された値を上回らないように Pod 数をコントロールする  
HPA がターゲットの追跡に利用する CPU 使用率やメモリ使用率などのメトリクスは Metrics Server から取得するが、EKS のばあいは Metrics Server がデフォルトでインストールされていないため、Deployment として別途デプロイ氏、kubelet から収集された Pod のメトリクスを api-server 経由で Metrics Server から取得する  
Metrics Server がない状態で`kubectl top`コマンドを実行していみるとメトリクスが取得できない

```bash
# Metrics ServerをClusterへインストールする
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
# 確認
kubectl get deployment metrics-server -n kube-system
# Nodeのメトリクスを取得
kubectl top node

# HPAをデプロイ
## loadかけるdeploymentを追加
kubectl apply -f https://k8s.io/examples/application/php-apache.yaml
kubectl get pods
## 次にphp-apache Deploymentへ、HPAの設定を行う
kubectl apply -f hpa-php-apache.yaml

## 確認
kubectl get hpa
NAME         REFERENCE               TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
php-apache   Deployment/php-apache   <unknown>/50%   1         10        0          4s
## 負荷をかける
kubectl run -i --tty load-generator --rm --image=busybox --restart=Never \
    -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done"
## 監視
kubectl get hpa -w
```

```bash
# Clean Up
kubectl delete hpa php-apache
kubectl delete deployment php-apache
kubectl delete service php-apache
kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

## VPA

Pod の CPU とメモリの Resource Requests/Limits を調整し、垂直スケーリングさせることでスケールアップ・ダウンを行う.  
また、VPA はスケール前に設定されていた Resource Requests/Limits の比率を維持する

\*VPA と HPA の CPU 使用率もしくはメモリ使用率をターゲットとして併用しないように警告されている

## Cluster Autoscaler

Node のリソースが枯渇していて Pod が起動できないときや、Cluster 内の特定の Node のリソース使用率が低く Pod を別の Node に再スケジュールできる場合に、Cluster 内のノード数を自動的に調整する.  
Cluster Autoscaler 自体は K8s の機能だが、各クラウドプロバイダごとに最適化されていて、AWS では Cluster Autoscaler が EKS node Group(EC2 Autoscaling Group)を直接調整することで、Cluster 内の Node 数を増減させる　　
Cluster Autoscaler はデフォルトで使用できないので、Deployment として Cluster へインストールする必要がある  
Cluster Autoscaler をデプロイするには、EKSCluster 内に OIDC プロバイダが存在している必要と、対象の EC2 AutoScaling Group に下表のタグが付与されている必要がある

| Key                                       | Value |
| ----------------------------------------- | ----- |
| k8s.io/cluster-autoscaler/${CLUSTER_NAME} | owned |
| k8s.io/cluster-autoscaler/enabled         | TRUE  |

```bash
# Cluster Autoscalerをインストールする手順
## Node GroupのARNを取得
aws eks describe-nodegroup --cluster-name eks-cluster --nodegroup-name \
    $(aws eks list-nodegroups --cluster-name eks-cluster \
        --query 'nodegroups[]' --output text) \
    --query 'nodegroup.nodegroupArn'

## まずNode GroupへCluster Autoscalerが使用するタグを付与
aws eks tag-resource --resource-arn arn:aws:eks:ap-northeast-1:528163014577:nodegroup/eks-cluster/bottlerocket/68c47331-2f94-3079-d61e-7e9b11b20180 --tags "k8s.io/cluster-autoscaler/eks-cluster=owned,k8s.io/cluster-autoscaler/enabled=TRUE"

# タグの確認
aws eks describe-nodegroup --cluster-name eks-cluster --nodegroup-name \
    $(aws eks list-nodegroups --cluster-name eks-cluster \
        --query 'nodegroups[]' --output text) \
    --query 'nodegroup.tags'

## NodeGroupのスケーリング設定の最大数を3へ変更
aws eks update-nodegroup-config --cluster-name eks-cluster \
    --nodegroup-name $(aws eks list-nodegroups --cluster-name eks-cluster \
        --query 'nodegroups[]' --output text) \
    --scaling-config "minSize=1,maxSize=3,desiredSize=1"

aws eks describe-nodegroup --cluster-name eks-cluster --nodegroup-name \
    $(aws eks list-nodegroups --cluster-name eks-cluster \
        --query 'nodegroups[]' --output text) \
    --query 'nodegroup.status'
"ACTIVE"

## Cluster AutoscalerのService Accountにマッピングする, IAM Policy とService Accountを作成
aws iam create-policy \
    --policy-name AmazonEKSClusterAutoscalerPolicy \
    --policy-document \
'{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "autoscaling:DescribeAutoScalingGroups",
        "autoscaling:DescribeAutoScalingInstances",
        "autoscaling:DescribeLaunchConfigurations",
        "autoscaling:DescribeTags",
        "autoscaling:SetDesiredCapacity",
        "autoscaling:TerminateInstanceInAutoScalingGroup",
        "ec2:DescribeLaunchTemplateVersions"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}'

eksctl create iamserviceaccount \
    --cluster=eks-cluster \
    --namespace=kube-system \
    --name=cluster-autoscaler \
    --attach-policy-arn="arn:aws:iam::528163014577:policy/AmazonEKSClusterAutoscalerPolicy" \
    --approve

## Cluster AutoScalerのマニフェストをダウンロードし、Deploymentリソースへ変更を適用する
curl -o cluster-autoscaler-autodiscover.yaml https://raw.githubusercontent.com/kubernetes/autoscaler/master/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml

## spec.template.metadata.annotations → cluster-autoscaler.kubernetes.io/safe-to-evict: 'false'
## spec.template.spec.containers.image → k8s.gcr.io/autoscaling/cluster-autoscaler: v1.22.2
## spec.template.spec.containers.command の --node-group-auto-discoveryオプションの末尾をEKS Cluster名に変更
## spec.template.spec.containers.commandに「--balance-similar-node-groups」「--skip-nodes-with-system-pods=false」を付与

# 修正したらapply
kubectl apply -f cluster-autoscaler-autodiscover.yaml

# Cluster Autoscalerがデプロイされてログが出力されていることを確認
kubectl get deployment -n kube-system cluster-autoscaler
kubectl -n kube-system logs deployment.apps/cluster-autoscaler | tail

# 実際にNodeのリソースを枯渇させてNodeをスケールアウト・インしてみる
kubectl apply -f deployment.yaml
# Nodeのリソースが枯渇しているため1つしかpodが作れない
kubectl get pods
nginx-deployment-54f758864f-4bfdw   0/1     Pending   0          48s
nginx-deployment-54f758864f-8k8hk   0/1     Pending   0          48s
nginx-deployment-54f758864f-gqm99   0/1     Pending   0          48s
nginx-deployment-54f758864f-ktb9x   1/1     Running   0          48s
nginx-deployment-54f758864f-ngm99   0/1     Pending   0          48s
# Nodeがmaxの3台スケーリングされている
kubectl get nodes -w
NAME                                                STATUS   ROLES    AGE     VERSION
ip-192-168-59-109.ap-northeast-1.compute.internal   Ready    <none>   97m     v1.22.15-eks-fb459a0
ip-192-168-64-195.ap-northeast-1.compute.internal   Ready    <none>   2m25s   v1.22.15-eks-fb459a0
ip-192-168-9-212.ap-northeast-1.compute.internal    Ready    <none>   2m21s   v1.22.15-eks-fb459a0
# 再度podのステータスを確認すると2台のpodが追加が起動するようになっている
kubectl get pods -w
nginx-deployment-54f758864f-4bfdw   0/1     Pending   0          5m2s
nginx-deployment-54f758864f-8k8hk   1/1     Running   0          5m2s
nginx-deployment-54f758864f-gqm99   1/1     Running   0          5m2s
nginx-deployment-54f758864f-ktb9x   1/1     Running   0          5m2s
nginx-deployment-54f758864f-ngm99   0/1     Pending   0          5m2s


# Nodeのスケールインの挙動を確認
kubectl delete -f deployment.yaml

kubectl get nodes
NAME                                                STATUS   ROLES    AGE    VERSION
ip-192-168-59-109.ap-northeast-1.compute.internal   Ready    <none>   114m   v1.22.15-eks-fb459a0
```

```bash
# Clean up
kubectl delete -f cluster-autoscaler-autodiscover.yaml
eksctl delete iamserviceaccount \
    --cluster=eks-cluster \
    --namespace=kube-system \
    --name=cluster-autoscaler

aws iam delete-policy --policy-arn arn:aws:iam::528163014577:policy/AmazonEKSClusterAutoscalerPolicy

aws eks update-nodegroup-config --cluster-name eks-cluster \
    --nodegroup-name $(aws eks list-nodegroups --cluster-name eks-cluster \
        --query 'nodegroups[]' --output text) \
    --scaling-config "minSize=1,maxSize=1,desiredSize=1"

```

# レジリエンスの維持

## Pod の Priority

Priority は Pod がスケジューリングされる際や Node から eviction(退避)させられるばあいに考慮され、Node 上に何らかの理由でスケジューリングできない Pod がある場合は、その Pod より Priority の低い Pod が eviction される  
Pod へ Priority を設定するには PriorityClass を作成する必要がある  
value に 10 億以下の任意の 32 ビットの整数値を用いて Priority 値を指定.  
Cluster 内で**globalDefault**は 1 つの PriorityClass だけが持つことができるが、**globalDefault**が一つもなく Pod でも PriorityClass が指定 s れていないときは Priority 値が 0 となる

同様の概念として QoS クラスがある  
Priority と QoS はどちらも Pod の優先度を定義するものだが、次の相違点がある

- QoS クラスはスケジューリングの優先順位に影響を与えない
- Node のリソース不足による evition が実行されるばあいのみ QoS クラスが考慮される

  - Pod の Resource Requests が Node のキャパを超えているばあいに**QoS の優先度順**で eviction される(_ここでは Priority は考慮されない_)
  - 実際の使用量が Requests を超えていない場合は**Priority の優先度順**に eviction される
  - eviction 対象の Priority が同じ場合は Resource Requests に対する使用率の相対値を用いて eviction される Pod が決定される

\*重要度の高い Pod がより湯煎的にスケジューリングされるためにも、Pod への Priority の設定が必要となる

## Pod の Anti-Affinity

Pod Anti-Affinity でクラウドプロバイダの ZOne を指定することで Pod の Replica をゾーン間に均等に配置することができる.  
EKS の場合も Node へ**topology.kubernetes.io/zone=ap-northeast-1c**などのトポロジキーが自動で付与される

## Probe によるコンテナの正常性確認

Pod 内のコンテナの正常性を確認する機能として、Liveness Probe(コンテナの死活監視), Readiness Probe(コンテナがサービスを提供できる状態にあるかのチェック), Startup Probe(起動に時間のかかるコンテナの起動が完了したかのチェック)の 3 つが存在する.  
これらの 3 つの Probe はコマンドの実行、HTTP リクエスト、TCP コネクションの確立のいずれかを使用してコンテナの正常性を確認できる

### コマンド実行

**spec.containers.livenessProbe**の**exec.command**で実行させない Linux コマンドを指定する. このコマンドの実行結果が 0(正常終了)でリターンされると、Liveness Probe はチェック結果を成功とする.  
**spec.containers.livenessProbe**の**initialDelaySeconds**でコンテナが起動してから初回のチェックが行われるまでの猶予時間を**periodSeconds**で Liveness Probe の実行間隔を指定する

### HTTP リクエスト

**spec.containers.livenessProbe**の**httpGet.path**で HTTP リクエストを行う URL パスを、**httpGet.port**で対象のポートを指定する.  
HTTP ステータスコードが 200 から 400 未満で返却されるとチェックが成功と判断され、それ以外のステータスが返却されると＾ LIveness Probe はコンテナを**再起動**する.

### TCP コネクション

**spec.containers.livenessProbe**の**tcpSocket.port**で TCP コネクションの確立を行う Port を指定.  
TCP コネクションが確立できた場合はチェックが成功と判断され、コネクションが確立できないばあいはコンテナが**再起動**される.

- Readiness Probe
  成功したばあいは、K8s Service からのリクエストを割り振られるが、失敗した場合は、そのコンテナはまだリクエストを処理する準備が整っていないと判断され、Service からのリクエストが割り振られない.しかし、コンテナが**再起動することはない**.

- Startup Probe
  起動に時間のかかるコンテナの起動が完了したかをチェックする用途で利用される  
  Liveness Probe と同じチェック方法、設定がなされるのが一般的  
  Startup Probe でコンテナが起動したかのチェックし、コンテ起動後は Liveness Probe へその確認を引き継ぐ形となる.

## Pod のターミネーション

実行中の Pod が突然強制終了されると、その Pod で処理していたリクエストも強制終了され、Web アプリケーションの場合はクライアントに 500 系のエラーが帰ってしまう.K8s には Pod を安全に終了させる機能がいくつか用意されている

- preStop フック
  preStop フックは exec によるコマンドなどを実行し、コンテナを安全に終了させるための処理を実行.  
  exec のほかにも tcpSocket や httpGet を指定でき、kubelet により Pod の終了プロセスが開始されたタイミング実行される
- `terminationGracePeriodSeconds`
  Pod が終了されるまでの猶予期間の定義を行う.デフォルトは 30s  
  pod は`kubectl delete pod`コマンドなどにより Preemption(終了)されると、次の流れで Pod を終了させる
  1. `kubectl delete pod`などにより Pod を終了すると、Pod のステータスが Terminating に変更され、猶予期間（terminationGracePeriodSeconds）が Pod へ設定される
  2. 猶予期間が設定されたことを kubelet が検知し、Pod の終了プロセスを開始する.
  3. prestop フックが設定されているコンテナが存在する場合は「preStop」で定義された処理が実行される
  4. 2 と同時にコントロールプレーンが終了中の＾ Pod を Endpoints から削除する(終了中のコンテナへの Service からのリクエストが止まる)
  5. コンテナから preStop フックの正常終了を受け取った場合は「SIGTERM」シグナルが送られ、コンテナが正常終了される.
  6. 猶予期間が経過しても終了されていないコンテナが存在する場合は、1 度だけ猶予期間を 2 秒延長する
  7. 延長後の猶予期間内に終了しないコンテナへ「SIGKILL」シグナルが送られ、Pod を強制終了

## PodDisruptionBudget

PodDisruptionBudget(PDB)は Pod が最小限起動しておく必要のある数と、一度に停止できる数の最大数を指定することで、`kubectl drain`コマンドなどによる Eviction API 経由で Pod が eviction される際に、サービスの維持に必要な稼働数を維持した状態で安全に Pod を退避することができる.

注意事項として、Deployment や replicaset などによる、ローリングアップデート時は PDB の設定が考慮されない  
また、PriorityClass による Pod eviction が実行される際も、場合によっては PDB が考慮されない

## レジリエンスの高い Cluster 構成まとめ

- QoS クラスと PriorityClass によって、より優先度の高い Pod が優先的に Node にスケジューリングされる
- ゾーンをトポロジキーにした Pod の Anti-Affinity によって、同一 Pod がゾーン間で均等にスケジューリングされる.
- Node のリソース不足が発生したときには、優先度の低い Pod から eviction される
- 万が一 Pod 内で障害が発生しても、Liveness Probe により自動で再起動される.
- Readiness Probe により Service からのリクエストを処理できる状態になったタイミングでリクエストを受け取ることができる
- Pod が終了される場合も preStop フックと PodDisruptionBudget により、サービスの維持に必要な Pod 数を確保した状態で Pod を安全に終了させることができる.
- HPA や Cluster Autoscaler を利用することで、急激な負荷増加や特定のゾーン障害が発生した際のリソース不足などにも自動で解決することができる

# K8s のモニタリング

## メトリクスの収集と可視化

Prometheus は exporter と呼ばれるメトリクス収集エージェントに対して、サーバー側から取得を行う Pull 型のアーキテクチャを採用しており、KeyValue の多次元データモデルとして時系列データで保存するので、コンテナ環境下で必要とされる複雑なメトリクス収集にも対応できる  
また、PromQL でメトリクスを直感的にクエリすることもできる  
アラート設定やメトリクスの可視化を行う際には、PromQL でかんたんにこれらの設定を行うことができる  
Prometheus を EKS Cluster へ直接インストールして利用することもできるが、Prometheus はメトリクスを永続保存しないため、Thanos や Cortex などの別サービスを利用して永続化させる必要がある。さらにクラスタ管理にかかるコストが肥大化してしまう　　
これらの課題を解決すべく、Amazon Managed Service for Prometheus(AMP)がリリースされた  
AMP は既存の Prometheus と互換性を持ったまま、Workload のスケールアップ及びスケールダウンにおうじて、運用メトリクスの取り込み、ストレージ、クエリを自動的にスケーリングする機能をマネージドサービスとして提供  
また、AMP はマルチ AZ としてデプロイされ、取り込まれたメトリクスデータなどはリージョン内の 3 つの AZ へレプリケートされるので、高い可用性も備えている. AMP のワークスペースに取り込まれたメトリクスは 150 日間保存され、その後に自動的に削除される.  
AMP では Prometheus のダッシュボードへのアクセスを IAM を用いて認証することができるので、すでに AWS を利用している場合は手軽に始めることができる

```bash
# EKS Cluster内のメトリクスをAMPへ取り込む手順
aws amp create-workspace --alias eks-cluster --region ap-northeast-1
## 次にPrometheus ZServerがAMPへメトリクスを書き込む際のIRSAを作成
sh ./create_irsa.sh

## 次にHelmを使って、EKS ClusterへPrometheus Serverをインストール
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

kubectl create ns prometheus

export SERVICE_ACCOUNT_IAM_ROLE=EKS-AMP-ServiceAccount-Role
export SERVICE_ACCOUNT_IAM_ROLE_ARN=$(aws iam get-role \
    --role-name $SERVICE_ACCOUNT_IAM_ROLE --query 'Role.Arn' --output text)

WORKSPACE_ID=$(aws amp list-workspaces --alias eks-cluster \
    --query 'workspaces[].workspaceId' --output text)

helm install prometheus-for-amp prometheus-community/prometheus -n prometheus -f ./amp_ingest_override_values.yaml \
    --set serviceAccounts.server.annotations."eks\.amazonaws\.com/role-arn"="${SERVICE_ACCOUNT_IAM_ROLE_ARN}" \
    --set 'server.remoteWrite[0].url'="https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/${WORKSPACE_ID}/api/v1/remote_write" \
    --set 'server.remoteWrite[0].sigv4.region'=ap-northeast-1

# インストールされて起動していることを確認
kubectl get all -n prometheus

```

## メトリクスの可視化

Prometheus は Grafana と呼ばれる OSS の分析、可視化サービスとともによく使われる  
AWS は AMP と合わせて Grafana のマネージドサービスである、Amazon Managed Service for Grafana(AMG)を提供  
AMG は AMP や CloudWatch、Opensearch などをデータソースとして統合することができ、また、収集した PRometheus メトリクスを PromQL でクエリして、ダッシュボードやアラートを作成できる
AMG を利用するに当たり、AWS アカウントで AWS Organizations と Single Sign on(SSO)が有効になっている必要がある

## メトリクスのアラート設定

Prometheus メトリクスに対してアラートを設定し、特定の条件を満たした場合に通知を行う方法  
なお、AMP や AMG は AWS SNS と統合されているので、SNS 経由で Slack や pagerDuty などに対してアラートを通知することが可能.  
ここでは AMP で Prometheus のアラートルールへ　アラート設定を行う.  
Alertmanager によってアラートルールの通知などが行われる.  
るーるは　 2 種類存在し、アラートルールの他にも記憶ルールと呼ばれるルールが存在する. 記憶ルールは頻繁に実行される場合うや計算量が多いクエリを事前に計算し、別のメトリクスとして　保存しておくことができる.  
AMP へアラート設定を行う手順.  
なお、ここでは Pod 内のコンテナで OOM Killer によるコンテナプロセスが Kill された場合にトリガーされるアラートを PromQL を利用して設定する

```bash
# 記憶ルールとアラートルールを作成
## base64でエンコードしてAMPへ取り込む
base64 -i rules.yaml -o rules_base64.yaml

aws amp create-rule-groups-namespace \
    --data file://rules_base64.yaml --name rules \
    --workspace-id $(aws amp list-workspaces \
        --query 'workspaces[].workspaceId' --output text)

## 取り込みが完了していることを確認
aws amp describe-rule-groups-namespace \
    --workspace-id $(aws amp list-workspaces \
        --query 'workspaces[].workspaceId' --output text) --name rules \
    --query 'ruleGroupsNamespace.status.statusCode'
"ACTIVE"
# アラート先のAWS SNSトピックを作成
aws sns create-topic --name eks-cluster
## AMPからのアクセスを許可するためSNSトピックへポリシーをアタッチ
aws sns set-topic-attributes \
    --topic-arn 'arn:aws:sns:ap-northeast-1:528163014577:eks-cluster' \
    --attribute-name 'Policy' --attribute-value file://sns-policy.json
## 通知を行うサブスクリプションの追加
aws sns subscribe \
    --topic-arn 'arn:aws:sns:ap-northeast-1:528163014577:eks-cluster' \
    --protocol email --notification-endpoint 'matsuya@h01.itscom.net'

# AMPへAlertmanager定義を追加し、アラートがトリガーされた際にSNSへ通知を行う
base64 -i alertmanager.yaml -o alertmanager_base64.yaml

aws amp create-alert-manager-definition \
    --data file://alertmanager_base64.yaml \
    --workspace-id $(aws amp list-workspaces \
        --query 'workspaces[].workspaceId' --output text)
# 有効になっているか確認
aws amp describe-alert-manager-definition \
    --workspace-id $(aws amp list-workspaces \
        --query 'workspaces[].workspaceId' --output text) \
    --query 'alertManagerDefinition.status.statusCode'
"ACTIVE"

# アラートをトリガーさせる
MemoryLimitsを設定したコンテナを含むPodへstressコマンドでLimits以上のメモリ負荷をかけてOOM Killを発生させる
kubectl apply -f pod.yaml
```

```bash
# CleanUp
kubectl delete pod oom-triger
helm uninstall prometheus-for-amp -n prometheus
kubectl delete ns prometheus
aws sns delete-topic \
    --topic-arn "arn:aws:sns:ap-northeast-1:528163014577:eks-cluster"
aws amp delete-workspace --workspace-id $(aws amp list-workspaces \
    --alias eks-cluster --query 'workspaces[].workspaceId' --output text)
aws grafana delete-workspace --workspace-id $(aws grafana list-workspaces \
    --query 'workspaces[].id' --output text)
SERVICE_ACCOUNT_IAM_ROLE=EKS-AMP-ServiceAccount-Role
SERVICE_ACCOUNT_IAM_POLICY=AWSManagedPrometheusWriteAccessPolicy
SERVICE_ACCOUNT_IAM_POLICY_ARN=arn:aws:iam::528163014577:policy/$SERVICE_ACCOUNT_IAM_POLICY
aws iam detach-role-policy --role-name $SERVICE_ACCOUNT_IAM_ROLE \
    --policy-arn $SERVICE_ACCOUNT_IAM_POLICY_ARN
aws iam delete-role --role-name $SERVICE_ACCOUNT_IAM_ROLE
```

# K8s ロギング

## ログの収集と転送

K8s では Control Plane, Node, コンテナの 3 つでログが出力されるため、これらのログを収集し安全に保存する必要がある

- Control Plane のログ  
  EKS Cluster でログの出力を有効にすると、CWL に出力される
- Data Plane のログ  
   Node のログと Workload(pod)のログが存在する.
  Node のログは　 Node 上で稼働する kubelet や kube-proxy のログ、さらに Node のホスト OS 上に出力される Syslog など.  
  Workload のログは、コンテナ内で動作するアプリケーションが出力するログ

Node とコンテナのログは k8s によって、それぞれホスト OS の下表の場所へ出力される

| ログ種別         | 出力先                                                 |
| ---------------- | ------------------------------------------------------ |
| Node(Data Plane) | /var/log/journal                                       |
| Node(ホスト OS)  | /var/log/dmesg<br>/var/log/secure<br>/var/log/messages |
| コンテナ         | /var/log/containers<br>/var/log/pods                   |

ログを収集する方法として、サイドカーコンテナとしてログエージェントを同梱させる方法雨や、アプリケーション自身が外部の＾ログサービスへ直接書き込む方法などが存在するが、一般的な方法では、Fluentd などのログ収集ツールを DaemonSet として各 Node 上に起動し、ホスト OS 上に出力される Node やコンテナのログを収集する手法が利用される

クラウドのコンテナ環境下では Fluent Bit が利用されることが多い.  
FluentBit は Fluentd をベースに作られており、Fluentd は Ruby の Gem などの依存関係が必要とされるのに対して、FluentBit は C のみで書かれているため、特別なプラグインなどが必要としない限り依存関係を考慮する必要がない. また FluentBit は Fluentd よりも高速なため大量のログが出力されるコンテナ環境下でも高いパフォーマンスを破棄することができ、Elasticsearch や CWL,
AWS Kinesis Data Firehorse など、さまざまなログサービスにログを転送することもできるため、高度なクエリを用いた分析も可能.

```bash
# amazon-cloudwatch namespaceを作成
kubectl apply -f https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/cloudwatch-namespace.yaml

# FluentBit Daemonsetが利用するConfigMapを作成
ClusterName=eks-cluster
RegionName=ap-northeast-1
FluentBitHttpPort='2020'
FluentBitReadFromHead='Off'
[[ ${FluentBitReadFromHead} = 'On' ]] && FluentBitReadFromTail='Off' || \
FluentBitReadFromTail='On'
[[ -z ${FluentBitHttpPort} ]] && FluentBitHttpServer='Off' || \
FluentBitHttpServer='On'
kubectl create configmap fluent-bit-cluster-info \
    --from-literal=cluster.name=${ClusterName} \
    --from-literal=http.server=${FluentBitHttpServer} \
    --from-literal=http.port=${FluentBitHttpPort} \
    --from-literal=read.head=${FluentBitReadFromHead} \
    --from-literal=read.tail=${FluentBitReadFromTail} \
    --from-literal=logs.region=${RegionName} -n amazon-cloudwatch

# FluentBit最適化設定のFluentBit daemonsetをデプロイ
## FluentBitのマニフェストファイルをダウンロード
curl -O https://raw.githubusercontent.com/aws-samples/amazon-cloudwatch-container-insights/latest/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml

## ConfigMapを修正してFluentBitをデプロイ
kubectl apply -f fluent-bit.yaml

## 確認
kubectl get pods -n amazon-cloudwatch
```

収集されたログは以下のロググループに出力される

- `/aws/containerinsights/${CLUSTER_NAME}`/application
- `/aws/containerinsights/${CLUSTER_NAME}`/host
- `/aws/containerinsights/${CLUSTER_NAME}`/dataplane

## ログの可視化と分析

CWL に出力されたログを CWL insights を使用して分析.  
分析用のログとして Nginx コンテナのログを JSON で出力し、それを CWL へ転送する

```bash
kubectl apply -f nginx.yaml
# 次に作成したNginxのServiceへ、ローカルからポートフォワードする.
# これによりhttp://localhost:8080/にアクセスすると、Cluster上のNginx Podへリクエストが転送される。
kubectl port-forward svc/nginx-service 8080:80 > /dev/null &

# Nginx Podへ30回リクエスト
hey -n 30 -c 1 http://localhost:8080/
# /aws/containerinsights/eks-cluster/application内にNginxログがJSONで出力されている
# nginx-プレフィックスを含むログストリームを対象として、Nginxのリクエストが200したログのHTTPメソッド名、りくえすとパス、HTTPステータスコード、Nginx Pod名を出力するように指定
filter @logStream like "nginx-" | filter log_processed.status = "200" | fields log_processed.method, log_processed.path, log_processed.status, kubernetes.pod_name
#ステータスコードが200を返却したログの総数を確認(bin(1h)により1時間単位で集計)
filter @logStream like "nginx-" | filter log_processed.status = "200" | stats count(*) as statusCount by bin(1h)
```

```bash
# Clean Up
ps aux | grep 'kubectl port-forward svc/nginx' | grep -v 'grep' | \
    awk '{print $2}' | xargs kill -9
kubectl delete service nginx-service
kubectl delete deployment nginx
kubectl delete configmap nginx-conf
kubectl delete ns amazon-cloudwatch
```
