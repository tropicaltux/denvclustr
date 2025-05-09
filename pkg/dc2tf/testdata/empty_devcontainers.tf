provider "aws" {
  region = "us-west-2"
  alias  = "infrastructure1"
}

module "node1" {
  source        = "github.com/tropicaltux/terraform-devcontainers"
  name          = "node1"
  instance_type = "t3.micro"
  providers     = {
    aws = aws.infrastructure1
  }

  devcontainers = []

  public_ssh_key = {
    local_key_path = "~/.ssh/id_rsa.pub"
  }
}

output "node1_output" {
  value     = {
    module = module.node1
  }
  sensitive = true
} 