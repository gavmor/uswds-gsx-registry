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
