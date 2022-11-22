# kustomize sample

Kustomize は kubeectl のバージョンが 1.14 以上であれば、kubectl に組み込み済み　　
`kubectl apply -k ~`のようにすぐ使用可能

- base
  共通リソースの作成  
  kustomization.yaml には`resources[]`に共通リソースのマニフェスト名の一覧を記載する

- overlays
  base ディレクトリ配下にあるリソースを元に、差分があるリソースのみを記述

- develop
  base の deployment リソースの`replicas`を 1 から 3 に

- production
  `replicas`を 5 に上書き、Nginx のコンテナイメージタグを`1.20`から`1.19`
