# ホスト OS のセキュリティ Concepts

## アタックサーフェス

アタックサーフェスとは、攻撃される可能性のある箇所を意味する言葉で、攻撃対象領域とも言われる  
コンテナを動作させるために最低限必要なライブラリやツールのみを提供する、コンテナ最適化 OS が開発された。代表的なものとして Amazon が開発した Bottlerocket などがある

## Bottlerocket を利用

K8s Node のホスト OS に不要なツールは含まれていない。Node のホスト OS に利用することで以下のメリットを享受できる

- OS 軽量化
- 起動の高速化
- リソース使用率の削減
- アタックサーフェスの軽減
- アトミック k な OS 自動更新による運用上のオーバーヘッド削減

EKS Node のホストとして default では Amazon Linux 2 が利用されているが、AL2 が必要となるケースは限られている。

- どうしてもホスト OS 上で完全な root 権限が必要な Pod を起動したい
- ホスト OS が使用するライブラリをパッケージで厳密に管理したいなど

可能な限り Bottlerocket を利用することが推奨される

```bash
eksctl create cluster -f eksctl_create_cluster.yaml

# ホストOSがBottlerocketかどうかの確認
kubectl get nodes -o=custom-columns=NODE:.metadata.name,ARCH:.status.nodeInfo.architecture,OS-Image:.status.nodeInfo.osImage,OS:.status.nodeInfo.operatingSystem

# ハンズオンをやめる場合
eksctl delete cluster -f eksctl_create_cluster.yaml

```

## 共有カーネル

すべてのコンテナとホストがカーネルを制限なく共有していると、特定のコンテナからに任意のコンテナやホストへ通信できたり、そのコンテナやホストが使用しているメモリ空間の内容を参照できてしまう。

これを制限するために、K8s では Linux の Namespace 機能を利用してコンテナ間通信などを制限している。

＊K8s の Namespace と Linux の Namespace は直接的な関係はなく、互換性もない

- コンテナを root で実行しない

  コンテナで作成された root ユーザーが、そのままホスト OS でも root として認識されてしまう  
  root ユーザーは UID 0 なので、ホスト OS とコンテナの root ユーザーを同一と判断してしまうため。  
  したがって、コンテナや Pod を実行する際は、一般ユーザーで j ９っこうするように指定を行い,root ユーザーで実行させないようにする

- ホスト OS の NameSpace を共有しない

  特権コンテナにホスト OS の Namespcae を共有すると、特権コンテナ内からホスト OS を管理者権限で操作できてしまう。

- コンテナ化されていないぷろせすを　実行させない
  ホスト OS 上でコンテナ化されていないサービスを直接実行させるとアタックサーフェスを提供しかねないので、EKS Node から切り離すべき

## AMi を定期的に更新

```bash
eksctl upgrade nodegroup --name=bottlerocket --cluster=eks-cluster
```

## 不適切なユーザーアクセス権

Bottlerocket ではデフォルトで提供されるコントロールコンテナがある。

このコントロールコンテナでは SSM Agent が実行されるので、AWS SystemsManger の SessionManager 経由でコントロールコンテナにログインできる

```bash
aws ssm start-session --target $(aws ec2 describe-instances \
    --filters Name=tag:nodegroup-type,Values=Bottlerocket \
    --query 'Reservations[].Instances[].InstanceId' \
    --output text | awk '{print $1}')
```

## ホスト OS のファイルの改竄対策

コンテナからホスト OS の特定の directory をマウントすることで、コンテナブレイクアウト（コンテナを経由したホストへの侵入）のリスクが高まる

- ホストの気密性の高いディレクトリをマウントさせない
  `spec.volumes.hostPath`を指定することで、ホストの directory がマウントできてしまうのでコンテナからホストの directory をマウントさせないようにする。  
  どうしてもホストのディレクトリをマウントさせたいときは、Pod の`spec.containers.securityContext.readOnlyRootFilesystem`と`true`に設定して、読み取り専用でマウントすること

# コンテナのセキュリティコンセプト

## ランタイムソフトウェア内の脆弱性

runc や containerd などのコンテナランタイムに脆弱性が存在するとコンテナブレイクアウトなどの攻撃を可能にしてしまう

- ホスト OS の AMI を更新
- 脆弱性情報を管理する

## コンテナからの無制限のネットワークアクセス

K8s では異なる Namespace 間の Pod や、Pod からホストネットワークへの通信がデフォルトでできてしまう。万が一コンテナに侵入されて攻撃者がネットワークないをスキャンすることにより、さらなる脆弱性や設定不備を発見され悪用されてしまう。

Pod のアウトバウンド通信を制限することが望ましいが、K8s のネットワークは高度に仮想化されており、従来のツールでは Pod 間の通信を細かく制御することは難しいので、最適化されたツールを使用するのが推奨されている

- NetworkPolicy を使用
  NetworkPolicy を実装するための CNI プラグインが別途必要になる  
  EKS ではデフォルトで AWS VPC CNI for K8s が CNI プラグインとして利用できる。このプラグインがないと、NetworkPolicy の設定はできても何も制御されない

  EKS では NetworkPolicy を実装する際に、Calico という CNI プラグインを使用することが推奨されている。

```bash
# calicoのインストール
helm repo add projectcalico https://projectcalico.docs.tigera.io/charts

echo '{installation: {kubernetesProvider: EKS}}' > values.yaml

kubectl create namespace tigera-operator

helm install calico projectcalico/tigera-operator \
    --version v3.23.1 -f values.yaml --namespace tigera-operator

# NetworkPolicyを作成し, Podのアウトバウンド通信を制限
kubectl apply -f default-deny-all-egress-policy.yaml

kubectl get networkpolicy

NAME                      POD-SELECTOR   AGE
default-deny-all-egress   <none>         6s

# Nginx Deployment
kubectl apply -f deployment.yaml
# PodのIPアドレスの確認
kubectl get pods -l app=nginx -o "custom-columns=podIP:status.podIP"
192.168.190.183
# Nginxへリクエストを行うクライアントPod
kubectl apply -f pod.yaml
kubectl exec -it test-pd -- bash

[root@test-pd /]# curl http://192.168.190.183/
curl: (7) Failed connect to 192.168.190.183:80; Connection timed out

# NetworkPolicyを更新してアウトバウンドを許可する
kubectl apply -f allow-to-nginx-egress-policy.yaml
kubectl get networkpolicy
NAME                      POD-SELECTOR   AGE
allow-to-nginx-egress     app=test-pd    4s
default-deny-all-egress   <none>         8m9s
```

