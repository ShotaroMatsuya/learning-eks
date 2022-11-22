variable "repository_name" {
  description = "GitHubのリポジトリ名"
  type        = string
}

variable "github_user" {
  description = "GitHubのユーザー名"
  type        = string
}

variable "aws_account_id" {
  description = "AWS ACCOUNT ID"
  type        = string
}

variable "my_ip_adress" {
  description = "ALBにアクセスする際の使用中のIPアドレス"
  type        = string
}
