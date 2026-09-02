output "decisions_table" {
  value       = module.decisions_table.table_name
  description = "DynamoDB table for the enforcement DecisionStore adapter"
}
