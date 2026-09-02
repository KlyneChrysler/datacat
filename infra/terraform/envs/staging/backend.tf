# Remote state lives in S3 with native lockfile locking. The state bucket is
# the one resource created BY HAND, once, before the first init (the
# chicken-and-egg of Terraform state):
#
#   aws s3 mb s3://<your-state-bucket> --region ap-southeast-1
#   terraform init -backend-config="bucket=<your-state-bucket>"
#
# Until the AWS session, state stays local: leave this block commented.
#
# terraform {
#   backend "s3" {
#     key          = "staging/datacat.tfstate"
#     region       = "ap-southeast-1"
#     use_lockfile = true
#   }
# }
