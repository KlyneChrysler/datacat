output "vpc_id" {
  value       = aws_vpc.this.id
  description = "VPC id"
}

output "public_subnet_ids" {
  value       = aws_subnet.public[*].id
  description = "Public subnet ids (load balancers)"
}

output "private_subnet_ids" {
  value       = aws_subnet.private[*].id
  description = "Private subnet ids (EKS nodes, brokers)"
}
