package risk_test

import (
	"testing"

	"helm-values-diff-explainer/internal/risk"
)

func TestParseRisk(t *testing.T) {
	tests := []struct {
		input string
		want  risk.Level
	}{
		{"none", risk.None},
		{"low", risk.Low},
		{"medium", risk.Medium},
		{"high", risk.High},
		{"unknown", risk.None},
		{"", risk.None},
	}
	for _, tt := range tests {
		got := risk.Parse(tt.input)
		if got != tt.want {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestShouldFailNone(t *testing.T) {
	if risk.ShouldFail("none", []string{"high", "medium", "low"}) {
		t.Error("threshold=none should never fail")
	}
}

func TestShouldFailHighThresholdWithHigh(t *testing.T) {
	if !risk.ShouldFail("high", []string{"low", "medium", "high"}) {
		t.Error("threshold=high should fail when high risk present")
	}
}

func TestShouldFailHighThresholdWithMedium(t *testing.T) {
	if risk.ShouldFail("high", []string{"low", "medium"}) {
		t.Error("threshold=high should not fail when only low/medium risk present")
	}
}

func TestShouldFailMediumThresholdWithMedium(t *testing.T) {
	if !risk.ShouldFail("medium", []string{"low", "medium"}) {
		t.Error("threshold=medium should fail when medium risk present")
	}
}

func TestShouldFailMediumThresholdWithLow(t *testing.T) {
	if risk.ShouldFail("medium", []string{"low"}) {
		t.Error("threshold=medium should not fail when only low risk present")
	}
}

func TestShouldFailLowThresholdWithLow(t *testing.T) {
	if !risk.ShouldFail("low", []string{"low"}) {
		t.Error("threshold=low should fail when low risk present")
	}
}

func TestShouldFailEmptyChanges(t *testing.T) {
	if risk.ShouldFail("high", []string{}) {
		t.Error("should not fail with no changes")
	}
}
