# hvediff — Helm Values Diff Explainer

Explain Helm values changes before they surprise production.

## Quick Demo

```bash
hvediff diff --old examples/values-old.yaml --new examples/values-new.yaml --component loki
```

```
Helm Values Diff Explainer

Old: examples/values-old.yaml
New: examples/values-new.yaml
Component: loki

Summary:
- Added:       4
- Removed:     0
- Changed:     13
- High risk:   3
- Medium risk: 4
- Low risk:    6

Changes:

[HIGH] [CHANGED] loki.schemaConfig.configs[0].schema
old: v11
new: v12

Impact:
- Loki schema configuration changed.
- Incorrect schema changes may break ingestion or querying.
- Review Loki migration requirements carefully.

[HIGH] [CHANGED] loki.schemaConfig.configs[0].object_store
old: filesystem
new: s3

Impact:
- Loki object storage backend changed.
- Ensure the new backend is accessible before deploying.
- Data migration may be required.

[HIGH] [CHANGED] loki.compactor.retention_enabled
old: false
new: true

Impact:
- Loki compactor retention setting changed.
- Enabling retention may affect stored log volume.
- Ensure retention rules are configured before enabling.

[MEDIUM] [CHANGED] loki.limits_config.retention_period
old: 168h
new: 720h

Impact:
- Loki log retention increased.
- Object storage usage may increase.
- Query range and storage cost may increase.

[MEDIUM] [CHANGED] service.type
old: ClusterIP
new: LoadBalancer

Impact:
- Kubernetes Service type changed.
- Changing to LoadBalancer may expose the service externally.
- Verify firewall and network policies.

... (8 more changes)
```

## Demo

GIF demo coming soon.

---

## Quick Start