NetworkPolicy は Namespace 単位で作成する必要があるため、kube-system 以外の Namespace でデフォルトですべての Egress 通信を制限する NetworkPolicy を適用し、疎通が必要な Workload 間の通信のみを許可する NetworkPolicy を適用することが望ましい

kube-system では基本的に K8s クラスタ￥を運用する上で重要な Workload がデプロイされるので、Egress の制限は行わないほうが良い

## Pod の SG を利用

K8s ないの通信制御は NetworkPolicy を使用するが、Pod から AWS リソースへの通信制業を行う場合は、SG を使う

しかし、以下の制約があることに留意

- SG が適用された Pod は NetworkPolicy が適用されない
- Windows ノードでは使用できない
- IPv6 アドレスを利用する EC2 ノードでは利用できない
- t3 インスタンスファミリーはサポート外

```bash
VPCID=$(aws eks describe-cluster --name eks-cluster \
      --query "cluster.resourcesVpcConfig.vpcId" \
      --output text)

# SGを作成し、デフォルトで追加されるアウトバウンドルールを削除
PODSG=$(aws ec2 create-security-group --group-name PodToRDSAccessSG \
      --description "Security group to apply to apps that need access to RDS" \
      --vpc-id $VPCID \
      --query "GroupId" --output text)
aws ec2 revoke-security-group-egress \
    --group-id $PODSG \
    --cidr 0.0.0.0/0 \
    --protocol all \
    --output text
# RDS用SGを作成し、デフォルトで追加されるアウトバウンドルールを削除
RDSSG=$(aws ec2 create-security-group --group-name RDSSG \
      --description "Security groups to apply to the RDS cluster" \
      --vpc-id $VPCID \
      --query "GroupId" --output text)
aws ec2 revoke-security-group-egress \
    --group-id $RDSSG \
    --cidr 0.0.0.0/0 \
    --protocol all \
    --output text
```

```bash
# POD用SGのアウトバウンドルールにMySQL通信を許可させる
aws ec2 authorize-security-group-egress \
    --group-id $PODSG \
    --ip-permissions IpProtocol=tcp,FromPort=3306,ToPort=3306,UserIdGroupPairs="[{GroupId=$RDSSG}]" \
    --query "Return" --output text

# RDS用SGのインバウンとルールにPODSGを許可させる
aws ec2 authorize-security-group-ingress \
    --group-id $RDSSG \
    --ip-permissions IpProtocol=tcp,FromPort=3306,ToPort=3306,UserIdGroupPairs="[{GroupId=$PODSG}]" \
    --query "Return" --output text

```

```bash
# RDSの作成
# プライベートサブネットの確認
aws ec2 describe-subnets \
    --filters \
    "Name=vpc-id,Values=$VPCID" "Name=tag:kubernetes.io/role/internal-elb,Values=1" \
    --query 'Subnets[].SubnetId' --output text
subnet-0559e9842ebc901ea        subnet-0e631d391c1081a2b        subnet-0ee3f439d51d79c5e

# サブネットグループ作成
aws rds create-db-subnet-group \
    --db-subnet-group-name book-kubernetes-pod-sg \
    --db-subnet-group-description "DB subnet group for testing the pod security group" \
    --subnet-ids '["subnet-0559e9842ebc901ea","subnet-0e631d391c1081a2b","subnet-0ee3f439d51d79c5e"]' \
    --query DBSubnetGroup.SubnetGroupStatus \
    --output text
# クラスター作成
aws rds create-db-cluster --db-cluster-identifier pod-sg \
    --engine aurora-mysql \
    --engine-version 8.0 --master-username admin \
    --master-user-password Pf9qKvbn \
    --db-subnet-group-name book-kubernetes-pod-sg \
    --query DBCluster.Status --output text

# インスタンス作成
aws rds create-db-instance --db-instance-identifier pod-sg-instance \
    --db-cluster-identifier pod-sg --engine aurora-mysql \
    --db-instance-class db.t3.medium \
    --query DBInstance.DBInstanceStatus --output text
```

```bash
# aws-node DaemonSetの環境変数「ENABLE_POD_ENI」をtrueに設定し、PodのENI利用を有効化
kubectl set env daemonset -n kube-system aws-node ENABLE_POD_ENI=true

# PodへSGを適用するためにSecurityGroupPolicy CRDへ定義を行う
# EKSクラスタが使用するSGを取得し、SecurityGroupPolicyを適用
CLUSTERSG=$(aws eks describe-cluster --name eks-cluster \
      --query "cluster.resourcesVpcConfig.clusterSecurityGroupId" \
      --output text)

kubectl apply -f sgp-policy.yaml
kubectl get sgp

NAME           SECURITY-GROUP-IDS
my-sg-policy   ["sg-05052b6526ce16947","sg-0188682d538cdca68"]

# bastion用podのデプロイ
kubectl apply -f aurora-mysql-test-pod.yaml

# RDSのエンドポイント確認
aws rds describe-db-clusters --db-cluster-identifier pod-sg \
    --query "DBClusters[].[Endpoint]" --output text
pod-sg.cluster-cmuonzjucukc.ap-northeast-1.rds.amazonaws.com

kubectl exec -it aurora-mysql-test -- bash

[root@aurora-mysql-test /]# rpm -Uvh https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
[root@aurora-mysql-test /]# rpm --import https://repo.mysql.com/RPM-GPG-KEY-mysql-2022
[root@aurora-mysql-test /]# yum install -y mysql-community-client
[root@aurora-mysql-test /]# mysql -u admin --password=Pf9qKvbn -h pod-sg.cluster-cmuonzjucukc.ap-northeast-1.rds.amazonaws.com

# 接続の確認
mysql> show databases;
```

