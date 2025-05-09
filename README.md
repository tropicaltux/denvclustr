# denvclustr — Development Environments Cluster

🚧 The project is currently under active development and not ready for production use.

## CLI Tool

The denvclustr CLI tool allows you to generate Terraform HCL files from denvclustr files and deploy devcontainers.

### Installation

```bash
go install github.com/tropicaltux/denvclustr/cmd/denvclustr@latest
```

### Usage

The CLI provides the following commands:

1. Generate Terraform HCL from a denvclustr configuration file:

```bash
# Specify input and output files
denvclustr generate path/to/config.json -o terraform/output.tf

# Use default input file (denvclustr.json in current directory)
denvclustr generate

# Use specific input file with default output (./denvclustr.tf)
denvclustr generate path/to/config.json
```

2. Show deployment plan without applying changes:

```bash
# Specify input file
denvclustr deploy path/to/config.json --plan

# Use default input file (denvclustr.json in current directory)
denvclustr deploy --plan

# Specify custom working directory
denvclustr deploy path/to/config.json --plan -w custom-dir
```

The plan command will:
- Show a summary of infrastructure, nodes, and devcontainers to be deployed
- Use Terraform to generate a detailed plan showing all resource changes
- Display which resources will be created, updated, or deleted
- Store Terraform files in the specified working directory for inspection and reuse

3. Deploy devcontainers based on a denvclustr configuration file:

```bash
# Specify input file
denvclustr deploy path/to/config.json

# Use default input file (denvclustr.json in current directory)
denvclustr deploy

# Specify custom working directory
denvclustr deploy path/to/config.json -w custom-dir
```

The deploy command will:
- Convert the denvclustr configuration to Terraform HCL
- Display a plan of what will be deployed before execution
- Execute Terraform init and apply commands directly
- Deploy the resources defined in your configuration
- Display the outputs from the Terraform deployment in a formatted, easy-to-read structure
- Store Terraform files in the specified working directory for future reference
- Automatically clean up resources if deployment fails

#### Deployment Outputs

After successful deployment, the tool will display all outputs from Terraform in a structured format:

```
Deployment Outputs:
-------------------

primary_node_output:
  module: module.primary_node
  name: primary_node
  public_ip: 203.0.113.10
  instance_id: i-0abc123def456789

  Devcontainers:

  [0] ID: devcontainer1
      Source: https://github.com/user/repo
      Remote Access:
        OpenVSCode Server: https://devcontainer1.example.com/?tkn=abcdef123456
        SSH: ssh -p 2222 root@devcontainer1.example.com

  [1] ID: devcontainer2
      Source: https://github.com/user/another-repo
      Remote Access:
        SSH: ssh -p 2223 root@203.0.113.10
```

The outputs are formatted for readability, with special focus on devcontainer remote access information that makes it easy to connect to your deployed environments.

4. Destroy previously deployed resources:

```bash
# Destroy resources using default working directory
denvclustr destroy

# Specify input file (optional as state is loaded from working directory)
denvclustr destroy path/to/config.json

# Show destroy plan without destroying resources
denvclustr destroy --plan

# Specify custom working directory where resources were deployed
denvclustr destroy -w custom-dir
```

The destroy command will:
- Check if resources exist in the specified working directory
- Display a plan of what will be destroyed before execution
- Show a confirmation prompt before destroying resources
- Run Terraform destroy to remove all deployed resources
- Preserve the Terraform files in the working directory

### Command Options

#### Generate Command

- `-o, --output`: Specify the output Terraform file (default: `./denvclustr.tf`)
  - If not specified, the output will be written to `denvclustr.tf` in the current directory
  - If the output file already exists, it will be overwritten

#### Deploy Command

- `-p, --plan`: Show deployment plan without applying changes
- `-w, --working-dir`: Specify the working directory for Terraform operations (default: `output`)

#### Destroy Command

- `-p, --plan`: Show destroy plan without applying changes
- `-w, --working-dir`: Specify the working directory where resources were deployed (default: `output`)

### Default Files and Directories

- If no input file is specified, the tool will look for `denvclustr.json` in the current directory
- If no output file is specified for the generate command, the tool will create `denvclustr.tf` in the current directory
- By default, the `plan`, `deploy`, and `destroy` commands use the `./output` directory in your current working directory
- You can specify a custom working directory with the `-w` or `--working-dir` flag for all commands that use Terraform

### Help

For more information about the available commands, use:

```bash
denvclustr --help
denvclustr generate --help
denvclustr deploy --help
denvclustr destroy --help
```

## Requirements

- Terraform CLI must be installed and available in your PATH
