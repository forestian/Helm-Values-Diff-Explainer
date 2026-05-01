package explain

var mimirRules = []Rule{
	// ── block storage ─────────────────────────────────────────────────────────
	{
		PathContains: []string{"blocks_storage", "bucket", "storage"},
		Risk:         "high",
		Impact: []string{
			"Mimir block storage configuration changed.",
			"Metrics durability and query availability may be affected.",
		},
	},

	// ── limits / ingestion ────────────────────────────────────────────────────
	{
		PathContains: []string{"limits", "ingestion", "max_global_series_per_user", "ingestion_rate"},
		Risk:         "medium",
		Impact: []string{
			"Mimir limits or ingestion behavior changed.",
			"Tenant ingestion or cardinality behavior may be affected.",
		},
	},

	// ── compactor ─────────────────────────────────────────────────────────────
	{
		PathContains: []string{"compactor"},
		Risk:         "medium",
		Impact: []string{
			"Mimir compactor behavior changed.",
			"Block compaction or retention may be affected.",
		},
	},
}