```bash
# SGが適用されていないPodで検証
kubectl apply -f non-sg-pod.yaml
kubectl exec -it aurora-mysql-test-non-sg -- bash

[root@aurora-mysql-test-non-sg /]# rpm -Uvh https://dev.mysql.com/get/mysql80-community-release-el7-3.noarch.rpm
[root@aurora-mysql-test-non-sg /]# rpm --import https://repo.mysql.com/RPM-GPG-KEY-mysql-2022
[root@aurora-mysql-test-non-sg /]# yum install -y mysql-community-client
[root@aurora-mysql-test-non-sg /]# mysql -u admin --password=Pf9qKvbn -h pod-sg.cluster-cmuonzjucukc.ap-northeast-1.rds.amazonaws.com

```

```bash
# Cleanup
kubectl delete sgp my-sg-policy
kubectl delete pod --force aurora-mysql-test
kubectl delete pod --force aurora-mysql-test-non-sg
aws rds delete-db-instance --db-instance-identifier pod-sg-instance  --skip-final-snapshot --delete-automated-backups
aws rds delete-db-cluster --db-cluster-identifier pod-sg --skip-final-snapshot
aws rds delete-db-subnet-group --db-subnet-group-name book-kubernetes-pod-sg
aws ec2 revoke-security-group-egress --group-id $PODSG --ip-permissions IpProtocol=tcp,FromPort=3306,ToPort=3306,UserIdGroupPairs="[{GroupId=$RDSSG}]"
aws ec2 revoke-security-group-ingress --group-id $RDSSG --ip-permissions IpProtocol=tcp,FromPort=3306,ToPort=3306,UserIdGroupPairs="[{GroupId=$PODSG}]"
aws ec2 delete-security-group --group-id $PODSG
aws ec2 delete-security-group --group-id $RDSSG
```

## セキュアではないコンテナランタイムの設定

- root ユーザーで実行させない
  コンテナはデフォルトで root ユーザーで実行される。コンテナは原則一般ユーザーで実行委させること。

  K8s 環境でコンテナの実行ユーザーに、一般ユーザーを指定する方法は 2 通り。1 つは Dockefile で指定する方法。  
  もう一つは、k8s マニフェストの`spec.containers.securityContext.runAsUser`と`spec.containers.securityContext.runAsNonRoot`で指定する方法
  なお、`runAsUser`はユーザー名ではなく、UID で指定する必要がある

- 特権コンテナを実行させない
  特権コンテナとは、Privileged(特権)が付与されたコンテナを指す。この privileged なコンテナは、ホスト OS のすべてのデバイスへアクセス可能  
  root ユーザーで実行するだけでは本来セキュリティリスクになりえないというのは、この Privilieged や強い権限が root に付与されていない前提の話であって、一般ユーザーでも Privileged が付与されていればリスクになり得る。

## アプリの脆弱性

- WebApplication Firewall の利用
  WAF は DDoS や SQL インジェクションといった L7 れべるの攻撃を防御することができる。

## 未承認コンテナと対策

- 未承認イメージのデプロイを制限する

# コンテナオーケストレータのセキュリティコンセプト

## 制限のない管理者アクセス

- aws-auth で最小の権限を付与

EKS は IAM Entity と K8s Entity のマッピング譲歩湯を、K8s の ConfigMap で aws-auth として保持。  
EKS クラスターを作成した IAM Entity は自動的に system:masters(K8s RBAC の管理者権限)が付与されるが、これは aws-auth では管理されていない。  
aws-auth へ認証設定を追加するまでは作成時に自動で付与された権限を削除することはできないので、この自動で付与される権限から管理権限を引き継ぐ IAM Entity を aws-auth へ登録しオペレーションを行う必要がある  
自動で付与された IAM Entity は aws-auth に破壊的な変更が加わったときなど aws-auth に認証できなくなったときの復旧用として用いるができる

```bash
# aws-authの設定を確認
kubectl get cm -n kube-system aws-auth -o yaml

```

Managed NodeGroup の場合は、デフォルトで Node が利用する IAM Role が K8s RBAC の system:bootstrappers と system:nodes Group に登録される  
aws-auth は kubectl コマンドで直接変更を行うことができるが、誤った内容で登録してしまうと Cluster へアクセスできなくなるので、eksctl で編集を行うことが推奨されている

```bash
# 特定のIAMユーザーをclusterに付与する例
eksctl create iamidentitymapping --cluster eks-cluster \
    --arn 'arn:aws:iam::528163014577:user/eks-cluster' \
    --group system:masters --username admin
# aws-authの内容確認
kubectl get cm -n kube-system aws-auth -o yaml

# 権限削除
eksctl delete iamidentitymapping --cluster eks-cluster \
    --arn 'arn:aws:iam::528163014577:user/eks-cluster'

# aws-auth認証できないことを確認
kubectl get cm -n kube-system aws-auth -o yaml

error: You must be logged in to the server (Unauthorized)
```

*EKS Cluster へのアクセス*は AWS IAM と aws-auth、*権限の設定*は 次で紹介する K8s RBAC を用いる

- RBAC で最小の権限を付与
  個々の User,Group, Service Account などの Entity に対して, Role をベースに Pod などのオブジェクトへのアクセス権限を付与  
  K8s おけるユーザーは`K8sによって管領されるServiceAccount`と`通常のユーザー`がある　　
  通常のユーザーは K8s クラスターから独立したサービスで管理し、K8s RBAC の Subjects としてマッピングすることではじめて K8s RBAC でユーザーを識別することができる.

  EKS では aws-auth により、IAM Entity を K8s の通常のユーザーとしてマッピングすることでコレを実現する

