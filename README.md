# terraform-provider-digesttracker

A minimal Terraform provider that tracks a user-supplied `digest` value and automatically increments a `version` number whenever the digest changes.


## Inputs

| Field     | Type   | Description                      |
| --------- | ------ | -------------------------------- |
| `digest`  | string | Required. Any identifier string. |


## Outputs

| Field     | Type   | Description                      |
| --------- | ------ | -------------------------------- |
| `version` | int    | Computed. Increments on change.  |


## Example usage

Useful when using the secret_string_wo option in AWS Secrets Manager, e.g.:

```
terraform {
  required_providers {
    digesttracker = {
      source  = "local/digesttracker"
      version = "0.2.0"
    }
  }
}

provider "digesttracker" {}
provider "aws" {
  region = "us-west-2"
}

data "aws_kms_secrets" "example" {
  secret {
    name    = "example"
    payload = file("${path.module}/example_encrypted_secret.enc")
  }
}

resource "digesttracker_tracker" "example" {
  digest = filesha256("${path.module}/example_encrypted_secret.enc")
}

resource "aws_secretsmanager_secret" "example" {
  name = "example"
}

resource "aws_secretsmanager_secret_version" "example" {
  secret_id                = aws_secretsmanager_secret.example.id
  secret_string_wo         = data.aws_kms_secrets.example.plaintext["example"]
  secret_string_wo_version = digesttracker_tracker.example.version
}
```


## License

MIT
