package assets

import "embed"

var (
	//go:embed dist/*
	FS embed.FS
)
