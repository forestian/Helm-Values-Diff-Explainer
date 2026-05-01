package explain

var alloyRules = []Rule{
	// ── export destinations ───────────────────────────────────────────────────
	{
		PathContains: []string{"remote_write", "loki.write", "otelcol.exporter"},
		Risk:         "high",
		Impact: []string{
			"Telemetry export destination changed.",
			"Data may stop reaching the backend if misconfigured.",
		},
	},

	// ── runtime config ────────────────────────────────────────────────────────
	{
		PathContains: []string{"alloy.config", "configMap", "controller"},
		Risk:         "medium",
		Impact: []string{
			"Alloy runtime configuration changed.",
			"Metrics, logs, or traces pipeline behavior may change.",
		},
	},
}
