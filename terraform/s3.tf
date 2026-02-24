resource "aws_s3_bucket" "invoice_storage" {
  bucket = "art-central-facturas-${var.environment}"
}

resource "aws_s3_bucket_cors_configuration" "invoice_cors" {
  bucket = aws_s3_bucket.invoice_storage.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "GET"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
