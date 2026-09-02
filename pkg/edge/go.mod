module github.com/KlyneChrysler/datacat/pkg/edge

go 1.26

require (
	github.com/KlyneChrysler/datacat/pkg/events v0.0.0
	github.com/KlyneChrysler/datacat/pkg/guard v0.0.0
	github.com/KlyneChrysler/datacat/pkg/hashx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/httpx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/obsx v0.0.0
)

replace (
	github.com/KlyneChrysler/datacat/pkg/events => ../events
	github.com/KlyneChrysler/datacat/pkg/guard => ../guard
	github.com/KlyneChrysler/datacat/pkg/hashx => ../hashx
	github.com/KlyneChrysler/datacat/pkg/httpx => ../httpx
	github.com/KlyneChrysler/datacat/pkg/obsx => ../obsx
)

require github.com/KlyneChrysler/datacat/pkg/policy v0.0.0

replace github.com/KlyneChrysler/datacat/pkg/policy => ../policy