```bash
# kube-systemのPodへの読み取りアクセス権を付与するRole
kubectl apply -f role.yaml

# ops_userへ付与(role bindingするまえにroleを作成しておくこと)
kubectl apply -f role_binding.yaml

# aws-authでIAM のeks-cluster UserをK8sのops_userへマッピング
eksctl create iamidentitymapping --cluster eks-cluster \
    --arn 'arn:aws:iam::528163014577:user/eks-cluster' \
    --username ops_user

# ops_userでkube-system以外のpodを確認
kubectl get pods

Error from server (Forbidden): pods is forbidden: User "ops_user" cannot list resource "pods" in API group "" in the namespace "default"

# ops_userでkube-systemのpodsを確認
kubectl get pods -n kube-system
NAME                       READY   STATUS    RESTARTS   AGE
aws-node-9r6n9             1/1     Running   0          72m
coredns-5b6d4bd6f7-dk565   1/1     Running   0          83m
coredns-5b6d4bd6f7-q7l4x   1/1     Running   0          83m
kube-proxy-n5qck           1/1     Running   0          72m
```

**RBAC Manager** と **RBAC Lookup** を利用することにより
RBAC を直感的に管理できる  
**RBAC Lookup** は一致する User, Group, ServiceAccount に
Binding されている、**Role や ClusterRole を一覧で出力**
Role は存在するが、Binding が　設定されていないものは出力されないため、
Cluster 内で有効な Role のみを間 t なんでより視覚的に確認することができる

```bash
# install
brew install FairwindsOps/tap/rbac-lookup

# ops_userのRBACを一覧出力
rbac-lookup ops_user

SUBJECT     SCOPE         ROLE
ops_user    kube-system   Role/pod-reader
```

**RBAC Manager** は **RoleBinding** を一元管理することができる

```bash
# RBAC Managerのinstall
helm repo add fairwinds-stable https://charts.fairwinds.com/stable
helm install fairwinds-stable/rbac-manager --generate-name \
    --namespace rbac-manager --version 1.11.1 --create-namespace

# RBACDefinitionをデプロイ
kubectl apply -f rbac-definition.yaml
# RoleBindingが作成されているか確認
kubectl get rbacdefinition
rbac-lookup jane
SUBJECT             SCOPE          ROLE
jane@example.com    cluster-wide   ClusterRole/cluster-admin
```

```bash
# Clean up
kubectl delete rolebinding read-pods -n kube-system
kubectl delete role pod-reader -n kube-system
eksctl delete iamidentitymapping --cluster eks-cluster \
    --arn 'arn:aws:iam::528163014577:user/eks-cluster'
kubectl delete rbacdefinition rbac-manager-users-example
kubectl delete ns rbac-manager
```

- IRSA を利用
  IAM Roles for ServiceAccounts(IRSA)は、K8s の ServiceAccount と AWS IAM Role を紐付けることにより、その Service Account を利用する Pod へ IAM Role の権限を付与する EKS の機能。  
  K8s の ServiceAccount は K8s により管理される User オブジェクト。  
  ServiceAccount へ RBAC 認証を設定し、その ServiceAccount を利用するように Pod を作成することで ServiceAccount 経由で Pod が認証され、RBAC により認可を受けることができる  
  Pod 単位で必要となる最小の権限のみを付与することが可能

```bash
# IRSAを利用するにはEKS Cluster内にOIDC Providerを作成する必要がある
eksctl utils associate-iam-oidc-provider --cluster eks-cluster --approve

# OIDC Providerが存在するか確認
aws eks describe-cluster --name eks-cluster \
    --query "cluster.identity.oidc.issuer" --output text
https://oidc.eks.ap-northeast-1.amazonaws.com/id/361DFBCB67DA6FE59991F3C6AD4C864D

# ServiceAccountにマッピングさせるIAMロールを作成
aws iam create-role --role-name batch-sa-role --assume-role-policy-document \
'{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::528163014577:oidc-provider/oidc.eks.ap-northeast-1.amazonaws.com/id/361DFBCB67DA6FE59991F3C6AD4C864D"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "oidc.eks.ap-northeast-1.amazonaws.com/id/361DFBCB67DA6FE59991F3C6AD4C864D:sub": "system:serviceaccount:default:batch"
        }
      }
    }
  ]
}'
# 作成したRoleにAmazonS3ReadOnlyAccessをアタッチ
aws iam attach-role-policy --role-name batch-sa-role \
    --policy-arn 'arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess'

# IAM Roleを使用するServiceAccountを作成
kubectl apply -f sa.yaml
# ServiceAccountを利用するPodをデプロイし、ログを確認
kubectl apply -f pod.yaml
kubectl logs batch

# ポリシーをdetachして確認
kubectl delete -f pod.yaml
aws iam detach-role-policy --role-name batch-sa-role \
    --policy-arn 'arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess'
kubectl apply -f pod.yaml
kubectl logs batch
An error occurred (AccessDenied) when calling the ListObjectsV2 operation: Access Denied
```

また、Pod から Node(EC2)のインスタンスメタデータへアクセスできると、Node に付与されている IAM Role を Pod から利用できてしまう  
EKS Node には基本的に ECR リポジトリの参照権限が付与されているため、それらの権限を Pod が引き継いでしまう.

```bash
kubectl apply -f pod-node-role.yaml
kubectl logs pod-node-role
```

コレを防ぐためにも Workload 単位で IRSA による IAM Role の付与を行う必要がある

```bash
# clean up
kubectl delete pod pod-node-role
kubectl delete pod batch
kubectl delete sa batch
aws iam delete-role --role-name batch-sa-role
```

## Workload の機微性レベルの混同

