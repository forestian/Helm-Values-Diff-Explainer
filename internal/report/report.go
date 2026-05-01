package report

import "helm-values-diff-explainer/internal/diff"

// Report holds all data needed to render output in any format.
type Report struct {
	OldFile   string
	NewFile   string
	Component string
	Changes   []diff.Change
}

// Summary aggregates change counts by type and risk.
type Summary struct {
	Added    int
	Removed  int
	Changed  int
	HighRisk int
	MedRisk  int
	LowRisk  int
}

// Summarize computes counts over the report's change list.
func (r *Report) Summarize() Summary {
	var s Summary
	for _, c := range r.Changes {
		switch c.Type {
		case diff.ChangeAdded:
			s.Added++
		case diff.ChangeRemoved:
			s.Removed++
		case diff.ChangeChanged:
			s.Changed++
		}
		switch c.Risk {
		case "high":
			s.HighRisk++
		case "medium":
			s.MedRisk++
		case "low":
			s.LowRisk++
		}
	}
	return s
}
