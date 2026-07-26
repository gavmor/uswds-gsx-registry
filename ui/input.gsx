package uswds

import "github.com/gsxhq/gsx"

// Input is the USWDS text input (usa-input): 1px base-dark border, square
// corners, units(5) = 40px tall, with the USWDS focus convention (4px solid
// focus-color outline at offset 0). type="text" is an overridable default —
// pass type="email" etc. at the call site. Void, childless element: the
// explicit { attrs... } spread opts it into fallthrough. Styling lives in
// css/uswds-components.css under @layer components, so caller utility
// classes still win per property.
component Input(attrs gsx.Attrs) {
	<input
		data-slot="input"
		type="text"
		class="usa-input"
		{ attrs... }
	/>
}
