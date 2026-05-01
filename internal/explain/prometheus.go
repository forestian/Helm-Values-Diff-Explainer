package explain

var prometheusRules = []Rule{
	{
		PathContains: []string{"scrape_configs", "remote_write", "rule_files", "alerting"},
		Risk:         "high",
		Impact: []string{
			"Prometheus collection, remote write, or alerting configuration changed.",
			"Metrics ingestion or alert delivery may be affected.",
		},
	},
}
