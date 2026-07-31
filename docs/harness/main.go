// Command harness renders each registry component in isolation to a
// self-contained HTML page, for the README screenshots in
// docs/screenshots/. Stock USWDS tokens only — no theme overrides —
// so the images show exactly what a bare `add` renders.
//
// Regenerate the pages with:
//
//	go run ./docs/harness -out /tmp/uswds-iso
//
// then screenshot each page at ~800px viewport width and copy the
// images over docs/screenshots/<component>.png.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	registry "github.com/gavmor/uswds-gsx-registry"
	uswds "github.com/gavmor/uswds-gsx-registry/ui"
	"github.com/gsxhq/gsx"
)

func main() {
	out := flag.String("out", "/tmp/uswds-iso", "directory to write per-component HTML pages into")
	flag.Parse()
	if err := os.MkdirAll(*out, 0o755); err != nil {
		panic(err)
	}

	tokens, err := registry.Files.ReadFile("css/uswds-tokens.css")
	if err != nil {
		panic(err)
	}
	comps, err := registry.Files.ReadFile("css/uswds-components.css")
	if err != nil {
		panic(err)
	}
	css := string(tokens) + string(comps)

	for name, node := range pages() {
		var body strings.Builder
		if err := node.Render(context.Background(), &body); err != nil {
			panic(err)
		}
		html := `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>` + name + `</title><style>
` + css + `
body { font-family: "Public Sans", system-ui, sans-serif; margin: 0; padding: 24px; background: #fff; color: #1c1d1f; }
</style></head><body>` + body.String() + `</body></html>`
		path := filepath.Join(*out, name+".html")
		if err := os.WriteFile(path, []byte(html), 0o644); err != nil {
			panic(err)
		}
		fmt.Println(path)
	}
}

// pages returns one representative, domain-free render per component,
// exercising the variants worth seeing rather than every permutation.
func pages() map[string]gsx.Node {
	stack := func(nodes ...gsx.Node) gsx.Node {
		var wrapped []gsx.Node
		for _, n := range nodes {
			wrapped = append(wrapped, gsx.Raw(`<div style="margin-bottom:12px">`), n, gsx.Raw(`</div>`))
		}
		return gsx.Fragment(wrapped...)
	}
	row := func(nodes ...gsx.Node) gsx.Node {
		var spaced []gsx.Node
		for i, n := range nodes {
			if i > 0 {
				spaced = append(spaced, gsx.Raw(" "))
			}
			spaced = append(spaced, n)
		}
		return gsx.Fragment(spaced...)
	}

	return map[string]gsx.Node{
		"button": row(
			uswds.Button("", "", false, gsx.Text("Primary"), nil),
			uswds.Button("outline", "", false, gsx.Text("Outline"), nil),
			uswds.Button("base", "", false, gsx.Text("Base"), nil),
			uswds.Button("secondary", "", false, gsx.Text("Secondary"), nil),
			uswds.Button("big", "", false, gsx.Text("Big"), nil),
			uswds.Button("unstyled", "", false, gsx.Text("Unstyled"), nil),
			uswds.Button("", "", true, gsx.Text("Disabled"), nil),
		),

		"input": stack(
			uswds.Input(gsx.AttrMap{"name": "email", "type": "email", "aria-label": "Email"}.ToAttrs()),
			uswds.Input(gsx.AttrMap{"name": "code", "aria-invalid": "true", "aria-label": "Code", "value": "not-a-code"}.ToAttrs()),
		),

		"alert": stack(
			uswds.Alert("info", gsx.Fragment(
				uswds.AlertHeading(gsx.Text("Information"), nil),
				uswds.AlertText(gsx.Text("The poll closes Friday at noon."), nil),
			), nil),
			uswds.Alert("success", gsx.Fragment(
				uswds.AlertHeading(gsx.Text("Success"), nil),
				uswds.AlertText(gsx.Text("Your ballot was recorded."), nil),
			), nil),
			uswds.Alert("warning", gsx.Fragment(
				uswds.AlertHeading(gsx.Text("Warning"), nil),
				uswds.AlertText(gsx.Text("Two options are tied."), nil),
			), nil),
			uswds.Alert("error", gsx.Fragment(
				uswds.AlertHeading(gsx.Text("Error"), nil),
				uswds.AlertText(gsx.Text("That code has already been used."), nil),
			), nil),
		),

		"table": stack(
			uswds.Table("striped", false, true, gsx.Raw(`
				<caption>Striped, scrollable</caption>
				<thead><tr><th scope="col">row over col</th><th scope="col">Alpha</th><th scope="col">Beta</th><th scope="col">Gamma</th></tr></thead>
				<tbody>
					<tr><th scope="row">Alpha</th><td>—</td><td>5</td><td>4</td></tr>
					<tr><th scope="row">Beta</th><td>2</td><td>—</td><td>6</td></tr>
					<tr><th scope="row">Gamma</th><td>3</td><td>1</td><td>—</td></tr>
				</tbody>`), gsx.AttrMap{"aria-label": "Pairwise preferences"}.ToAttrs()),
			uswds.Table("borderless", true, false, gsx.Raw(`
				<caption>Borderless, compact</caption>
				<thead><tr><th scope="col">Option</th><th scope="col">Cost</th></tr></thead>
				<tbody>
					<tr><th scope="row">New login flow</th><td>120 hours</td></tr>
					<tr><th scope="row">Dark mode</th><td>80 hours</td></tr>
				</tbody>`), nil),
		),

		"tag": row(
			uswds.Tag(false, gsx.Text("Winner"), nil),
			uswds.Tag(false, gsx.Text("Funded"), nil),
			uswds.Tag(true, gsx.Text("Big tag"), nil),
		),

		"ranked-list": uswds.RankedList([]uswds.RankedItem{
			{Body: gsx.Text("New login flow — 120 hours"), Status: "success", Trailing: uswds.Tag(false, gsx.Text("Funded"), nil)},
			{Body: gsx.Text("Dark mode — 80 hours"), Status: "success", Trailing: uswds.Tag(false, gsx.Text("Funded"), nil)},
			{Body: gsx.Text("Realtime sync — 300 hours"), Muted: true},
			{Body: gsx.Text("Emoji reactions — 40 hours"), Status: "success", Trailing: uswds.Tag(false, gsx.Text("Funded"), nil)},
			{Body: gsx.Text("Plugin system — 500 hours"), Muted: true},
		}, 2, "budget line — 240 hours", nil),

		"identicon": row(
			uswds.Identicon("alice", nil),
			uswds.Identicon("bob", nil),
			uswds.Identicon("carol", nil),
			uswds.Identicon("dave", nil),
			uswds.Identicon("erin", nil),
			uswds.Identicon("frank", nil),
			uswds.Identicon("grace", nil),
			uswds.Identicon("heidi", nil),
			uswds.Identicon("ivan", nil),
			uswds.Identicon("judy", nil),
			uswds.Identicon("mallory", nil),
			uswds.Identicon("niaj", nil),
		),
	}
}
