terraform {
  required_version = ">= 1.10"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Provider >= 6.13 polls DescribeTable's WarmThroughput.Status after
      # create; DynamoDB Local never returns that field, so the waiter loops
      # forever (localstack/localstack#13140). Staging is unpinned - real
      # AWS returns the field.
      version = ">= 6.0, < 6.13"
    }
  }
}

# DynamoDB Local: fake credentials, API calls redirected to the container
# from docker-compose (port 8000). This env exercises the dynamo module for
# $0 with no accounts (LocalStack now requires an auth token even for
# community usage, so it is not used here). S3/network/bootstrap modules are
# exercised in staging against real AWS.
provider "aws" {
  region     = "us-east-1"
  access_key = "local"
  secret_key = "local"

  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    dynamodb = "http://localhost:8000"
  }
}
