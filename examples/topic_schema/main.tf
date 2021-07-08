terraform {
  required_providers {
    protogit = {
      source = "curve-technology/protogit"
      version = "~> 0.1.0"
    }
  }
}


variable "git_password" {
  description = "Git password"
  type        = string
  sensitive   = true
}


provider "protogit" {
  url         = "github.com/curve-technology/terraform-provider-protogit"
  tag_version = "v0.1.0"
  password    = var.git_password
}


data "protogit_schemas" "schemas_collection" {
  entries {
    topic    = "topic1"
    section  = "value"
    filepath = "messaging/domain1/v1/event1.proto"
  }

  entries {
    topic    = "topic2"
    section  = "value"
    filepath = "messaging/domain2/v1/event1.proto"
  }
}


output "protogit_output" {
  value = data.protogit_schemas.schemas_collection
}