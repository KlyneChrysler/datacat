variable "environment" {
  type        = string
  description = "Environment name for tagging"
}

variable "cidr_block" {
  type        = string
  description = "VPC CIDR, /16 recommended"

  validation {
    condition     = can(cidrhost(var.cidr_block, 0))
    error_message = "cidr_block must be a valid CIDR."
  }
}
