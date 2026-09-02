variable "region" {
  type        = string
  description = "AWS region"
  default     = "ap-southeast-1"
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR block"
  default     = "10.20.0.0/16"
}

variable "monthly_budget_usd" {
  type        = number
  description = "Hard monthly spend expectation; alarms fire at 80/100%"
}

variable "budget_email" {
  type        = string
  description = "Email receiving budget alerts"
}

variable "archive_bucket_name" {
  type        = string
  description = "Globally-unique S3 bucket for the event archive"
}
