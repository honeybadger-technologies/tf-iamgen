output "app_bucket_id" {
  description = "ID of the application S3 bucket"
  value       = aws_s3_bucket.app_bucket.id
}

output "app_bucket_arn" {
  description = "ARN of the application S3 bucket"
  value       = aws_s3_bucket.app_bucket.arn
}

output "app_bucket_region" {
  description = "Region of the application S3 bucket"
  value       = aws_s3_bucket.app_bucket.region
}

output "log_bucket_id" {
  description = "ID of the logging S3 bucket"
  value       = aws_s3_bucket.log_bucket.id
}

output "log_bucket_arn" {
  description = "ARN of the logging S3 bucket"
  value       = aws_s3_bucket.log_bucket.arn
}

output "bucket_versioning_enabled" {
  description = "Whether versioning is enabled on the app bucket"
  value       = aws_s3_bucket_versioning.app_bucket_versioning.versioning_configuration[0].status == "Enabled"
}
