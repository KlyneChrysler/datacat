module github.com/KlyneChrysler/datacat/services/enforcement

go 1.26

require (
	github.com/KlyneChrysler/datacat/pkg/events v0.0.0-00010101000000-000000000000
	github.com/KlyneChrysler/datacat/pkg/httpx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/kafkax v0.0.0-00010101000000-000000000000
	github.com/KlyneChrysler/datacat/pkg/obsx v0.0.0
	github.com/aws/aws-sdk-go-v2 v1.45.1
	github.com/aws/aws-sdk-go-v2/config v1.33.2
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.21.2
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.66.0
	golang.org/x/sync v0.22.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.20.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.19.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.39.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.36.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.41.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.48.0 // indirect
	github.com/aws/smithy-go v1.28.1 // indirect
	github.com/klauspost/compress v1.18.7 // indirect
	github.com/pierrec/lz4/v4 v4.1.26 // indirect
	github.com/twmb/franz-go v1.21.6 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.13.1 // indirect
)

replace (
	github.com/KlyneChrysler/datacat/pkg/httpx => ../../pkg/httpx
	github.com/KlyneChrysler/datacat/pkg/obsx => ../../pkg/obsx
)

replace (
	github.com/KlyneChrysler/datacat/pkg/events => ../../pkg/events
	github.com/KlyneChrysler/datacat/pkg/kafkax => ../../pkg/kafkax
)
