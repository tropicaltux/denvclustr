# Denvclustr Schema Generator

This utility generates a JSON schema for the denvclustr configuration. The schema can be used for validation and auto-completion in editors that support JSON schemas.

## Usage

The schema generator can be run from the project root with the following command:

```bash
go run cmd/denvclustr_schema/generate_schema.go [options]
```

### Options

- `-out <filename>`: Specifies the output file for the JSON schema. If not provided, the schema will be printed to stdout.

## Example

1. Generate the schema and print it to stdout:

```bash
go run cmd/denvclustr_schema/generate_schema.go
```

2. Generate the schema and save it to a file:

```bash
go run cmd/denvclustr_schema/generate_schema.go -out denvclustr-schema.json
```

## VS Code Integration

A VS Code task has been configured to automatically update the schema:

1. Press `Ctrl+Shift+P` (or `Cmd+Shift+P` on macOS)
2. Select "Tasks: Run Task"
3. Choose "Update denvclustr schema"

This will update the schema at `/path/to/repo/api/denvclustr/configuration-schema.json`.
