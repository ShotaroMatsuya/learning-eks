resource "helm_release" "aws-load-balancer-controller" {
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = "1.4.2"
  namespace  = "kube-system"
  depends_on = [
    kubernetes_service_account.aws_loadbalancer_controller,
    null_resource.kubeconfig
  ]

  set {
    name  = "clusterName"
    value = data.aws_eks_cluster.eks.name
  }

  set {
    name  = "serviceAccount.create"
    value = false
  }

  set {
    name  = "serviceAccount.name"
    value = "aws-load-balancer-controller"
  }

  set {
    name  = "image.repository"
    value = "602401143452.dkr.ecr.ap-northeast-1.amazonaws.com/amazon/aws-load-balancer-controller"
  }

  set {
    name  = "image.tag"
    value = "v2.4.2"
  }
}