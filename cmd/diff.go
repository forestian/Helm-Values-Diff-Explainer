package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"helm-values-diff-explainer/internal/diff"
	"helm-values-diff-explainer/internal/explain"
	"helm-values-diff-explainer/internal/parser"
	"helm-values-diff-explainer/internal/report"
	"helm-values-diff-explainer/internal/risk"
)

var (
	flagOld        string
	flagNew        string
	flagFormat     string
	flagComponent  string
	flagOutput     string
	flagFailOnRisk string
)

var diffCmd = &cobra.Command{
	Use:          "diff",
	Short:        "Compare two Helm values YAML files",
	SilenceUsage: true,
	RunE:         runDiff,
}

func init() {
	diffCmd.Flags().StringVar(&flagOld, "old", "", "Path to old values YAML file (required)")
	diffCmd.Flags().StringVar(&flagNew, "new", "", "Path to new values YAML file (required)")
	diffCmd.Flags().StringVar(&flagFormat, "format", "text", "Output format: text, json, markdown")
	diffCmd.Flags().StringVar(&flagComponent, "component", "generic", "Component: generic, loki, mimir, tempo, alloy, prometheus, grafana")
	diffCmd.Flags().StringVar(&flagOutput, "output", "", "Write output to file instead of stdout")
	diffCmd.Flags().StringVar(&flagFailOnRisk, "fail-on-risk", "none", "Exit non-zero if risk threshold met: none, low, medium, high")
}

func runDiff(cmd *cobra.Command, args []string) error {
	if flagOld == "" {
		return fmt.Errorf("--old is required")
	}
	if flagNew == "" {
		return fmt.Errorf("--new is required")
	}

	validFormats := map[string]bool{"text": true, "json": true, "markdown": true}
	if !validFormats[flagFormat] {
		return fmt.Errorf("--format must be text, json, or markdown")
	}

	validComponents := map[string]bool{
		"generic": true, "loki": true, "mimir": true, "tempo": true,
		"alloy": true, "prometheus": true, "grafana": true,
	}
	if !validComponents[flagComponent] {
		return fmt.Errorf("--component must be one of: generic, loki, mimir, tempo, alloy, prometheus, grafana")
	}

	validRisks := map[string]bool{"none": true, "low": true, "medium": true, "high": true}
	if !validRisks[flagFailOnRisk] {
		return fmt.Errorf("--fail-on-risk must be none, low, medium, or high")
	}

	oldValues, err := parser.ParseFile(flagOld)
	if err != nil {
		return fmt.Errorf("error reading --old file: %w", err)
	}
	newValues, err := parser.ParseFile(flagNew)
	if err != nil {
		return fmt.Errorf("error reading --new file: %w", err)
	}

	changes := diff.Compare(oldValues, newValues)

	explainer := explain.NewExplainer(flagComponent)
	for i := range changes {
		explainer.Explain(&changes[i])
	}

	r := &report.Report{
		OldFile:   flagOld,
		NewFile:   flagNew,
		Component: flagComponent,
		Changes:   changes,
	}

	var content string
	switch flagFormat {
	case "text":
		content = report.RenderText(r)
	case "json":
		content, err = report.RenderJSON(r)
		if err != nil {
			return fmt.Errorf("error generating JSON report: %w", err)
		}
	case "markdown":
		content = report.RenderMarkdown(r)
	}

	if flagOutput != "" {
		if err := os.WriteFile(flagOutput, []byte(content), 0o644); err != nil {
			return fmt.Errorf("error writing output file: %w", err)
		}
	} else {
		fmt.Print(content)
	}

	if flagFailOnRisk != "none" {
		riskLevels := make([]string, len(changes))
		for i, c := range changes {
			riskLevels[i] = c.Risk
		}
		if risk.ShouldFail(flagFailOnRisk, riskLevels) {
			os.Exit(1)
		}
	}

	return nil
}
