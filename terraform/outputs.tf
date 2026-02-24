output "api_url" {
  value = "https://${var.subdomain}"
}

output "s3_bucket_name" {
  value = aws_s3_bucket.invoice_storage.id
}

output "frontend_bucket_name" {
  value = aws_s3_bucket.frontend_web.id
}
