package cmd

import (
	_loacl "github.com/licat233/goss/modules/local"
	"github.com/licat233/goss/utils"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:     "local",
	GroupID: "modules",
	Short:   "checkout: " + _loacl.Name,
}

var localStartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "run " + _loacl.Name,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			utils.Success("done.")
		}()
		run(_loacl.Run, _loacl.CheckoutFileExt)
	},
}

func init() {
	localCmd.AddCommand(htmlStartCmd)
	localCmd.SetHelpTemplate(setColorizeHelp(htmlCmd.HelpTemplate()))
}
