package explain_test

import (
	"testing"

	"helm-values-diff-explainer/internal/diff"
	"helm-values-diff-explainer/internal/explain"
)

func makeChange(ct diff.ChangeType, path string, old, new any) *diff.Change {
	return &diff.Change{
		Type:     ct,
		Path:     path,
		OldValue: old,
		NewValue: new,
	}
}

// ── generic rules ─────────────────────────────────────────────────────────────

func TestGenericReplicaIncreased(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "replicaCount", 2, 3)
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "low" {
		t.Errorf("expected risk=low, got %q", c.Risk)
	}
	if len(c.Impact) == 0 {
		t.Error("expected impact messages")
	}
}

func TestGenericReplicaDecreased(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "replicaCount", 3, 1)
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "medium" {
		t.Errorf("expected risk=medium, got %q", c.Risk)
	}
}

func TestGenericReplicaAdded(t *testing.T) {
	c := makeChange(diff.ChangeAdded, "replicaCount", nil, 2)
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "low" {
		t.Errorf("expected risk=low for added replica, got %q", c.Risk)
	}
}

func TestGenericResourceRequestDecreased(t *testing.T) {
	// Kubernetes CPU strings can't be compared numerically — fallback rule fires.
	c := makeChange(diff.ChangeChanged, "resources.requests.cpu", "200m", "100m")
	explain.NewExplainer("generic").Explain(c)
	if len(c.Impact) == 0 {
		t.Error("expected impact messages for resource change")
	}
}

func TestGenericResourceRequestIncreased(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "resources.requests.memory", 128, 256)
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "low" {
		t.Errorf("expected risk=low for increased memory, got %q", c.Risk)
	}
}

func TestGenericResourceRequestDecreaseNumeric(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "resources.requests.memory", 256, 128)
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "medium" {
		t.Errorf("expected risk=medium for decreased memory, got %q", c.Risk)
	}
}

func TestGenericImageTagLatest(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "image.tag", "1.0.0", "latest")
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "high" {
		t.Errorf("expected risk=high for latest tag, got %q", c.Risk)
	}
}

func TestGenericImageTagChange(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "image.tag", "1.0.0", "1.1.0")
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "medium" {
		t.Errorf("expected risk=medium for image tag change, got %q", c.Risk)
	}
}

func TestGenericServiceTypeToLoadBalancer(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "service.type", "ClusterIP", "LoadBalancer")
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "medium" {
		t.Errorf("expected risk=medium, got %q", c.Risk)
	}
}

func TestGenericPrivilegedTrue(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "securityContext.privileged", false, true)
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "high" {
		t.Errorf("expected risk=high for privileged=true, got %q", c.Risk)
	}
}

func TestGenericNoMatchGivesNone(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "someObscureKey.deep.value", "a", "b")
	explain.NewExplainer("generic").Explain(c)
	if c.Risk != "none" {
		t.Errorf("expected risk=none for unknown path, got %q", c.Risk)
	}
}

// ── loki rules ────────────────────────────────────────────────────────────────

func TestLokiRetentionIncreased(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "loki.limits_config.retention_period", "168h", "720h")
	explain.NewExplainer("loki").Explain(c)
	if c.Risk != "medium" {
		t.Errorf("expected risk=medium for retention increase, got %q", c.Risk)
	}
}

func TestLokiRetentionDecreased(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "loki.limits_config.retention_period", "720h", "168h")
	explain.NewExplainer("loki").Explain(c)
	if c.Risk != "high" {
		t.Errorf("expected risk=high for retention decrease, got %q", c.Risk)
	}
}

func TestLokiSchemaHighRisk(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "loki.schemaConfig.configs[0].schema", "v11", "v12")
	explain.NewExplainer("loki").Explain(c)
	if c.Risk != "high" {
		t.Errorf("expected risk=high for schema change, got %q", c.Risk)
	}
	if len(c.Impact) == 0 {
		t.Error("expected impact messages")
	}
}

func TestLokiStorageHighRisk(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "loki.storage.bucketNames.chunks", "old-bucket", "new-bucket")
	explain.NewExplainer("loki").Explain(c)
	if c.Risk != "high" {
		t.Errorf("expected risk=high for storage change, got %q", c.Risk)
	}
}

func TestLokiCompactorMediumRisk(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "loki.compactor.retention_enabled", false, true)
	explain.NewExplainer("loki").Explain(c)
	if c.Risk != "medium" {
		t.Errorf("expected risk=medium for compactor change, got %q", c.Risk)
	}
}

// ── mimir rules ───────────────────────────────────────────────────────────────

func TestMimirBlockStorageHighRisk(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "mimir.blocks_storage.backend", "s3", "gcs")
	explain.NewExplainer("mimir").Explain(c)
	if c.Risk != "high" {
		t.Errorf("expected risk=high for block storage change, got %q", c.Risk)
	}
}

// ── tempo rules ───────────────────────────────────────────────────────────────

func TestTempoStorageHighRisk(t *testing.T) {
	c := makeChange(diff.ChangeChanged, "tempo.storage.trace.backend", "local", "s3")
	explain.NewExplainer("tempo").Explain(c)
	if c.Risk != "high" {
		t.Errorf("expected risk=high for tempo storage change, got %q", c.Risk)
	}
}
