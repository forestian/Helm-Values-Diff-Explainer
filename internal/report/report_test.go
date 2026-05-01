package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"helm-values-diff-explainer/internal/diff"
	"helm-values-diff-explainer/internal/report"
)

func sampleReport() *report.Report {
	return &report.Report{
		OldFile:   "old.yaml",
		NewFile:   "new.yaml",
		Component: "loki",
		Changes: []diff.Change{
			{
				Type:     diff.ChangeChanged,
				Path:     "loki.schemaConfig.configs[0].schema",
				OldValue: "v11",
				NewValue: "v12",
				Risk:     "high",
				Impact: []string{
					"Loki schema configuration changed.",
					"Incorrect schema changes may break ingestion or querying.",
					"Review Loki migration requirements carefully.",
				},
			},
			{
				Type:     diff.ChangeChanged,
				Path:     "loki.limits_config.retention_period",
				OldValue: "168h",
				NewValue: "720h",
				Risk:     "medium",
				Impact: []string{
					"Loki log retention increased.",
					"Object storage usage may increase.",
				},
			},
			{
				Type:     diff.ChangeAdded,
				Path:     "extraEnv",
				NewValue: "some-value",
				Risk:     "low",
				Impact:   []string{"Environment configuration changed."},
			},
		},
	}
}

func TestSummarize(t *testing.T) {
	r := sampleReport()
	s := r.Summarize()
	if s.Changed != 2 {
		t.Errorf("expected Changed=2, got %d", s.Changed)
	}
	if s.Added != 1 {
		t.Errorf("expected Added=1, got %d", s.Added)
	}
	if s.Removed != 0 {
		t.Errorf("expected Removed=0, got %d", s.Removed)
	}
	if s.HighRisk != 1 {
		t.Errorf("expected HighRisk=1, got %d", s.HighRisk)
	}
	if s.MedRisk != 1 {
		t.Errorf("expected MedRisk=1, got %d", s.MedRisk)
	}
	if s.LowRisk != 1 {
		t.Errorf("expected LowRisk=1, got %d", s.LowRisk)
	}
}

func TestRenderText(t *testing.T) {
	r := sampleReport()
	out := report.RenderText(r)

	if !strings.Contains(out, "Helm Values Diff Explainer") {
		t.Error("text output missing header")
	}
	if !strings.Contains(out, "old.yaml") {
		t.Error("text output missing old file")
	}
	if !strings.Contains(out, "[HIGH] [CHANGED]") {
		t.Error("text output missing HIGH CHANGED entry")
	}
	if !strings.Contains(out, "[MEDIUM] [CHANGED]") {
		t.Error("text output missing MEDIUM CHANGED entry")
	}
	if !strings.Contains(out, "Impact:") {
		t.Error("text output missing Impact section")
	}
}

func TestRenderTextNoChanges(t *testing.T) {
	r := &report.Report{OldFile: "a.yaml", NewFile: "b.yaml", Component: "generic"}
	out := report.RenderText(r)
	if !strings.Contains(out, "No changes detected") {
		t.Error("expected 'No changes detected' message")
	}
}

func TestRenderJSON(t *testing.T) {
	r := sampleReport()
	content, err := report.RenderJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, content)
	}

	if parsed["old_file"] != "old.yaml" {
		t.Errorf("expected old_file=old.yaml, got %v", parsed["old_file"])
	}
	if parsed["new_file"] != "new.yaml" {
		t.Errorf("expected new_file=new.yaml, got %v", parsed["new_file"])
	}
	if parsed["component"] != "loki" {
		t.Errorf("expected component=loki, got %v", parsed["component"])
	}

	summary, ok := parsed["summary"].(map[string]any)
	if !ok {
		t.Fatal("expected summary object in JSON")
	}
	if summary["high_risk"].(float64) != 1 {
		t.Errorf("expected high_risk=1, got %v", summary["high_risk"])
	}

	changes, ok := parsed["changes"].([]any)
	if !ok {
		t.Fatal("expected changes array in JSON")
	}
	if len(changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(changes))
	}
}

func TestRenderMarkdown(t *testing.T) {
	r := sampleReport()
	out := report.RenderMarkdown(r)

	if !strings.Contains(out, "# Helm Values Diff Explainer") {
		t.Error("markdown output missing H1 header")
	}
	if !strings.Contains(out, "## Summary") {
		t.Error("markdown output missing Summary section")
	}
	if !strings.Contains(out, "## Changes") {
		t.Error("markdown output missing Changes section")
	}
	if !strings.Contains(out, "HIGH / CHANGED") {
		t.Error("markdown output missing HIGH / CHANGED entry")
	}
	if !strings.Contains(out, "MEDIUM / CHANGED") {
		t.Error("markdown output missing MEDIUM / CHANGED entry")
	}
	if !strings.Contains(out, "**Impact:**") {
		t.Error("markdown output missing Impact section")
	}
}

func TestRenderMarkdownNoChanges(t *testing.T) {
	r := &report.Report{OldFile: "a.yaml", NewFile: "b.yaml", Component: "generic"}
	out := report.RenderMarkdown(r)
	if !strings.Contains(out, "No changes detected") {
		t.Error("expected 'No changes detected' message in markdown")
	}
}