- NodeSelector を利用
  機密情報を扱う必要があり限定的に公開している Workload と、パブリックに公開している Workload が同一 Node 上で稼働していると、万が一コンテナブレイク・アウトされた場合に、Node のホスト OS 経由で限定公開の Workload へ侵害される可能性がある.  
   K8s ではデフォルトですべての Workload が同一 Node 上で稼働できるようになっているため、機微性の高い Workload と低い Workload はそれぞれ専用の Node Group を用意するなど、機微製の異なる workload をセグメント化し、同一 Node 上で稼働しないようにする
  | 機能名 | 特徴 |
  | -------------------------- | ------------------------------------------------------------------------------------------------ |
  | NodeSelector | Pod を稼働させる Node を、Node に付与されている Label を選択してスケジュールを行う |
  | Node Affinity | 基本的には NodeSelector と同等の機能だが、優先条件と必須条件を指定することができる |
  | Pod Affinity/Anti-Affinity | Node に付与された Label ではなく、すでに Node で稼働している Pod の Label を選択してスケジュールを行う |
  | Taint/Toleration | Node へ Taint(汚れ)を付与し、Toleration(汚れの受け入れ)を認めない Pod はスケジュール及び実行させない |

- NodeAffinity を利用
  Node に Label を付与し、それを Pod の`spec.containers.nodeSelector`で選択する

```bash
# Nodeにラベルを付与
kubectl get nodes
kubectl label nodes ip-192-168-167-61.ap-northeast-1.compute.internal disktype=ssd

# Labelに一致するNodeが存在する場合は、デプロイされる
kubectl apply -f node-selector-pod.yaml
kubectl get pods
kubectl delete pod nginx
# Labelに一致するNodeが存在しないばあいは、デプロイされない
kubectl apply -f node-selector-nfs-pod.yaml
kubectl get pods

NAME    READY   STATUS    RESTARTS   AGE
nginx   0/1     Pending   0          7s
```

```bash
# Clean up
kubectl delete pod nginx
kubectl label nodes ip-192-168-167-61.ap-northeast-1.compute.internal disktype-
```

- Node Affinity を利用
  NodeAffinity ではさらに必須条件と優先条件を指定することができるため、より柔軟なスケジューリング制御を行うことができる  
  `spec.affinity.nodeAffinity`に指定した条件を満たす場合のみスケジューリングを行う必須条件か、必須条件を満たさないばあいでもスケジューリングは実行される優先条件、もしくはその両方を指定。  
  必須条件は`requiredDuringSchedulingIgnoredDuringExecution`を、優先条件は`preferredDuringSchedulingIgnoredDuringExecution`

- PodAffinity/Anti-Affinity を利用する
  NodeAffinity は Node を Labe で指定する機能だったが、Pod Affinity/Anti-Affinity は Node 上で稼働する別の＾ Pod が持つ Label に対して、指定した条件と一致するばあいのみスケジューリングするといった制御を行うことができる  
  PodAffinity/Anti-Affinity は、大規模な Cluster 上で使用する際にスケジューリングを＾非常に遅くする恐れのある多くの処理を要するため、数百台以上の Node からなる Cluster では使用を推奨されない

- Taint/ Toleration を利用
  ここまでは優先条件と必須条件をもとに、一致するばあいに Node へスケジューリングする方法を解説したが、これから解説する Taint は逆に、**条件を満たさないばあいは Pod を Node から排除する**  
  Taint(汚れ)を Node へ Label することにより Pod はその Node へスケジューリングできなくなる。Taint を Toleration(許容)する Pod のみ、その Node 上へスケジューリングされることが許可される

```bash
# nodeにtaintを付与
kubectl get nodes
kubectl taint nodes ip-192-168-167-61.ap-northeast-1.compute.internal key1=value1:NoSchedule

# TolerationせずにPodをデプロイする
kubectl apply -f not-toleration-pod.yaml
kubectl get pods
NAME    READY   STATUS    RESTARTS   AGE
nginx   0/1     Pending   0          3s

kubectl describe pods nginx

# TolerationしたPodをデプロイ
kubectl delete pod nginx
kubectl apply -f toleration-pod.yaml
kubectl get pods
NAME    READY   STATUS    RESTARTS   AGE
nginx   1/1     Running   0          3s
```

```bash
# CleanUP
kubectl delete pod nginx
kubectl taint nodes ip-192-168-167-61.ap-northeast-1.compute.internal key1=:NoSchedule-
```

## オーケストレータオーケストレータ Node の信頼性

K8s の Control Plane に不正な Node が接続されるなど、K8s Node が侵害されると Control Plane や他の Node、その上で稼働している Container へ侵害される可能性があるため、信頼できる Node のみが Cluster に参加できる状態となっていることが望ましい。

- エンドツーエンドの通信を暗号化する
  EKS では、DataPlane から Control Plane への接続は、PKI 証明書を利用し TLS で暗号化された状態で通信される。  
  Control Plane から Data Plane への通信は HTTPS で暗号化して通信されるが、K8s の default では K8s API は kubelet が提供する証明書を検証しないため、中間者攻撃を受けやすくなる。  
  また,Control Plane の K8s API から Node, Pod, Service へのアクセスは、K8s のデフォルトでは平文の HTTP で接続されるため、認証も暗号化もされない。  
  そこで**Service Mesh**により提供される、mutual TLS(mTLS)を利用し、通信元と通信先がお互いに認証を行うことで、Cluster 内の Workload 間でエンドツーエンドの略号化を実現することができる。  
  AWS ではマネージドの Service Mesh サービスとして、App Mesh が提供されている。  
  Data Plane はサイドカープロ棋士が Workload 内に出入りする通信の制御を行い、Control Plane では Data Plane のトラフィックルール管理や証明書管理など、Service Mesh 全体の管理を行う.

