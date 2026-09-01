package uswds

import (
	"context"
	"strings"
	"testing"

	"github.com/gsxhq/gsx"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var b strings.Builder
	if err := n.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestButtonVariants(t *testing.T) {
	for variant, wantClass := range map[string]string{
		"":          `class="usa-button"`,
		"outline":   "usa-button--outline",
		"base":      "usa-button--base",
		"secondary": "usa-button--secondary",
		"big":       "usa-button--big",
		"unstyled":  "usa-button--unstyled",
	} {
		html := render(t, Button(variant, "", false, gsx.Text("Go"), nil))
		if !strings.Contains(html, wantClass) {
			t.Errorf("variant %q: %s missing %q", variant, html, wantClass)
		}
		if !strings.Contains(html, "<button") || !strings.Contains(html, `type="button"`) {
			t.Errorf("variant %q should render a type=button <button>: %s", variant, html)
		}
	}
}

func TestButtonHrefRendersAnchor(t *testing.T) {
	html := render(t, Button("outline", "/about", false, gsx.Text("About"), nil))
	if !strings.Contains(html, "<a") || !strings.Contains(html, `href="/about"`) {
		t.Errorf("href button should render an anchor: %s", html)
	}
	if strings.Contains(html, "<button") {
		t.Errorf("href button must not also render a <button>: %s", html)
	}
}

func TestButtonDisabledRendersRealDisabledButton(t *testing.T) {
	html := render(t, Button("", "/ignored", true, gsx.Text("Nope"), nil))
	if !strings.Contains(html, "<button") || !strings.Contains(html, "disabled") {
		t.Errorf("disabled button should render a disabled <button> even with href: %s", html)
	}
}

func TestInput(t *testing.T) {
	html := render(t, Input(gsx.AttrMap{"name": "email", "type": "email"}.ToAttrs()))
	for _, want := range []string{`class="usa-input"`, `name="email"`, `type="email"`} {
		if !strings.Contains(html, want) {
			t.Errorf("input missing %q: %s", want, html)
		}
	}
}

func TestAlertFamilies(t *testing.T) {
	for variant, want := range map[string]string{
		"":        "usa-alert--info",
		"info":    "usa-alert--info",
		"success": "usa-alert--success",
		"warning": "usa-alert--warning",
		"error":   "usa-alert--error",
	} {
		html := render(t, Alert(variant, AlertText(gsx.Text("hi"), nil), nil))
		for _, cls := range []string{"usa-alert", want, `role="alert"`, "usa-alert__body", "usa-alert__text"} {
			if !strings.Contains(html, cls) {
				t.Errorf("alert %q missing %q: %s", variant, cls, html)
			}
		}
	}
}

func TestTableVariants(t *testing.T) {
	for variant, wantClass := range map[string]string{
		"":           `class="usa-table"`,
		"striped":    "usa-table--striped",
		"borderless": "usa-table--borderless",
	} {
		html := render(t, Table(variant, false, false, gsx.Text("cells"), nil))
		if !strings.Contains(html, wantClass) {
			t.Errorf("variant %q: %s missing %q", variant, html, wantClass)
		}
		if strings.Contains(html, "usa-table-container--scrollable") {
			t.Errorf("non-scrollable table must not render the container: %s", html)
		}
	}
}

func TestTableCompact(t *testing.T) {
	html := render(t, Table("striped", true, false, gsx.Text("cells"), nil))
	for _, want := range []string{"usa-table--striped", "usa-table--compact"} {
		if !strings.Contains(html, want) {
			t.Errorf("compact striped table missing %q: %s", want, html)
		}
	}
}

func TestTableScrollableContainer(t *testing.T) {
	html := render(t, Table("", false, true, gsx.Text("cells"), gsx.AttrMap{"aria-label": "Wide data"}.ToAttrs()))
	for _, want := range []string{
		"usa-table-container--scrollable",
		`tabindex="0"`,
		`role="region"`,
		`aria-label="Wide data"`,
		`class="usa-table"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("scrollable table missing %q: %s", want, html)
		}
	}
}

func TestTag(t *testing.T) {
	html := render(t, Tag(false, gsx.Text("Winner"), nil))
	if !strings.Contains(html, `class="usa-tag"`) || !strings.Contains(html, "Winner") {
		t.Errorf("tag missing anatomy: %s", html)
	}
	if strings.Contains(html, "usa-tag--big") {
		t.Errorf("small tag must not carry the big modifier: %s", html)
	}
	if big := render(t, Tag(true, gsx.Text("New"), nil)); !strings.Contains(big, "usa-tag--big") {
		t.Errorf("big tag missing usa-tag--big: %s", big)
	}
}

func TestRankedList(t *testing.T) {
	items := []RankedItem{
		{Body: gsx.Text("Alpha"), Status: "success", Trailing: Tag(false, gsx.Text("Funded"), nil)},
		{Body: gsx.Text("Beta"), Status: "success"},
		{Body: gsx.Text("Gamma"), Muted: true},
	}
	html := render(t, RankedList(items, 2, "budget line — 200 hours", nil))
	for _, want := range []string{
		"usa-ranked-list",
		"usa-ranked-list__item--success",
		"usa-ranked-list__item--muted",
		"usa-ranked-list__waterline",
		"budget line — 200 hours",
		">1.<", ">2.<", ">3.<",
		"usa-tag",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("ranked list missing %q: %s", want, html)
		}
	}
	// The waterline sits above item 3: Beta renders before the rule, Gamma after.
	if strings.Index(html, "Beta") > strings.Index(html, "usa-ranked-list__waterline") ||
		strings.Index(html, "usa-ranked-list__waterline") > strings.Index(html, "Gamma") {
		t.Errorf("waterline not between items 2 and 3: %s", html)
	}
}

func TestRankedListWithoutWaterline(t *testing.T) {
	items := []RankedItem{{Body: gsx.Text("Alpha")}, {Body: gsx.Text("Beta")}}
	for _, waterline := range []int{-1, len(items)} {
		html := render(t, RankedList(items, waterline, "unused", nil))
		if strings.Contains(html, "usa-ranked-list__waterline") {
			t.Errorf("waterline %d should render no rule: %s", waterline, html)
		}
	}
}

func TestRankedListWaterlineAboveFirstItem(t *testing.T) {
	items := []RankedItem{{Body: gsx.Text("Alpha"), Muted: true}}
	html := render(t, RankedList(items, 0, "budget line", nil))
	if strings.Index(html, "usa-ranked-list__waterline") > strings.Index(html, "Alpha") {
		t.Errorf("waterline 0 should render before the first item: %s", html)
	}
}

func TestIdenticonDeterministic(t *testing.T) {
	a := render(t, Identicon("seed-1", nil))
	if b := render(t, Identicon("seed-1", nil)); a != b {
		t.Errorf("same seed must render byte-identical markup:\n%s\n%s", a, b)
	}
	if c := render(t, Identicon("seed-2", nil)); a == c {
		t.Errorf("different seeds should draw different faces: %s", a)
	}
	for _, want := range []string{"usa-identicon", `aria-hidden="true"`, "<polygon", "data-hue="} {
		if !strings.Contains(a, want) {
			t.Errorf("identicon missing %q: %s", want, a)
		}
	}
}

func TestIdenticonTwoDistinctHues(t *testing.T) {
	for _, seed := range []string{"a", "b", "c", "voter-42", "long-opaque-id"} {
		if identiconHue(seed) == identiconSecondHue(seed) {
			t.Errorf("seed %q: primary and secondary hue collapsed to %d", seed, identiconHue(seed))
		}
	}
}

func TestModalRendersAnatomyOnANativeDialog(t *testing.T) {
	html := render(t, Modal("", gsx.Text("body"), nil))
	for _, want := range []string{
		"<dialog", `class="usa-modal"`, "usa-modal__content", "usa-modal__main",
		"usa-modal__close", `method="dialog"`, "body",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("modal missing %q: %s", want, html)
		}
	}
	if !strings.Contains(html, `aria-label="Close"`) {
		t.Errorf("modal's close button should default aria-label to Close: %s", html)
	}
}

func TestModalCloseLabelOverride(t *testing.T) {
	html := render(t, Modal("Dismiss this notice", gsx.Text("body"), nil))
	if !strings.Contains(html, `aria-label="Dismiss this notice"`) {
		t.Errorf("modal should render the given closeLabel: %s", html)
	}
}

func TestModalAttrsPassThroughToTheDialogElement(t *testing.T) {
	attrs := gsx.AttrMap{"id": "example-modal", "aria-labelledby": "example-modal-heading"}.ToAttrs()
	html := render(t, Modal("", gsx.Text("body"), attrs))
	for _, want := range []string{`id="example-modal"`, `aria-labelledby="example-modal-heading"`} {
		if !strings.Contains(html, want) {
			t.Errorf("modal missing passed-through attr %q: %s", want, html)
		}
	}
}

func TestModalHeadingAndFooter(t *testing.T) {
	html := render(t, ModalHeading(gsx.Text("Are you sure?"), nil))
	for _, want := range []string{"usa-modal__heading", "Are you sure?"} {
		if !strings.Contains(html, want) {
			t.Errorf("modal heading missing %q: %s", want, html)
		}
	}

	html = render(t, ModalFooter(gsx.Text("actions"), nil))
	for _, want := range []string{"usa-modal__footer", "actions"} {
		if !strings.Contains(html, want) {
			t.Errorf("modal footer missing %q: %s", want, html)
		}
	}
}
