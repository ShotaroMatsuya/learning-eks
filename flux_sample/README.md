# flux sample

## Terraform 実行

```bash
terraform init
terraform apply -auto-approve
```

## サンプルアプリのデプロイ

```bash
# Dockerfimageの作成
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 528163014577.dkr.ecr.ap-northeast-1.amazonaws.com
docker build -t sample-app .
docker tag sample-app:latest 528163014577.dkr.ecr.ap-northeast-1.amazonaws.com/sample-app:latest
docker push 528163014577.dkr.ecr.ap-northeast-1.amazonaws.com/sample-app:latest
```

```bash
# AWS LoadBalancer Controllerのインストールを確認
kubectl get deployment -n kube-system aws-load-balancer-controller
```

1. Flux のインストール

```bash
brew install fluxcd/tap/flux

# 使っているEKSクラスターがFluxコンポーネントをインストールするための要件を見対しているかを確認
flux check --pre
```

2. Personal Access Token を作成

略

3. Flux ツールキットコンポーネントのインストール

```bash
# 環境変数の設定
export GITHUB_TOKEN=ghpxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export GITHUB_USER=ShotaroMatsuya

# コンポーネントinstall
flux bootstrap github \
    --owner=${GITHUB_USER} \
    --repository=flux-repo \
    --branch=master \
    --path=eks-cluster \
    --personal
```

4. リソース確認

```bash
# flux-systemが生成されていることを確認
kubectl get ns

# Deploymentリソースにツールキットコンポーネントがインストールされていることを確認
kubectl get deployment -n flux-system

NAME                      READY   UP-TO-DATE   AVAILABLE   AGE
helm-controller           1/1     1            1           2m24s
kustomize-controller      1/1     1            1           2m24s
notification-controller   1/1     1            1           2m24s
source-controller         1/1     1            1           2m24s

# さらにprivateリポジトリとして「flux-repo」が生成される
```

5. リポジトリの修正(@flux-repository)

```bash
# /eks-cluster/環境/配下にGitRepositoryリソースを作成するマニフェストを自動的に生成

# staging環境
mkdir eks-cluster/staging
flux create source git eks-sample-app-stg \
    --url=https://github.com/${GITHUB_USER}/learning-eks \
    --branch=staging \
    --interval=30s \
    --export >  ./eks-cluster/staging/app-source.yaml

# production環境
mkdir eks-cluster/production
flux create source git eks-sample-app-prd \
    --url=https://github.com/${GITHUB_USER}/learning-eks \
    --branch=master \
    --interval=30s \
    --export > ./eks-cluster/production/app-source.yaml

```

6. kustomization リソースを作成

```bash
# staging 環境
flux create kustomization eks-sample-app-stg \
    --target-namespace=staging \
    --source=eks-sample-app-stg \
    --path="./flux_sample/sample-app/kustomize/overlays/staging" \
    --prune=true \
    --interval=1m \
    --export > ./eks-cluster/staging/app-sync.yaml

# production環境
flux create kustomization eks-sample-app-prd \
    --target-namespace=production \
    --source=eks-sample-app-prd \
    --path="./flux_sample/sample-app/kustomize/overlays/production" \
    --prune=true \
    --interval=1m \
    --export > ./eks-cluster/production/app-sync.yaml
```

7. flux-repository を push して Cluster を deploy する

```bash
# kustomizationリソースのステータスを確認
watch flux get kustomizations

# Ingressリソースの情報を取得
kubectl get ing -n staging

```

8. ソースコードを変更して push すると、Cluster 内に変更が瞬時に反映される、さらに、再構築された image が`{環境}/kustomization.yaml`に反映される
9. 逆に kubectl コマンドでイメージを書き換えると、Flux によって Github のマニフェストに強制的にロールバックされる

```bash
# imageの変更
kubectl set image -n staging deployment/deployment app=nginx:1.21

# 確認
kubectl get deployment -n staging deployment -o yaml | grep image:

```

10. clean up

```bash
flux uninstall --namespace=flux-system

kubectl delete -k kustomize/overlays/staging
kubectl delete -k kustomize/overlays/production

terraform destroy -auto-approve
```
