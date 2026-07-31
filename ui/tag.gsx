package uswds

import "github.com/gsxhq/gsx"

// Tag is the USWDS tag (usa-tag): the small uppercase label USWDS uses to
// mark an item's state — "Winner", "Funded", "New". It is a label, not a
// control: a <span>, no interaction states. big renders usa-tag--big (the
// 15px reading-size variant for body copy; the default is deliberately
// small and quiet). Status stays a caller concern — USWDS tags are one
// neutral gray; tint per-site with a utility class if you must, but the
// stock look is the accessible default.
component Tag(big bool, children gsx.Node, attrs gsx.Attrs) {
	<span data-slot="tag" class={ "usa-tag", tagBigClass(big) } { attrs... }>
		{ children }
	</span>
}

func tagBigClass(big bool) string {
	if big {
		return "usa-tag--big"
	}
	return ""
}
