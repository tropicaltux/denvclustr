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

  devcontainers = [
    {
      id = "dev1"
      source = {
        url = "https://github.com/example/repo"
        branch = "feature-branch"
        devcontainer_path = ".devcontainer"
        ssh_key = {
          ref = "github-key"
          src = "secrets_manager"
        }
      }
      remote_access = {
        openvscode_server = {
          port = 3000
        }
        ssh = {
          port = 2222
          public_ssh_key = {
            local_key_path = "~/.ssh/custom_key.pub"
          }
        }
      }
    }
  ]

  public_ssh_key = {
    local_key_path = "~/.ssh/id_rsa.pub"
  }

  dns = {
    high_level_domain = "example.com"
  }
}

output "node1_output" {
  value     = {
    module = module.node1
  }
  sensitive = true
} 