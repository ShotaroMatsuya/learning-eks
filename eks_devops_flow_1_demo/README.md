# EKS devops flow

## Part0. Create Cluster Via EKSCTL

```bash
$ eksctl create cluster --name ekscicd-test
```

---

## Part1. IAM and Setup

- create IAM role for

1. codePipeline - via cloudFormation
   (ekscicdIAMstack-CodePipelineServiceRole-19FPHVGS4SE5X)
2. codeBuild - via cloudFormation
   (ekscicdIAMstack-CodeBuildServiceRole-1UV013P5YBKW7å)
3. Kubectl role - via AWS Management console
   (arn:aws:iam::528163014577:role/EKSkubectlRole)

- Edit ConfigMap to give access to Kubectl role

---

## Part2. Pipeline

- Create CodePipeline pipeline
- Discuss relevant files
  ① Dockerfile
  ② Buildspec
- Run The pipeilne
