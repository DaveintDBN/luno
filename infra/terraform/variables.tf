variable "aws_region" {
  type        = string
  description = "AWS region to deploy resources"
  default     = "us-east-1"
}

variable "ami_id" {
  type        = string
  description = "AMI ID for the EC2 instance"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t3.medium"
}

variable "public_key_path" {
  type        = string
  description = "Path to SSH public key"
}

variable "repo_url" {
  type        = string
  description = "Git repo URL of the trading bot code"
}

variable "api_key_id" {
  type        = string
  description = "Luno API Key ID"
}

variable "api_key_secret" {
  type        = string
  description = "Luno API Key Secret"
  sensitive   = true
}

variable "dockerhub_username" {
  type        = string
  description = "Docker Hub username for pulling images"
}

variable "dockerhub_token" {
  type        = string
  description = "Docker Hub access token"
  sensitive   = true
}
