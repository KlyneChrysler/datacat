variable "environment" {
  type        = string
  description = "Environment name for tagging"
}

variable "table_name" {
  type        = string
  description = "DynamoDB table name, prefixed by environment"

  validation {
    condition     = can(regex("^[a-z][a-z0-9.-]+$", var.table_name))
    error_message = "table_name must be lowercase kebab/dot case."
  }
}

variable "hash_key" {
  type        = string
  description = "Partition key attribute name (string type)"
}

variable "ttl_attribute" {
  type        = string
  description = "Attribute holding epoch-seconds expiry; empty disables TTL"
  default     = ""
}

variable "billing_mode" {
  type        = string
  description = "PAY_PER_REQUEST (default) or PROVISIONED (DynamoDB Local compatibility)"
  default     = "PAY_PER_REQUEST"

  validation {
    condition     = contains(["PAY_PER_REQUEST", "PROVISIONED"], var.billing_mode)
    error_message = "billing_mode must be PAY_PER_REQUEST or PROVISIONED."
  }
}

variable "tags" {
  type        = map(string)
  description = "Resource tags. Empty for DynamoDB Local, which does not support tagging"
  default     = {}
}
