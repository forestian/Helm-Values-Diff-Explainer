package explain

var grafanaRules = []Rule{
	{
		PathContains: []string{"adminPassword", "admin", "datasources", "dashboardProviders", "ingress"},
		Risk:         "medium",
		Impact: []string{
			"Grafana access, datasource, or dashboard configuration changed.",
			"Users may see different dashboards or data sources.",
		},
	},
}
