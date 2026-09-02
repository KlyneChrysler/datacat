module github.com/KlyneChrysler/datacat/services/traffic-sim

go 1.26

require (
	github.com/KlyneChrysler/datacat/pkg/envx v0.0.0-00010101000000-000000000000
	github.com/KlyneChrysler/datacat/pkg/events v0.0.0-00010101000000-000000000000
	github.com/KlyneChrysler/datacat/pkg/obsx v0.0.0
	golang.org/x/sync v0.22.0
)

replace github.com/KlyneChrysler/datacat/pkg/obsx => ../../pkg/obsx

replace (
	github.com/KlyneChrysler/datacat/pkg/envx => ../../pkg/envx
	github.com/KlyneChrysler/datacat/pkg/events => ../../pkg/events
)
