# Terraform (AWS phase — priority 1 and 4)

Layout per docs/standards/infra.md:

```
modules/          reusable, environment-agnostic
  bootstrap/      budget alarm — created BEFORE any other AWS resource
  network/        VPC, subnets, endpoints
  eks/            cluster, node groups, IRSA
  kafka/          MSK or Redpanda-on-EKS
  dynamo/         tables, TTL
  storage/        S3 buckets, lifecycle
  verdict-api/    Lambda + API Gateway
  agent-verification/  Step Functions + IAM
envs/
  local/          LocalStack targets
  staging/
  prod/           backend.tf + main.tf (module calls only) + terraform.tfvars
```

Rules (binding): envs differ only in tfvars and backend config; remote state in
S3 with locking; `plan` on PR, `apply` from CI on merge only; every service
gets its own least-privilege IAM role via IRSA; `terraform destroy` between
learning sessions; the budget alarm module is applied first.

Modules are implemented in the Terraform phase, after the local stack
(Docker Compose + kind) is exercised end to end.
