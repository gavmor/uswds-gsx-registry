package uswds

import (
	"strconv"

	"github.com/gsxhq/gsx"
)

// RankedItem is one rung of a RankedList. Body is the row's content (it
// flexes to fill the row); Trailing is an optional right-aligned slot for
// annotations — a cost, a Tag, an action — kept separate so the body can
// wrap without pushing them off-row. Status draws the row's left status
// bar: "" (none — the gutter stays, transparent, so rows column-align) or
// "success". Muted grays the whole row — an item that is out of the
// running but keeps its place (rank is a fact even when the item didn't
// make the cut).
type RankedItem struct {
	Body     gsx.Node
	Trailing gsx.Node
	Status   string
	Muted    bool
}

// RankedList is a numbered standing. The <ol> carries the order
// semantically (screen readers announce position); the printed rank span
// is aria-hidden typography, styled in place of the list marker so the
// number sits inside the row's flex line. Not a stock USWDS component
// (registry extension): the anatomy borrows the alert's status left bar
// and the list's quiet row dividers, all on the usa-* tokens.
//
// waterline draws a labeled cutoff rule above items[waterline] — the
// budget line, the pass mark, the top-N fold. Pass -1 (or len(items)) for
// no line. It is an index, not a property of the items, because the line
// and the row status are independent facts: a greedy budget walk can fund
// a cheap item below the line, and that row keeps its success bar there.
component RankedList(items []RankedItem, waterline int, waterlineLabel string, attrs gsx.Attrs) {
	<ol data-slot="ranked-list" class="usa-ranked-list" { attrs... }>
		{ for i, it := range items {
			{ if i == waterline {
				<li data-slot="ranked-list-waterline" class="usa-ranked-list__waterline" role="presentation">
					{waterlineLabel}
				</li>
			} }
			<li class={ "usa-ranked-list__item", rankedStatusClass(it) }>
				<span class="usa-ranked-list__rank" aria-hidden="true">{strconv.Itoa(i + 1)}.</span>
				<span class="usa-ranked-list__body">{it.Body}</span>
				{ if it.Trailing != nil {
					{it.Trailing}
				} }
			</li>
		} }
	</ol>
}

func rankedStatusClass(it RankedItem) string {
	cls := ""
	if it.Status == "success" {
		cls = "usa-ranked-list__item--success"
	}
	if it.Muted {
		if cls != "" {
			cls += " "
		}
		cls += "usa-ranked-list__item--muted"
	}
	return cls
}
