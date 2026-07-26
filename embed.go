// Package uswdsgsx embeds the USWDS-for-gsx component registry the
// uswds-gsx CLI vendors from. The embedded tree is the single source of
// truth: the CLI derives the component list from it at runtime — a
// component is a .gsx file basename under ui/, exactly gsxui's model.
package uswdsgsx

import "embed"

//go:embed ui css NOTICE.md
var Files embed.FS
