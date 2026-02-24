# 📦 Inventory Backend

REST API for a business inventory management system, built with Go and deployed on AWS Lambda via Terraform.

## 🚀 Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.21+ |
| HTTP Router | [chi](https://github.com/go-chi/chi) |
| Database | PostgreSQL (Supabase-compatible) |
| Auth | JWT (HS256) + bcrypt |
| Storage | AWS S3 |
| PDF Generation | Go `embed` + custom renderer |
| Deployment | AWS Lambda + API Gateway v2 |
| Infrastructure | Terraform |

## 📋 Features in this Repository

- Product catalog with categories and pagination
- Sales orders with PDF quote generation
- JWT authentication

## 📋 Other Features in this Repository full project

- Purchase orders with S3 invoice upload
- Stock movements & adjustments
- Customer & supplier management (with RUT/email validation)

---

## ⚙️ Local Setup

### Prerequisites
- Go 1.21+
- PostgreSQL database (or [Supabase](https://supabase.com) project)
- AWS credentials configured (`~/.aws/credentials` or environment)

### 1. Clone the repo
```bash
git clone https://github.com/benjoks/inventory.git
cd inventory
```

### 2. Configure environment variables
```bash
cp .env.example .env
# Edit .env and fill in your values
```

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Secret key for signing JWT tokens |
| `COMPANY_NAME` | Company name shown in generated PDFs |
| `S3_BUCKET_NAME` | AWS S3 bucket for invoice storage |
| `FRONTEND_URL` | Allowed CORS origin (empty = allow all) |
| `LAMBDA_FUNCTION_NAME` | Set when running in AWS Lambda mode |

### 3. Run locally
```bash
make run
# Server starts at http://localhost:8080
```

### 4. Run tests
```bash
make test
```

---

## 🏗️ Project Structure

```
.
├── cmd/api/main.go          # Entry point (Lambda or local)
├── internal/
│   ├── core/domain.go       # Domain models
│   ├── database/            # DB connection (pgx)
│   ├── repository/          # Data access layer
│   ├── service/             # Business logic + HTTP handlers
│   ├── middleware/          # Auth & logging middleware
│   └── utils/               # Auth (JWT/bcrypt), validation, response helpers
├── terraform/               # AWS infrastructure as code
├── tests/                   # Integration tests
└── Makefile                 # Build & dev commands
```

---

## ☁️ Deployment (AWS Lambda)

### Prerequisites
- AWS CLI configured with appropriate IAM permissions
- Terraform >= 1.0
- Domain name configured in Route 53

### 1. Configure Terraform variables
```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
```

### 2. Build the Lambda binary
```bash
make build-lambda
# Creates build/main.zip
```

### 3. Deploy infrastructure
```bash
cd terraform
terraform init
terraform plan
terraform apply
```

### Terraform Variables

| Variable | Description |
|----------|-------------|
| `db_url` | PostgreSQL connection string |
| `jwt_secret` | JWT signing secret |
| `company_name` | Company name for PDFs |
| `environment` | `prod` / `staging` |
| `region` | AWS region (default: `sa-east-1`) |
| `domain_name` | Your domain (e.g. `example.com`) |
| `subdomain` | API subdomain (e.g. `api.example.com`) |

---

## 🔒 Security

- All secrets are loaded from environment variables — never hardcoded
- `.env` and `terraform.tfvars` are in `.gitignore` and never committed
- JWT tokens validated on every protected route
- Passwords hashed with bcrypt (cost 14)
- CORS locked to `FRONTEND_URL` in production

---

## 📄 API Routes

All routes are prefixed with `/v1`.

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/v1/login` | Authenticate and receive JWT |

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v1/products` | List products (paginated, with search) |
| GET | `/v1/products/{id}` | Get product by ID |
| POST | `/v1/products` | Create product |
| PUT | `/v1/products/{id}` | Update product |
| DELETE | `/v1/products/{id}` | Delete product |

### This Repository includes: Orders, Categories
Standard CRUD endpoints follow the same pattern as above.

### Full project includes: Orders, Purchases, Movements, Customers, Suppliers, Categories
Standard CRUD endpoints follow the same pattern as above.

---

## 🛠️ Makefile Commands

```bash
make run            # Run locally (loads .env automatically)
make build          # Compile binary
make build-lambda   # Build & zip for AWS Lambda
make test           # Run tests
make tidy           # go mod tidy
make clean          # Remove build artifacts
```