**Download a prebuilt binary** from [GitHub Releases](https://github.com/forestian/Helm-Values-Diff-Explainer/releases):

```bash
# Linux / macOS
tar -xzf hvediff_<version>_<os>_<arch>.tar.gz
chmod +x hvediff
./hvediff version

# Windows
.\hvediff.exe version
```

**Build from source** (requires Go 1.22+):

```bash
git clone https://github.com/forestian/Helm-Values-Diff-Explainer.git
cd Helm-Values-Diff-Explainer
go build -o hvediff .
```

---

## Use Cases

- Catch high-risk config changes (schema migrations, storage backend switches) before they reach the cluster
- Fail a CI pipeline if changes exceed a risk threshold (`--fail-on-risk high`)
- Generate Markdown diff reports for automated PR comments (`--format markdown --output report.md`)
- Review Loki, Mimir, Prometheus, or Grafana values changes with component-aware context
- Audit what changed between environments without reading raw YAML diffs

---

## Why this tool exists

Plain YAML diffs are hard to review before a Helm deployment. A line like `-retention_period: 168h / +retention_period: 720h` tells you a value changed, but not that:

- Log storage usage will grow 4×
- Older queries will scan more data
- You may have compliance implications

`hvediff` adds that context — deterministically, without AI, and without any network calls.

**Who it is for:** DevOps engineers, SREs, Kubernetes operators, platform engineers, and anyone who reviews Helm values changes before deployment.

---

## Commands

### `hvediff version`

```
hvediff version
# hvediff version 0.1.0
```

### `hvediff diff`

```
hvediff diff --old <path> --new <path> [flags]
```

| Flag | Default | Description |
|---|---|---|
| `--old` | — | Path to old values YAML file **(required)** |
| `--new` | — | Path to new values YAML file **(required)** |
| `--format` | `text` | Output format: `text`, `json`, `markdown` |
| `--component` | `generic` | Component hint for richer impact messages |
| `--output` | stdout | Write output to a file instead of stdout |
| `--fail-on-risk` | `none` | Exit non-zero when changes reach a risk threshold |

**Supported components:** `generic` · `loki` · `mimir` · `tempo` · `alloy` · `prometheus` · `grafana`

**Fail-on-risk levels:**

| Value | Exits non-zero when... |
|---|---|
| `none` | Never |
| `high` | Any `high` risk change is present |
| `medium` | Any `medium` or `high` risk change is present |
| `low` | Any `low`, `medium`, or `high` risk change is present |

The report is always printed/written **before** exiting non-zero.

---

## Usage Examples

```bash
# Basic diff (text output to stdout)
hvediff diff --old values-old.yaml --new values-new.yaml

# Loki-specific impact messages
hvediff diff --old values-old.yaml --new values-new.yaml --component loki

# Write a Markdown report for a PR comment
hvediff diff --old values-old.yaml --new values-new.yaml --format markdown --output report.md

# JSON report
hvediff diff --old values-old.yaml --new values-new.yaml --format json

# Fail CI if high-risk changes are present
hvediff diff --old values-old.yaml --new values-new.yaml --component loki --fail-on-risk high
```

---

## Example Output

### Text

```
Helm Values Diff Explainer

Old: values-old.yaml
New: values-new.yaml
Component: loki

Summary:
- Added: 1
- Removed: 0
- Changed: 3
- High risk: 1
- Medium risk: 2
- Low risk: 1

Changes:

[HIGH] [CHANGED] loki.schemaConfig.configs[0].schema
old: v11
new: v12

Impact:
- Loki schema configuration changed.
- Incorrect schema changes may break ingestion or querying.
- Review Loki migration requirements carefully.

[MEDIUM] [CHANGED] loki.limits_config.retention_period
old: 168h
new: 720h

Impact:
- Loki log retention increased.
- Object storage usage may increase.
- Query range and storage cost may increase.
```

### JSON

```json
{
  "old_file": "values-old.yaml",
  "new_file": "values-new.yaml",
  "component": "loki",
  "summary": {
    "added": 1,
    "removed": 0,
    "changed": 3,
    "high_risk": 1,
    "medium_risk": 2,
    "low_risk": 1
  },
  "changes": [
    {
      "type": "CHANGED",
      "path": "loki.schemaConfig.configs[0].schema",
      "old_value": "v11",
      "new_value": "v12",
      "risk": "high",
      "impact": [
        "Loki schema configuration changed.",
        "Incorrect schema changes may break ingestion or querying.",
        "Review Loki migration requirements carefully."
      ]
    }
  ]
}
```

### Markdown

```markdown
# Helm Values Diff Explainer

**Old:** `values-old.yaml`
**New:** `values-new.yaml`
**Component:** `loki`

## Summary

| Type | Count |
|---|---:|
| Added | 1 |
| Removed | 0 |
| Changed | 3 |
| High risk | 1 |
| Medium risk | 2 |
| Low risk | 1 |

## Changes

### HIGH / CHANGED / `loki.schemaConfig.configs[0].schema`

**Old:** `v11`
**New:** `v12`

**Impact:**

- Loki schema configuration changed.
- Incorrect schema changes may break ingestion or querying.
- Review Loki migration requirements carefully.
```

---

## Component-specific modes

When you pass `--component <name>`, `hvediff` applies richer rules tuned to that component.

| Component | Key paths monitored with extra context |
|---|---|
| `loki` | `retention_period`, `schemaConfig`, `storage`, `compactor` |
| `mimir` | `blocks_storage`, `limits`, `ingestion_rate`, `compactor` |
| `tempo` | `storage`, `traces`, `retention` |
| `alloy` | `alloy.config`, `remote_write`, `otelcol.exporter` |
| `prometheus` | `scrape_configs`, `remote_write`, `rule_files`, `alerting` |
| `grafana` | `adminPassword`, `datasources`, `dashboardProviders` |

When no component is given (`generic`), only general Kubernetes/Helm impact rules apply (replicas, resources, service type, ingress, image, security context, etc.).

---

## Risk levels

| Level | Meaning |
|---|---|
| `high` | Change may cause data loss, outage, or security regression |
| `medium` | Change affects availability, access, or resource use — review carefully |
| `low` | Minor change with limited blast radius |
| `none` | No matching impact rule found |

---

## Path notation

Nested keys use dot notation. List indexes use bracket notation.

```
server.replicas
loki.limits_config.retention_period
resources.requests.cpu
ingress.hosts[0].host
extraEnv[1].name
```

---

## Limitations

- **Read-only and local-only.** `hvediff` reads YAML files and prints a report. It does not connect to any cluster, API, or external service.
- **Path matching is substring-based.** A path like `limits_config` will match rules that look for the keyword `config`. Component-specific rules are applied first to reduce false positives.
- **Kubernetes resource strings** (`100m`, `256Mi`) cannot be compared as numbers. Direction-aware rules (increased/decreased) will not fire for these values; a generic resource-changed message is shown instead.
- **List diffs are index-based.** If a list is reordered, every moved element will appear as changed rather than as a move operation.
- **No Helm execution.** `hvediff` reads raw YAML files and does not template, merge, or lint Helm charts.
- **No schema validation.** The tool does not know the schema of any specific Helm chart.
- **Generated output is not a substitute for review.** Always verify the report before applying changes to production.

---

## Roadmap

- GitHub Actions integration
- GitHub PR comment bot
- ArgoCD diff integration
- Support for more observability components
- Helm chart-aware schema validation

---

*Part of the [Forestian Cloud Native Toolkit](https://github.com/forestian) — small CLI tools for Kubernetes, observability, GitOps, and platform engineering.*

## License

MIT
