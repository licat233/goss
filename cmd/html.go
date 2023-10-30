package cmd

import (
	"fmt"

	"github.com/licat233/goss/config"
	_html "github.com/licat233/goss/modules/html"
	"github.com/licat233/goss/utils"
	"github.com/spf13/cobra"
)

var htmlCmd = &cobra.Command{
	Use:     "html",
	GroupID: "modules",
	Short:   "checkout: " + _html.Name,
}

var htmlStartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "run " + _html.Name,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			utils.Success("done.")
		}()
		run(_html.Run, _html.CheckoutFileExt)
	},
}

func init() {
	htmlCmd.PersistentFlags().StringSliceVar(&config.HtmlTags, "tags", []string{"*"}, fmt.Sprintf("Select the HTML tags to process,current support: %s", _html.SupportTags))
	htmlCmd.AddCommand(htmlStartCmd)
	htmlCmd.SetHelpTemplate(setColorizeHelp(htmlCmd.HelpTemplate()))
}
