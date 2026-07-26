# Third-party attribution

The visual design implemented by these components — anatomy, tokens, and
interaction conventions — is the **U.S. Web Design System**
(https://designsystem.digital.gov, public domain as a work of the United
States Government; CC0 1.0 Universal for non-government contributions). No
USWDS CSS or JS source is copied; the styles are written from the published
design tokens and component specifications.

The CLI's vendoring model (copy-in components you own, gsxui.json
discovery, the derived-registry-over-embed idea, and the writeVendored
contract) follows **gsxui** ([gsxhq/gsxui](https://github.com/gsxhq/gsxui),
MIT). Component API shapes (attrs fallthrough, variant props, data-slot
markers, href-renders-an-anchor) match gsxui's conventions so the two
component sets compose in one project.

Public Sans — the USWDS typeface these components assume for full fidelity
— is not bundled; vendor it in your project (OFL,
https://public-sans.digital.gov).
