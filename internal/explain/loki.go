package explain

var lokiRules = []Rule{
	// ── retention ─────────────────────────────────────────────────────────────
	{
		PathContains: []string{"retention", "retention_period", "limits_config.retention_period"},
		Condition:    isIncreased,
		Risk:         "medium",
		Impact: []string{
			"Loki log retention increased.",
			"Object storage usage may increase.",
			"Query range and storage cost may increase.",
		},
	},
	{
		PathContains: []string{"retention", "retention_period", "limits_config.retention_period"},
		Condition:    isDecreased,
		Risk:         "high",
		Impact: []string{
			"Loki log retention decreased.",
			"Older logs may expire sooner.",
			"Confirm compliance and troubleshooting requirements.",
		},
	},
	{
		PathContains: []string{"retention_period", "limits_config.retention_period"},
		Risk:         "medium",
		Impact: []string{
			"Loki log retention configuration changed.",
			"Object storage usage and log availability may be affected.",
		},
	},

	// ── schema ────────────────────────────────────────────────────────────────
	{
		PathContains: []string{"schemaConfig", "schema_config"},
		Risk:         "high",
		Impact: []string{
			"Loki schema configuration changed.",
			"Incorrect schema changes may break ingestion or querying.",
			"Review Loki migration requirements carefully.",
		},
	},

	// ── object storage ────────────────────────────────────────────────────────
	{
		PathContains: []string{"object_store", "bucket"},
		Risk:         "high",
		Impact: []string{
			"Loki object storage configuration changed.",
			"Data availability may be affected if misconfigured.",
		},
	},
	{
		PathContains: []string{"storage"},
		Risk:         "high",
		Impact: []string{
			"Loki object storage configuration changed.",
			"Data availability may be affected if misconfigured.",
		},
	},

	// ── compactor ─────────────────────────────────────────────────────────────
	{
		PathContains: []string{"compactor"},
		Risk:         "medium",
		Impact: []string{
			"Loki compactor settings changed.",
			"Retention and deletion behavior may be affected.",
		},
	},
}
