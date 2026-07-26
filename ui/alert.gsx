package uswds

import "github.com/gsxhq/gsx"

// Alert is the USWDS alert (usa-alert), no-icon anatomy: a status-tinted
// panel with a thick status-colored left bar; the copy stays regular
// foreground ink (USWDS alerts never color their text — the tint and bar
// carry the status, keeping AA contrast on every variant). variant: "info"
// (default) | "success" | "warning" | "error". No icons are shipped —
// keeping the registry dependency-free — so the CSS renders the
// usa-alert--no-icon layout.
component Alert(variant string, children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="alert"
		role="alert"
		class={ "usa-alert", familyClass(variant) }
		{ attrs... }
	>
		<div class="usa-alert__body">
			{ children }
		</div>
	</div>
}

// AlertHeading is usa-alert__heading. It renders a div, not a fixed h4 —
// heading level belongs to the page's outline, not the component; wrap it
// in (or replace it with) the level your document needs.
component AlertHeading(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="alert-heading" class="usa-alert__heading" { attrs... }>
		{ children }
	</div>
}

component AlertText(children gsx.Node, attrs gsx.Attrs) {
	<p data-slot="alert-text" class="usa-alert__text" { attrs... }>
		{ children }
	</p>
}

func familyClass(variant string) string {
	switch variant {
	case "success":
		return "usa-alert--success"
	case "warning":
		return "usa-alert--warning"
	case "error":
		return "usa-alert--error"
	default:
		return "usa-alert--info"
	}
}
