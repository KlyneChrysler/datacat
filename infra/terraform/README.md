# Terraform (AWS phase - priority 1 and 4)

Layout per docs/standards/infra.md:

```
modules/          reusable, environment-agnostic
  bootstrap/      budget alarm - created BEFORE any other AWS resource   [built]
  network/        VPC, 2 AZs, public/private subnets, single NAT         [built]
  dynamo/         tables, TTL, on-demand                                 [built, applied locally]
  storage/        S3 bucket, SSE, lifecycle, public-access block         [built]
  eks/            cluster, node groups, IRSA                             [AWS session]
  kafka/          MSK                                                    [AWS session]
  verdict-api/    Lambda + API Gateway                                   [AWS session]
envs/
  local/          DynamoDB Local target - applied and verified for $0
  staging/        real AWS shape - validated; apply awaits credentials
```

Rules (binding): envs differ only in tfvars and backend config; env roots
contain module calls only; remote state in S3 with locking (see
`envs/staging/backend.tf` for the state-bucket bootstrap); `plan` on PR,
`apply` from CI on merge only; `terraform destroy` between learning
sessions; the budget alarm is in the first apply.

## Local loop ($0, no accounts)

```bash
docker compose up -d dynamodb
cd infra/terraform/envs/local
terraform init && terraform apply     # creates local-datacat-decisions
terraform destroy                     # leave nothing behind
```

Known emulator gaps, all documented in-code: DynamoDB Local has no tagging
API and no TTL emulation (both enabled only in staging), and AWS provider
>= 6.13 hangs on its missing WarmThroughput field, so the local env pins
`< 6.13` (localstack/localstack#13140). LocalStack itself is not used - it
now requires an auth token even for community usage.

## First real AWS apply (when credentials exist)

```bash
aws s3 mb s3://<state-bucket> --region ap-southeast-1
cd infra/terraform/envs/staging
cp terraform.tfvars.example terraform.tfvars   # fill in; git-ignored
# uncomment backend.tf, then:
terraform init -backend-config="bucket=<state-bucket>"
terraform apply -target=module.bootstrap       # budget alarm FIRST
terraform apply
terraform destroy                              # between sessions
```
