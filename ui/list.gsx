package uswds

import "github.com/gsxhq/gsx"

// ListBlockItem is the standard USWDS block layout. It uses left borders
// and background colors to indicate state.
component ListBlockItem(variant string, number string, title string, desc string, modifiers []string, attrs gsx.Attrs) {
	<li
		class={
			"relative p-4 rounded border flex items-start gap-4 hover:shadow-md",
			blockBaseClass(variant, modifiers),
			dimmedClass(variant, modifiers),
		}
		{ attrs... }
	>
		<div class={ "flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm bg-background border", colorClass(variant, "border"), colorClass(variant, "text") }>
			{ if variant == "success" {
				<i data-lucide="check" class="w-4 h-4"></i>
			} else {
				{ number }
			} }
		</div>
		<div>
			<div class={ "font-bold", if variant == "disabled" { "line-through text-muted-foreground" } else { "text-foreground" } }>{ title }</div>
			<div class={ "text-sm mt-0.5", if variant == "disabled" { "line-through text-muted-foreground/70" } else { "text-muted-foreground" } }>{ desc }</div>
		</div>
	</li>
}

// ListMinimalRow is the tighter, table-like layout.
component ListMinimalRow(variant string, number string, title string, modifiers []string, attrs gsx.Attrs) {
	<li
		class={
			"flex items-center justify-between p-3 hover:bg-muted/50",
			minimalBaseClass(variant, modifiers),
			dimmedClass(variant, modifiers),
		}
		{ attrs... }
	>
		<div class="flex items-center gap-3">
			<span class={ "font-mono text-xs font-bold w-6 text-right", colorClass(variant, "text") }>{ number }.</span>
			<span class={ "text-sm", if variant == "disabled" { "line-through text-muted-foreground" } else { "font-medium text-foreground" } }>{ title }</span>
		</div>
		{ if variant != "" && variant != "disabled" {
			<i data-lucide={ iconForVariant(variant) } class={ "w-4 h-4", colorClass(variant, "text") }></i>
		} }
	</li>
}

// ListTaggedRow uses explicit tag components and icons.
component ListTaggedRow(variant string, number string, title string, desc string, modifiers []string, attrs gsx.Attrs) {
	<li
		class={
			"flex flex-col sm:flex-row sm:items-center justify-between p-4 gap-3 bg-background hover:bg-muted/50",
			taggedBaseClass(modifiers),
			dimmedClass(variant, modifiers),
		}
		{ attrs... }
	>
		<div>
			<div class={ "font-bold text-sm", if variant == "disabled" { "text-muted-foreground" } else { "text-foreground" } }>{ number }. { title }</div>
			<div class="text-xs text-muted-foreground mt-1">{ desc }</div>
		</div>
		<div class="flex-shrink-0 flex items-center">
			{ renderTag(variant) }
			{ if hasModifier(modifiers, "archived") {
				<span class="px-2 py-1 text-xs font-bold rounded bg-muted text-muted-foreground uppercase tracking-wider ml-1">Archived</span>
			} }
		</div>
	</li>
}

// ListProcessStep simulates a workflow step.
component ListProcessStep(variant string, number string, title string, desc string, isLast bool, modifiers []string, attrs gsx.Attrs) {
	<li
		class={
			"relative pl-12 py-2 mb-4",
			if hasModifier(modifiers, "dimmed") || hasModifier(modifiers, "archived") { "opacity-50" },
		}
		{ attrs... }
	>
		{ if !isLast {
			<div class={ "absolute left-4 top-8 bottom-[-1.5rem] w-0.5 -ml-px z-0", processLineClass(variant) }></div>
		} }
		<div class="absolute left-0 top-2">
			<div class={ "w-8 h-8 rounded-full border-2 flex items-center justify-center font-bold text-sm bg-background z-10 relative", processCircleClass(variant) }>
				{ if variant == "success" {
					<i data-lucide="check" class="w-4 h-4"></i>
				} else {
					{ number }
				} }
			</div>
		</div>
		<div class={ if variant == "disabled" { "text-muted-foreground" } }>
			<div class={ "font-bold text-sm", if variant == "info" { "text-info" } else { "text-foreground" } }>{ title }</div>
			<div class="text-xs text-muted-foreground mt-1">{ desc }</div>
			{ if variant == "error" {
				<div class="mt-2 text-xs text-destructive bg-destructive-muted p-2 rounded inline-flex items-center gap-1 border border-destructive">
					<i data-lucide="alert-circle" class="w-3 h-3"></i> Action required to proceed.
				</div>
			} }
		</div>
	</li>
}

