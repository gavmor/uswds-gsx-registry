package uswds

import "github.com/gsxhq/gsx"

// Button is the USWDS button (usa-button), typed for gsx. The variant
// vocabulary is USWDS's own modifier set — "" (primary), "outline", "base",
// "secondary", "big", "unstyled" — mapped to the corresponding
// usa-button--* class; unknown variants fall back to primary. A non-empty
// href on an enabled Button renders an <a> (same contract as gsxui's
// Button); disabled always renders a real disabled <button>. type="button"
// is an overridable default — pass type="submit" at the call site.
//
// Styling lives in css/uswds-components.css, keyed entirely off the
// usa-button classes; no Tailwind utilities are baked in, so callers can
// still adjust per-site with their own utility classes (the component CSS
// sits in @layer components, below Tailwind's utilities layer).
//
// Caution for shadcn/gsxui muscle memory: USWDS "secondary" is the RED
// secondary palette, not a quiet gray — for a subdued button use "base"
// (gray) or "outline".
component Button(variant string, href string, disabled bool, children gsx.Node, attrs gsx.Attrs) {
	{ if href != "" && !disabled {
		<a
			data-slot="button"
			data-variant={variant |> default("default")}
			href={href}
			class={ "usa-button", modifierClass(variant) }
			{ attrs... }
		>
			{ children }
		</a>
	} else {
		<button
			data-slot="button"
			data-variant={variant |> default("default")}
			type="button"
			class={ "usa-button", modifierClass(variant) }
			disabled={disabled}
			{ attrs... }
		>
			{ children }
		</button>
	} }
}

func modifierClass(variant string) string {
	switch variant {
	case "outline":
		return "usa-button--outline"
	case "base":
		return "usa-button--base"
	case "secondary":
		return "usa-button--secondary"
	case "big":
		return "usa-button--big"
	case "unstyled":
		return "usa-button--unstyled"
	default:
		return ""
	}
}
