provider "aws" {
  region = "us-west-2"
  alias  = "infrastructure1"
}

module "node1" {
  source        = "github.com/tropicaltux/terraform-devcontainers"
  name          = "node1"
  instance_type = "t3.micro"
  provider      = "aws.infrastructure1"

  devcontainers = [
    {
      id = "dev1"
      source = {
        url = "https://github.com/example/repo"
      }
      remote_access = {}
    }
  ]

  public_ssh_key = {
    local_key_path = "~/.ssh/id_rsa.pub"
  }
}

module "node2" {
  source        = "github.com/tropicaltux/terraform-devcontainers"
  name          = "node2"
  instance_type = "t3.large"
  provider      = "aws.infrastructure1"

  devcontainers = [
    {
      id = "dev2"
      source = {
        url = "https://github.com/example/app"
        branch = "main"
      }
      remote_access = {}
    }
  ]

  public_ssh_key = {
    local_key_path = "~/.ssh/id_rsa.pub"
  }
}

output "node1_output" {
  value     = {
    module = "module.node1"
  }
  sensitive = true
}

output "node2_output" {
  value     = {
    module = "module.node2"
  }
  sensitive = true
} 