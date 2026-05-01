package report

import (
	"encoding/json"

	"helm-values-diff-explainer/internal/diff"
)

type jsonReport struct {
	OldFile   string        `json:"old_file"`
	NewFile   string        `json:"new_file"`
	Component string        `json:"component"`
	Summary   jsonSummary   `json:"summary"`
	Changes   []diff.Change `json:"changes"`
}

type jsonSummary struct {
	Added      int `json:"added"`
	Removed    int `json:"removed"`
	Changed    int `json:"changed"`
	HighRisk   int `json:"high_risk"`
	MediumRisk int `json:"medium_risk"`
	LowRisk    int `json:"low_risk"`
}

// RenderJSON serialises the report to indented JSON.
func RenderJSON(r *Report) (string, error) {
	s := r.Summarize()
	jr := jsonReport{
		OldFile:   r.OldFile,
		NewFile:   r.NewFile,
		Component: r.Component,
		Summary: jsonSummary{
			Added:      s.Added,
			Removed:    s.Removed,
			Changed:    s.Changed,
			HighRisk:   s.HighRisk,
			MediumRisk: s.MedRisk,
			LowRisk:    s.LowRisk,
		},
		Changes: r.Changes,
	}
	data, err := json.MarshalIndent(jr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}
