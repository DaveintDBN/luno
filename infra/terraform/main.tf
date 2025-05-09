terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

resource "aws_key_pair" "deployer" {
  key_name   = "luno-deployer"
  public_key = file(var.public_key_path)
}

resource "aws_security_group" "luno_sg" {
  name_prefix = "luno-sg-"
  description = "Security group for Luno Bot"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 8081
    to_port     = 8081
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 9091
    to_port     = 9091
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 3002
    to_port     = 3002
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "luno" {
  ami                    = var.ami_id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.deployer.key_name
  vpc_security_group_ids = [aws_security_group.luno_sg.id]

  user_data = templatefile("${path.module}/user_data.sh", {
    repo_url           = var.repo_url
    api_key_id         = var.api_key_id
    api_key_secret     = var.api_key_secret
    dockerhub_username = var.dockerhub_username
    dockerhub_token    = var.dockerhub_token
  })

  tags = {
    Name = "luno-bot"
  }
}

output "instance_public_ip" {
  description = "Public IP of the Luno Bot EC2 instance"
  value       = aws_instance.luno.public_ip
}
