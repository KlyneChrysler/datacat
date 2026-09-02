module github.com/KlyneChrysler/datacat/services/enforcement

go 1.26

require (
	github.com/KlyneChrysler/datacat/pkg/httpx v0.0.0
	github.com/KlyneChrysler/datacat/pkg/obsx v0.0.0
)

replace (
	github.com/KlyneChrysler/datacat/pkg/httpx => ../../pkg/httpx
	github.com/KlyneChrysler/datacat/pkg/obsx => ../../pkg/obsx
)
