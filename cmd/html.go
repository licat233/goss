package cmd

import (
	"fmt"
	"strings"

	"github.com/licat233/goss/config"
	_html "github.com/licat233/goss/modules/html"
	"github.com/licat233/goss/utils"
	"github.com/spf13/cobra"
)

var htmlCmd = &cobra.Command{
	Use:     "html",
	GroupID: "modules",
	Short:   _html.Name,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			utils.Success("done.")
		}()
		run(_html.Run, _html.CheckoutFileExt)
	},
}

func init() {
	htmlCmd.PersistentFlags().StringSliceVar(&config.HtmlTags, "tags", []string{"*"}, fmt.Sprintf("Select the HTML tags to process,current support: %s", _html.SupportTags))
	// htmlCmd.AddCommand(htmlStartCmd)
	htmlCmd.SetHelpTemplate(setColorizeHelp(htmlCmd.HelpTemplate()))
	rootCmd.AddCommand(htmlCmd)
}

func checkHtmlConfig() error {
	if len(config.HtmlTags) != 0 {
		//如果用户指定了，则以用户的为准
		for index := range config.HtmlTags {
			tag := config.HtmlTags[index]
			if tag == "*" {
				config.HtmlTags = _html.SupportTags
				break
			}
			config.HtmlTags[index] = strings.ToLower(tag)
		}
	} else {
		//否则以系统默认的为准
		config.HtmlTags = _html.SupportTags
	}
	return nil
}
