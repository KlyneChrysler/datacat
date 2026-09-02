variable "environment" {
  type        = string
  description = "Environment name prefixing the budget (staging, prod)"
}

variable "monthly_limit_usd" {
  type        = number
  description = "Monthly spend cap in USD; alerts fire at 80% actual and 100% forecasted"

  validation {
    condition     = var.monthly_limit_usd > 0
    error_message = "monthly_limit_usd must be positive."
  }
}

variable "notification_email" {
  type        = string
  description = "Email receiving budget alerts"

  validation {
    condition     = can(regex("^[^@\\s]+@[^@\\s]+$", var.notification_email))
    error_message = "notification_email must be a valid email address."
  }
}
