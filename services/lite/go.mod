module github.com/KlyneChrysler/datacat/services/lite

go 1.26

require (
	github.com/KlyneChrysler/datacat/pkg/edge v0.0.0
	github.com/KlyneChrysler/datacat/pkg/envx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/guard v0.0.0
	github.com/KlyneChrysler/datacat/pkg/httpx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/obsx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/policy v0.0.0
)

require (
	github.com/KlyneChrysler/datacat/pkg/events v0.0.0 // indirect
	github.com/KlyneChrysler/datacat/pkg/hashx v0.0.0 // indirect
)

replace (
	github.com/KlyneChrysler/datacat/pkg/edge => ../../pkg/edge
	github.com/KlyneChrysler/datacat/pkg/envx => ../../pkg/envx
	github.com/KlyneChrysler/datacat/pkg/events => ../../pkg/events
	github.com/KlyneChrysler/datacat/pkg/guard => ../../pkg/guard
	github.com/KlyneChrysler/datacat/pkg/hashx => ../../pkg/hashx
	github.com/KlyneChrysler/datacat/pkg/httpx => ../../pkg/httpx
	github.com/KlyneChrysler/datacat/pkg/obsx => ../../pkg/obsx
	github.com/KlyneChrysler/datacat/pkg/policy => ../../pkg/policy
)
