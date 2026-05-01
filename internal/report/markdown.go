package report

import (
	"fmt"
	"strings"
)

// RenderMarkdown produces a GitHub-friendly Markdown diff report.
func RenderMarkdown(r *Report) string {
	var b strings.Builder
	s := r.Summarize()

	b.WriteString("# Helm Values Diff Explainer\n\n")
	fmt.Fprintf(&b, "**Old:** `%s`  \n", r.OldFile)
	fmt.Fprintf(&b, "**New:** `%s`  \n", r.NewFile)
	fmt.Fprintf(&b, "**Component:** `%s`\n\n", r.Component)

	b.WriteString("## Summary\n\n")
	b.WriteString("| Type | Count |\n")
	b.WriteString("|---|---:|\n")
	fmt.Fprintf(&b, "| Added | %d |\n", s.Added)
	fmt.Fprintf(&b, "| Removed | %d |\n", s.Removed)
	fmt.Fprintf(&b, "| Changed | %d |\n", s.Changed)
	fmt.Fprintf(&b, "| High risk | %d |\n", s.HighRisk)
	fmt.Fprintf(&b, "| Medium risk | %d |\n", s.MedRisk)
	fmt.Fprintf(&b, "| Low risk | %d |\n", s.LowRisk)
	b.WriteString("\n")

	if len(r.Changes) == 0 {
		b.WriteString("No changes detected.\n")
		return b.String()
	}

	b.WriteString("## Changes\n\n")
	for _, c := range r.Changes {
		riskLabel := strings.ToUpper(c.Risk)
		if riskLabel == "" || riskLabel == "NONE" {
			riskLabel = "NONE"
		}
		fmt.Fprintf(&b, "### %s / %s / `%s`\n\n", riskLabel, string(c.Type), c.Path)

		if c.OldValue != nil {
			fmt.Fprintf(&b, "**Old:** `%v`  \n", c.OldValue)
		}
		if c.NewValue != nil {
			fmt.Fprintf(&b, "**New:** `%v`  \n", c.NewValue)
		}

		if len(c.Impact) > 0 {
			b.WriteString("\n**Impact:**\n\n")
			for _, line := range c.Impact {
				fmt.Fprintf(&b, "- %s\n", line)
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}
