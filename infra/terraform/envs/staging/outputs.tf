output "vpc_id" {
  value       = module.network.vpc_id
  description = "Staging VPC id (EKS module consumes this)"
}

output "private_subnet_ids" {
  value       = module.network.private_subnet_ids
  description = "Private subnets for EKS nodes and brokers"
}

output "decisions_table" {
  value       = module.decisions_table.table_name
  description = "DynamoDB table for the enforcement DecisionStore adapter"
}
