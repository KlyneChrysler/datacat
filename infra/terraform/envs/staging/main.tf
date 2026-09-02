# Module calls only. Apply order note: bootstrap first is enforced socially,
# not mechanically - the budget alarm is in the same apply, and `terraform
# apply -target=module.bootstrap` exists for the very first run.

module "bootstrap" {
  source             = "../../modules/bootstrap"
  environment        = "staging"
  monthly_limit_usd  = var.monthly_budget_usd
  notification_email = var.budget_email
}

module "network" {
  source      = "../../modules/network"
  environment = "staging"
  cidr_block  = var.vpc_cidr
}

module "decisions_table" {
  source        = "../../modules/dynamo"
  environment   = "staging"
  table_name    = "staging-datacat-decisions"
  hash_key      = "session_id"
  ttl_attribute = "expires_at"

  tags = {
    environment = "staging"
    project     = "datacat"
  }
}

module "archive_bucket" {
  source         = "../../modules/storage"
  environment    = "staging"
  bucket_name    = var.archive_bucket_name
  retention_days = 30
}
