# Environment roots contain module calls ONLY (the composition root
# pattern).

module "decisions_table" {
  source      = "../../modules/dynamo"
  environment = "local"
  table_name  = "local-datacat-decisions"
  hash_key    = "session_id"
  # Local divergences, forced by DynamoDB Local gaps (staging uses the real
  # defaults): no TTL emulation, and no tagging API (TagResource throws
  # UnknownOperationException, wedging the provider's post-create wait).
  ttl_attribute = ""
  tags          = {}
}
