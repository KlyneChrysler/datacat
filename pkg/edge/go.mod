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

require github.com/KlyneChrysler/datacat/pkg/agentsig v0.0.0

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/KlyneChrysler/datacat/pkg/agentsig => ../../pkg/agentsig
