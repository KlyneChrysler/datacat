variable "environment" {
  type        = string
  description = "Environment name for tagging"
}

variable "bucket_name" {
  type        = string
  description = "Globally-unique bucket name"

  validation {
    condition     = can(regex("^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$", var.bucket_name))
    error_message = "bucket_name must satisfy S3 naming rules."
  }
}

variable "retention_days" {
  type        = number
  description = "Days before archived objects expire"
  default     = 30

  validation {
    condition     = var.retention_days >= 1
    error_message = "retention_days must be at least 1."
  }
}
