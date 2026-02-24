variable "region" { default = "sa-east-1" } # Sao Paulo para Chile
variable "environment" { default = "prod" }
variable "domain_name" { type = string }
variable "subdomain" { type = string }

# These are passed when running or by .tfvars file
variable "db_url" { type = string }
variable "jwt_secret" { type = string }
variable "company_name" {
  type    = string
  default = "Mi Empresa"
}
