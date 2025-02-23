package sharedconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	FormatJSON = "json"
	FormatText = "text"
)

type FormatConfig struct {
	Output string `koanf:"output" short:"o" description:"output format, one of 'json' or 'text'"`
}

func NewFormatConfig() FormatConfig {
	return FormatConfig{
		Output: FormatText,
	}
}

func (cfg *FormatConfig) Validate() error {
	allowed := []string{FormatJSON, FormatText}
	lOutput := strings.ToLower(cfg.Output)
	if !isOneOf(lOutput, allowed...) {
		return fmt.Errorf("invalid output format %q: must be one of %v", cfg.Output, allowed)
	}
	cfg.Output = lOutput
	return nil
}

func isOneOf(s string, values ...string) bool {
	for _, v := range values {
		if s == v {
			return true
		}
	}
	return false
}

func (cfg *FormatConfig) Print(cmd *cobra.Command, a any) error {
	switch cfg.Output {
	case FormatText:
		return printText(cmd, a)
	case FormatJSON:
		return printJSON(cmd, a)
	default:
		// should never happen
		return fmt.Errorf("unsupported output format: %s", cfg.Output)
	}
}

func printText(cmd *cobra.Command, a any) error {
	s := a.(fmt.Stringer) // will panic if used incorrectly
	_, err := fmt.Fprintln(cmd.OutOrStdout(), s.String())
	return err
}

func printJSON(cmd *cobra.Command, a any) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json result: %w", err)
	}

	_, err = cmd.OutOrStdout().Write(data)
	if err != nil {
		return fmt.Errorf("failed to print json result: %w", err)
	}
	fmt.Fprint(cmd.OutOrStdout(), "\n")
	return nil
}