```bash
# App Mesh の mTLS を利用し、エンドツーエンドの暗号化を行う手順
# リポジトリのクローン
git clone https://github.com/aws/aws-app-mesh-examples.git
# AppMesh controller for K8s のCRDをClusterへインストール
helm repo add eks https://aws.github.io/eks-charts
kubectl apply \
  -k "https://github.com/aws/eks-charts/stable/appmesh-controller/crds?ref=master"
# インストールされたか確認
kubectl api-resources --api-group=appmesh.k8s.aws -o wide
# 次にコントローラ本体をインストール
helm upgrade -i appmesh-controller eks/appmesh-controller \
    --namespace appmesh-system \
    --set region=ap-northeast-1 \
    --set serviceAccount.create=false \
    --set serviceAccount.name=appmesh-controller \
    --set sds.enabled=true \
    --version 1.5.0

# コントローラのバージョンの確認
kubectl -n appmesh-system get deployment appmesh-controller -o json | \
    jq -r ".spec.template.spec.containers[].image" | cut -f2 -d ':'
v1.5.0
# SPIREサーバーをStatefulSetとして、
# SPIREエージェントをDaemonSetとしてNodeごとに1つずつ信頼ドメインであるhowto-k8s-mtls-sds-based.awsを使用してインストール
kubectl apply -f aws-app-mesh-examples/walkthroughs/howto-k8s-mtls-sds-based/spire/spire_setup.yaml

# SPIREサーバーとSPIREエージェントが正常に稼働しているか確認
kubectl -n spire get all

# フロントアプリケーション(blueとred)を登録
./aws-app-mesh-examples/walkthroughs/howto-k8s-mtls-sds-based/spire/register_server_entries.sh register

# 接続性をテストするサンプルアプリをdeploy
export ENVOY_IMAGE=840364872350.dkr.ecr.ap-northeast-1.amazonaws.com/aws-appmesh-envoy:v1.22.0.0-prod
export AWS_ACCOUNT_ID=528163014577
export AWS_DEFAULT_REGION=ap-northeast-1

./aws-app-mesh-examples/walkthroughs/howto-k8s-mtls-sds-based/deploy_app.sh

# リソースの確認
kubectl -n howto-k8s-mtls-sds-based get all
```

mTLS による相互認証が適用されていることを確認

```bash
# 動作確認用podへログイン
kubectl run -i --tty curler \
    --image=public.ecr.aws/k8m1l3p1/alpine/curler:latest --rm

# mTLSが有効なblueへアクセス
/ # curl -H "color_header: blue" front.howto-k8s-mtls-sds-based.svc.cluster.local:8080; echo;
blue

# mTLSが有効でないgreenへアクセス
/ # curl -H "color:_header: green" front.howto-k8s-mtls-sds-based.svc.cluster.local:8080; echo;
```

```bash
# Clean Up
kubectl delete ns howto-k8s-mtls-sds-based
./aws-app-mesh-examples/walkthroughs/howto-k8s-mtls-sds-based/cleanup.sh
helm uninstall -n appmesh-system appmesh-controller
kubectl delete -k \
    "https://github.com/aws/eks-charts/stable/appmesh-controller/crds?ref=master"
aws ecr delete-repository --repository-name howto-k8s-mtls-sds-based
unset AWS_ACCOUNT_ID
unset AWS_DEFAULT_REGION
unset ENVOY_IMAGE
```

# コンテナイメージのセキュリティコンセプト

## コンテナイメージの脆弱性

- イメージを最小化
  イメージ内に使用しないライブラリやツールなどのバイナリが存在していると、脆弱性対応を行う対象が増えてしまい、アタックサーフェスを拡大してしまう。  
  イメージを最小化することでイメージのビルドにかかる時間短縮にもなる

  1. 最小化 s れたベースイメージを利用する
  2. マルチステージビルドを利用する

- パッケージを Update

- コンテナイメージの脆弱性をスキャン

  1. AWS Security Hub のセットアップ：コンソールから有効化
  2. 「統合」から**aqua security**を入力
  3. GitHub Actions が利用する IAM Role を作成

```bash
aws iam create-role --role-name GitHubActionRole \
    --assume-role-policy-document \
'{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::528163014577:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:ShotaroMatsuya/learning-eks:*"
        }
      }
    }
  ]
}'

aws iam attach-role-policy --role-name GitHubActionsRole --policy-arn 'arn:aws:iam::aws:policy/AdministratorAccess'
```

どうしても特定の脆弱性を許容せざるを得ないばあいは, **.trivyignore**ファイルに許容する脆弱性の脆弱性 ID を記載する

## コンテナイメージの設定不備

- hadolint を利用
- Polaris を利用

## 埋め込まれた平文の秘密情報

- Secret を利用
  K8s では秘密情報を管理し、必要になったタイミングで提供する機能として、Secrets が提供される。  
   Secrets には下表の 8 つのタイプがあり、Secret 作成時に type フィールドで指定することで、Secret 情報登録時の型を定義できる(デフォルトは*Opaque*)
  | タイプ | 利用用途 | 特徴 |
  | ----------------------------------- | --------------------------------------------------- | ------------------------------------------------------------------------------------------- |
  | Opaque | 任意設定可能なユーザー定義データ | タイプを省略した際のデフォルトで、任意のデータを登録できる |
  | kubernetes.io/service-account-token | ServiceAccount Token | K8s の ServiceAccount Token を格納できる |
  | kubernetes.io/dockerconfig | シリアライズされた「~/.dockercfg」ファイル | Docker レジストリにアクセスするための資格情報を格納、「~/.dockercfg」ファイルを登録 |
  | kubernetes.io/dockerconfigjson | シリアライズされた「~/.docker/config.json」ファイル | Docker レジストリにアクセスするための資格情報を格納、「~/.docker/config.json」ファイルを登録 |
  | kubernetes.io/basic-auth | ベーシック認証のための認証情報 | 「data」フィールドに「username」と「password」キー指定する |
  | kubernetes.io/ssh-auth | SSH 認証のための認証情報 | 「data」フィールドに「ssh-private key」キーとその値のペアを指定する |
  | kubernetes.io/tls | TLS クライアントまたはサーバーの鍵データ | 「data」フィールドに「tls.key」と「tls.crt」キーを指定する |
  | bootstrap.kubernetes.io/token | ブートストラップトークンデータ | Node のブートストラッププロセス中に使用されるトークンで、ConfigMap に署名するために使用 |

  Secret 内の**data**フィールドのキーへ値を指定するばあいは、base64 でエンコードされている必要がある。  
  data フィールドの代わりに、**stringData**フィールドを利用することで、base64 のエンコードは不要で平文で値を指定可能になる  
  Secrets の Pod へのまうんとほうほうとして　、次の 3 つがある

  1. ボリューム内のファイルとして、Pod の単一または複数のコンテナにマウントする
  2. コンテナの環境変数として利用
  3. Pod を生成するために、kubelet がイメージを pull するときに使用

