module github.com/KlyneChrysler/datacat/services/edge-proxy

go 1.26

require (
	github.com/KlyneChrysler/datacat/pkg/events v0.0.0-00010101000000-000000000000
	github.com/KlyneChrysler/datacat/pkg/httpx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/kafkax v0.0.0-00010101000000-000000000000
	github.com/KlyneChrysler/datacat/pkg/obsx v0.0.0
	golang.org/x/sync v0.22.0
)

require (
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
