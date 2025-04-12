package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jaytaylor/html2text"
	"github.com/spf13/cobra"

	"github.com/ccbrown/cloud-snitch/backend/app"
)

var renderTemplateCmd = &cobra.Command{
	Use:   "render-template",
	Short: "renders a template",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := app.New(rootConfig.App)
		if err != nil {
			return err
		}

		template, _ := cmd.Flags().GetString("template")

		var params map[string]any
		if p, _ := cmd.Flags().GetString("params"); p != "" {
			if err := json.Unmarshal([]byte(p), &params); err != nil {
				return fmt.Errorf("failed to unmarshal params: %w", err)
			}
		}

		html, err := a.RenderTemplate(template, params)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("format")
		switch format {
		case "text":
			text, err := html2text.FromString(html)
			if err != nil {
				return err
			}
			os.Stdout.WriteString(text)
		case "html":
			os.Stdout.WriteString(html)
		default:
			return fmt.Errorf("unknown format %q", format)
		}

		return nil
	},
}

func init() {
	renderTemplateCmd.Flags().StringP("template", "t", "", "the template to render")
	renderTemplateCmd.MarkFlagRequired("template")

	renderTemplateCmd.Flags().StringP("format", "f", "html", "the format to render the template in (html, text)")

	renderTemplateCmd.Flags().StringP("params", "p", "", "the parameters to pass to the template as json")

	rootCmd.AddCommand(renderTemplateCmd)
}
