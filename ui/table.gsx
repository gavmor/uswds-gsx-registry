package uswds

import "github.com/gsxhq/gsx"

// Table is the USWDS table (usa-table). Unlike gsxui's table, there are no
// per-element subcomponents: USWDS styles thead/th/td by descendant
// selector, so callers write plain <thead>/<tbody>/<tr>/<th>/<td> (and an
// optional <caption>) as children and the anatomy just applies. variant:
// "" (bordered default) | "striped" | "borderless". compact tightens the
// cell padding (usa-table--compact, stackable with either variant).
// scrollable wraps the table in usa-table-container--scrollable — the USWDS
// pattern for wide tables: a focusable overflow region so keyboard users
// can pan it (hence tabindex="0" and role="region"; give it an aria-label
// via attrs when the table has no caption).
component Table(variant string, compact bool, scrollable bool, children gsx.Node, attrs gsx.Attrs) {
	{ if scrollable {
		<div class="usa-table-container--scrollable" tabindex="0" role="region" { attrs... }>
			<table data-slot="table" class={ "usa-table", tableVariantClass(variant), tableCompactClass(compact) }>
				{ children }
			</table>
		</div>
	} else {
		<table data-slot="table" class={ "usa-table", tableVariantClass(variant), tableCompactClass(compact) } { attrs... }>
			{ children }
		</table>
	} }
}

func tableVariantClass(variant string) string {
	switch variant {
	case "striped":
		return "usa-table--striped"
	case "borderless":
		return "usa-table--borderless"
	default:
		return ""
	}
}

func tableCompactClass(compact bool) string {
	if compact {
		return "usa-table--compact"
	}
	return ""
}
