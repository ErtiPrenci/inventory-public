# Addtional Provider for Virginia (Required for CloudFront)
provider "aws" {
  alias  = "virginia"
  region = "us-east-1"
}

# 1. Bucket for static files
resource "aws_s3_bucket" "frontend_web" {
  bucket = var.domain_name
}

# 2. Public access block (required for hosting web)
resource "aws_s3_bucket_public_access_block" "frontend_web" {
  bucket = aws_s3_bucket.frontend_web.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

# 3. Policy to allow public read access
resource "aws_s3_bucket_policy" "public_read" {
  bucket = aws_s3_bucket.frontend_web.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Sid       = "PublicReadGetObject"
      Effect    = "Allow"
      Principal = "*"
      Action    = "s3:GetObject"
      Resource  = "${aws_s3_bucket.frontend_web.arn}/*"
    }]
  })
  depends_on = [aws_s3_bucket_public_access_block.frontend_web]
}

# 4. Hosting Web Configuration
resource "aws_s3_bucket_website_configuration" "frontend_web" {
  bucket = aws_s3_bucket.frontend_web.id
  index_document { suffix = "index.html" }
  error_document { key = "index.html" } # Important for React Router
}

# 1. Get the certificate (assuming it already exists in us-east-1)
data "aws_acm_certificate" "cert" {
  domain   = var.domain_name
  statuses = ["ISSUED"]
  provider = aws.virginia
}

# 2. CloudFront Distribution
resource "aws_cloudfront_distribution" "frontend_cdn" {
  origin {
    domain_name = aws_s3_bucket.frontend_web.bucket_regional_domain_name
    origin_id   = "S3-Frontend"
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  aliases = [var.domain_name, "www.${var.domain_name}"]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-Frontend"

    viewer_protocol_policy = "redirect-to-https"

    forwarded_values {
      query_string = false
      cookies { forward = "none" }
    }
  }

  viewer_certificate {
    acm_certificate_arn      = data.aws_acm_certificate.cert.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  restrictions {
    geo_restriction { restriction_type = "none" }
  }
}
