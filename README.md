# uswds-gsx-registry

[USWDS](https://designsystem.digital.gov)-flavored components for
[gsx](https://github.com/gsxhq/gsx), distributed
[gsxui](https://github.com/gsxhq/gsxui)-style — copy-in, type-checked,
server-rendered, zero dependency on the `uswds` npm package.

Components carry real `usa-*` classes; the styles ship as plain CSS written
against `--usa-*` design-token custom properties. Every color reads a
shadcn-style semantic token first and falls back to the authentic USWDS
primitive, so a bare project renders stock USWDS and a themed gsxui project
(including dark-mode re-grades) restyles the components automatically.

## Install

In a gsxui-initialized Go module (the CLI reads your `gsxui.json`):

    go run github.com/gavmor/uswds-gsx-registry/cmd/uswds-gsx@latest add button input alert
    go run github.com/gavmor/uswds-gsx-registry/cmd/uswds-gsx@latest list

This vendors `.gsx` sources into `<ui>/uswds/` (their own `uswds` package —
no collisions with your gsxui components), the two stylesheets into
`<cssdir>/uswds/`, and runs `gsx generate`. Then add to the top of your CSS
entry point:

    @import "./uswds/uswds-tokens.css";
    @import "./uswds/uswds-components.css";

and use them beside your gsxui components:

    import "yourmodule/ui/uswds"

    <uswds.Button type="submit">Sign</uswds.Button>
    <uswds.Button variant="outline" href="/about">About</uswds.Button>
    <uswds.Input name="email" type="email"/>
    <uswds.Alert variant="warning">
        <uswds.AlertHeading>Heads up</uswds.AlertHeading>
        <uswds.AlertText>This form closes Friday.</uswds.AlertText>
    </uswds.Alert>

## Components

| Component | USWDS anatomy | Variants |
|---|---|---|
| `Button` | `usa-button` | `""` (primary), `outline`, `base`, `secondary`*, `big`, `unstyled` |
| `Input` | `usa-input` | — (`aria-invalid` styles the error state) |
| `Alert` (+`AlertHeading`, `AlertText`) | `usa-alert`, no-icon anatomy | `info` (default), `success`, `warning`, `error` |

\* USWDS "secondary" is the **red** secondary palette, not a quiet gray —
use `base` or `outline` for a subdued button.

## Design contract

- **You own the code.** `add` copies source into your tree; edit freely.
  A locally modified file is never overwritten without `--overwrite`.
- **Cascade position:** all component CSS sits in `@layer components`.
  Tailwind v4 orders `theme < base < components < utilities`, so your
  utility classes win per property — same contract as gsxui's components.
- **Tokens:** `uswds-tokens.css` is the canonical `--usa-*` primitive set
  (verbatim USWDS values). If your project already defines them, import
  exactly one source.
- **No JS.** Nothing here needs behavior; gsxui's runtime is untouched.

## Development

`.gsx` sources live in `ui/`; committed `.x.go` files keep the module
`go build`-able (consumers regenerate with their own gsx toolchain and
class merger on `add`). Regenerate here with `go tool gsx generate`.
