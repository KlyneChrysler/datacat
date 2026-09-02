output "bucket_name" {
  value       = aws_s3_bucket.this.bucket
  description = "Created bucket name"
}

output "bucket_arn" {
  value       = aws_s3_bucket.this.arn
  description = "Bucket ARN for least-privilege IAM policies"
}