```bash
# ボリュームとしてマウントする方法
# apply
kubectl apply -f secret-file.yaml
kubectl exec -it secret-file-pod bash
# コンテナ内ではデコードされている
root@secret-file-pod:/data# ls /etc/secret/
password  username
root@secret-file-pod:/data# cat /etc/secret/username
admin
root@secret-file-pod:/data# cat /etc/secret/password
12345678
```

kubelet は定期的な動機のたびにマウントされた Secret が新しいかどうかを確認。  
Secret の値が更新されると、マウントとされている Secret の値も更新される。  
この動作は**KubeletConfiguration struct**内の**ConfigMapAndSecretChangeDetectionStrategy**によっていくかの動作モードが定義されていて、デフォルトで**watch**(Secret の値の変更を監視)

```bash
# 環境変数として利用
# apply
kubectl apply -f secret-env.yaml
# 確認
kubectl exec -it secret-env-pod bash
root@secret-env-pod:/data# printenv | grep -i username
SECRET_USERNAME=admin
root@secret-env-pod:/data# printenv | grep -i password
SECRET_PASSWORD=12345678
```

セキュリティの観点から言うと、環境変数として Secret を利用すると、コンテナランタイムやアプリケーションの仕様によりログへ環境変数の情報が出力される場合があるため、ボリュームからマウントして Secret を利用する胃ことが望ましい

```bash
kubectl delete pod secret-env-pod
kubectl delete pod secret-file-pod
kubectl delete secret test-secret
```

- ASCP を利用する

  Secret をマニフェストで管理するばあいは、マニフェスト内に機密情報をハードコーディングするわけに行かない上、kubectl で直接作成すると GitOps でステート管理できなくなってしまう  
  この問題に対応するため、主要なクラウドベンダーでは Secrets Store CSI ドライ g バー向けのプロバイダが提供されている  
  AWS Secrets and Configuration Provider(ASCP)がそれに当たる  
  ASCP を使えば、AWS Secrets Manager でシークレットを安全に保管及び管理し、カスタムコードを記述しなくても K8s 上で動作するアプリケーションからシークレットを取得することができる。  
  また、IAM とリソースポリシーを Secret に適用することで, K8s クラスタ内の特定の Pod にアクセス制御をかけることができる.

```bash
# EKS ClusterへASCPをインストールする手順
## まずSecrets Store CSIドライバーをインストール
helm repo add secrets-store-csi-driver https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts
helm install -n kube-system csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver
# 確認
kubectl --namespace=kube-system get pods -l "app=secrets-store-csi-driver"

## ASCPをインストール
kubectl apply -f https://raw.githubusercontent.com/aws/secrets-store-csi-driver-provider-aws/main/deployment/aws-provider-installer.yaml

kubectl get daemonset -n kube-system -l app=csi-secrets-store-provider-aws

## AWS Secrets ManagerのSecretを作成
aws secretsmanager create-secret --name User \
    --secret-string '{"username": "admin", "password": "12345678"}'

## IRSA用のIAM Policy
POLICY_ARN=$(aws --region ap-northeast-1 --query Policy.Arn \
    --output text iam create-policy --policy-name nginx-deployment-policy \
    --policy-document \
    '{
      "Version": "2012-10-17",
      "Statement": [{
        "Effect": "Allow",
        "Action": ["secretsmanager:GetSecretValue", "secretsmanager:DescribeSecret"],
        "Resource": ["arn:aws:secretsmanager:ap-northeast-1:528163014577:secret:User-izurlu"]
      }]
    }')
## IRSA用のServiceAccount
eksctl create iamserviceaccount --name nginx-deployment-sa \
    --cluster "eks-cluster" --attach-policy-arn $POLICY_ARN \
    --approve --override-existing-serviceaccounts

## ASCPにSecrets ManagerのSecretとK8sのSecretをマッピングさせる
kubectl apply -f secret_provider_class.yaml
kubectl get secretproviderclass

## SecretProviderClassで作成したSecretを利用するNginx
kubectl apply -f nginx.yaml
kubectl get pod
## SecretがただしくNginxにマウントされていることを確認
kubectl exec -it $(kubectl get pods | awk '/nginx/{print $1}' | head -1) \
    -- cat /secrets/User; echo
{"username": "admin", "password": "12345678"}
```

```bash
kubectl delete deployment nginx
kubectl delete secretproviderclass user-aws-secrets
kubectl delete -f https://raw.githubusercontent.com/aws/secrets-store-csi-driver-provider-aws/main/deployment/aws-provider-installer.yaml
eksctl delete iamserviceaccount --name nginx-deployment-sa \
    --cluster "eks-cluster"
helm uninstall -n kube-system csi-secrets-store
```

# コンテナレジストリのセキュリティコンセプト

## コンテナレジストリへのセキュアでない接続

