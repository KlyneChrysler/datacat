# One table per module call. On-demand capacity by default: a learning
# environment has no baseline load worth provisioning for. DynamoDB Local
# cannot emulate on-demand's WarmThroughput status (the provider's waiter
# hangs on it), so the local env overrides to PROVISIONED.

resource "aws_dynamodb_table" "this" {
  name           = var.table_name
  billing_mode   = var.billing_mode
  hash_key       = var.hash_key
  read_capacity  = var.billing_mode == "PROVISIONED" ? 1 : null
  write_capacity = var.billing_mode == "PROVISIONED" ? 1 : null

  attribute {
    name = var.hash_key
    type = "S"
  }

  dynamic "ttl" {
    for_each = var.ttl_attribute == "" ? [] : [var.ttl_attribute]
    content {
      attribute_name = ttl.value
      enabled        = true
    }
  }

  tags = var.tags
}
