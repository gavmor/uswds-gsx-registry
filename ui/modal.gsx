package uswds

import "github.com/gsxhq/gsx"

// Modal is USWDS's dialog (usa-modal), its documented anatomy —
// usa-modal__content > usa-modal__main, a heading, a footer, a corner
// close button — rendered onto a native <dialog> rather than USWDS's own
// usa-modal-wrapper/usa-modal-overlay pair and the JS behavior module
// that drives them. That JS exists to fake exactly what <dialog> already
// gives the platform for free: top-layer stacking above the rest of the
// page, an inert background, a real ::backdrop, and (via
// dialog.showModal()/.close(), or a plain <form method="dialog">) focus
// trapping and dismissal. This registry ships no JS at all — every other
// component here is server-rendered markup plus CSS, and Modal keeps
// that contract rather than becoming the one component that needs a
// script tag.
//
// A caller opens it with dialog.showModal() (an autoopen-style attribute
// plus a small event listener is the usual shape — see the README's
// Modal section) or wires open/close entirely through <form
// method="dialog"> buttons, which need no JS of the caller's either. The
// close button below already does the latter.
//
// closeLabel is the corner close button's aria-label ("Close" if empty).
// There's no icon glyph — USWDS's own close button is normally an SVG
// sprite reference, and this registry ships no icon assets (Alert's own
// "no icon" anatomy makes the same call) — just a plain "×" glyph.
component Modal(closeLabel string, children gsx.Node, attrs gsx.Attrs) {
	<dialog data-slot="modal" class="usa-modal" { attrs... }>
		<div class="usa-modal__content">
			<div class="usa-modal__main">
				{ children }
			</div>
			<form method="dialog">
				<button
					type="submit"
					class="usa-modal__close"
					aria-label={closeLabel |> default("Close")}
				>
					×
				</button>
			</form>
		</div>
	</dialog>
}

// ModalHeading is usa-modal__heading. Give it an id and point the Modal's
// own aria-labelledby at it — native <dialog> has no implicit label
// wiring the way usa-modal's own markup convention assumes, so the
// caller still owns that link (aria-describedby on the body text, the
// same way).
component ModalHeading(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-slot="modal-heading" class="usa-modal__heading" { attrs... }>
		{ children }
	</h2>
}

// ModalFooter is usa-modal__footer — the button, or button group, that
// closes or acts on the dialog.
component ModalFooter(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="modal-footer" class="usa-modal__footer" { attrs... }>
		{ children }
	</div>
}
