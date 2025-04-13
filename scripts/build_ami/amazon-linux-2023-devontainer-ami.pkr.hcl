packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "region" {
  type    = string
  default = "eu-central-1"
}

variable "architecture" {
  type    = string
  default = "x86_64"
}

source "amazon-ebs" "amazon_linux_2023" {
  region                      = var.region
  instance_type               = "t2.micro"
  ami_name                    = "amazon-linux-2023-devcontainers-${var.architecture}-{{timestamp}}"
  ssh_username                = "ec2-user"
  associate_public_ip_address = true

  source_ami_filter {
    filters = {
      name                = "al2023-ami-2023*-kernel-6.1-${var.architecture}"
      virtualization-type = "hvm"
      root-device-type    = "ebs"
    }
    most_recent = true
    owners      = ["137112412989"]
  }

  tags = {
    "Name" = "Amazon Linux 2023 Dev Containers (${var.architecture})"
    # "Project"     = "Project-X"
    "Owner"       = "tropicaltux@proton.me"
    "CreatedBy"   = "Packer"
    "Application" = "Dev Container"
    "Version"     = "0.1.0"
  }
}

build {
  name    = "amazon-linux-2023-devcontainers-${var.architecture}"
  sources = ["source.amazon-ebs.amazon_linux_2023"]

  provisioner "shell" {
    inline = [
      "sudo dnf update -y",
      "sudo dnf install -y git docker nodejs",
      "sudo systemctl enable docker",
      "sudo usermod -aG docker ec2-user",
      "sudo npm install -g @devcontainers/cli"
    ]
  }
}
