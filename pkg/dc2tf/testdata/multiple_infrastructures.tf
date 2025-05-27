provider "aws" {
  region = "us-west-2"
  alias  = "infrastructure1"
}

provider "aws" {
  region = "eu-west-1"
  alias  = "infrastructure2"
}

module "node1" {
  source        = "github.com/tropicaltux/terraform-devcontainers"
  name          = "node1"
  instance_type = "t3.micro"
  providers     = {
    aws = aws.infrastructure1
  }

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
  instance_type = "t3.small"
  providers     = {
    aws = aws.infrastructure2
  }

  devcontainers = [
    {
      id = "dev2"
      source = {
        url = "https://github.com/example/app"
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
    module = module.node1
  }
  
}

output "node2_output" {
  value     = {
    module = module.node2
  }
  
} 