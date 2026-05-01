package explain

import (
	"strconv"
	"strings"
	"time"

	"helm-values-diff-explainer/internal/diff"
)

// Rule maps a path pattern (and optional condition) to a risk level and impact messages.
type Rule struct {
	PathContains []string
	Condition    func(c *diff.Change) bool
	Risk         string
	Impact       []string
}

// Explainer annotates a Change with risk and impact.
type Explainer interface {
	Explain(c *diff.Change)
}

type componentExplainer struct {
	rules []Rule
}

// NewExplainer returns an Explainer for the given component name.
func NewExplainer(component string) Explainer {
	var componentRules []Rule
	switch component {
	case "loki":
		componentRules = lokiRules
	case "mimir":
		componentRules = mimirRules
	case "tempo":
		componentRules = tempoRules
	case "alloy":
		componentRules = alloyRules
	case "prometheus":
		componentRules = prometheusRules
	case "grafana":
		componentRules = grafanaRules
	}
	return &componentExplainer{rules: append(componentRules, genericRules...)}
}

func (e *componentExplainer) Explain(c *diff.Change) {
	for _, rule := range e.rules {
		if matchesRule(c, rule) {
			c.Risk = rule.Risk
			c.Impact = rule.Impact
			return
		}
	}
	c.Risk = "none"
	c.Impact = nil
}

func matchesRule(c *diff.Change, rule Rule) bool {
	matched := false
	for _, kw := range rule.PathContains {
		if strings.Contains(c.Path, kw) {
			matched = true
			break
		}
	}
	if !matched {
		return false
	}
	if rule.Condition != nil {
		return rule.Condition(c)
	}
	return true
}

// ── value comparison helpers ──────────────────────────────────────────────────

func isIncreased(c *diff.Change) bool {
	if c.Type != diff.ChangeChanged {
		return false
	}
	if of, nf, ok := asFloats(c.OldValue, c.NewValue); ok {
		return nf > of
	}
	if od, nd, ok := asDurations(c.OldValue, c.NewValue); ok {
		return nd > od
	}
	return false
}

func isDecreased(c *diff.Change) bool {
	if c.Type != diff.ChangeChanged {
		return false
	}
	if of, nf, ok := asFloats(c.OldValue, c.NewValue); ok {
		return nf < of
	}
	if od, nd, ok := asDurations(c.OldValue, c.NewValue); ok {
		return nd < od
	}
	return false
}

func asFloats(a, b any) (float64, float64, bool) {
	fa, oka := toFloat64(a)
	fb, okb := toFloat64(b)
	return fa, fb, oka && okb
}

func asDurations(a, b any) (time.Duration, time.Duration, bool) {
	da, oka := toDuration(a)
	db, okb := toDuration(b)
	return da, db, oka && okb
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case string:
		f, err := strconv.ParseFloat(n, 64)
		return f, err == nil
	}
	return 0, false
}

func toDuration(v any) (time.Duration, bool) {
	s, ok := v.(string)
	if !ok {
		return 0, false
	}
	d, err := time.ParseDuration(s)
	return d, err == nil
}

func newValueIs(s string) func(*diff.Change) bool {
	return func(c *diff.Change) bool {
		v, ok := c.NewValue.(string)
		return ok && v == s
	}
}

func newValueBoolTrue(c *diff.Change) bool {
	v, ok := c.NewValue.(bool)
	return ok && v
}

func oldValueIs(s string) func(*diff.Change) bool {
	return func(c *diff.Change) bool {
		v, ok := c.OldValue.(string)
		return ok && v == s
	}
}
