package explain

var genericRules = []Rule{
	// ── replicas ──────────────────────────────────────────────────────────────
	{
		PathContains: []string{"replicas", "replicaCount"},
		Condition:    isIncreased,
		Risk:         "low",
		Impact: []string{
			"Workload replica count increased.",
			"Resource usage may increase.",
			"Availability may improve.",
		},
	},
	{
		PathContains: []string{"replicas", "replicaCount"},
		Condition:    isDecreased,
		Risk:         "medium",
		Impact: []string{
			"Workload replica count decreased.",
			"Availability or redundancy may be reduced.",
		},
	},
	{
		PathContains: []string{"replicas", "replicaCount"},
		Risk:         "low",
		Impact: []string{
			"Workload replica count changed.",
			"Resource usage and availability may be affected.",
		},
	},

	// ── resources ─────────────────────────────────────────────────────────────
	{
		PathContains: []string{"resources.requests.cpu", "resources.requests.memory", "resources.limits.cpu", "resources.limits.memory"},
		Condition:    isIncreased,
		Risk:         "low",
		Impact: []string{
			"Resource allocation increased.",
			"Scheduling requirements may increase.",
			"Cluster capacity usage may increase.",
		},
	},
	{
		PathContains: []string{"resources.requests.cpu", "resources.requests.memory", "resources.limits.cpu", "resources.limits.memory"},
		Condition:    isDecreased,
		Risk:         "medium",
		Impact: []string{
			"Resource allocation decreased.",
			"Risk of throttling or OOM may increase if the value is too low.",
		},
	},
	{
		PathContains: []string{"resources.requests.", "resources.limits.", "resources"},
		Risk:         "low",
		Impact: []string{
			"Resource allocation changed.",
			"Review scheduling and cluster capacity requirements.",
		},
	},

	// ── persistence / storage ─────────────────────────────────────────────────
	{
		PathContains: []string{"persistence", "storageClass", "storage", "size"},
		Risk:         "medium",
		Impact: []string{
			"Storage-related configuration changed.",
			"Persistent data behavior or storage usage may be affected.",
			"Review PVC and retention implications before deployment.",
		},
	},

	// ── service type ──────────────────────────────────────────────────────────
	{
		PathContains: []string{"service.type"},
		Condition:    newValueIs("LoadBalancer"),
		Risk:         "medium",
		Impact: []string{
			"Service may become externally reachable depending on cloud/network configuration.",
			"Cloud load balancer cost may be introduced.",
		},
	},
	{
		PathContains: []string{"service.type"},
		Condition:    oldValueIs("LoadBalancer"),
		Risk:         "medium",
		Impact: []string{
			"External access behavior may change.",
			"Existing clients may lose access if they depend on the previous service type.",
		},
	},
	{
		PathContains: []string{"service.type"},
		Risk:         "medium",
		Impact: []string{
			"Service type changed.",
			"External access behavior may change.",
		},
	},

	// ── ingress ───────────────────────────────────────────────────────────────
	{
		PathContains: []string{"ingress"},
		Risk:         "medium",
		Impact: []string{
			"Ingress routing or external exposure changed.",
			"Review host, TLS, and annotation changes.",
		},
	},

	// ── security context ──────────────────────────────────────────────────────
	{
		PathContains: []string{"privileged", "allowPrivilegeEscalation"},
		Condition:    newValueBoolTrue,
		Risk:         "high",
		Impact: []string{
			"Container privilege level increased.",
			"Security risk may increase.",
		},
	},
	{
		PathContains: []string{"runAsUser", "runAsNonRoot"},
		Risk:         "medium",
		Impact: []string{
			"Runtime user configuration changed.",
			"File permissions or security posture may be affected.",
		},
	},
	{
		PathContains: []string{"securityContext", "podSecurityContext"},
		Risk:         "medium",
		Impact: []string{
			"Runtime user configuration changed.",
			"File permissions or security posture may be affected.",
		},
	},

	// ── image ─────────────────────────────────────────────────────────────────
	{
		PathContains: []string{"image.tag"},
		Condition:    newValueIs("latest"),
		Risk:         "high",
		Impact: []string{
			"Floating image tag detected.",
			"Deployments may become non-deterministic.",
		},
	},
	{
		PathContains: []string{"image.repository", "image.tag", "image.digest"},
		Risk:         "medium",
		Impact: []string{
			"Container image changed.",
			"Runtime behavior may change.",
			"Review release notes and compatibility.",
		},
	},

	// ── environment variables ─────────────────────────────────────────────────
	{
		PathContains: []string{"extraEnv", "env"},
		Risk:         "low",
		Impact: []string{
			"Environment configuration changed.",
			"Application behavior may change depending on the variable.",
		},
	},

	// ── config ────────────────────────────────────────────────────────────────
	{
		PathContains: []string{"extraConfig", "structuredConfig", "config"},
		Risk:         "medium",
		Impact: []string{
			"Application configuration changed.",
			"Review chart-specific behavior.",
		},
	},

	// ── scheduling ────────────────────────────────────────────────────────────
	{
		PathContains: []string{"tolerations", "nodeSelector", "affinity"},
		Risk:         "medium",
		Impact: []string{
			"Scheduling behavior changed.",
			"Workloads may move to different nodes or fail to schedule.",
		},
	},
}