// Helpers

func hasModifier(modifiers []string, m string) bool {
	for _, x := range modifiers {
		if x == m {
			return true
		}
	}
	return false
}

func dimmedClass(variant string, modifiers []string) string {
	if hasModifier(modifiers, "dimmed") || variant == "disabled" {
		return "opacity-60"
	}
	return ""
}

func blockBaseClass(variant string, modifiers []string) string {
	base := ""
	if hasModifier(modifiers, "archived") {
		base = "striped-archived border-muted-foreground "
	}
	switch variant {
	case "info":
		return base + "bg-info-muted border-l-4 border-info border-y-border border-r-border"
	case "success":
		return base + "bg-success-muted border-l-4 border-success border-y-border border-r-border"
	case "warning":
		return base + "bg-warning-muted border-l-4 border-warning border-y-border border-r-border"
	case "error":
		return base + "bg-destructive-muted border-l-4 border-destructive border-y-border border-r-border"
	case "disabled":
		return base + "bg-muted border-l-4 border-muted-foreground border-y-border border-r-border"
	default:
		if base != "" {
			return base + "border-border"
		}
		return "bg-background border-border"
	}
}

func minimalBaseClass(variant string, modifiers []string) string {
	base := ""
	if hasModifier(modifiers, "archived") {
		base = "striped-archived shadow-[inset_4px_0_0_0_var(--muted-foreground)] "
	}
	switch variant {
	case "info":
		return base + "bg-info-muted shadow-[inset_4px_0_0_0_var(--info)]"
	case "success":
		return base + "bg-success-muted shadow-[inset_4px_0_0_0_var(--success)]"
	case "warning":
		return base + "bg-warning-muted shadow-[inset_4px_0_0_0_var(--warning)]"
	case "error":
		return base + "bg-destructive-muted shadow-[inset_4px_0_0_0_var(--destructive)]"
	case "disabled":
		return base + "bg-muted shadow-[inset_4px_0_0_0_var(--muted-foreground)]"
	default:
		if base != "" {
			return base
		}
		return "bg-background"
	}
}

func taggedBaseClass(modifiers []string) string {
	if hasModifier(modifiers, "archived") {
		return "striped-archived"
	}
	return ""
}

func processLineClass(variant string) string {
	if variant == "success" {
		return "bg-success"
	}
	return "bg-border"
}

func processCircleClass(variant string) string {
	switch variant {
	case "success":
		return "bg-success border-success text-success-foreground"
	case "info":
		return "border-info text-info ring-4 ring-info-muted"
	case "warning":
		return "border-warning text-warning"
	case "error":
		return "border-destructive text-destructive"
	case "disabled":
		return "border-border text-muted-foreground bg-muted"
	default:
		return "border-border text-muted-foreground bg-muted"
	}
}

func colorClass(variant string, prop string) string {
	switch variant {
	case "info":
		return prop + "-info"
	case "success":
		return prop + "-success"
	case "warning":
		return prop + "-warning"
	case "error":
		return prop + "-destructive"
	case "disabled":
		return prop + "-muted-foreground"
	default:
		return prop + "-foreground"
	}
}

func iconForVariant(variant string) string {
	switch variant {
	case "info":
		return "chevron-right"
	case "success":
		return "check"
	case "warning":
		return "alert-circle"
	case "error":
		return "octagon-alert"
	default:
		return "circle"
	}
}

func renderTag(variant string) gsx.Node {
	switch variant {
	case "success":
		return <span class="px-2 py-1 text-xs font-bold rounded bg-success text-success-foreground uppercase tracking-wider">Completed</span>
	case "warning":
		return <span class="px-2 py-1 text-xs font-bold rounded bg-warning text-warning-foreground uppercase tracking-wider flex items-center gap-1"><i data-lucide="alert-triangle" class="w-3 h-3"></i> Draft</span>
	case "error":
		return <span class="px-2 py-1 text-xs font-bold rounded bg-destructive text-destructive-foreground uppercase tracking-wider flex items-center gap-1"><i data-lucide="octagon-alert" class="w-3 h-3"></i> Alert</span>
	case "info":
		return <span class="px-2 py-1 text-xs font-bold rounded bg-info text-info-foreground uppercase tracking-wider">Current</span>
	case "disabled":
		return <span class="px-2 py-1 text-xs font-bold rounded border border-muted-foreground text-muted-foreground uppercase tracking-wider">Disabled</span>
	default:
		return nil
	}
}
