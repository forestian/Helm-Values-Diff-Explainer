package explain

var tempoRules = []Rule{
	// ── trace storage ─────────────────────────────────────────────────────────
	{
		PathContains: []string{"storage", "bucket", "traces"},
		Risk:         "high",
		Impact: []string{
			"Tempo trace storage configuration changed.",
			"Trace durability or query behavior may be affected.",
		},
	},

	// ── retention ─────────────────────────────────────────────────────────────
	{
		PathContains: []string{"retention"},
		Risk:         "medium",
		Impact: []string{
			"Tempo trace retention changed.",
			"Storage usage and trace availability may change.",
		},
	},
}