ECR ではデフォルトで HTTPS 経由で暗号化されて通信する。 しかしデフォルトの状態だと、インターネット経由での通信となってしまう。  
データは暗号化されていても、インターネット上の攻撃者は通信自体の傍受は可能になってしまうというリスクがある  
暗号化に加えてプライベートにレジストリと通信を行うことでよりセキュリティを高めることができる

- VPC エンドポイントを利用
  Interface 型はサブネット内にそのサービスの VPC エンドポイントの ENI が生成される  
  あた VPC エンドポイントのプライベート DNS を有効化することで、プライベートサブネットはもちろんパブリックサブネットからの対象サービスへの通信も自動的にプライベート IP で通信される  
  Gateway 型は VPC エンドポイント経由でプライベートに通信を行いたいサブネットのルートテーブルにルートの設定を行うことで、プライベート IP で通信される  
  Gateway 型の VPC エンドポイント咲く世辞に指定したルートテーブルに pl-xxxxxxx のルートが自動的に設定される

```bash
# Interface型のECR VPCエンドポイントを設定
## private subnetを指定
aws ec2 create-vpc-endpoint --vpc-id vpc-0d72437fcc093711c  \
    --vpc-endpoint-type Interface \
    --service-name "com.amazonaws.ap-northeast-1.ecr.api" \
    --subnet-ids subnet-04f68f76fe45966e7 subnet-005fa6168c9d3f5a7 subnet-04adf4e0f3e8f44b8 \
    --security-group-ids sg-0aef696c7ee151180 \
    --private-dns-enabled \
    --tag-specifications 'ResourceType=vpc-endpoint,Tags=[{Key=Name,Value=test-endpoint1}]'
VpcEndpointId=vpce-068b9b8937203cebe

aws ec2 create-vpc-endpoint --vpc-id  vpc-0d72437fcc093711c \
    --vpc-endpoint-type Interface \
    --service-name "com.amazonaws.ap-northeast-1.ecr.dkr" \
    --subnet-ids subnet-04f68f76fe45966e7 subnet-005fa6168c9d3f5a7 subnet-04adf4e0f3e8f44b8  \
    --security-group-ids sg-0aef696c7ee151180 \
    --private-dns-enabled \
    --tag-specifications 'ResourceType=vpc-endpoint,Tags=[{Key=Name,Value=test-endpoint2}]'
VpcEndpointId=vpce-00f65096fda03d002

# ECRは内部的にS3が使用されているため、Gateway型のS3のVPCエンドポイントも作成する
aws ec2 create-vpc-endpoint --vpc-id  vpc-0d72437fcc093711c \
    --vpc-endpoint-type Gateway \
    --service-name "com.amazonaws.ap-northeast-1.s3" \
    --route-table-ids rtb-0284019d965f3436e rtb-08e5d6360397a4bbb rtb-052e699109d639772 \
    --tag-specifications 'ResourceType=vpc-endpoint,Tags=[{Key=Name,Value=test-endpoint3}]'
VpcEndpointId=vpce-069ff0e0ff8b0607d

# NAT Gatewayが存在しないprivate subnetでは前述の3つのVPCエンドポイントに加えてEC2のVPCエンドポイントが必要
# Interface型のVPCエンドポイントを作成
aws ec2 create-vpc-endpoint --vpc-id  vpc-0d72437fcc093711c \
    --vpc-endpoint-type Interface \
    --service-name "com.amazonaws.ap-northeast-1.ec2" \
    --subnet-ids subnet-04f68f76fe45966e7 subnet-005fa6168c9d3f5a7 subnet-04adf4e0f3e8f44b8  \
    --security-group-ids sg-0aef696c7ee151180 \
    --private-dns-enabled \
    --tag-specifications 'ResourceType=vpc-endpoint,Tags=[{Key=Name,Value=test-endpoint4}]'
VpcEndpointId=vpce-001fe2cd291bad80d
```

```bash
# VPCエンドポイント
## VPCエンドポイント経由でECRに対するアクションを制限できる
## なおPrincipalはEKSノードのIAM Roleにするのが望ましい
aws ec2 modify-vpc-endpoint --vpc-endpoint-id ${VPC_ENDPOINT_ID} \
    --policy-document \
    '{
      "Statement": [
        {
          "Sid": "AllowAll",
          "Effect": "Allow",
          "Principal": "*",
          "Action": "*",
          "Resource": "*"
        },
        {
          "Sid": "AllowPull",
          "Principal": {
            "AWS": "arn:aws:iam::${AWS_ACCOUT_ID}:role/${ROLE_NAME}"
          },
          "Action": [
            "ecr:BatchGetImage",
            "ecr:GetDownloadUrlForLayer",
            "ecr:GetAuthorizationToken"
          ],
          "Effect": "Allow",
          "Resource": "*"
        }
      ]
    }'
## ECRが内部に使用するS3バケットからのみオブジェクトの取得を許可するVPC Endopointポリシーの一例
aws ec2 modify-vpc-endpoint --vpc-endpoint-id vpce-069ff0e0ff8b0607d \
    --policy-document \
    '{
      "Statement": [
        {
          "Sid": "Access-to-specific-bucket-only",
          "Principal": "*",
          "Action": [
            "s3:GetObject"
          ],
          "Effect": "Allow",
          "Resource": ["arn:aws:s3:::prod-ap-northeast-1-starport-layer-bucket/*"]
        }
      ]
    }'
```

## コンテナレジストリ内の古いイメージ

- タグのイミュータビリティを有効にする
  イミュータビリティを有効化することで同一タグを使用したイメージの上書きを防止できる.  
  これにより最新バージョンのイメージをデプロイする際は、目地的にバージョンを指定する運用を行うことが可能  
  マニフェスト側でイメージタグを指定する際は、latest タグではなく、そのイメージのタグ名を用いて指定を行うことが誤デプロイ防止の観点から望ましい

- 古いイメージを自動的に削除
  特定期間が経過した古いイメージが定期的に削除されることで、古いイメージが誤ってデプロイされることを防ぐことができる
