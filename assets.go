// Package assets exposes the compiled frontend as an embed.FS.
// This file lives at the repository root alongside webapp/ so that
// the //go:embed directive can reference webapp/dist.
package assets

import "embed"

// FS holds the entire webapp/dist directory, baked into the binary at compile time.
//
//go:embed all:webapp/dist
var FS embed.FS
