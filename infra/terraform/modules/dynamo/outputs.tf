output "table_name" {
  value       = aws_dynamodb_table.this.name
  description = "Created table name (services receive it via env)"
}

output "table_arn" {
  value       = aws_dynamodb_table.this.arn
  description = "Table ARN for least-privilege IAM policies"
}
