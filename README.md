# hvediff — Helm Values Diff Explainer

`hvediff` compares two Helm values YAML files and explains the differences in plain English with operational impact hints.

## Why this tool exists

Plain YAML diffs are hard to review before a Helm deployment. A line like `-retention_period: 168h / +retention_period: 720h` tells you a value changed, but not that:

- Log storage usage will grow 4×
- Older queries will scan more data
- You may have compliance implications

`hvediff` adds that context — deterministically, without AI, and without any network calls.

## Who it is for

DevOps engineers, SREs, Kubernetes operators, platform engineers, and anyone who reviews Helm values changes before deployment.

---

## Installation

### Build from source

```bash
git clone <repo-url> helm-values-diff-explainer
cd helm-values-diff-explainer
go build -o hvediff .
```

Or run without installing:

```bash
go run . diff --old values-old.yaml --new values-new.yaml
```

### Install from GitHub Releases

Download a prebuilt binary from the [GitHub Releases page](https://github.com/forestian/Helm-Values-Diff-Explainer/releases).

**Linux / macOS:**
```bash
tar -xzf hvediff_<version>_<os>_<arch>.tar.gz
chmod +x hvediff
./hvediff version
```

**Windows:**
```powershell
# Extract the archive, then:
.\hvediff.exe version
```

### Requirements

- Go 1.22 or later (build from source only)

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

#### Flags

| Flag | Default | Description |
|---|---|---|
| `--old` | — | Path to old values YAML file **(required)** |
| `--new` | — | Path to new values YAML file **(required)** |
| `--format` | `text` | Output format: `text`, `json`, `markdown` |
| `--component` | `generic` | Component hint for richer impact messages |
| `--output` | stdout | Write output to a file instead of stdout |
| `--fail-on-risk` | `none` | Exit non-zero when changes reach a risk threshold |

#### Supported components

`generic` · `loki` · `mimir` · `tempo` · `alloy` · `prometheus` · `grafana`

#### Fail-on-risk levels

| Value | Exits non-zero when... |
|---|---|
| `none` | Never |
| `high` | Any `high` risk change is present |
| `medium` | Any `medium` or `high` risk change is present |
| `low` | Any `low`, `medium`, or `high` risk change is present |

The report is always printed/written **before** exiting non-zero.

---

## Usage examples

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

## Output examples

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

- **Path matching is substring-based.** A path like `limits_config` will match rules that look for the keyword `config`. Component-specific rules are applied first to reduce false positives.
- **Kubernetes resource strings** (`100m`, `256Mi`) cannot be compared as numbers. The direction-aware rules (increased/decreased) will not fire for these values; a generic resource-changed message is shown instead.
- **List diffs are index-based.** If a list is reordered, every moved element will appear as changed rather than as a move operation.
- **No Helm execution.** `hvediff` reads raw YAML files and does not template, merge, or lint Helm charts.
- **No schema validation.** The tool does not know the schema of any specific Helm chart.

---

## Roadmap (not implemented)

- GitHub Actions integration
- GitHub PR comment bot
- ArgoCD diff integration
- Helm chart-aware schema validation
- AI-generated natural-language explanations
- Web UI / SaaS backend

---

## License

MIT
