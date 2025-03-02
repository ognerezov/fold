package configurator

import "embed"

var (
	//go:embed data/*
	dataOs embed.FS
)
